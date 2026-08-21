package slotsmongo

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/comm/ut"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type OpDailyDoc struct {
	ID             string `json:"-" bson:"_id"`
	BetAmount      int64  `bson:"betamount"` // 支出
	WinAmount      int64  `bson:"winamount"` // 赚取
	SpinCount      int    `bson:"spincount"` // 转动次数
	RoundCount     int    `bson:"roundcount"`
	LoginPlrCount  int    `bson:"loginplrcount"`
	RegistPlrCount int    `bson:"registplrcount"`

	Date  string `bson:"date"`
	AppID string `bson:"appid"`
}

// 运营商报表
// 投注日期 运营商名称 活跃玩家数量 新注册玩家数量	投注次数 回合次数 投注金额 总赢得金额 玩家输赢 公司输赢 公司输赢占比
func addOperatorReport(dateStr, appId string, incDoc bson.D) {
	operator := comm.Operator_V2{}
	collOperator := db.Collection2("GameAdmin", "AdminOperator")
	collOperator.FindOne(context.TODO(), bson.M{"AppID": appId}).Decode(&operator)
	coll := ReportsCollection("OperatorDailyReport")
	// dateStr := time.Now().Format("20060102")
	id := ut.JoinStr(':', dateStr, appId)

	// doc := &OpDailyDoc{
	// 	ID:    id,
	// 	Date:  dateStr,
	// 	AppID: appId,
	// }
	// coll.InsertOne(context.TODO(), doc)

	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"appid", appId,
			"OperatorBalance", operator.Balance,
		),
		"$inc", incDoc,
	)

	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}

func IncRegistCount(appId string) {
	coll := ReportsCollection("OperatorDailyReport")
	dateStr := time.Now().UTC().Format("20060102")
	id := ut.JoinStr(':', dateStr, appId)

	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"appid", appId,
		),
		"$inc", db.D(
			"registplrcount", 1,
		),
	)

	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}

func IncLoginCount(appId string) {
	operator := comm.Operator_V2{}
	collOperator := db.Collection2("GameAdmin", "AdminOperator")
	collOperator.FindOne(context.TODO(), bson.M{"AppID": appId}).Decode(&operator)
	coll := ReportsCollection("OperatorDailyReport")
	dateStr := time.Now().UTC().Format("20060102")
	id := ut.JoinStr(':', dateStr, appId)

	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"appid", appId,
			"OperatorBalance", operator.Balance,
		),
		"$inc", db.D(
			"loginplrcount", 1,
		),
	)

	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}
