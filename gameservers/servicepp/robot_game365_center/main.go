package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"
	"serve/comm/mux"
	"strings"
	"sync"
	"time"
)

func main() {
	if len(os.Args) > 1 {
		PGHack(10000, "tester0114", os.Args[1])
	}
	time.Sleep(5 * time.Second)
}

type Token struct {
	Token string `json:"token"`
}

const fakeAppUrl = "http://127.0.0.1:27000/GetLink"

const PGReqSpinUrl = "https://api.ppply.site/game-api/%v/v2/Spin"
const PGReqInfoUrl = "https://api.ppply.site/game-api/%v/v2/GameInfo/Get"

type GetlinkPs struct {
	UserID     string `json:"user_id"`
	GameID     string `json:"game_id"`
	Language   string `json:"language"`
	GamePreFix string `json:"game_pre_fix"`
	BuyEX      int    `json:"BuyEX"`
}

// // RoundC 随机取权重，下注
//
//	func RoundC() {
//		allWeight := 0
//		for _, v := range ppr.arrC {
//			allWeight += v
//		}
//		randNum := rand.Int() % allWeight
//		for k, v := range ppr.arrC {
//			randNum -= v
//			if randNum < 0 {
//				ppr.c = k
//				break
//			}
//		}
//	}

func genTraceId() string {
	// return "VWDLLC27"
	var buf [6]byte
	for i := 0; i < len(buf); i++ {
		c := rand.IntN(26) + 'A'
		buf[i] = byte(c)
	}

	return string(buf[:]) + "27"
}

func PGHack(count int, userName, GameId string) {
	numThreads := runtime.NumCPU() // 获取 CPU 核心数
	runtime.GOMAXPROCS(numThreads) // 设置最大可同时使用的 CPU 核心数
	var wg sync.WaitGroup
	//numThreads = 1
	// 每个协程处理的迭代次数
	perThread := count / numThreads
	remainder := count % numThreads

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			if userName == "" {
				userName = "test_user"
			}
			uid := fmt.Sprintf("%s_%d", userName, threadID)
			spinNum := perThread
			if threadID == numThreads-1 {
				spinNum += remainder
			}
			fmt.Printf("协程 id：%d ,分配次数：%d\n", threadID, spinNum)
			var err error
			var token Token
			var getlinkPs GetlinkPs
			//str := "action=doSpin&symbol=vs20olympx&c=%d&l=20&bl=0&index=5&counter=183&repeat=0&mgckey=eyJQIjoxMDAyMjYsIkUiOjE3MzI2Mjk1MTAsIlMiOjEwMDAsIkQiOiI2NSJ9.kGBj8OA4x1r9fKaqrwqD6pLNFRIGS8uAYFKGsVhkDNs"
			//str = fmt.Sprintf(str, ps.Bet)
			//bytes := []byte(str)
			//cs := Cs[0]
			//ml := Ml[0]
			getlinkPs.UserID = uid
			getlinkPs.GameID = GameId
			getlinkPs.Language = "en"
			getlinkPs.GamePreFix = ""
			//不足1000000 自动设置
			err = mux.HttpInvoke(fakeAppUrl, getlinkPs, &token)
			if err != nil {
				slog.Error("spin_stat_err", "HttpInvoke", err)
				return
			}
			//token换pid

			//请求info
			sid, cs := PGInfo(strings.TrimPrefix(GameId, "pg_"), token.Token)
			for i := 0; i < spinNum; i++ {
				//ret.Count++
				for {
					sid = pgSpin(strings.TrimPrefix(GameId, "pg_"), token.Token, sid, cs)

					sid1 := strings.Split(sid, "_")
					if len(sid1) == 1 || sid1[1] == sid1[2] { // 表示是最后一盘
						break
					}
				}
			}

		}(i)
	}

	wg.Wait()
}

func pgSpin(GameId, token, lastSid string, cs float64) string {
	payload := strings.NewReader(fmt.Sprintf("cs=%v&ml=1&pf=1&id=%s&wk=0_C&btt=1&atk=%s", cs, lastSid, token))

	client := &http.Client{}

	SpinUrl := fmt.Sprintf(PGReqSpinUrl, GameId)
	SpinUrl = SpinUrl + "?" + genTraceId()
	req, err := http.NewRequest(http.MethodPost, SpinUrl, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("origin", "https://m.ppply.site")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://m.ppply.site/")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	type Response struct {
		DT struct {
			SI struct {
				Sid string `json:"sid"`
			} `json:"si"`
		} `json:"dt"`
	}
	fmt.Println(string(body))
	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	// 获取 sid
	sid := resp.DT.SI.Sid
	fmt.Println("SID:", sid)
	return sid
}

