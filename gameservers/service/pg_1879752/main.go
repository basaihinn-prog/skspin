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
	"serve/service/pg_1879752/internal/gamedata"
	"serve/service/pg_1879752/internal/gendata"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "serve/service/pg_1879752/internal/rpc"

	"github.com/samber/lo"
)

func main() {
	//defer func() {
	//	if r := recover(); r != nil {
	//		slog.Error("HttpInvoke", "HttpInvoke", r)
	//	}
	//}()
	lazy.Init("pg_1879752")
	gamedata.Load()

	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	slotsmongo.RestoreSpinData()

	// 先准备好数据, 再监听消息
	if len(os.Args) > 1 && os.Args[1] == "cache" {
		slotsmongo.InitData("simulate", "pg_1879752")
		lo.Must0(gendata.LoadCombineData())
		os.Exit(0)
	}

	lo.Must0(gendata.LoadCombineData())

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)
	//mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)

	httpAddr, _ := lazy.RouteFile.Get("pg_1879752.http.test")
	if len(httpAddr) != 0 {
		mux.StartHttpServer(httpAddr)
	}
	lazy.Serve()
}

func test() {
	coll := db.Collection("simulate")
	findFilter := bson.M{}
	findId := bson.M{"_id": true, "DropPan": true}
	cursor, _ := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(findId))

	var ids []primitive.ObjectID
	for cursor.Next(context.TODO()) {
		var result struct {
			Id      primitive.ObjectID `bson:"_id,omitempty"`
			DropPan []struct {
				Sid  string `bson:"sid"`
				Psid string `bson:"psid"`
			}
		}
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		sid := strings.Split(result.DropPan[0].Sid, "_")
		if sid[0] != result.DropPan[0].Psid {
			ids = append(ids, result.Id)
		}
	}
	findFilterId := bson.M{"_id": bson.M{"$in": ids}}
	many, err := coll.DeleteMany(context.TODO(), findFilterId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(many)
}

func test1() {
	coll := db.Collection("simulate")
	findFilter := bson.M{}
	findId := bson.M{"_id": true, "times": true, "DropPan": true}
	cursor, _ := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(findId))

	for cursor.Next(context.TODO()) {
		var result struct {
			Id      primitive.ObjectID `bson:"_id,omitempty"`
			Times   float64            `bson:"times"`
			DropPan []struct {
				Tw  float64 `bson:"tw"`
				Tbb float64 `bson:"tbb"`
			}
		}
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		num := len(result.DropPan) - 1
		tw := result.DropPan[num].Tw
		tbb := result.DropPan[num].Tbb
		times := tw / tbb
		if result.Times != times {
			fmt.Println(result.Times, times, result.Id)
			//updateFilter := bson.M{"_id": result.Id}
			//update := bson.M{"$set": bson.M{"times": times}}
			//_, err = coll.UpdateOne(context.TODO(), updateFilter, update)
			//if err != nil {
			//	fmt.Println(err)
			//	break
			//}
		}
	}
}

//func test2() {
//	coll := db.Collection("simulate")
//	GBuckets := gendata.NewBuckets()
//	for k, v := range GBuckets.Bounds {
//		if v.Group >= 7 && v.Group <= 11 {
//			continue
//		}
//		fmt.Println(k, v)
//		updateFilter := bson.M{
//			"$and": []bson.M{
//				{"type": v.Type},
//				{"selected": true},
//				{"times": bson.M{"$gt": v.Min}},
//				{"times": bson.M{"$lte": v.Max}},
//			},
//		}
//		updateData := bson.M{"$set": bson.M{"bucketid": k}}
//		_, err := coll.UpdateMany(context.TODO(), updateFilter, updateData)
//		fmt.Println(err)
//	}
//}
