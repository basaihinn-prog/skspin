package slotsmongo

import (
	"context"
	"game/comm/db"
	"game/comm/ut"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type OpDailyDoc struct {
	ID             map[string]any `json:"-" bson:"_id"`
	BetAmount      int64          `bson:"betamount"` // 支出
	WinAmount      int64          `bson:"winamount"` // 赚取
	SpinCount      int            `bson:"spincount"` // 转动次数
	RoundCount     int            `bson:"roundcount"`
	LoginPlrCount  int            `bson:"loginplrcount"`
	RegistPlrCount int            `bson:"registplrcount"`

	Date  string `bson:"date"`
	AppID string `bson:"appid"`
}

// 运营商报表
// 投注日期 运营商名称 活跃玩家数量 新注册玩家数量	投注次数 回合次数 投注金额 总赢得金额 玩家输赢 公司输赢 公司输赢占比
func addOperatorReport(dateStr, appId string, incDoc bson.D) {
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
		),
		"$inc", incDoc,
	)

	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}

func IncRegistCount(appId string) {
	coll := ReportsCollection("OperatorDailyReport")
	dateStr := time.Now().Format("20060102")
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
	coll := ReportsCollection("OperatorDailyReport")
	dateStr := time.Now().Format("20060102")
	id := ut.JoinStr(':', dateStr, appId)

	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"appid", appId,
		),
		"$inc", db.D(
			"loginplrcount", 1,
		),
	)

	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}
