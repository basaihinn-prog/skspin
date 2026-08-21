package main

import (
	"github.com/samber/lo"
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/servicehacksaw/hacksaw_1067/internal"
	"serve/servicehacksaw/hacksaw_1067/internal/gendata"
	_ "serve/servicehacksaw/hacksaw_1067/internal/rpc"
	"serve/servicehacksaw/hacksawcomm"
)

func main() {
	lazy.Init(internal.GameID)

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()
	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("simulate", "")
		lo.Must0(gendata.LoadCombineData())
		return
	}

	// 先准备好数据, 再监听消息
	lo.Must0(gendata.LoadCombineData())
	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	hacksawcomm.RegistRpcToMQ(mqconn)

	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)
	lazy.Serve()
}
