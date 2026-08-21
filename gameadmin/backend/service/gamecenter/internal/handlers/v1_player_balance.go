package handlers

import (
	"errors"
	"game/comm"
	"game/comm/mux"
	"game/comm/ut"
	"game/service/gamecenter/internal/operator"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/balance",
		Handler:      v1_player_balance,
		Desc:         "获取玩家余额",
		Kind:         "api/v1",
		ParamsSample: v1PlayerBalancePs{"user_id"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v1PlayerBalancePs struct {
	UserID string
}

type v1PlayerBalanceRet struct {
	Balance          float64
	Unsettled        float64
	UnsettledDetails map[string]float64
}

func v1_player_balance(app *operator.MemApp, ps v1PlayerBalancePs, ret *v1PlayerBalanceRet) (err error) {
	if app.WalletMode != comm.WalletModeTransfer {
		err = errors.New("Not transfer Mode")
		return
	}

	memplr, err := operator.AppMgr.GetPlr2(app, ps.UserID)
	if err != nil {
		return
	}

	err = getBalanceByMemplr(memplr, ret)
	return
}

func getBalanceByMemplr(memplr *operator.MemPlr, ret *v1PlayerBalanceRet) (err error) {
	balance, err := memplr.Balance()
	ret.Balance = ut.Gold2Money(balance)

	ret.UnsettledDetails = map[string]float64{}

	// var unSettlementRet struct {
	// 	Balance int64
	// }

	// mq.Invoke("/lotter/admin/unSettlement", map[string]any{"Pid": memplr.Pid}, &unSettlementRet)
	// if unSettlementRet.Balance != 0 {
	// 	b := ut.Gold2Money(unSettlementRet.Balance)
	// 	ret.UnsettledDetails["lottery"] = b
	// 	ret.Unsettled += b
	// }

	// unSettlementRet.Balance = 0
	// mq.Invoke("/hilo/player/getUnsettledBalance", map[string]any{"Pid": memplr.Pid}, &unSettlementRet)
	// if unSettlementRet.Balance != 0 {
	// 	b := ut.Gold2Money(unSettlementRet.Balance)
	// 	ret.UnsettledDetails["Hilo"] = b
	// 	ret.Unsettled += b
	// }

	return
}
