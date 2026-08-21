package handlers

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/ut"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"sort"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/pendinggames",
		Handler:      v1_player_pending_games,
		Desc:         "获取玩家未结算游戏列表&当前玩家余额(单一钱包为0)",
		Kind:         "api/v1",
		ParamsSample: v1PlayerBalancePs{"user_id"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type pendingGames struct {
	Balance      float64
	PendingGames []string
}

func v1_player_pending_games(app *operator.MemApp, ps v1PlayerBalancePs, ret *pendingGames) (err error) {
	coll := gcdb.CollPlayers.Coll()

	ret.PendingGames = []string{}

	memplr, err := operator.AppMgr.GetPlr2(app, ps.UserID)
	if err != nil {
		return
	}

	if app.WalletMode == comm.WalletModeTransfer {
		balance, _ := memplr.Balance()
		ret.Balance = ut.Gold2Money(balance)
	}

	opts := options.FindOne().SetProjection(db.D("CompletedGames", 1))

	var doc operator.DocPlayer
	err = coll.FindOne(context.TODO(), db.ID(memplr.Pid), opts).Decode(&doc)
	if err != nil {
		return
	}
	for game, completed := range doc.CompletedGames {
		if !completed {
			ret.PendingGames = append(ret.PendingGames, game)
		}
	}

	sort.Strings(ret.PendingGames)
	return
}
