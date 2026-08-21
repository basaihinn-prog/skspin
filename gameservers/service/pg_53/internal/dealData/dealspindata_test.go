package dealData

import (
	"context"
	"fmt"
	"serve/comm/db"
	"serve/comm/ut"
	"serve/service/pg_53/internal/gendata"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_53")

	collNames, _ := db.Client().Database("pg_53").ListCollectionNames(context.TODO(), bson.M{})
	for i := 0; i < len(collNames); i++ {
		if strings.HasPrefix(collNames[i], "pgSpinData_") {
			name := collNames[i]
			go dealNormal(name)
		}

	}
	select {}
}

func TestCaclBucketId(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_53")
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

func dealNormal(collName string) {
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
