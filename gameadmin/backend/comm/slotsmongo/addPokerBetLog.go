package slotsmongo

import (
	"context"
	"game/comm/db"
	"game/comm/ut"
	"go.mongodb.org/mongo-driver/bson"
)

func AddPokerBetDailyReport(dateStr, game string, appId string, incDoc bson.D) {
	coll := ReportsCollection("BetPokerDailyReport")
	id := ut.JoinStr(':', dateStr, game, appId)
	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"game", game,
			"appid", appId,
		),
		"$inc", incDoc,
	)
	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}
