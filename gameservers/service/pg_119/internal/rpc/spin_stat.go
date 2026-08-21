package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"runtime"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_119/internal/models"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {

	mux.RegRpc("/game-rtp-api/119/game/spinstat", "spinstat", "game-rtp-api", db.WrapRpcRtpPlayer(playStatistics), nil)
	mux.RegRpc("/game-rtp-api/119/game/querysimulate", "spinstat", "query-simulate", db.WrapRpcRtpPlayer(QuerySimulate), nil)

	mux.RegHttpWithSample("/test/statistics", "多次转动,统计回报", "test", GoRTP, ut.SpinStatPs{
		Count: 1000,
	})
}

//const fakeAppUrl = "https://free.ppply.site/GetLink"

const fakeAppUrl = "http://127.0.0.1:27000/GetLink"

type Token struct {
	Token string `json:"token"`
}
type GetlinkPs struct {
	UserID     string `json:"user_id"`
	GameID     string `json:"game_id"`
	Language   string `json:"language"`
	GamePreFix string `json:"game_pre_fix"`
}
type PostFrom struct {
	Count    int
	GameId   string
	GameName string
	UName    string
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
			cs := Cs[0]
			ml := Ml[0]
			getlinkPs.UserID = uid
			getlinkPs.GameID = postFrom.GameId
			getlinkPs.Language = "en"
			getlinkPs.GamePreFix = postFrom.GameName
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
			spinMsg := define.PGParams{
				Form:   url.Values{},
				Header: map[string]any{},
			}

			spinMsg.Form.Set("id", fmt.Sprintf("%v", 0))
			spinMsg.Form.Set("cs", fmt.Sprintf("%v", cs))
			spinMsg.Form.Set("ml", fmt.Sprintf("%v", ml))

			spinMsg.Header["stat_bet"] = int64(0)
			spinMsg.Header["stat_win"] = int64(0)
			spinMsg.Header["jump_log"] = int64(1)
			spinMsg.Header["LastSid"] = ""

			coll := db.Collection("players")
			player, err := db.GetPlr[models.Player](pid)
			if err != nil {
				slog.Error("spin_stat_err", "GetPlr", err)
				return
			}
			//缓存
			isEnd, _ := player.IsEndO()
			redisx.GetPlayerCs(player.AppID, player.PID, isEnd)

			respon := &json.RawMessage{}
			for i := 0; i < spinNum; i++ {
				ret.Count++
				for {
					player, err = db.GetPlr[models.Player](pid)
					if err != nil {
						slog.Error("spin_stat_err", "GetPlr", err)
						return
					}
					err = spin(player, spinMsg, respon)
					if err != nil {
						slog.Error("spin_stat_err", "spin", err)
					}

					update := bson.M{
						"$set": player, // 更新所有字段，假设 player 只包含需要更新的字段
					}

					// 执行更新操作
					_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": player.PID}, update, options.Update().SetUpsert(true))

					if err != nil {
						fmt.Println(player)
						slog.Error("ReplaceOne", "ReplaceOne", err)
					}
					bet := spinMsg.Header["stat_bet"].(int64)
					ret.Bet += bet
					win := spinMsg.Header["stat_win"].(int64)
					ret.Win += win
					lastId := spinMsg.Header["LastSid"].(string)
					if lastId == "" {
						break
					}
					spinMsg.Form.Set("id", fmt.Sprintf("%v", lastId))
					spinMsg.Header["stat_bet"] = int64(0)
					spinMsg.Header["stat_win"] = int64(0)
					params := strings.Split(lastId, "_")
					//一轮游戏结束
					if len(params) != 3 || params[1] == params[2] {
						break
					}
				}
			}

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

func QuerySimulate(plr *models.Player, ps define.PGParams, ret *json.RawMessage) (err error) {
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
	tmp.Ml = Ml[0]
	tmp.Cs = Cs[0]
	tmp.Mxl = Mxl
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

func playStatistics(plr *models.Player, ps define.PGParams, ret *json.RawMessage) (err error) {
	ss := ut.SpinStatPs{}
	rets := ut.SpinStatRet{}

	ss.Bet = 0.05
	ss.BuyBonus, _ = strconv.ParseBool(ps.Form.Get("BuyBonus"))
	ss.Count, _ = strconv.Atoi(ps.Form.Get("Count"))

	var postFrom PostFrom
	postFrom.GameId = ps.Form.Get("GameId")
	postFrom.GameName = ps.Form.Get("GameName")
	postFrom.UName = ps.Form.Get("UName")
	postFrom.Count, _ = strconv.Atoi(ps.Form.Get("Count"))

	GoRTP(ss, &rets, postFrom)

	//rtp := rets.Win / (rets.Bet / 100)
	marshal, err := json.Marshal(rets)
	if err != nil {
		return
	}

	*ret = marshal

	//fmt.Println(rtp, "%")

	//fmt.Println(time.Now().Format(time.DateTime))

	return err
}
