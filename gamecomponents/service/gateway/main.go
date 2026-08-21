package main

import (
	"fmt"
	"game/comm/db"
	"game/comm/mq"
	"game/duck/lazy"
	"game/duck/logger"
	"game/duck/mongodb"
	"game/duck/ut2"
	"strconv"
	"strings"

	"net/http"
	"net/url"

	"github.com/caddyserver/certmagic"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var client gamepb.GameRpcClient

func ConnectMongoByConfig() {
	mongoAddr := lo.Must(lazy.RouteFile.Get("mongo"))

	logger.Info("连接mongodb", mongoAddr)

	// err := DB.Connect("mongodb://" + mongoAddr)

	DB := mongodb.NewDB(lazy.ServiceName)
	err := DB.Connect(mongoAddr)

	if err != nil {
		logger.Info("连接mongodb失败", err)
	} else {
		logger.Info("连接mongodb成功")
		db.SetupClient(DB.Client)
	}
}
func main() {
	lazy.Init("gateway")

	// mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	// db.DialToMongo(mgoaddr, lazy.ServiceName)

	ConnectMongoByConfig()

	gw := &Gateway{players: ut2.NewSyncMap[int64, *Connection]()}

	mq.ConnectServerMust()
	mq.Subscribe("gateway", func(msg *mq.Forward) {
		for _, pid := range msg.Pids {
			player, ok := gw.players.Load(pid)
			if ok {
				player.SendMsg(msg.Data)
			}
		}
	})

	// client = lazy.NewGrpcClient("gamecenter", gamepb.NewGameRpcClient)

	lazy.ServeFn(func() { httpLoginRun(gw) })
}

func httpLoginRun(g *Gateway) {

	mux := http.NewServeMux()
	mux.HandleFunc("/login", g.route_login)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	addr := lazy.GetAddr("gateway.http")
	logger.Info("启动", addr)

	var err error

	if !strings.HasPrefix(addr, "https") {
		// gateway.http = localhost:8080

		port := lazy.GetPortMust("gateway.http")
		err = http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
	} else {
		// gateway.http = https://baidu.com?https=8443&http=8080&email=auto@baidu.com
		url, err := url.Parse(addr)
		if err != nil {
			panic(err)
		}
		https, err := strconv.Atoi(url.Query().Get("https"))
		if err != nil {
			panic(err)
		}
		http, err := strconv.Atoi(url.Query().Get("http"))
		if err != nil {
			panic(err)
		}

		domains := []string{url.Host}
		email := url.Query().Get("email")

		// Tips: 验证https 需要从80端口转发到 监听端口
		certmagic.HTTPSPort = https
		certmagic.HTTPPort = http

		certmagic.DefaultACME.Email = email
		certmagic.DefaultACME.Agreed = true

		// 把etcd的log 写入到我们的logger中
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core := zapcore.NewCore(encoder, zapcore.AddSync(logger.DefaultLogger), zapcore.DebugLevel)
		certmagic.DefaultACME.Logger = zap.New(core)

		logger.Info("HTTPSPort", https, "HTTPPort", http, email, domains)

		err = certmagic.HTTPS(domains, mux)
	}

	if err != nil {
		panic(err)
	}
}
