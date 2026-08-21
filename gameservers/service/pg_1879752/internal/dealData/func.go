package dealData

import (
	"context"
	"fmt"

	"serve/comm/db"
	"serve/comm/ut"
	"serve/service/pg_1879752/internal/gendata"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DealNormal(collName string) {
	coll := db.Collection(collName)
	newColl := db.Collection("simulate")
	var lastId primitive.ObjectID
	db.MiscGet(collName, &lastId)
	match := bson.M{}
	if lastId != primitive.NilObjectID {
		match["_id"] = bson.M{"$gte": lastId}
	}
	cursor, _ := coll.Find(context.TODO(), match)
	defer cursor.Close(context.TODO())

	pans := []bson.M{}
	for cursor.Next(context.TODO()) {
		var pan map[string]any
		err := cursor.Decode(&pan)
		if err != nil {
			fmt.Println("err:", err,
				"_id:", pan["_id"],
				"sid:", pan["sid"],
				"psid:", pan["psid"])
			break
		}
		_id := pan["_id"]
		pan["orignid"] = fmt.Sprintf("%s_%s", pan["_id"], collName)
		delete(pan, "_id")
		if pan["sid"].(string) == pan["psid"].(string) && len(pans) > 0 {
			if pans[0]["sid"].(string) == pans[0]["psid"].(string) { // 第一盘是否有效
				objectId := primitive.NewObjectID()
				hasGame := false
				for j := 0; j < len(pans); j++ {
					pans[j]["psid"] = objectId.Hex()
					pans[j]["sid"] = fmt.Sprintf("%s_%d_%d", objectId.Hex(), j+1, len(pans))
					if pans[j]["fs"] != nil {
						hasGame = true
					}
				}
				times := ut.GetFloat(pans[len(pans)-1]["aw"])
				data := gendata.SimulateData{
					Id:       objectId,
					DropPan:  pans,
					HasGame:  hasGame,
					Times:    times,
					BucketId: gendata.GBuckets.GetBucket(times, db.BoundType(gendata.GameTypeNormal)),
					Type:     gendata.GameTypeNormal,
					Selected: true,
				}
				lo.Must(newColl.InsertOne(context.TODO(), data))
				fmt.Println("success", objectId.Hex())
			}
			db.MiscSet(collName, _id)
			pans = []bson.M{}
		}
		pans = append(pans, pan)
	}
}

func DealNormal2(collName string) {
	coll := db.Collection(collName)
	newColl := db.Collection("simulate")
	var lastId primitive.ObjectID
	db.MiscGet(collName, &lastId)
	match := bson.M{}
	if lastId != primitive.NilObjectID {
		match["_id"] = bson.M{"$gte": lastId}
	}
	cursor, _ := coll.Find(context.TODO(), match)
	defer cursor.Close(context.TODO())

	pans := []bson.M{}
	for cursor.Next(context.TODO()) {
		var pan map[string]any
		err := cursor.Decode(&pan)
		if err != nil {
			fmt.Println("err:", err,
				"_id:", pan["_id"],
				"sid:", pan["sid"],
				"psid:", pan["psid"])
			break
		}
		_id := pan["_id"]
		pan["orignid"] = fmt.Sprintf("%s_%s", pan["_id"], collName)
		delete(pan, "_id")
		if pan["sid"].(string) == pan["psid"].(string) && len(pans) > 0 {
			if pans[0]["sid"].(string) == pans[0]["psid"].(string) { // 第一盘是否有效
				objectId := primitive.NewObjectID()
				for j := 0; j < len(pans); j++ {
					pans[j]["psid"] = objectId.Hex()
					pans[j]["sid"] = fmt.Sprintf("%s_%d_%d", objectId.Hex(), j+1, len(pans))
				}
				times := ut.GetFloat(pans[len(pans)-1]["aw"])
				data := gendata.SimulateData{
					Id:       objectId,
					DropPan:  pans,
					HasGame:  true,
					Times:    times,
					BucketId: gendata.GBuckets.GetBucket(times, db.BoundType(gendata.GameTypeDouble)),
					Type:     gendata.GameTypeDouble,
					Selected: true,
				}
				lo.Must(newColl.InsertOne(context.TODO(), data))
				fmt.Println("success", objectId.Hex())
			}
			db.MiscSet(collName, _id)
			pans = []bson.M{}
		}
		pans = append(pans, pan)
	}
}
