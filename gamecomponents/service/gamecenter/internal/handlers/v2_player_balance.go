package handlers

import (
	"game/comm"
	"game/comm/define"
	"game/service/gamecenter/internal/operator"
)

func init() {
	// mux.DefaultRpcMux.Add(&mux.PHandler{
	// 	Path:         "/api/v2/player/balance",
	// 	Handler:      v1_player_balance,
	// 	Desc:         "获取玩家余额&待完成的游戏",
	// 	Kind:         "api/v1",
	// 	ParamsSample: v1PlayerBalancePs{"user_id"},
	// 	Class:        "operator",
	// 	GetArg0:      getArg0,
	// })
}

// type v1PlayerBalancePs struct {
// 	UserID string
// }

type v2PlayerBalanceRet struct {
	Balance        float64
	UnsettledGames []string
}

func v2_player_balance(app *operator.MemApp, ps v1PlayerBalancePs, ret *v1PlayerBalanceRet) (err error) {
	if app.WalletMode != comm.WalletModeTransfer {
		err = define.NewErrCode("Not transfer Mode", 1009)
		return
	}

	memplr, err := operator.AppMgr.GetPlr2(app, ps.UserID)
	if err != nil {
		return
	}

	err = getBalanceByMemplr(memplr, ret)

	return
}
