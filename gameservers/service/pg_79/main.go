// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/pg_79 ./service/pg_79
package main

import (
	"fmt"
	"log/slog"
	"os"
	"serve/comm/redisx"

	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/service/pg_79/internal/config"
	"serve/service/pg_79/internal/gendata"
	_ "serve/service/pg_79/internal/rpc"

	"github.com/samber/lo"
)

// 1.金砖模式的时候不能连线
// 2.连线的时候可以最多出现4个金砖块
// hitLines没中的时候给一空数组
// 5个金砖 其它全是1了

const GameName = "pg_79"

// func main() {
// 	lazy.Init(GameName)

// 	core.LoadConfig(func() {
// 		fmt.Println("config加载完成了++++++++++++")
// 		// 先准备好数据, 再监听消息
// 		if gendata.NeedLoad() {
// 			fmt.Println("NeedLoad gendata.LoadCombineData++++++++++++")
// 			lo.Must0(gendata.LoadCombineData())
// 		}
// 	})

// 	mqconn := lo.Must(mq.ConnectServerMust())
// 	mux.RegistRpcToMQ(mqconn)
// 	mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())

// 	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
// 	db.DialToMongo(mgoaddr, lazy.ServiceName)

// 	httpAddr, _ := lazy.RouteFile.Get("pg_79.http.test")
// 	if len(httpAddr) != 0 {
// 		mux.StartHttpServer(httpAddr)
// 	}
// 	lazy.ServeWithoutGrpc()
// }

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("HttpInvoke", "HttpInvoke", r)
		}
	}()
	lazy.Init(GameName)
	// lazy.InitWithoutGrpc(GameName)
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()

	config.LoadConfig(func() {
		fmt.Println("config加载完成了++++++++++++")
		// 先准备好数据, 再监听消息
		if len(os.Args) > 1 && os.Args[1] == "cache" {
			slotsmongo.InitData("simulate", "")
			lo.Must0(gendata.LoadCombineData())
			os.Exit(0)
		}
		lo.Must0(gendata.LoadCombineData())
	})

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)
	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)

	httpAddr, _ := lazy.RouteFile.Get("pg_79.http.test")
	if len(httpAddr) != 0 {
		mux.StartHttpServer(httpAddr)
	}

	lazy.Serve()
}
