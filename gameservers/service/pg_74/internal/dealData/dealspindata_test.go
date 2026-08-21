package dealData

import (
	"context"
	"fmt"
	"testing"

	"serve/comm/db"
	"serve/comm/ut"
	"serve/service/pg_74/internal/gendata"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_74")

	for i := 0; i < 50; i++ {
		num := i
		go dealNormal(num)
		//go dealGame(num)
	}
	select {}
}

func TestCaclBucketId(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_74")
	coll := db.Collection("simulate")

	cur, _ := coll.Find(context.TODO(), db.D(), options.Find().SetProjection(db.D("times", 1, "hasgame", 1, "type", 1)))

	models := []mongo.WriteModel{}
	for cur.Next(context.TODO()) {
		var doc struct {
			ID      primitive.ObjectID `bson:"_id"`
			Times   float64            `bson:"times"`
			HasGame bool               `bson:"hasgame"`
			Type    db.BoundType       `bson:"type"`
		}
		assert.Nil(t, cur.Decode(&doc))

		bucketId := gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, doc.Type)

		model := mongo.NewUpdateOneModel().SetFilter(db.ID(doc.ID)).SetUpdate(db.D("$set", db.D("bucketid", bucketId, "selected", true)))
		models = append(models, model)
	}

	coll.BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(false))
}

func dealNormal(i int) {
	collName := fmt.Sprintf("pgSpinData_%d", i)
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
					BucketId: gendata.GBuckets.GetBucket(times, hasGame, db.BoundType(gendata.GameTypeNormal)),
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

func dealGame(i int) {
	collName := fmt.Sprintf("pgSpinDataGame_%d", i)
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
					BucketId: gendata.GBuckets.GetBucket(times, hasGame, db.BoundType(gendata.GameTypeGame)),
					Type:     gendata.GameTypeGame,
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

func TestAddLimitCount(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_74")

	coll := db.Collection("simulate")
	match := bson.M{
		"bucketid": 0,
	}

	for i := 0; i < 1000; i++ {
		cursor, err := coll.Aggregate(context.TODO(), bson.A{
			bson.D{{"$match", match}},
			bson.D{
				{"$sample", bson.D{{"size", 10}}},
			},
		})
		// id, _ := primitive.ObjectIDFromHex("6604d1cce04f4526f342ac8a")
		// cursor, _ := coll.Find(context.TODO(), db.ID(id))
		if err != nil {
			fmt.Println(err)
			return
		}

		// cursor.Next()
		var docs []*gendata.SimulateData
		err = cursor.All(context.TODO(), &docs)
		if err != nil {
			fmt.Println(err)
			return
		}
		var bsonDocs []interface{}
		for j := range docs {
			docs[j].Id = primitive.NewObjectID()
			bsonDocs = append(bsonDocs, docs[j])
		}
		_, err = coll.InsertMany(context.TODO(), bsonDocs)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}
}
