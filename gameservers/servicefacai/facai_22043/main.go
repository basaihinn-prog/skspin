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
	"serve/servicefacai/facai_22043/internal"
	"serve/servicefacai/facai_22043/internal/gendata"
	_ "serve/servicefacai/facai_22043/internal/rpc"
	"serve/servicefacai/facaicomm"
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
	facaicomm.RegistRpcToMQ(mqconn)

	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)
	lazy.Serve()
}
