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
	"serve/servicepp/pp_vs20daydead/internal"
	"serve/servicepp/pp_vs20daydead/internal/gamedata"
	"serve/servicepp/pp_vs20daydead/internal/gendata"
	_ "serve/servicepp/pp_vs20daydead/internal/rpc"
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

func test1() {
	coll := db.Collection("simulate")
	GBuckets := gendata.GBuckets
	for k, v := range GBuckets.Bounds {
		fmt.Println(k, v)
		var updateFilter primitive.M
		if v.Type == 0 {
			updateFilter = bson.M{
				"$and": []bson.M{
					{"type": v.Type},
					{"selected": true},
					{"hasgame": v.HasGame},
					{"times": bson.M{"$gt": v.Min}},
					{"times": bson.M{"$lte": v.Max}},
				},
			}
		}

		if v.Type == 1 {
			updateFilter = bson.M{
				"$and": []bson.M{
					{"type": v.Type},
					{"selected": true},
					//{"hasgame": v.HasGame},
					{"times": bson.M{"$gt": v.Min}},
					{"times": bson.M{"$lte": v.Max}},
				},
			}
		}

		updateData := bson.M{"$set": bson.M{"bucketid": k}}
		_, err := coll.UpdateMany(context.TODO(), updateFilter, updateData)
		fmt.Println(err)
	}

	return
}

func test2() {
	coll := db.Collection("simulate")
	findFilter := bson.M{
		"$and": []bson.M{
			{"selected": true},
			{"hasgame": false},
			{"type": 0},
		},
	}
	findId := bson.M{"_id": true, "droppan": true}
	cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(findId))
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(context.TODO()) {
		var result struct {
			Id      primitive.ObjectID `bson:"_id,omitempty"`
			DropPan []any              `bson:"droppan"`
		}
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		if len(result.DropPan) > 1 {
			fmt.Println(result.Id)
			updateFilter := bson.M{"_id": result.Id}
			updateData := bson.M{"$set": bson.M{"hasgame": true}}
			_, err := coll.UpdateMany(context.TODO(), updateFilter, updateData)
			fmt.Println(err)
		}
	}
}

func test3() {
	coll := db.Collection("simulate")
	findFilter := bson.M{
		"$and": []bson.M{
			{"selected": true},
			{"type": 0},
		},
	}
	findId := bson.M{"_id": true, "droppan": true}
	cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(findId))
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(context.TODO()) {
		var result struct {
			Id      primitive.ObjectID `bson:"_id,omitempty"`
			DropPan []struct {
				TW string `bson:"tw"`
			} `bson:"droppan"`
			Times float64 `bson:"times"`
		}
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		if len(result.DropPan) > 1 {
			s := result.DropPan[len(result.DropPan)-1]
			if s.TW == "0.00" {
				fmt.Println(result.Id, result.Times)
			}
		}
	}
}
