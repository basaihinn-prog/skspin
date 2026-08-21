package main

import (
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/service/pg_89/internal/gamedata"
	"serve/service/pg_89/internal/gendata"

	_ "serve/service/pg_89/internal/rpc"

	"github.com/samber/lo"
)

func main() {
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println(r)
	//		slog.Error("HttpInvoke", "HttpInvoke", r)
	//	}
	//}()

	lazy.Init("pg_89")
	gamedata.Load()

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()

	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("simulate", "")
		lo.Must0(gendata.LoadCombineData())
		return
	}
	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)
	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)
	httpAddr, _ := lazy.RouteFile.Get("pg_89.http.test")
	if len(httpAddr) != 0 {
		mux.StartHttpServer(httpAddr)
	}

	lazy.Serve()
}
