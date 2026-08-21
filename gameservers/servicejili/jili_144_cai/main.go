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
	"serve/servicejili/jili_144_cai/internal"
	"serve/servicejili/jili_144_cai/internal/gamedata"
	"serve/servicejili/jili_144_cai/internal/gendata"
	_ "serve/servicejili/jili_144_cai/internal/gendata"
	_ "serve/servicejili/jili_144_cai/internal/rpc"
	"serve/servicejili/jiliut"

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
	slotsmongo.RestoreClientCacheJILI(internal.GameShortName)

	// 先准备好数据, 再监听消息
	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("rawSpinData", "")
		lo.Must0(gendata.LoadCombineData())
		return
	}
	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	jiliut.RegistRpcToMQ(mqconn)
	// mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())
	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)

	// httpAddr, _ := lazy.RouteFile.Get("jili_2.http.test")
	// if httpAddr != "" {
	// 	mux.StartHttpServer(httpAddr)
	// }

	lazy.Serve()
}

func test() {
	coll := db.Collection("rawSpinData")

	arr := []int{48, 49, 50, 51, 52, 53, 54, 55, 56, 57}
	for _, i := range arr {
		findFilter := bson.M{
			"$and": []bson.M{
				{"type": 2},
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
		num := int(float64(len(ids)) * 0.1)
		findFilterId := bson.M{"_id": bson.M{"$in": ids[:num]}}
		many, err := coll.DeleteMany(context.TODO(), findFilterId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(many)
	}
}
