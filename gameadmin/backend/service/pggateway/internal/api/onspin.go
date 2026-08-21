package api

import (
	"context"
	"encoding/json"
	"game/comm/db"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"log/slog"
	"strings"
)

func onSpin(r *PGParams, respdata json.RawMessage) {
	game, ok := isSpinPath(r.Path)
	if !ok {
		return
	}

	var data struct {
		Si struct {
			Sid string
		}
	}

	err := json.Unmarshal(respdata, &data)
	if err != nil {
		return
	}

	// 6618db8d8a3189557283f886_28_28
	arr := strings.Split(data.Si.Sid, "_")
	if len(arr) != 3 {
		return
	}

	completed := arr[1] == arr[2]

	udpatePlayerGameStatus(&slotsmongo.DocBetLog{
		Pid:       r.Pid,
		GameID:    game,
		Completed: completed,
	})
}

func isSpinPath(pth string) (game string, ok bool) {
	pth = strings.ToLower(pth)
	arr := strings.Split(pth, "/")
	ok = arr[0] == "" &&
		arr[1] == "game-api" &&
		arr[4] == "spin"

	if ok {
		game = "pg_" + arr[2]
	}

	return
}

func udpatePlayerGameStatus(betlog *slotsmongo.DocBetLog) {
	coll := db.Collection2("game", "Players")
	update := db.D("$set", db.D("CompletedGames."+betlog.GameID, betlog.Completed))
	_, err := coll.UpdateByID(context.TODO(), betlog.Pid, update)
	if err != nil {
		// mongo.WriteError
		filter := db.D("_id", betlog.Pid, "CompletedGames", nil)
		update := db.D(
			"$set", db.D("CompletedGames", db.D(
				betlog.GameID, betlog.Completed),
			),
		)
		_, err := coll.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			slog.Error("udpatePlayerGameStatus", "err", ut.ErrString(err))
		}
	}
}
