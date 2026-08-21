package gendata

import (
	"context"

	"serve/comm/db"
	"serve/comm/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.RegRpc("/pg_39/AdminSlotsRpc/combine/gen", "生成组合数据", "gendata", genCombineData, newCombine()).SetOnlyDev()
}

func syncCombine(ps Combine) (err error) {
	coll := db.Collection("combine")

	// save back the combine.count to mongo
	var models []mongo.WriteModel
	for _, c := range ps {
		md := mongo.NewUpdateOneModel().SetFilter(db.D("_id", c.ID)).SetUpdate(db.D("$set", db.D("count", c.Count))).SetUpsert(true)

		models = append(models, md)
	}

	if len(models) != 0 {
		_, err = coll.BulkWrite(context.TODO(), models)
		if err != nil {
			return
		}
	}

	return
}

type posinfo struct {
	combineID  int
	combineSeq int
}

type PlayId struct {
	ID         primitive.ObjectID
	CombineID  int
	CombineSeq int
	Pos        int
	Padding    bool
}

func genCombineData(ps Combine, ret *mux.EmptyResult) (err error) {
	err = syncCombine(ps)
	if err != nil {
		return
	}

	err = LoadCombineData()
	return
}
