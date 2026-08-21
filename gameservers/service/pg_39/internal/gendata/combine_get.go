package gendata

import (
	"context"

	"serve/comm/db"
	"serve/comm/mux"

	"github.com/samber/lo"
)

func init() {
	mux.RegRpc("/pg_39/AdminSlotsRpc/combine/get", "获取组合配置", "gendata", httpGetCombine, nil).SetOnlyDev()
}

func httpGetCombine(_ mux.EmptyParams, ret *Combine) (err error) {
	combine, err := getCombine()
	if err != nil {
		return
	}

	*ret = combine
	return
}

func getCombine() (combine Combine, err error) {
	coll := db.Collection("combine")

	cursor, err := coll.Find(context.TODO(), db.D())
	if err != nil {
		return
	}

	var tmpcombine Combine
	err = cursor.All(context.TODO(), &tmpcombine)
	if err != nil {
		return
	}

	combine = newCombine()
	if len(combine) != len(tmpcombine) {
		coll.DeleteMany(context.TODO(), db.D())
		syncCombine(combine)
	} else {
		for i, c := range combine {
			lo.Must0(c.ID == tmpcombine[i].ID)
			c.Count = tmpcombine[i].Count
		}
	}

	return
}
