package dealData

import (
	"context"
	"fmt"
	"serve/comm/db"
	"serve/service/pg_124/internal/gendata"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:aajohsujie9vecieSohqu4weivai7oxayei9rie5Yoh4vuojohwaothee0waethi@47.238.67.30:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_124")

	collNames, _ := db.Client().Database("pg_124").ListCollectionNames(context.TODO(), bson.M{})
	for i := 0; i < len(collNames); i++ {
		if strings.HasPrefix(collNames[i], "pgSpinData_") {
			name := collNames[i]
			go DealNormal(name)
		} else if strings.HasPrefix(collNames[i], "pgSpinDataGame_") {
			name := collNames[i]
			go DealGame(name)
		}

	}
	select {}
}

func TestCaclBucketId(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:aajohsujie9vecieSohqu4weivai7oxayei9rie5Yoh4vuojohwaothee0waethi@47.238.67.30:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_124")
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

func TestAddLimitCount(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:aajohsujie9vecieSohqu4weivai7oxayei9rie5Yoh4vuojohwaothee0waethi@47.238.67.30:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_123")

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
