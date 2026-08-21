package rpc

import (
	"game/comm/mq"
	"game/comm/mux"
	"game/comm/slotsmongo"
)

func init() {
	mux.RegRpc("/alerter/betlog/add", "下注结算完成调用", "betlog", addBetLog, nil)

}

func addBetLog(ps *slotsmongo.DocBetLog, _ *mux.EmptyResult) (err error) {
	// err = alert.onAddBetLog_old(ps)
	return
}

func RegSub() {
	mq.JsonNC.Subscribe("/games/betlog", func(ps *slotsmongo.DocBetLog) {
		alert.onAddBetLog(ps)
	})

	mq.JsonNC.Subscribe("/games/minibetlog", func(ps *slotsmongo.DocBetLogAviator) {
		alert.onAddMiniBetLog(ps)
	})

	mq.JsonNC.Subscribe("/games/lottery/openPrize", func(lotteryAlert *slotsmongo.LotteryAlert) {
		alert.onOpenLottery(lotteryAlert)
	})

	mq.JsonNC.Subscribe("/alerter/onTransferOut", onTransferOut)
	mq.JsonNC.Subscribe("/alerter/OperatorBalanceAlert", operatorBalanceAlert)
}
