package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/service/betlog/internal/logStation"
	"syscall"
	"time"

	"github.com/samber/lo"
)

func main() {
	lazy.ServiceName = "betlog"
	lazy.RouteFile = lazy.NewFilePortProvider()

	clickHouseAddr := lo.Must(lazy.RouteFile.Get("clickhouse"))

	//conn, err := sql.Open("clickhouse", clickHouseAddr)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer conn.Close()

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)
	ctx := context.Background()
	logService, err := logStation.NewLogService(ctx, clickHouseAddr, 100000, 15, 10*time.Second, 500*time.Millisecond)
	if err != nil {
		log.Fatalf("Failed to initialize log service: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("log Service started")
	logService.Start(c)
	// Block until a signal is received.
	os.Exit(0)

	//lazy.SignalProc()

}
