package main

import (
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	"game/service/minigateway/internal/gamedata"
	"game/service/minigateway/internal/staticproxy"
	"log"

	"github.com/samber/lo"
)

func main() {
	lazy.InitWithoutGrpc("minigateway")
	gamedata.InitConfig()

	mqconn := lo.Must(mq.ConnectServerMust())
	mux.RegistRpcToMQ(mqconn)

	httpAddrProxy := lo.Must(lazy.RouteFile.Get("minigateway.http.proxy"))

	staticproxy.StartProxy(httpAddrProxy)

	log.Println("start...")
	lazy.ServeWithoutGrpc()
}
