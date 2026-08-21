package main

import (
	"errors"
	"game/duck/lazy"
	"time"

	"game/duck/logger"

	"github.com/nats-io/nats-server/v2/server"
)

func main() {
	lazy.Init("proxy")

	err := StartNats()
	if err != nil {
		logger.Info("启动失败", err)
	}

	// lazy.Serve()
	lazy.SignalProc()
}

// go get -d github.com/nats-io/nats-server/v2
// https://dev.to/karanpratapsingh/embedding-nats-in-go-19o

func StartNats() error {
	port := lazy.GetPortMust("proxy.mq")

	opts := &server.Options{
		Host: "0.0.0.0",
		Port: port,
		// 开启日志
		// Debug: true,
		// Trace: true,
	}
	ns, err := server.NewServer(opts)
	ns.ConfigureLogger()

	if err != nil {
		return err
	}

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		return errors.New("not ready for connection")
	}

	logger.Info("nats 启动成功", ns.ClientURL())

	return nil
}
