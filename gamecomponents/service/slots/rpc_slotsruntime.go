package main

import (
	"context"
	"fmt"
	"game/duck/logger"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/slots"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocSlotsRuntimeDay struct {
	ID              string `bson:"_id"`
	Expense         int    // 支出
	Income          int    // 赚取
	NetProfit       int    // 净利润
	SelfPoolReward  int
	Hit             int
	AwardHit        int
	SmallGameHit    int
	SmallGameIncome int
	FreeGameHit     int
	FreeGameIncome  int
	Day             string
	Game            string
}

type DocSlotsRuntimePlayer struct {
	ID      string `bson:"_id"`
	Expense int    // 所有玩家的单日消费
	Income  int    // 所有玩家的总消费
	Hit     int
}

type DocSlotsRuntimeSummary struct {
	Type string `bson:"_id"`
	Val  int
}

type RuntimeRpcServer struct {
}

func RuntimeDay_GetID(day, game string) string {
	return fmt.Sprintf("%v-%v", day, game)
}

func RuntimePlayer_GetID(day string, game string, pid int) string {
	return fmt.Sprintf("%v-%v-%v", day, game, pid)
}

func RuntimeSummary_PlayerCount(day, game string) string {
	return fmt.Sprintf("%v-%v", day, game)
}

func RuntimeSummary_TotalIncome(game string) string {
	return fmt.Sprintf("%v-TotalIncome", game)
}

func RuntimeSummary_TotalExpense(game string) string {
	return fmt.Sprintf("%v-TotalExpense", game)
}

func (RuntimeRpcServer) AddRunTimeData(ctx context.Context, req *slots.AddRunTimeDataReq) (*pb.Empty, error) {
	resp := &pb.Empty{}
	day := time.Now().Format("2006-01-02")

	// 统计单日
	{
		id := RuntimeDay_GetID(day, req.Game)

		var doc *DocSlotsRuntimeDay
		err := CollRuntimeDay.FindId(id, &doc)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.Err("CollSlotsRuntimeDay.FindId err:", day, err, req)
			return resp, err
		}
		if err == mongo.ErrNoDocuments {
			CollRuntimeDay.InsertOne(&DocSlotsRuntimeDay{ID: id, Game: req.Game, Day: day})
		}

		incDoc := bson.M{
			"Expense":        req.Expense,
			"Income":         req.Income,
			"NetProfit":      req.Income - req.Expense,
			"SelfPoolReward": req.SelfPoolReward,
			"Hit":            1,
		}

		if req.Income > 0 {
			incDoc["AwardHit"] = 1
		}

		if req.IsSmallGame {
			incDoc["SmallGameHit"] = 1
			incDoc["SmallGameIncome"] = req.Income
		}

		if req.IsFreeGameTurn {
			incDoc["FreeGameHit"] = 1
			incDoc["FreeGameIncome"] = req.Income
		}

		CollRuntimeDay.UpdateId(id, bson.M{"$inc": incDoc})
	}

	// 统计玩家
	{
		id := RuntimePlayer_GetID(day, req.Game, int(req.Pid))
		var doc *DocSlotsRuntimePlayer
		err := CollRuntimePlayer.FindId(id, &doc)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.Err("CollSlotsRuntimeDay.FindId err:", id, err, req)
			return resp, err
		}
		if err == mongo.ErrNoDocuments {
			CollRuntimePlayer.InsertOne(&DocSlotsRuntimePlayer{ID: id})

			// 递增统计进入
			CollRuntimeSummary.UpsertId(RuntimeSummary_PlayerCount(day, req.Game), bson.M{"$inc": bson.M{"Val": 1}})
		}

		CollRuntimePlayer.UpdateId(id, bson.M{
			"$inc": bson.M{
				"Expense": req.Expense,
				"Income":  req.Income,
				"Hit":     1,
			},
		})
	}

	// 统计总和
	{
		CollRuntimeSummary.UpsertId(RuntimeSummary_TotalIncome(req.Game), bson.M{"$inc": bson.M{"Val": req.Income}})
		CollRuntimeSummary.UpsertId(RuntimeSummary_TotalExpense(req.Game), bson.M{"$inc": bson.M{"Val": req.Expense}})
	}

	return resp, nil
}

func (sf RuntimeRpcServer) FixRunTimeData(ctx context.Context, req *slots.FixRunTimeDataReq) (*slots.ThreeMonkeyRunTimeData, error) {
	id := RuntimeDay_GetID(req.Day, req.Game)
	var doc = &DocSlotsRuntimeDay{}
	err := CollRuntimeDay.FindId(id, &doc)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	result := slots.ThreeMonkeyRunTimeData{
		Date:            req.Day,
		Expense:         int64(doc.Expense),
		Income:          int64(doc.Income),
		NetProfit:       int64(doc.NetProfit),
		SelfPoolReward:  int64(doc.SelfPoolReward),
		Hit:             int64(doc.Hit),
		AwardHit:        int64(doc.AwardHit),
		SmallGameHit:    int64(doc.SmallGameHit),
		SmallGameIncome: int64(doc.SmallGameIncome),
		FreeGameHit:     int64(doc.FreeGameHit),
		FreeGameIncome:  int64(doc.FreeGameIncome),
	}

	result.PlrCount, _ = sf.getSummaryVal(RuntimeSummary_PlayerCount(req.Day, req.Game))

	return &result, nil
}

func (sf RuntimeRpcServer) GetTotal(ctx context.Context, req *slots.GetTotalReq) (*slots.GetTotalResp, error) {
	var ret = &slots.GetTotalResp{}

	ret.Income, _ = sf.getSummaryVal(RuntimeSummary_TotalIncome(req.Game))
	ret.Expense, _ = sf.getSummaryVal(RuntimeSummary_TotalExpense(req.Game))

	return ret, nil
}

func (RuntimeRpcServer) getSummaryVal(key string) (int64, error) {
	var doc *DocSlotsRuntimeSummary

	err := CollRuntimeSummary.FindId(key, &doc)
	if err == nil {
		return int64(doc.Val), nil
	}

	if err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Err("RuntimeRpcServer.getSummaryVal err:", key, err)
		}
	}

	return 0, err
}