func PGInfo(GameId, token string) (string, float64) {

	payload := strings.NewReader(fmt.Sprintf("btt=1&atk=%s&pf=1", token))

	client := &http.Client{}
	InfoUrl := fmt.Sprintf(PGReqInfoUrl, GameId)
	InfoUrl = InfoUrl + "?" + genTraceId()
	req, err := http.NewRequest(http.MethodPost, InfoUrl, payload)
	if err != nil {
		fmt.Println(err)
		return "", 0
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("origin", "https://m.ppply.site")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://m.ppply.site/")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", 0
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", 0
	}
	type Response struct {
		DT struct {
			CS []float64 `json:"cs"` // cs是一个float64数组
			LS struct {
				SI struct {
					Sid string `json:"sid"`
					// 其他字段可以根据需要添加
				} `json:"si"`
			} `json:"ls"`
		} `json:"dt"`
		Err interface{} `json:"err"`
	}
	fmt.Println(string(body))
	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	// 获取 sid
	sid := resp.DT.LS.SI.Sid
	cs := resp.DT.CS[0]
	return sid, cs
}

//func PPHack(count int, userName, GameId string) {
//	numThreads := runtime.NumCPU() // 获取 CPU 核心数
//	runtime.GOMAXPROCS(numThreads) // 设置最大可同时使用的 CPU 核心数
//	var wg sync.WaitGroup
//
//	// 每个协程处理的迭代次数
//	perThread := count / numThreads
//	remainder := count % numThreads
//
//	for i := 0; i < numThreads; i++ {
//		wg.Add(1)
//		go func(threadID int) {
//			defer wg.Done()
//			if userName == "" {
//				userName = "test_user"
//			}
//			uid := fmt.Sprintf("%s_%d", userName, threadID)
//			spinNum := perThread
//			if threadID == numThreads-1 {
//				spinNum += remainder
//			}
//			fmt.Printf("协程 id：%d ,分配次数：%d\n", threadID, spinNum)
//			var err error
//			var token Token
//			var getlinkPs GetlinkPs
//			//str := "action=doSpin&symbol=vs20olympx&c=%d&l=20&bl=0&index=5&counter=183&repeat=0&mgckey=eyJQIjoxMDAyMjYsIkUiOjE3MzI2Mjk1MTAsIlMiOjEwMDAsIkQiOiI2NSJ9.kGBj8OA4x1r9fKaqrwqD6pLNFRIGS8uAYFKGsVhkDNs"
//			//str = fmt.Sprintf(str, ps.Bet)
//			//bytes := []byte(str)
//			//cs := Cs[0]
//			//ml := Ml[0]
//			getlinkPs.UserID = uid
//			getlinkPs.GameID = GameId
//			getlinkPs.Language = "en"
//			getlinkPs.GamePreFix = ""
//			//不足1000000 自动设置
//			err = mux.HttpInvoke(fakeAppUrl, getlinkPs, &token)
//			if err != nil {
//				slog.Error("spin_stat_err", "HttpInvoke", err)
//				return
//			}
//			//token换pid
//			pid, _, _ := jwtutil.ParseTokenData(token.Token)
//
//			params := fmt.Sprintf("mgckey=%s&c=%v", token.Token, 1)
//
//			spinMsg := &nats.Msg{
//				Data: []byte(params),
//				Header: nats.Header{
//					"Token": []string{token.Token},
//				},
//			}
//
//			coll := db.Collection("players")
//			player, err := db.GetPlr[ppcomm.Player](pid)
//			if err != nil {
//				slog.Error("db GetPlr", "GetPlr", err)
//				return
//			}
//
//			//缓存
//			appkey := slotsmongo.MakeKey(player.AppID, GameId)
//			info2, ok := slotsmongo.Load(appkey)
//			if !ok {
//				info2, err = slotsmongo.GetAppInfo(player.AppID, GameId)
//				if err != nil {
//					slog.Error("slotsmongo.GetAppInfo", "GetAppInfo", err)
//					return
//				}
//				info2 = slotsmongo.StoreAndLoad(appkey, info2, player.PID)
//			} else {
//				sid := strings.Split(player.LastData.Str("gid"), "_")
//				c, _ := slotsmongo.GetPlayerCs(appkey, player.PID)
//
//				if len(player.LastSid) == 0 || sid[1] == sid[2] ||
//					len(c) == 0 || c == nil { // 表示是最后一盘
//					//直接覆盖
//					info2 = slotsmongo.StoreAndLoad(appkey, info2, player.PID)
//				}
//			}
//
//			for i := 0; i < spinNum; i++ {
//				//ret.Count++
//				for {
//
//					_, err = doSpin(spinMsg)
//					if err != nil {
//						slog.Error("doSpin", "doSpin", err)
//						return
//					}
//					player, err = db.GetPlr[ppcomm.Player](pid)
//					if err != nil {
//						slog.Error("spin_stat_err", "GetPlr", err)
//						return
//					}
//					bet, _ := strconv.Atoi(spinMsg.Header.Get("stat_bet"))
//					ret.Bet += int64(bet)
//					win, _ := strconv.Atoi(spinMsg.Header.Get("stat_win"))
//					ret.Win += int64(win)
//					//isEnd := player.IsEnd
//					//lastId := spinMsg.Header.Get("lastId")
//
//					spinMsg.Header.Set("stat_bet", "")
//					spinMsg.Header.Set("stat_win", "")
//
//					params2 := strings.Split(player.LastData.Str("gid"), "_")
//					//一轮游戏结束
//					if len(params2) != 3 || params2[1] == params2[2] || player.IsEnd {
//						spinMsg.Header.Set("lastId", "")
//						//结束
//						spinMsg2 := &nats.Msg{
//							Data: []byte(params),
//							Header: nats.Header{
//								"mgckey":   []string{token.Token},
//								"counter":  []string{strconv.Itoa(ret.Count + 1)},
//								"index":    []string{strconv.Itoa(ret.Count)},
//								"jump_log": []string{"1"},
//							},
//						}
//						doCollect(spinMsg2)
//						break
//					}
//				}
//
//			}
//
//			ret.Balance_after, _ = slotsmongo.GetBalance(ps.Pid)
//			_, err = coll.ReplaceOne(context.TODO(), db.ID(ps.Pid), player, options.Replace().SetUpsert(true))
//
//			//c, err := slotsmongo.GetPlayerCs(slotsmongo.MakeKey(player.AppID, lazy.ServiceName), player.PID)
//			//if err != nil {
//			//	slog.Error("spin_stat_err", "GetPlayerCs", err)
//			//}
//			//
//			//respon := &json.RawMessage{}
//			//for i := 0; i < spinNum; i++ {
//			//	ret.Count++
//			//	for {
//			//		player, err = db.GetPlr[models.Player](pid)
//			//		if err != nil {
//			//			slog.Error("spin_stat_err", "GetPlr", err)
//			//			return
//			//		}
//			//		err = spin(player, spinMsg, respon)
//			//		if err != nil {
//			//			slog.Error("spin_stat_err", "spin", err)
//			//		}
//			//
//			//		update := bson.M{
//			//			"$set": player, // 更新所有字段，假设 player 只包含需要更新的字段
//			//		}
//			//
//			//		// 执行更新操作
//			//		_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": player.PID}, update, options.Update().SetUpsert(true))
//			//
//			//		if err != nil {
//			//			fmt.Println(player)
//			//			slog.Error("ReplaceOne", "ReplaceOne", err)
//			//		}
//			//		bet := spinMsg.Header["stat_bet"].(int64)
//			//		ret.Bet += bet
//			//		win := spinMsg.Header["stat_win"].(int64)
//			//		ret.Win += win
//			//		lastId := spinMsg.Header["LastSid"].(string)
//			//		if lastId == "" {
//			//			break
//			//		}
//			//		spinMsg.Form.Set("id", fmt.Sprintf("%v", lastId))
//			//		spinMsg.Header["stat_bet"] = int64(0)
//			//		spinMsg.Header["stat_win"] = int64(0)
//			//		params := strings.Split(lastId, "_")
//			//		//一轮游戏结束
//			//		if len(params) != 3 || params[1] == params[2] {
//			//			break
//			//		}
//			//	}
//			//}
//			//
//			//_, err = coll.ReplaceOne(context.TODO(), db.ID(ps.Pid), player, options.Replace().SetUpsert(true))
//		}(i)
//	}
//
//	wg.Wait()
//	for _, statRet := range rets {
//		ret.Count += statRet.Count
//		ret.Win += statRet.Win
//		ret.Bet += statRet.Bet
//		ret.Balance_after += statRet.Balance_after
//		ret.Balance_before += statRet.Balance_before
//	}
//}
