package main

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/servicejdb/jdb_14010/internal"
	"serve/servicejdb/jdb_14010/internal/gendata"
	_ "serve/servicejdb/jdb_14010/internal/rpc"
	"serve/servicejdb/jdbcomm"
	"sync"
)

func main() {
	lazy.Init(internal.GameID)

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
	jdbcomm.RegistRpcToMQ(mqconn)

	mux.RegistRpcToMQ(mqconn)

	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)
	lazy.Serve()
}

func test() {
	var err error
	coll := db.Collection("simulate")
	findFilter := bson.M{"selected": true}
	findData := bson.M{"_id": true, "times": true}
	cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(findData))
	if err != nil {
		log.Fatal(err)
	}

	var i int
	var wg sync.WaitGroup
	sem := make(chan struct{}, 12)
	for cursor.Next(context.TODO()) {
		var result struct {
			Id    primitive.ObjectID `bson:"_id"`
			Times float64            `bson:"times"`
		}
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}

			updateFilter := bson.M{"_id": result.Id}
			updateData := bson.M{"$set": bson.M{"times": result.Times * 5 / 6}}
			_, err = coll.UpdateMany(context.TODO(), updateFilter, updateData)

			i++
			fmt.Println(i, err)
			<-sem
		}()
	}
	wg.Wait()

	for k, v := range jdbcomm.GBuckets.Bounds {
		updateFilter := bson.M{
			"$and": []bson.M{
				{"type": v.Type},
				{"selected": true},
				{"hasgame": v.HasGame},
				{"times": bson.M{"$gt": v.Min}},
				{"times": bson.M{"$lte": v.Max}},
			},
		}
		updateData := bson.M{"$set": bson.M{"bucketid": k}}
		_, err = coll.UpdateMany(context.TODO(), updateFilter, updateData)
		fmt.Println(err)
	}
}
