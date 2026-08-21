package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"serve/comm/db"
	"serve/comm/robot"
	"serve/servicepp/pp_vs1024mahjwins/internal"
	"serve/servicepp/pp_vs1024mahjwins/internal/gendata"
	"serve/servicepp/ppcomm"
)

func main() {
	var (
		game      = internal.GameID // 游戏ID
		line      = internal.Line   // 游戏line, 参照游戏发送doSpin请求里面的‘l’值
		insertC   = 0.05            // 插入数据库时的投注，参照游戏发送doSpin请求里面的‘c’值
		normalCnt = 5000            // 需要拉取普通盘的次数
		gameCnt   = 100             // 需要拉取游戏的次数
		// 使用最低3个挡位
		arrC = map[float64]int{ // 游戏下注挡位，k：下注值，v:权重
			0.05: 50,
			0.1:  10,
			0.15: 10,
			0.2:  5,
			0.25: 2,
		}
		gameWeight = 50 // 如果需要更多机器人拉数据，可以把这个值改大点，建议不修改
		GBuckets   = gendata.GBuckets
	)
	//mongoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@127.0.0.1:27017/?authSource=admin"
	fmt.Println(mongoaddr)
	db.DialToMongo(mongoaddr, game)

	coll := db.Collection2("robot", "game")
	coll.InsertOne(context.TODO(), robot.GameInfo{
		Id:        primitive.NewObjectID().Hex(),
		GameId:    game,
		Weight:    int64(gameWeight),
		StartTime: time.Now().Unix(),
		IsEnd:     false,
	})

	go ppcomm.FetchData(&ppcomm.FetchDataParam{
		Game:     game,
		ArrC:     arrC,
		GBuckets: GBuckets,
		InsertC:  insertC,
		L:        line,
		Ty:       ppcomm.GameTypeNormal,
	})
	go ppcomm.FetchData(&ppcomm.FetchDataParam{
		Game:     game,
		ArrC:     arrC,
		GBuckets: GBuckets,
		InsertC:  insertC,
		L:        line,
		Ty:       ppcomm.GameTypeGame,
	})
	needMap := map[int]int64{
		ppcomm.GameTypeNormal: int64(normalCnt),
		ppcomm.GameTypeGame:   int64(gameCnt),
	}
	go ppcomm.CheckGameCnt(fmt.Sprintf("pp_%s", game), needMap)
	select {}
}
