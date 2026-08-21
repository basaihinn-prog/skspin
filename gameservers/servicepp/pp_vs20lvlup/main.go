package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/servicepp/pp_vs20lvlup/internal"
	"serve/servicepp/pp_vs20lvlup/internal/gamedata"
	"serve/servicepp/pp_vs20lvlup/internal/gendata"
	_ "serve/servicepp/pp_vs20lvlup/internal/rpc"
	"serve/servicepp/ppcomm"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	lazy.Init(internal.GameID)
	gamedata.Load()
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()

	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("simulate", "")
		lo.Must0(gendata.LoadCombineData())
		return
	}
	// 先准备好数据, 再监听消息
	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	ppcomm.RegistRpcToMQ(mqconn)
	// mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())
	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)
	lazy.Serve()
}

func test() {
	coll := db.Collection("simulate")

	for i := 66; i <= 85; i++ {
		findFilter := bson.M{
			"$and": []bson.M{
				{"bucketid": i},
				{"selected": true},
			},
		}
		findId := bson.M{"_id": true}
		cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(findId))
		if err != nil {
			log.Fatal(err)
		}
		var ids []primitive.ObjectID
		for cursor.Next(context.TODO()) {
			var result struct {
				Id primitive.ObjectID `bson:"_id,omitempty"`
			}
			err = cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			ids = append(ids, result.Id)
		}
		num := int(float64(len(ids)) * 0.10)
		if num != 0 {
			findFilterId := bson.M{"_id": bson.M{"$in": ids[:num]}}
			many, err := coll.DeleteMany(context.TODO(), findFilterId)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(many)
		}
	}
}
