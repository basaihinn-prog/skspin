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
	"serve/servicepp/pp_vs15godsofwar/internal"
	"serve/servicepp/pp_vs15godsofwar/internal/gamedata"
	"serve/servicepp/pp_vs15godsofwar/internal/gendata"
	_ "serve/servicepp/pp_vs15godsofwar/internal/rpc"
	"serve/servicepp/ppcomm"

	"github.com/samber/lo"
)

func main() {
	lazy.Init(internal.GameID)
	gamedata.Load()

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()
	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("simulate", "pp_vs15godsofwar")
		lo.Must0(gendata.LoadCombineData())
		return
	}

	// 先准备好数据, 再监听消息
	err := gendata.LoadCombineData()
	if err != nil {
		slog.Error("err")
	}
	//lo.Must0()

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	ppcomm.RegistRpcToMQ(mqconn)
	// mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())
	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)
	lazy.Serve()
}
