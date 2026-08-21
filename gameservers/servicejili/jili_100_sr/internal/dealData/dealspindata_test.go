package dealData

import (
	"context"

	"serve/comm/db"
	"serve/servicejili/jili_100_sr/internal"
	"serve/servicejili/jili_100_sr/internal/gendata"

	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCaclBucketId(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

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
