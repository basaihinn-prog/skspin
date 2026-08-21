package public

import (
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/slotsmongo"
)

func AddBetLog(ps *slotsmongo.AddBetLogParams) {
	if ps.ServiceName == "" {
		ps.ServiceName = lazy.ServiceName
	}
	mq.JsonNC.Publish("/betlog/addBetLog", ps)
}
