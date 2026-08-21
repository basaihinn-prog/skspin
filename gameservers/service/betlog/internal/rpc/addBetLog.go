package rpc

import (
	"serve/comm/mq"
	"serve/comm/slotsmongo"
)

func addBetLog(betLog *slotsmongo.AddBetLogParams) {

	//slotsmongo.AddBetLog(betLog)
}

func RegSub() {
	mq.JsonNC.Subscribe("/betlog/addBetLog", addBetLog)
}
