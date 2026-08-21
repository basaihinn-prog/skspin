package slotsmongo

import (
	"context"
	"game/comm/db"
	"game/comm/ut"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

type PlayerDailyDoc struct {
	ID         string `json:"-" bson:"_id"`
	BetAmount  int64  `bson:"betamount"` // 支出
	WinAmount  int64  `bson:"winamount"` // 赚取
	SpinCount  int    `bson:"spincount"` // 转动次数
	RoundCount int    `bson:"roundcount"`

	Date  string `bson:"date"`
	AppID string `bson:"appid"`
	Pid   int64  `bson:"pid"`
}

func addPayerReport(dateStr, appId string, pid int64, incDoc bson.D) {
	coll := ReportsCollection("PlayerDailyReport")

	// dateStr := time.Now().Format("20060102")
	id := ut.JoinStr(':', dateStr, appId, strconv.Itoa(int(pid)))

	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"appid", appId,
			"pid", pid,
		),
		"$inc", incDoc,
	)

	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}
