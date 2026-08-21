package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"runtime"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicepp/ppcomm"
	"strconv"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	//mux.RegRpc("/pp_vs20olympx/game/spinstat", "spin stat", "stat", spinstat, nil).SetOnlyDev()

	mux.RegRpc("/game-rtp-api/pp_vs40demonpots/game/spinstat", "spinstat", "game-rtp-api", db.WrapRpcRtpPlayer(playStatistics), nil)
	mux.RegRpc("/game-rtp-api/pp_vs40demonpots/game/querysimulate", "spinstat", "query-simulate", db.WrapRpcRtpPlayer(QuerySimulate), nil)

}

const fakeAppUrl = "http://127.0.0.1:27000/GetLink"

type Token struct {
	Token string `json:"token"`
}
type GetlinkPs struct {
	UserID     string `json:"user_id"`
	GameID     string `json:"game_id"`
	Language   string `json:"language"`
	GamePreFix string `json:"game_pre_fix"`
	BuyEX      int    `json:"BuyEX"`
}
type PostFrom struct {
	Count    int
	GameId   string
	GameName string
	UName    string
	BuyEX    int
}

func GoRTP(ps ut.SpinStatPs, ret *ut.SpinStatRet, postFrom PostFrom) {
	numThreads := runtime.NumCPU() // 获取 CPU 核心数
	runtime.GOMAXPROCS(numThreads) // 设置最大可同时使用的 CPU 核心数
	var wg sync.WaitGroup
	rets := make([]*ut.SpinStatRet, numThreads)
	for i, _ := range rets {
		rets[i] = &ut.SpinStatRet{}
	}
	// 每个协程处理的迭代次数
	perThread := ps.Count / numThreads
	remainder := ps.Count % numThreads
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			if postFrom.UName == "" {
				postFrom.UName = "test_user"
			}
			uid := fmt.Sprintf("%s_%d", postFrom.UName, threadID)
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
			getlinkPs.GameID = postFrom.GameId
			getlinkPs.Language = "en"
			getlinkPs.GamePreFix = ""
			//不足1000000 自动设置
			err = mux.HttpInvoke(fakeAppUrl, getlinkPs, &token)
			if err != nil {
				slog.Error("spin_stat_err", "HttpInvoke", err)
				return
			}
			//token换pid
			pid, _, _ := jwtutil.ParseTokenData(token.Token)

			//获取用户余额
			ret.Balance_before, err = slotsmongo.GetBalance(pid)

			if err != nil {
				slog.Error("spin_stat_err", "GetBalance", err)
				return
			}
			//token := lo.Must(jwtutil.NewToken(ps.Pid, time.Now().Add(time.Hour*12)))
			params := fmt.Sprintf("mgckey=%s&c=%v", token.Token, ps.Bet)

			if ps.BuyBonus {
				if ps.BuyInt == 0 {
					params = fmt.Sprintf("%s&pur=%v", params, 0)
				}
				if ps.BuyInt == 1 {
					params = fmt.Sprintf("%s&pur=%v", params, 1)
				}
			}

			if postFrom.BuyEX == 1 {
				params = fmt.Sprintf("%s&bl=%v", params, 1)
			}

			spinMsg := &nats.Msg{
				Data: []byte(params),
				Header: nats.Header{
					"isstat":   []string{"1"},
					"Token":    []string{token.Token},
					"stat_bet": []string{"0"},
					"stat_win": []string{"0"},
					"jump_log": []string{"1"},
				},
			}

			coll := db.Collection("players")
			player, err := db.GetPlr[ppcomm.Player](pid)
			if err != nil {
				slog.Error("db GetPlr", "GetPlr", err)
				return
			}

			//缓存
			isEnd, _ := player.IsEndO()
			redisx.GetPlayerCs(player.AppID, player.PID, isEnd)

			for i := 0; i < spinNum; i++ {
				ret.Count++
				for {
					_, err = doSpin(spinMsg)
					if err != nil {
						slog.Error("doSpin", "doSpin", err)
						return
					}
					player, err = db.GetPlr[ppcomm.Player](pid)
					if err != nil {
						slog.Error("spin_stat_err", "GetPlr", err)
						return
					}
					bet, _ := strconv.Atoi(spinMsg.Header.Get("stat_bet"))
					ret.Bet += int64(bet)
					win, _ := strconv.Atoi(spinMsg.Header.Get("stat_win"))
					ret.Win += int64(win)
					//isEnd := player.IsEnd
					//lastId := spinMsg.Header.Get("lastId")

					spinMsg.Header.Set("stat_bet", "")
					spinMsg.Header.Set("stat_win", "")

					params2 := strings.Split(player.LastData.Str("gid"), "_")
					//一轮游戏结束
					if len(params2) != 3 || params2[1] == params2[2] || player.IsEnd {
						spinMsg.Header.Set("lastId", "")
						//结束
						spinMsg2 := &nats.Msg{
							Data: []byte(params),
							Header: nats.Header{
								"mgckey":   []string{token.Token},
								"counter":  []string{strconv.Itoa(ret.Count + 1)},
								"index":    []string{strconv.Itoa(ret.Count)},
								"jump_log": []string{"1"},
							},
						}
						doCollect(spinMsg2)
						break
					}
				}
			}

			ret.Balance_after, _ = slotsmongo.GetBalance(ps.Pid)
			_, err = coll.ReplaceOne(context.TODO(), db.ID(ps.Pid), player, options.Replace().SetUpsert(true))
		}(i)
	}

	wg.Wait()
	for _, statRet := range rets {
		ret.Count += statRet.Count
		ret.Win += statRet.Win
		ret.Bet += statRet.Bet
		ret.Balance_after += statRet.Balance_after
		ret.Balance_before += statRet.Balance_before
	}
}

