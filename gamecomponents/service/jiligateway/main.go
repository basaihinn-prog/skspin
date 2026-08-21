package main

import (
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"

	"game/service/jiligateway/internal/gamedata"
	_ "game/service/jiligateway/internal/rpc"
	"game/service/jiligateway/internal/staticproxy"

	"github.com/samber/lo"
)

func main() {
	lazy.InitWithoutGrpc("jiligateway")
	gamedata.InitConfig()

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)

	mqconn := lo.Must(mq.ConnectServerMust())
	mux.RegistRpcToMQ(mqconn)

	httpAddrProxy := lo.Must(lazy.RouteFile.Get("jiligateway.http.proxy"))

	staticproxy.StartProxy(httpAddrProxy)

	httpAddr, _ := lazy.RouteFile.Get("jiligateway.http")
	if httpAddr != "" {
		mux.StartHttpServer(httpAddr)
	}

	// log.Println("start...")
	lazy.ServeWithoutGrpc()
}

// [JL2]
// AgentKey = "7a7aca504fa4dc16e1e0ad7292789c12180a6794"
// AgentId = "RT237_THB"
// Url = "https://uat-wb-api.jlfafafa2.com/api1/"
