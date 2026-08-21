package main

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/servicepp/ppcomm"
	"serve/servicepp/robot_dispatch_center_hgc/internal"

	"github.com/samber/lo"
)

func main() {
	lazy.Init("robots_center_hgc")
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	ppcomm.RegistRpcToMQ(mqconn)
	go internal.RobotsDispatch()
	go internal.TasksDispatchHGC()
	//go internal.LoopUpdatePlayerSpinCountTo100(mgoaddr)
	lazy.Serve()
}
