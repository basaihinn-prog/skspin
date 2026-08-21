package main

import (
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	"game/service/alerter/internal/gamedata"
	"game/service/alerter/internal/rpc"

	"github.com/samber/lo"
)

func main() {
	lazy.InitWithoutGrpc("alerter")
	gamedata.Load()
	gamedata.LoadLocalization()
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)

	mqconn := lo.Must(mq.ConnectServerMust())
	mux.RegistRpcToMQ(mqconn)
	rpc.RegSub()

	httpAddr, _ := lazy.RouteFile.Get("alerter.http.test")
	if len(httpAddr) != 0 {
		mux.StartHttpServer(httpAddr)
	}
	lazy.ServeFn(rpc.Start)
}
