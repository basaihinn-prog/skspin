package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MiscSet(key string, value any) (err error) {
	coll := Collection("misc")

	_, err = coll.UpdateByID(context.TODO(), key, D("$set", D("value", value)), options.Update().SetUpsert(true))
	return
}

type objdoc struct {
	ID    string        `bson:"_id"`
	Value bson.RawValue `bson:"value"`
}

func MiscGet(key string, ptr any) (err error) {
	coll := Collection("misc")

	// var obj struct {
	// 	ID    string `bson:"_id"`
	// 	Value bson.RawValue
	// 	// Value bson.Raw
	// }

	var obj objdoc

	err = coll.FindOne(context.TODO(), D("_id", key)).Decode(&obj)
	if err != nil {
		return
	}
	// err = bson.Unmarshal(obj.Value, ptr)
	err = obj.Value.Unmarshal(ptr)
	return
}

func MiscInc(key string, d any) (err error) {
	coll := Collection("misc")

	update := D(
		"$inc", D("value", d),
	)
	_, err = coll.UpdateByID(context.TODO(), key, update, options.Update().SetUpsert(true))
	return
}

func MiscGet2(dbname, key string, ptr any) (err error) {
	coll := Collection2(dbname, "misc")

	var obj objdoc
	err = coll.FindOne(context.TODO(), D("_id", key)).Decode(&obj)
	if err != nil {
		return
	}
	// err = bson.Unmarshal(obj.Value, ptr)
	err = obj.Value.Unmarshal(ptr)
	return
}
