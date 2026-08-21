package main

import (
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/servicejili/jili_26_sweetheart3/internal"
	"serve/servicejili/jili_26_sweetheart3/internal/gamedata"
	"serve/servicejili/jili_26_sweetheart3/internal/gendata"
	_ "serve/servicejili/jili_26_sweetheart3/internal/rpc"
	"serve/servicejili/jiliut"

	"github.com/samber/lo"
)

func main() {
	lazy.Init(internal.GameID)
	gamedata.Load()
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()
	slotsmongo.RestoreClientCacheJILI(internal.GameShortName)

	// 先准备好数据, 再监听消息
	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("rawSpinData", "")
		lo.Must0(gendata.LoadCombineData())
		return
	}
	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	jiliut.RegistRpcToMQ(mqconn)
	// mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())
	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)

	// httpAddr, _ := lazy.RouteFile.Get("jili_2.http.test")
	// if httpAddr != "" {
	// 	mux.StartHttpServer(httpAddr)
	// }

	lazy.Serve()
}
