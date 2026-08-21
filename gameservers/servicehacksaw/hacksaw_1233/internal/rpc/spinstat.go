package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"runtime"
	"serve/comm/jwtutil"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicehacksaw/hacksawcomm"
	"strconv"
	"sync"
)

func init() {
	hacksawcomm.RegRpc("rtp", RTP)
	hacksawcomm.RegRpc("rtp_count", RTP)
}

const fakeAppUrl = "http://127.0.0.1:27000/GetLink"

type GetlinkPs struct {
	Bet        string `json:"bet"`
	UserID     string `json:"user_id"`
	GameID     string `json:"game_id"`
	Language   string `json:"language"`
	GamePreFix string `json:"game_pre_fix"`
	Count      int
}

func GoRTP(ps *GetlinkPs, ret *ut.SpinStatRet) {
	numThreads := runtime.NumCPU()
	runtime.GOMAXPROCS(numThreads)
	perThread := ps.Count / numThreads
	remainder := ps.Count % numThreads

	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			spinNum := perThread
			if threadID == numThreads-1 {
				spinNum += remainder
			}
			fmt.Printf("协程 id：%d ,分配次数：%d\n", threadID, spinNum)

			var token struct {
				Token string `json:"token"`
			}
			err := mux.HttpInvoke(fakeAppUrl, GetlinkPs{
				UserID:     ps.UserID + strconv.Itoa(threadID),
				GameID:     ps.GameID,
				Language:   "en",
				GamePreFix: ps.GamePreFix,
			}, &token)
			if err != nil {
				fmt.Println(err)
				return
			}

			pid, _, _ := jwtutil.ParseTokenData(token.Token)
			ret.Balance_before, _ = slotsmongo.GetBalance(pid)

			data := hacksawcomm.CreateSpinTest(ps.Bet)
			binary, _ := data.ToBinary()
			marshal, _ := json.Marshal(hacksawcomm.GateWayData{
				Pid:  pid,
				Data: binary,
			})
			spinMsg := &nats.Msg{
				Data: marshal,
				Header: nats.Header{
					"isstat":   []string{"1"},
					"jump_log": []string{"0"},
				},
			}

			for j := 0; j < spinNum; j++ {
				_, err = doSpin(spinMsg)
				if err != nil {
					fmt.Println(err)
					return
				}

				bet, _ := strconv.Atoi(spinMsg.Header.Get("stat_bet"))
				ret.Bet += int64(bet)
				win, _ := strconv.Atoi(spinMsg.Header.Get("stat_win"))
				ret.Win += int64(win)
			}

		}(i)
	}
	wg.Wait()
}

func RTP(msg *nats.Msg) ([]byte, error) {
	var Data struct {
		Pid      int
		BuyEX    int
		Count    int
		BuyBonus bool
		Bet      string
		GameId   string
		GameName string
		UName    string
	}
	err := json.Unmarshal(msg.Data, &Data)
	if err != nil {
		return nil, err
	}

	rets := ut.SpinStatRet{}
	GoRTP(&GetlinkPs{
		Bet:        Data.Bet,
		UserID:     Data.UName,
		GameID:     Data.GameId,
		GamePreFix: Data.GameName,
		Count:      Data.Count,
	}, &rets)

	marshal, err := json.Marshal(rets)
	if err != nil {
		return nil, err
	}

	return marshal, nil
}
