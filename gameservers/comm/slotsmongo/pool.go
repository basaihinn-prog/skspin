package slotsmongo

import (
	"context"
	"fmt"

	"serve/comm/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SlotsPoolRpcServer_GetID(pid int64) string {
	return fmt.Sprintf("%v-1", pid) // 注意int(typ)，否则会是字符串。
}

type DocSlotsPool struct {
	ID   string `bson:"_id"`
	Pid  int64  `bson:"Pid"`
	Gold int64  `bson:"Gold"`
}

func GetSelfSlotsPool(pid int64) (gold int64, err error) {
	id := SlotsPoolRpcServer_GetID(pid)
	var doc DocSlotsPool

	coll := SlotsCollection("Pool")
	err = coll.FindOne(context.TODO(), db.ID(id)).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		err = nil
		return
	}

	gold = doc.Gold
	return
}

func IncSelfSlotsPool(pid int64, gold int64) (aftergold int64, err error) {
	id := SlotsPoolRpcServer_GetID(pid)

	coll := SlotsCollection("Pool")
	var doc DocSlotsPool
	update := db.D(
		"$inc", db.D("Gold", gold),
		"$set", db.D("Pid", pid, "Type", 1),
	)

	// defer func() {
	// 	now := time.Now()
	// 	hisdoc := db.D(
	// 		"_id", primitive.NewObjectIDFromTimestamp(now),
	// 		"poolid", id,
	// 		"pid", pid,
	// 		"change", gold,
	// 		"after", aftergold,
	// 		"game", lazy.ServiceName,
	// 		"time", now,
	// 		"func", "IncSelfSlotsPool",
	// 		"err", ut.ErrString(err),
	// 	)

	// 	his := db.Collection2("slots", "PoolHistory")
	// 	his.InsertOne(context.TODO(), hisdoc)
	// }()

	err = coll.FindOneAndUpdate(context.TODO(), db.ID(id), update, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)).Decode(&doc)
	if err != nil {
		return
	}

	aftergold = doc.Gold

	return
}

func SetSelfSlotsPool(pid int64, gold int64) (err error) {
	id := SlotsPoolRpcServer_GetID(pid)

	coll := SlotsCollection("Pool")

	update := db.D(
		"$set", db.D(
			"Pid", pid,
			"Type", 1,
			"Gold", gold,
		),
	)
	_, err = coll.UpdateByID(context.TODO(), id, update, options.Update().SetUpsert(true))
	// if ret.UpsertedCount != 0 {
	// 	coll.UpdateByID(context.TODO(), ret.UpsertedID, db.D("$set", db.D("Pid", pid, "Type", 1)))
	// }

	/*
		now := time.Now()
		hisdoc := db.D(
			"_id", primitive.NewObjectIDFromTimestamp(now),
			"poolid", id,
			"pid", pid,
			"change", gold,
			"after", gold,
			"game", lazy.ServiceName,
			"time", now,
			"func", "SetSelfSlotsPool",
			"err", ut.ErrString(err),
		)

		his := db.Collection2("slots", "PoolHistory")
		his.InsertOne(context.TODO(), hisdoc)
	*/
	return
}
