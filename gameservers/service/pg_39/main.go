package main

import (
	"log/slog"
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/service/pg_39/internal/gamedata"
	"serve/service/pg_39/internal/gendata"
	_ "serve/service/pg_39/internal/rpc"

	"github.com/samber/lo"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("HttpInvoke", "HttpInvoke", r)
		}
	}()
	lazy.Init("pg_39")
	gamedata.Load()
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()

	// 先准备好数据, 再监听消息

	if len(os.Args) > 1 && os.Args[1] == "cache" {
		//slotsmongo.InitClassData2("pgSpinData")
		slotsmongo.InitData("pgSpinData", "pg_39")
		lo.Must0(gendata.LoadCombineData())
		return
	}
	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	// mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())
	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)

	httpAddr, _ := lazy.RouteFile.Get("pg_39.http.test")
	if httpAddr != "" {
		mux.StartHttpServer(httpAddr)
	}

	lazy.Serve()
}
