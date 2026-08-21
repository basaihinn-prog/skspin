package slotsmongo

import (
	"context"
	"fmt"

	"serve/comm/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getId(pid int64, game string) string {
	return fmt.Sprintf("%v-%v", game, pid)
}

func IncRotate(pid int64, game string) (err error) {
	id := getId(pid, game)
	coll := SlotsCollection("PlayerRotateData")
	_, err = coll.UpdateByID(context.TODO(), id, db.D("$inc", db.D("RotateCount", 1)), options.Update().SetUpsert(true))
	// err := CollPlayerRotateData.UpsertId(id, bson.M{"$inc": bson.M{"RotateCount": 1}})
	return
}

func GetRotate(pid int64, game string) (count int64, err error) {
	id := getId(pid, game)

	var doc struct {
		ID          string `bson:"_id"`
		RotateCount int64  `bson:"RotateCount"`
	}

	coll := SlotsCollection("PlayerRotateData")
	err = coll.FindOne(context.TODO(), db.ID(id)).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	count = doc.RotateCount
	return
}
