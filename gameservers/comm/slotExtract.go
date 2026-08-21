package comm

import (
	"context"
	"log/slog"

	"serve/comm/db"
	"serve/comm/slotsmongo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExtractActionPsItem struct {
	Index int
	Count int
}

func ExtractAction(reqs []ExtractActionPsItem, _ *int) (err error) {
	coll := db.Collection("simulate")

	bucketIds := lo.Map(reqs, func(item ExtractActionPsItem, _ int) int {
		return item.Index
	})
	coll.UpdateMany(context.TODO(), bson.M{
		"bucketid": bson.M{"$in": bucketIds},
	}, bson.M{
		"$set": bson.M{"selected": false},
	})

	type dt struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	for _, req := range reqs {
		// slog.Info("ExtractAction", "flag", req.Flag, "count", req.DstCount)
		slog.Info("ExtractAction", "index", req.Index, "count", req.Count)
		if req.Count <= 0 {
			continue
		}

		cursor, err1 := coll.Aggregate(context.TODO(), mongo.Pipeline{
			db.D("$match", db.D(
				"bucketid", req.Index,
			)),

			db.D("$sample", db.D("size", req.Count)),
			//db.D("$set", db.D("selected", true)),
			//db.D("$project", db.D("_id", 1, "times", 1)),
			db.D("$project", db.D("_id", 1)),
		})
		if err1 != nil {
			return err1
		}

		var docs []dt
		err1 = cursor.All(context.TODO(), &docs)
		if err1 != nil {
			return err1
		}

		ids := lo.Map(docs, func(item dt, _ int) primitive.ObjectID {
			return item.ID
		})
		_, err1 = coll.UpdateMany(context.TODO(), db.D("_id", db.D("$in", ids)), db.D("$set", db.D("selected", true)))
		if err != nil {
			return err1
		}
		if req.Index == 0 {
			match := db.D("bucketid", req.Index)
			allCount, _ := coll.CountDocuments(context.TODO(), match)
			if allCount < int64(req.Count) {
				curs, err2 := coll.Aggregate(context.TODO(), bson.A{
					bson.D{{"$match", match}},
					bson.D{
						{"$sample", bson.D{{"size", int64(req.Count) - allCount}}},
					},
				})
				if err2 != nil {
					return err2
				}
				var docs1 []*slotsmongo.SimulateData
				err2 = curs.All(context.TODO(), &docs1)
				if err2 != nil {
					return err2
				}
				var bsonDocs []interface{}
				for j := range docs1 {
					docs1[j].Id = primitive.NewObjectID()
					bsonDocs = append(bsonDocs, docs1[j])
				}
				_, err2 = coll.InsertMany(context.TODO(), bsonDocs)
				if err2 != nil {
					return err2
				}
			}
		}
	}
	return err
}
