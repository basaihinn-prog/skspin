package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"serve/comm/db"
	"serve/comm/robot"
	"serve/servicepp/pp_vs10bbdoubled/internal"
	"serve/servicepp/pp_vs10bbdoubled/internal/gendata"
	"serve/servicepp/ppcomm"

	"go.mongodb.org/mongo-driver/bson"
)

func TestFetchSpin(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:aajohsujie9vecieSohqu4weivai7oxayei9rie5Yoh4vuojohwaothee0waethi@127.0.0.1:27017/?authSource=admin&directConnection=true"
	db.DialToMongo(mongoaddr, internal.GameID)

	coll := db.Collection2("robot", "game")
	coll.InsertOne(context.TODO(), robot.GameInfo{
		Id:        internal.GameID,
		Weight:    int64(10),
		StartTime: time.Now().Unix(),
		IsEnd:     false,
	})

	cnt, err := db.Collection("simulate").CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
		return
	}
	// 超过100条就不能在本地拉取数据
	if cnt >= 50 {
		fmt.Println("测试数据达到最大限制")
		return
	}
	go normal()
	go game()

	go func() {
		for ; ; time.Sleep(time.Second) {
			cnt, err := db.Collection("simulate").CountDocuments(context.TODO(), bson.M{})
			if err != nil {
				fmt.Println(err)
				return
			}
			if cnt >= 50 {
				fmt.Println("测试数据达到最大限制")
				os.Exit(-1)
			}
		}
	}()
	select {}
}

func normal() {
	fetcher := ppcomm.NewFetcher(internal.GameID, "testusertjh_01")
	r := &robot.Robot{}
	ppR := ppcomm.NewPPRobot(r, fetcher, &ppcomm.FetchDataParam{
		Game: internal.GameID,
		ArrC: map[float64]int{ // 游戏下注挡位，k：下注值，v:权重
			0.05: 50,
		},
		GBuckets: gendata.GBuckets,
		InsertC:  0.05,
		L:        internal.Line,
		Ty:       ppcomm.GameTypeNormal,
	})
	ppR.EndTime = time.Now().Unix() + 24*60*60
	ppR.Run()
}

func game() {
	fetcher := ppcomm.NewFetcher(internal.GameID, "testusertjhgame_01")
	r := &robot.Robot{}
	ppR := ppcomm.NewPPRobot(r, fetcher, &ppcomm.FetchDataParam{
		Game: internal.GameID,
		ArrC: map[float64]int{ // 游戏下注挡位，k：下注值，v:权重
			0.05: 50,
		},
		GBuckets: gendata.GBuckets,
		InsertC:  0.05,
		L:        internal.Line,
		Ty:       ppcomm.GameTypeGame,
	})
	ppR.EndTime = time.Now().Unix() + 24*60*60
	ppR.Run()
}
