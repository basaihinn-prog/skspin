package slotsmongo

import (
	"context"
	"fmt"
	"time"

	"serve/comm/db"
	"serve/comm/ut"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

type AddRunTimeDataReq struct {
	Pid int64 `json:"pid"`
	// AppId          string `json:"appId"`
	Game           string `json:"game"`
	Expense        int64  `json:"expense"`
	Income         int64  `json:"income"`
	SelfPoolReward int64  `json:"selfPoolReward"`
	IsSmallGame    bool   `json:"isSmallGame"`
	IsFreeGameTurn bool   `json:"isFreeGameTurn"`
}

type DocSlotsRuntimeDay struct {
	ID              string `bson:"_id"`
	Expense         int    `bson:"Expense"`   // 支出
	Income          int    `bson:"Income"`    // 赚取
	NetProfit       int    `bson:"NetProfit"` // 净利润
	SelfPoolReward  int    `bson:"SelfPoolReward"`
	Hit             int    `bson:"Hit"`
	AwardHit        int    `bson:"AwardHit"`
	SmallGameHit    int    `bson:"SmallGameHit"`
	SmallGameIncome int    `bson:"SmallGameIncome"`
	FreeGameHit     int    `bson:"FreeGameHit"`
	FreeGameIncome  int    `bson:"FreeGameIncome"`
	Day             string `bson:"Day"`
	Game            string `bson:"Game"`
	AppID           string `bson:"AppID"`
}

type DocSlotsRuntimePlayer struct {
	ID      string `bson:"_id"`
	Expense int    `bson:"Expense"` // 所有玩家的单日消费
	Income  int    `bson:"Income"`  // 所有玩家的总消费
	Hit     int    `bson:"Hit"`
}

func AddRunTimeData(req *AddRunTimeDataReq) (err error) {
	day := time.Now().Format("2006-01-02")

	// 统计单日
	{
		id := RuntimeDay_GetID(day, req.Game)
		coll := SlotsCollection("RuntimeDay")
		coll.InsertOne(context.TODO(), DocSlotsRuntimeDay{
			ID:   id,
			Game: req.Game,
			Day:  day,
		})

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

		coll.UpdateByID(context.TODO(), id, bson.M{"$inc": incDoc})

		///////////////////////////////////////////////////////
		appId := players.GetAppID(req.Pid)
		id = ut.JoinStr('-', day, req.Game, appId)
		coll.InsertOne(context.TODO(), DocSlotsRuntimeDay{
			ID:    id,
			Game:  req.Game,
			Day:   day,
			AppID: appId,
		})
		coll.UpdateByID(context.TODO(), id, bson.M{"$inc": incDoc})
	}

	// 统计玩家
	{
		id := RuntimePlayer_GetID(day, req.Game, int(req.Pid))
		CollRuntimePlayer := SlotsCollection("RuntimePlayer")
		_, err = CollRuntimePlayer.InsertOne(context.TODO(), &DocSlotsRuntimePlayer{ID: id})
		if err == nil {
			// 递增统计进入
			SlotsCollection("RuntimeSummary").UpdateByID(context.TODO(), RuntimeSummary_PlayerCount(day, req.Game), bson.M{"$inc": bson.M{"Val": 1}}, options.Update().SetUpsert(true))
		}

		if mongo.IsDuplicateKeyError(err) {
			err = nil
		}

		if err != nil {
			return
		}

		_, err = CollRuntimePlayer.UpdateByID(context.TODO(), id, bson.M{
			"$inc": bson.M{
				"Expense": req.Expense,
				"Income":  req.Income,
				"Hit":     1,
			},
		})

		if err != nil {
			return
		}
	}
	// 统计总和
	{
		CollRuntimeSummary := SlotsCollection("RuntimeSummary")

		CollRuntimeSummary.UpdateByID(context.TODO(), RuntimeSummary_TotalIncome(req.Game), bson.M{"$inc": bson.M{"Val": req.Income}}, db.UpsertOpt())
		CollRuntimeSummary.UpdateByID(context.TODO(), RuntimeSummary_TotalExpense(req.Game), bson.M{"$inc": bson.M{"Val": req.Expense}}, db.UpsertOpt())
	}

	return
}
