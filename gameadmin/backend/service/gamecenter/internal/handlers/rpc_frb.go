package handlers

import (
	"context"
	"game/comm/db"
	"game/comm/mux"
	"time"
)

func init() {
	mux.RegRpc("/gamecenter/frb/next", "获取一个可用的frb", "frb", NextFRB, NextFRBPs{
		AppID:  "faketrans",
		UserID: "testuser1",
		GameID: "pp_vs20olympx",
		Remove: false,
	})
}

type NextFRBPs struct {
	AppID  string
	UserID string
	GameID string

	// Remove this record after fetch success
	Remove bool
}

type NextFRBRet = FRBPlayer

func NextFRB(ps NextFRBPs, ret *NextFRBRet) (err error) {
	coll := db.Collection2("game", "freeRoundBonus")

	now := time.Now()
	filter := db.D(
		"AppID", ps.AppID,
		"UserID", ps.UserID,
		"GameID", ps.GameID,
		"ExpirationDate", db.D(
			"$gt", now.Unix(),
		),
	)

	if ps.Remove {
		err = coll.FindOneAndDelete(context.TODO(), filter).Decode(ret)
	} else {
		err = coll.FindOne(context.TODO(), filter).Decode(ret)
	}
	return
}
