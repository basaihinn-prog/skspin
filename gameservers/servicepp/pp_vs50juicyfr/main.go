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
	"serve/servicepp/pp_vs50juicyfr/internal"
	"serve/servicepp/pp_vs50juicyfr/internal/gamedata"
	"serve/servicepp/pp_vs50juicyfr/internal/gendata"
	_ "serve/servicepp/pp_vs50juicyfr/internal/rpc"
	"serve/servicepp/ppcomm"
	"strconv"
	"strings"

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

	//for i := 1; i <= 61; i++ {
	for i := 63; i <= 63; i++ {
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
		num := int(float64(len(ids)) * 0.2)
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
	findFilter := bson.M{
		"$and": []bson.M{
			{"times": 0},
			{"type": 1},
			{"selected": true},
		},
	}
	findId := bson.M{"_id": true, "droppan": true}
	cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(findId))
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		var result ppcomm.SimulateData
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		i := len(result.DropPan) - 1

		lint := result.DropPan[i].Int("l")
		cfloat := result.DropPan[i].Float("c")
		var twstr string
		twstr = result.DropPan[i].Str("tw")
		twstr = strings.ReplaceAll(twstr, ",", "")
		twfloat, err := strconv.ParseFloat(twstr, 64)
		if err != nil {
			fmt.Println("转换失败:", err)
			return
		}

		times := twfloat / float64(lint) / cfloat
		fmt.Println(result.Id, times)

		updateFilter := bson.M{"_id": result.Id}
		update := bson.M{"$set": bson.M{"times": times}}

		_, err = coll.UpdateMany(context.TODO(), updateFilter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("处理完毕")
	}
}
