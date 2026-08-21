package main

import (
	"game/comm/slotscall"
	"game/duck/lazy"
	"game/pb/_gen/pb/stats"
)

func main() {
	lazy.Init("stats")

	ConnectMongoByConfig()

	slotscall.RpcCallInit()

	stats.RegisterAdminStatsRpcServer(lazy.GrpcServer, &StatsRpcServer{})

	lazy.Serve()
}
