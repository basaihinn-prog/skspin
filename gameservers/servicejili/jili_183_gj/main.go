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
	"serve/servicejili/jili_183_gj/internal"
	"serve/servicejili/jili_183_gj/internal/gamedata"
	"serve/servicejili/jili_183_gj/internal/gendata"
	_ "serve/servicejili/jili_183_gj/internal/rpc"
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
		slotsmongo.InitData("rawSpinData", internal.GameID)
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

func tmp() {
	coll := db.Collection("rawSpinData")

	//resetFilter := bson.M{
	//	"$and": []bson.M{
	//		{"type": 2},
	//		{"selected": false},
	//		{"bucketid": bson.M{"$gte": 19}},
	//		{"bucketid": bson.M{"$lte": 26}},
	//	},
	//}
	//resetUpdate := bson.M{"$set": bson.M{"selected": true}}
	//_, err := coll.UpdateMany(context.TODO(), resetFilter, resetUpdate)
	//if err != nil {
	//	log.Fatal(err)
	//}

	findFilter := bson.M{
		"$and": []bson.M{
			{"type": 2},
			{"selected": true},
			{"bucketid": bson.M{"$gte": 30}},
			{"bucketid": bson.M{"$lte": 46}},
		},
	}
	projection := bson.M{"_id": true}

	cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(projection))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

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

	var cancelGather []primitive.ObjectID
	cancelGatherNum := (len(ids) * 3) / 100

	cancelGather = ids[:cancelGatherNum]

	fmt.Println("ids", ids, len(ids))
	fmt.Println("cancelGather", cancelGather, len(cancelGather))

	updateFilter := bson.M{"_id": bson.M{"$in": cancelGather}}
	updateSet := bson.M{"$set": bson.M{"selected": false}}
	_, err = coll.UpdateMany(context.TODO(), updateFilter, updateSet)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("处理完毕")
}
