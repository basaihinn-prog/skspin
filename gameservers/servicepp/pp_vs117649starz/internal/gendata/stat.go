package gendata

import (
	"context"
	"serve/comm/db"

	"go.mongodb.org/mongo-driver/bson"
)

// func init() {
// 	mux.RegHttpWithSample("/gendata/stat", "获取各个bucket生成数量", "gendata", stat, nil).SetOnlyDev()
// }

type ListUnit struct {
	Count int
	Range [2]float64
}

type statRet struct {
	List  []*ListUnit
	Count int
}

func stat() (ret statRet, err error) {
	coll := db.Collection("simulate")

	cursor, err := coll.Aggregate(context.TODO(), []bson.M{
		{
			"$group": bson.M{
				"_id":   "$bucketid",
				"count": bson.M{"$count": bson.M{}},
			},
		},
	})
	if err != nil {
		return
	}

	var ans []struct {
		BucketId int `bson:"_id"`
		Count    int
	}
	err = cursor.All(context.TODO(), &ans)
	if err != nil {
		return
	}

	ret.List = make([]*ListUnit, len(GBuckets.Bounds))
	for i := 0; i < len(ret.List); i++ {
		b := GBuckets.Bounds[i]
		ret.List[i] = &ListUnit{
			Range: [2]float64{b.Min, b.Max},
		}
	}

	for _, v := range ans {
		u := ret.List[v.BucketId]
		u.Count = v.Count
		ret.Count += v.Count
	}

	return
}
