package common

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Aggregate[T any](coll *mongo.Collection, Pipeline ...bson.D) (res []*T, err error) {
	showInfoCursor, err := coll.Aggregate(context.TODO(), Pipeline)
	if err != nil {
		return res, err
	}

	if err := showInfoCursor.All(context.TODO(), &res); err != nil {
		return res, err
	}
	defer showInfoCursor.Close(context.Background())

	return res, nil
}
