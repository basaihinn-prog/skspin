package main

import (
	"game/comm/db"
	"game/duck/lazy"
	"game/duck/logger"
	"game/duck/mongodb"
	"github.com/samber/lo"
)

var DB = mongodb.NewDB("reports")

func ConnectMongoByConfig() {
	mongoAddr := lo.Must(lazy.RouteFile.Get("mongo"))

	logger.Info("连接mongodb", mongoAddr)

	err := DB.Connect(mongoAddr)

	if err != nil {
		logger.Info("连接mongodb失败", err)
	} else {
		logger.Info("连接mongodb成功")
		db.SetupClient(DB.Client)
	}
}
