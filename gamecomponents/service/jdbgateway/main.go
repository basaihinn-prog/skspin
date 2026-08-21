package main

import (
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	"game/service/jdbgateway/internal/gamedata"
	"game/service/jdbgateway/internal/staticproxy"
	"log"

	"github.com/samber/lo"
)

func main() {
	lazy.InitWithoutGrpc("jdbgateway")
	gamedata.InitConfig()

	mqconn := lo.Must(mq.ConnectServerMust())
	mux.RegistRpcToMQ(mqconn)
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	httpAddrProxy := lo.Must(lazy.RouteFile.Get("jdbgateway.http.proxy"))
	httpAddrProxy2 := lo.Must(lazy.RouteFile.Get("jdbgateway.http.proxy2"))
	httpAddrProxy3 := lo.Must(lazy.RouteFile.Get("jdbgateway.http.proxy3"))

	staticproxy.StartProxy(httpAddrProxy)
	staticproxy.StartProxy(httpAddrProxy2)
	staticproxy.StartProxy(httpAddrProxy3)

	log.Println("start...")
	lazy.ServeWithoutGrpc()
}
