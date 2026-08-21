package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/service/pg_24/internal/gamedata"
	"serve/service/pg_24/internal/gendata"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "serve/service/pg_24/internal/rpc"

	"github.com/samber/lo"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			slog.Error("HttpInvoke", "HttpInvoke", r)
		}
	}()
	lazy.Init("pg_24")
	gamedata.Load()

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()

	if len(os.Args) > 1 && os.Args[1] == "cache" {
		//slotsmongo.InitClassData4("simulate")
		slotsmongo.InitData("simulate", "pg_24")
		lo.Must0(gendata.LoadCombineData())
		return
	}
	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)
	//mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())
	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)

	httpAddr, _ := lazy.RouteFile.Get("pg_24.http.test")
	if len(httpAddr) != 0 {
		mux.StartHttpServer(httpAddr)
	}

	lazy.Serve()
}

func tmp() {
	coll := db.Collection("simulate")

	resetFilter := bson.M{
		"$and": []bson.M{
			{"type": 2},
			{"selected": false},
			{"bucketid": bson.M{"$gte": 19}},
			{"bucketid": bson.M{"$lte": 26}},
		},
	}
	resetUpdate := bson.M{"$set": bson.M{"selected": true}}
	_, err := coll.UpdateMany(context.TODO(), resetFilter, resetUpdate)
	if err != nil {
		log.Fatal(err)
	}

	findFilter := bson.M{
		"$and": []bson.M{
			{"type": 2},
			{"selected": true},
			{"bucketid": bson.M{"$gte": 19}},
			{"bucketid": bson.M{"$lte": 26}},
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
	cancelGatherNum := (len(ids) * 16) / 100

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
