package api

import (
	"context"
	"errors"
	"fmt"
	"game/comm/db"
	"game/service/pggateway/pgcomm"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://api.pg-demo.com/web-api/game-proxy/v2/BetHistory/Get?traceId=KKQMKE21
// gid=39&dtf=1710432000000&dtt=1711036799999&bn=1&rc=15&atk=D663F48D-040E-46B5-B2C8-3FDD51F66316&pf=1&wk=0_C&btt=1
// gid=39&dtf=1710950400000&dtt=1711036799999&bn=2&rc=15&lbt=1711005490264&atk=4EBA5DB7-222C-4898-8F7C-36522DE464EB&pf=1&wk=0_C&btt=1

var (
	checkedM sync.Map
)

func ensureIndex(coll *mongo.Collection) {
	key := coll.Database().Name()
	if _, ok := checkedM.Load(key); ok {
		return
	}
	checkedM.Store(key, true)

	indexmodels := []mongo.IndexModel{
		{
			Keys: db.D("bt", 1),
			// Options: options.Index().SetUnique(false),
		},
		{
			Keys: db.D("pid", 1),
		},
	}
	coll.Indexes().CreateMany(context.TODO(), indexmodels)
}

func getBetHistory(ps *PGParams, ret *M) (err error) {
	gid := ps.Get("gid")
	coll := db.Collection2("pg_"+gid, "BetHistory")

	ensureIndex(coll)

	dtt, lbt := ps.GetInt("dtt"), ps.GetInt("lbt")
	if lbt != 0 {
		dtt = lbt
	}
	// dtt := min()
	filter := db.D(
		"bt", db.D(
			"$gt", ps.GetInt("dtf"),
			"$lt", dtt,
		),
		"pid", ps.Pid,
	)

	opts := options.Find().SetLimit(int64(ps.GetInt("rc"))).SetSort(db.D("bt", -1))
	cur, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}
	var bh []*pgcomm.BHItem
	err = cur.All(context.TODO(), &bh)
	if err != nil {
		return
	}

	if len(bh) == 0 {
		bh = []*pgcomm.BHItem{}
	}

	*ret = M{
		"bh": bh,
	}

	return
}

type Result struct {
	PSID string `bson:"psid"`
}

// /back-office-proxy/Report/GetBetHistory
// sid=1775083500903469568&gid=1489936
func boGetBetHistory(ps *PGParams, ret *M) (err error) {
	var bh pgcomm.BHItem
	gid := ps.Get("gid")
	coll := db.Collection2("pg_"+gid, "BetHistory")
	sid := ps.Get("sid")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		//return 原来逻辑需要return，因为调整游戏真伪代码，需要这么操作
		coll2 := db.Collection2("pg_"+gid, "psidMap")
		var result Result
		err = coll2.FindOne(context.TODO(), db.D("sid", sid), options.FindOne().SetProjection(db.D("psid", 1))).Decode(&result)
		if err != nil {
			//errors.Is()
			if errors.Is(err, mongo.ErrNoDocuments) {
				err = &PGError{
					Cd:  "5000",
					Msg: "mongo: no documents in result psidMap",
					//Tid: traceId,
				}
			}
			return
		}
		sid = result.PSID
		err = coll.FindOne(context.TODO(), db.D("tid", sid)).Decode(&bh)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				err = &PGError{
					Cd:  "5000",
					Msg: "mongo: no documents in result BetHistory",
					//Tid: traceId,
				}
			}
			return
		}
	} else {
		err = coll.FindOne(context.TODO(), db.ID(id)).Decode(&bh)
		if err != nil {
			return
		}
	}
	// if len(bh.Bd) == 0 {
	// 	bh.Bd = []*pgcomm.BDItem{}
	// }
	if len(bh.CC) > 3 {
		bh.CC = fmt.Sprintf(" %s ", bh.CC[:3])
	}
	//bh.CC = bh.CC[:3]

	*ret = M{
		"bh": &bh,
	}
	return
}
