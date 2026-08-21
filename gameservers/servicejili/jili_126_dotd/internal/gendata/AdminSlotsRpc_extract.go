package gendata

import (
	"context"
	"fmt"
	"log/slog"

	"serve/comm/db"
	"serve/comm/mux"
	"serve/servicejili/jili_126_dotd/internal"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.RegRpc(fmt.Sprintf("/%s/AdminSlotsRpc/ExtractAction", internal.GameID), "提取数据", "gendata", ExtractAction, []ExtractActionPsItem{
		{
			Index: 5,
			Count: 1000,
		},
		{
			Index: 6,
			Count: 2000,
		},
	}).SetOnlyDev()
}

type ExtractActionPsItem struct {
	Index int
	Count int
}

func ExtractAction(reqs []ExtractActionPsItem, _ *int) (err error) {
	coll := db.Collection("rawSpinData")

	bucketIds := lo.Map(reqs, func(item ExtractActionPsItem, _ int) int {
		return item.Index
	})
	coll.UpdateMany(context.TODO(), bson.M{
		"bucketid": bson.M{"$in": bucketIds},
	}, bson.M{
		"$set": bson.M{"selected": false},
	})

	// db.users.aggregate(
	// 	[ { $sample: { size: 3 } } ]
	//  )

	type dt struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	for _, req := range reqs {
		// slog.Info("ExtractAction", "flag", req.Flag, "count", req.DstCount)
		slog.Info("ExtractAction", "index", req.Index, "count", req.Count)
		if req.Count <= 0 {
			continue
		}

		cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{
			db.D("$match", db.D(
				"bucketid", req.Index,
			)),

			db.D("$sample", db.D("size", req.Count)),
			//db.D("$set", db.D("selected", true)),
			//db.D("$project", db.D("_id", 1, "times", 1)),
			db.D("$project", db.D("_id", 1)),
		})
		if err != nil {
			return err
		}

		var docs []dt
		err = cursor.All(context.TODO(), &docs)
		if err != nil {
			return err
		}

		ids := lo.Map(docs, func(item dt, _ int) primitive.ObjectID {
			return item.ID
		})
		_, err = coll.UpdateMany(context.TODO(), db.D("_id", db.D("$in", ids)), db.D("$set", db.D("selected", true)))
		if err != nil {
			return err
		}
	}

	return
}
