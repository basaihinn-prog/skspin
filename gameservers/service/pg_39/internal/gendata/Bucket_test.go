package gendata

import (
	"context"
	"testing"

	"serve/comm/db"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSetBucket(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_39")

	coll := db.Collection("pgSpinData")

	cur, _ := coll.Find(context.TODO(), db.D("selected", true), options.Find().SetProjection(db.D("aw", 1)))

	count, err := coll.EstimatedDocumentCount(context.TODO())
	assert.Nil(t, err)
	assert.NotZero(t, count)
	var models = make([]mongo.WriteModel, 0, count)
	for cur.Next(context.TODO()) {
		var doc struct {
			ID primitive.ObjectID `bson:"_id"`
			Aw float64
		}
		assert.Nil(t, cur.Decode(&doc))

		bucketId := GBuckets.GetBucket(doc.Aw)

		model := mongo.NewUpdateOneModel().SetFilter(db.ID(doc.ID)).SetUpdate(db.D("$set", db.D("bucketid", bucketId, "selected", true)))
		models = append(models, model)
	}

	coll.BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(false))
}
