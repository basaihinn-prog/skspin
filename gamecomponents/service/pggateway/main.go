package main

import (
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	"log/slog"

	"game/service/pggateway/internal/api"
	"game/service/pggateway/internal/gamedata"
	_ "game/service/pggateway/internal/rpc"
	"game/service/pggateway/internal/staticproxy"
	_ "game/service/pggateway/platpg"

	"github.com/samber/lo"
)

func main() {
	lazy.InitWithoutGrpc("pggateway")
	gamedata.InitConfig()
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)

	mqconn := lo.Must(mq.ConnectServerMust())
	mux.RegistRpcToMQ(mqconn)
	// mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())

	httpAddr := lo.Must(lazy.RouteFile.Get("pggateway.http"))
	if httpAddr != "" {
		mux.StartHttpServer(httpAddr)
	}
	httpAddrProxy := lo.Must(lazy.RouteFile.Get("pggateway.http.proxy"))
	staticproxy.StartProxy(httpAddrProxy)
	httpAddrVerify := lo.Must(lazy.RouteFile.Get("pggateway.http.verify"))
	staticproxy.StartProxyVerify(httpAddrVerify)
	httpAddrApi := lo.Must(lazy.RouteFile.Get("pggateway.http.api"))
	api.StartApi(httpAddrApi)
	slog.Info("start listen",
		"static proxy", httpAddrProxy,
		"api proxy", httpAddrApi,
	)
	lazy.ServeWithoutGrpc()
}
