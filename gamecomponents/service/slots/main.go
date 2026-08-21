package main

import (
	"game/duck/lazy"
	"game/pb/_gen/pb/slots"
)

func main() {
	lazy.Init("slots")

	initConfig()

	ConnectMongoByConfig()

	slots.RegisterSlotsPoolRpcServer(lazy.GrpcServer, &SlotsPoolRpcServer{})
	slots.RegisterSlotsConfigRpcServer(lazy.GrpcServer, &ConfigRpcServer{})
	slots.RegisterSlotsRuntimeRpcServer(lazy.GrpcServer, &RuntimeRpcServer{})
	slots.RegisterSlotsGameDayStatRpcServer(lazy.GrpcServer, &GameDayStatRpc{})
	slots.RegisterSlotsPlayerGameDayStatRpcServer(lazy.GrpcServer, &PlayerGameDayStatRpc{})
	slots.RegisterSlotsPlayerRotateDataRpcServer(lazy.GrpcServer, &PlayerRotateRpcServer{})
	slots.RegisterSlotsStatsBuyGameRpcServer(lazy.GrpcServer, &SlotsStatsBuyGameRpcServer{})
	slots.RegisterAdminSlotsRpcServer(lazy.GrpcServer, &AdminRpcServer{})

	lazy.Serve()
}
