package dealData

import (
	"context"
	"strings"
	"testing"

	"serve/comm/db"
	"serve/service/pg_1879752/internal/gendata"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_1879752")

	collNames, _ := db.Client().Database("pg_1879752").ListCollectionNames(context.TODO(), bson.M{})
	for i := 0; i < len(collNames); i++ {
		if strings.HasPrefix(collNames[i], "pgSpinData_") {
			name := collNames[i]
			go DealNormal(name)
		} else if strings.HasPrefix(collNames[i], "pgSpinData2_") {
			name := collNames[i]
			go DealNormal2(name)
		}

	}
	select {}
}

func TestCaclBucketId(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_1879752")
	coll := db.Collection("simulate")

	cur, _ := coll.Find(context.TODO(), db.D(), options.Find().SetProjection(db.D("times", 1, "hasgame", 1, "type", 1)))

	models := []mongo.WriteModel{}
	for cur.Next(context.TODO()) {
		var doc struct {
			ID    primitive.ObjectID `bson:"_id"`
			Times float64            `bson:"times"`
			// HasGame bool               `bson:"hasgame"`
			Type db.BoundType `bson:"type"`
		}
		assert.Nil(t, cur.Decode(&doc))

		bucketId := gendata.GBuckets.GetBucket(doc.Times, doc.Type)

		model := mongo.NewUpdateOneModel().SetFilter(db.ID(doc.ID)).SetUpdate(db.D("$set", db.D("bucketid", bucketId, "selected", true)))
		models = append(models, model)
	}

	coll.BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(false))
}
