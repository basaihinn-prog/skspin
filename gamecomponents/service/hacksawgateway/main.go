package main

import (
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	"log"

	"game/service/hacksawgateway/internal/gamedata"
	"game/service/hacksawgateway/internal/staticproxy"
	"github.com/samber/lo"
)

func main() {
	lazy.InitWithoutGrpc("hacksawgateway")
	gamedata.InitConfig()

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)

	mqconn := lo.Must(mq.ConnectServerMust())
	mux.RegistRpcToMQ(mqconn)

	httpAddrProxy := lo.Must(lazy.RouteFile.Get("hacksawgateway.http.proxy"))

	staticproxy.StartProxy(httpAddrProxy)

	//httpAddrProxy2 := lo.Must(lazy.RouteFile.Get("ppgateway.http.proxy2"))
	//if httpAddrProxy2 != "" {
	//	staticproxy.StartReverseProxy(httpAddrProxy2)
	//}

	//httpAddr := lo.Must(lazy.RouteFile.Get("ppgateway.http.test"))
	//if httpAddr != "" {
	//	mux.StartHttpServer(httpAddr)
	//}

	log.Println("start...")
	lazy.ServeWithoutGrpc()
}
