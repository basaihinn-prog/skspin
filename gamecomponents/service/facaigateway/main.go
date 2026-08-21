package main

import (
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	"game/service/facaigateway/internal/gamedata"
	"game/service/facaigateway/internal/staticproxy"
	"log"

	"github.com/samber/lo"
)

func main() {
	lazy.InitWithoutGrpc("facaigateway")
	gamedata.InitConfig()
	//bytes, _ := base64.StdEncoding.DecodeString("KJotb6gIlXTPYEHhnhnYee8Wq7IoQYVHmkfgT61oH28Gyv+JoBYkIhl9OyyI59a2E0EPYrd5EYX/wal3nmxPgIOAa/gOG4iYIQE0jzxigePRJHMCvCnlRrO4MVxKtheD+QbNthyrh0PbZfZzWRyS9Ot1j9G7w/J5eIrXVW/599yfFpKjtRgTauBtxh0kq8bb145Fec3I2x6fHQwjWscJRkjc8uxuyfjvQk7pB8JrDZLxht7IXv0OHzT9nrJDWaAr6uco87ANO8qLfXZHJ9uFsA==")
	//fmt.Printf("%x\n", bytes)
	//so, err := slotsmongo.NewFromBinaryData(bytes)
	//fmt.Println(so, err)
	mqconn := lo.Must(mq.ConnectServerMust())
	mux.RegistRpcToMQ(mqconn)
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	httpAddrProxy := lo.Must(lazy.RouteFile.Get("facaigateway.http.proxy"))
	staticproxy.StartProxy(httpAddrProxy)

	log.Println("start...")
	lazy.ServeWithoutGrpc()
}
