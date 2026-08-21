package main

import (
	"game/comm/db"
	"game/duck/lazy"
	"game/duck/logger"
	"game/duck/mongodb"

	"github.com/samber/lo"
)

var (
	DB = mongodb.NewDB("slots")

	CollPool              = DB.Collection("Pool")              // slots的公共奖池
	CollRuntimePlayer     = DB.Collection("RuntimePlayer")     // slots的运行时统计(每个玩家)
	CollRuntimeDay        = DB.Collection("RuntimeDay")        // slots的运行时统计(每日)
	CollRuntimeSummary    = DB.Collection("RuntimeSummary")    // slots的运行时统计(所有总和)
	CollPlayerGameDayStat = DB.Collection("PlayerGameDayStat") // slots的各个玩家的每日统计
	CollGameDayStat       = DB.Collection("GameDayStat")       // slots的所有玩家的每日统计
	CollStatsBuyGame      = DB.Collection("StatsBuyGame")      // slots的购买小游戏的统计
	CollPlayerRotateData  = DB.Collection("PlayerRotateData")  // slots的玩家在各个游戏的转动次数
)

// func ConnectMongoByConfig() {
// 	mongoAddr := lazy.GetAddr(lazy.ServiceName + ".mongo")
// 	if len(mongoAddr) == 0 {
// 		mongoAddr = "127.0.0.1:27017"
// 	}

// 	logger.Info("连接mongodb", mongoAddr)

// 	err := DB.Connect("mongodb://" + mongoAddr)

// 	if err != nil {
// 		logger.Info("连接mongodb失败", err)
// 	} else {
// 		logger.Info("连接mongodb成功")
// 	}

// 	createDBIndex()
// }

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
	createDBIndex()
}

// 创建数据库索引
func createDBIndex() {
	CollRuntimeDay.CreateIndexByKeys("Day", "Game")
}
