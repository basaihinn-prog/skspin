package gendata

import (
	"context"
	"fmt"
	"serve/comm/db"
	"serve/comm/mux"
	"serve/servicejili/jili_375_psd/internal"

	"github.com/samber/lo"
)

func init() {
	mux.RegRpc(fmt.Sprintf("/%s/AdminSlotsRpc/combine/get", internal.GameID), "获取组合配置", "gendata", hgetCombine, nil).SetOnlyDev()
}

func hgetCombine(_ mux.EmptyParams, ret *Combine) (err error) {
	combine, err := getCombine()
	if err != nil {
		return
	}

	*ret = combine
	return err
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
	if len(tmpcombine) <= 0 { // 没有初始数据的时候设置一下
		tmpcombine = combine
		var docs []interface{}
		for _, c := range combine {
			docs = append(docs, c)
		}
		coll.InsertMany(context.TODO(), docs)
	}
	for i, c := range combine {
		lo.Must0(c.ID == tmpcombine[i].ID)
		c.Count = tmpcombine[i].Count
	}

	return
}
