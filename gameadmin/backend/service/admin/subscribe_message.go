package main

import (
	"game/comm"
	"game/comm/mq"
)

func RegSub() {
	mq.SubscribeMsg("/games/operator_balance_update", func(ps *comm.OperatorBalance) {
		UpdateOperatorBalance(ps)
	})
}
