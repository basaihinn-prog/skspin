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
	"serve/service/pg_33/internal/gamedata"
	"serve/service/pg_33/internal/gendata"

	_ "serve/service/pg_33/internal/rpc"

	"github.com/samber/lo"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("HttpInvoke", "HttpInvoke", r)
		}
	}()
	lazy.Init("pg_33")
	gamedata.Load()

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()

	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("simulate", "pg_33")
		lo.Must0(gendata.LoadCombineData())
		return
	}
	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)
	//mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())
	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)

	httpAddr, _ := lazy.RouteFile.Get("pg_33.http.test")
	if len(httpAddr) != 0 {
		mux.StartHttpServer(httpAddr)
	}
	lazy.Serve()
}