func QuerySimulate(plr *ppcomm.Player, ps define.PGParams, ret *json.RawMessage) (err error) {
	coll := db.Collection("simulate")
	filter := db.D("selected", true, "bucketid", bson.D{{"$lte", 28}})

	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	AllCount, err := coll.CountDocuments(context.TODO(), db.D())
	if err != nil {
		log.Fatal(err)
	}

	// $match 阶段：筛选 selected 为 true 且 bucketid 小于等于 28 的文档
	matchStage := bson.D{{"$match", bson.D{
		{"selected", true},
		{"bucketid", bson.D{{"$lte", 28}}},
	}}}

	// $group 阶段：将结果分组，计算 times 字段的总和
	groupStage := bson.D{{"$group", bson.D{
		{"_id", nil},
		{"totalAmount", bson.D{{"$sum", "$times"}}},
	}}}

	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		log.Fatal(err)
	}

	// 存储结果
	var results []struct {
		TotalAmount float64 `bson:"totalAmount"`
	}

	// 将查询结果解码到 results 变量中
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}
	var tmp struct {
		Count    int64   `json:"count"`
		AllCount int64   `json:"all_count"`
		Ml       float64 `json:"ml"`
		Cs       float64 `json:"cs"`
		Mxl      int     `json:"mxl"`
		Times    float64 `json:"times"`
	}

	tmp.Count = count
	tmp.AllCount = AllCount
	//tmp.Ml = Ml[0]
	//tmp.Cs = Cs[0]
	//tmp.Mxl = Mxl
	if len(results) > 0 {
		tmp.Times = results[0].TotalAmount
	}
	if len(results) == 0 {
		slog.Error("results is empty")
	}
	marshal, err := json.Marshal(tmp)
	if err != nil {
		return
	}

	*ret = marshal

	return
}

func playStatistics(plr *ppcomm.Player, ps define.PGParams, ret *json.RawMessage) (err error) {
	ss := ut.SpinStatPs{}
	rets := ut.SpinStatRet{}

	ss.Bet = 0.05
	ss.Pid = ps.Pid
	ss.BuyBonus, _ = strconv.ParseBool(ps.Form.Get("BuyBonus"))
	ss.Count, _ = strconv.Atoi(ps.Form.Get("Count"))
	ss.Extra, _ = strconv.Atoi(ps.Form.Get("BuyEX"))
	ss.BuyInt, _ = strconv.Atoi(ps.Form.Get("BuyInt"))

	var postFrom PostFrom
	postFrom.GameId = ps.Form.Get("GameId")
	postFrom.GameName = ps.Form.Get("GameName")
	postFrom.UName = ps.Form.Get("UName")
	postFrom.Count, _ = strconv.Atoi(ps.Form.Get("Count"))
	postFrom.BuyEX, _ = strconv.Atoi(ps.Form.Get("BuyEx"))

	GoRTP(ss, &rets, postFrom)

	marshal, err := json.Marshal(rets)
	if err != nil {
		return
	}

	*ret = marshal
	return err
}
