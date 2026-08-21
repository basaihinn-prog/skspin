package main

import (
	"context"
	"fmt"
	"os"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/servicepp/pp_vs40rainbowr/internal"
	"serve/servicepp/pp_vs40rainbowr/internal/gamedata"
	"serve/servicepp/pp_vs40rainbowr/internal/gendata"
	_ "serve/servicepp/pp_vs40rainbowr/internal/rpc"
	"serve/servicepp/ppcomm"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
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
	GBuckets := gendata.GBuckets
	for k, v := range GBuckets.Bounds {
		fmt.Println(k, v)
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
		_, err := coll.UpdateMany(context.TODO(), updateFilter, updateData)
		fmt.Println(err)
	}

	return
}
