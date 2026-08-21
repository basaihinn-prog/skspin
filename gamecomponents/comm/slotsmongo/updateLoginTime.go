package slotsmongo

import (
	"context"
	"game/comm/db"
	"game/comm/ut"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdatePlrLoginTime(pid int64, gameID string) (err error) {
	// coll := gcdb.CollPlayers.Coll()
	coll := db.Collection2("game", "Players")

	type Player struct {
		Id       int64     `json:"Pid" bson:"_id"`
		Uid      string    `json:"Uid" bson:"Uid"`
		AppID    string    `json:"AppID" bson:"AppID" md:"运营商"`
		LoginAt  time.Time `json:"LoginAt" bson:"LoginAt"`
		CreateAt time.Time `json:"CreateAt" bson:"CreateAt"`
	}

	var (
		plr Player
		now = time.Now()
	)

	update := db.D("$set", db.D("LoginAt", now, "LastGame", gameID))

	err = coll.FindOneAndUpdate(context.TODO(), db.ID(pid), update, options.FindOneAndUpdate().SetReturnDocument(options.Before).SetProjection(db.D("LoginAt", 1, "AppID", 1))).Decode(&plr)

	if err != nil {
		return
	}

	if !ut.IsSameDate(plr.LoginAt.Local(), now) {
		IncLoginCount(plr.AppID)
	}

	updateEnterPlrCount(gameID, plr.AppID, pid)

	return nil
}

func updateEnterPlrCount(game string, appId string, pid int64) {
	if !IsFirstVisitToday(ut.JoinStr(':', "enter", game, strconv.Itoa(int(pid)))) {
		return
	}

	dateStr := time.Now().UTC().Format("20060102")
	incDoc := db.D("enterplrcount", 1)

	addBetDailyReport(dateStr, game, appId, incDoc)
	addOperatorReport(dateStr, appId, incDoc)
	//addPayerReport(dateStr, appId, pid, incDoc)
	addBetSummaryReport(game, appId, incDoc)
}
