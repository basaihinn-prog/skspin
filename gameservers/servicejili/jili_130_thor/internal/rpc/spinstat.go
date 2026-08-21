package rpc

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"runtime"
	"serve/comm/db"
	"serve/comm/ut"
	"serve/servicejili/jili_130_thor/internal/models"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"serve/comm/jwtutil"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/servicejili/jili_130_thor/internal"
	"serve/servicejili/jiliut/jiliUtMessage"

	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

func init() {
	mux.RegRpc(fmt.Sprintf("%s/game/spinstat", internal.GameID), "spinstat", "game-rtp-api", playStatistics, nil)
	mux.RegRpc(fmt.Sprintf("%s/game/querysimulate", internal.GameID), "spinstat", "query-simulate", QuerySimulate, nil)
}

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

const fakeAppUrl = "http://127.0.0.1:27000/GetLink"

func GoRTP(ps ut.SpinStatPs, ret *ut.SpinStatRet, postFrom PostFrom) {
	numThreads := runtime.NumCPU()
	runtime.GOMAXPROCS(numThreads)
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			var err error

			var getlinkPs GetlinkPs
			var token Token
			getlinkPs.UserID = fmt.Sprintf("%s_%d", postFrom.UName, threadID)
			getlinkPs.GameID = postFrom.GameId
			getlinkPs.Language = "en"
			err = mux.HttpInvoke(fakeAppUrl, getlinkPs, &token)
			pid, _, _ := jwtutil.ParseTokenData(token.Token)
			ret.Balance_before, _ = slotsmongo.GetBalance(pid)

			gameReq := &jiliUtMessage.Server_Request{
				Ack: proto.Int32(0),
				Req: jiliut.ProtoEncode(&jiliUtMessage.Server_SpinReq{
					BuyCustomType: proto.Int32(0),
					Bet:           proto.Float64(cmp.Or(ps.Bet, 1.0)),
					MallBet:       proto.Float64(lo.Ternary(ps.BuyBonus, 1.0, 0.0)),
				}),
			}

			spinMsg := &nats.Msg{
				Data: jiliut.ProtoEncode(gameReq),
				Header: nats.Header{
					"isstat":   []string{"1"},
					"Token":    []string{token.Token},
					"jump_log": []string{"1"},
				},
			}

			_, err = reqInfo(pid, spinMsg.Data, spinMsg)
			for j := 0; j < ps.Count; j++ {

				_, err = req(spinMsg)
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

func QuerySimulate(ps GameData, ret *json.RawMessage) (err error) {
	coll := db.Collection("rawSpinData")
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
		Bet      float64 `json:"bet"`
		Times    float64 `json:"times"`
	}

	tmp.Count = count
	tmp.AllCount = AllCount
	tmp.Bet, err = strconv.ParseFloat(ps.Bet, 64)
	if err != nil {
		slog.Error(err.Error())
	}
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

type GameData struct {
	Pid      int64  `json:"Pid"`
	Count    int    `json:"Count"`
	BuyBonus bool   `json:"BuyBonus"`
	Bet      string `json:"Bet"`
	GameId   string `json:"GameId"`
	GameName string `json:"GameName"`
	UName    string `json:"UName"`
}

func playStatistics(ps GameData, ret *json.RawMessage) (err error) {
	fmt.Println(ps)
	ss := ut.SpinStatPs{}
	rets := ut.SpinStatRet{}

	ss.Pid = ps.Pid
	ss.Bet, _ = strconv.ParseFloat(ps.Bet, 64)
	ss.BuyBonus = ps.BuyBonus
	ss.Count = ps.Count

	var postFrom PostFrom
	postFrom.GameId = ps.GameId
	postFrom.GameName = ps.GameName
	postFrom.UName = ps.UName
	postFrom.Count = ps.Count

	GoRTP(ss, &rets, postFrom)

	marshal, err := json.Marshal(rets)
	if err != nil {
		return
	}

	*ret = marshal
	return err
}

func infortp(plr *models.Player) {

}

func spinstat(ps ut.SpinStatPs, ret *ut.SpinStatRet) (err error) {
	ret.Balance_before, _ = slotsmongo.GetBalance(ps.Pid)
	token := lo.Must(jwtutil.NewToken(ps.Pid, time.Now().Add(time.Hour*12)))
	gameReq := &jiliUtMessage.Server_Request{
		Ack: proto.Int32(0),
		Req: jiliut.ProtoEncode(&jiliUtMessage.Server_SpinReq{
			BuyCustomType: proto.Int32(0),
			Bet:           proto.Float64(cmp.Or(ps.Bet, 1.0)),
			MallBet:       proto.Float64(lo.Ternary(ps.BuyBonus, 1.0, 0.0)),
		}),
	}

	spinMsg := &nats.Msg{
		Data: jiliut.ProtoEncode(gameReq),
		Header: nats.Header{
			"isstat": []string{"1"},
			"Token":  []string{token},
		},
	}

	for i := 0; i < ps.Count; i++ {
		_, err = req(spinMsg)
		if err != nil {
			return
		}

		ret.Count++
		bet, _ := strconv.Atoi(spinMsg.Header.Get("stat_bet"))
		ret.Bet += int64(bet)
		win, _ := strconv.Atoi(spinMsg.Header.Get("stat_win"))
		ret.Win += int64(win)
	}

	ret.Balance_after, _ = slotsmongo.GetBalance(ps.Pid)
	return
}
