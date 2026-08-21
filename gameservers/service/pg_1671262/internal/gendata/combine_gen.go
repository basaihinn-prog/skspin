package gendata

import (
	"context"
	"serve/comm/db"
	"serve/comm/mux"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.RegRpc("/pg_1671262/AdminSlotsRpc/combine/gen", "生成组合数据", "gendata", genCombineData, newCombine()).SetOnlyDev()
}

func genCombineData(ps Combine, ret *mux.EmptyResult) (err error) {
	err = syncCombine(ps)
	if err != nil {
		return
	}

	err = LoadCombineData()
	return err
}

func syncCombine(ps Combine) error {
	coll := db.Collection("combine")
	var models []mongo.WriteModel
	var md mongo.WriteModel
	for _, c := range ps {
		md = mongo.NewUpdateOneModel().SetFilter(db.D("_id", c.ID)).SetUpdate(db.D("$set", db.D("count", c.Count)))
		models = append(models, md)
	}

	if len(models) > 0 {
		_, err := coll.BulkWrite(context.TODO(), models)
		if err != nil {
			return err
		}
	}
	return nil
}
