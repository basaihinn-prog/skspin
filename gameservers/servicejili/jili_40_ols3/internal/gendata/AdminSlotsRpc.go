package gendata

import (
	"context"
	"fmt"

	"serve/comm/db"
	"serve/comm/mux"
	"serve/servicejili/jili_40_ols3/internal"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	mux.RegRpc(fmt.Sprintf("/%s/AdminSlotsRpc/GetPoolStatus", internal.GameID), "原始生成的数据", "gendata", GetPoolStatus, nil).SetOnlyDev()
}

type RateDistributionStatus struct {
	Index              int
	Flag               string
	MinRate            float64
	MaxRate            float64
	ExpectCount        int
	CurrentCount       int
	UseCombineCount    int
	NotUseCombineCount int
	DstCount           int
	DstTimes           float64
	DstTimesMinusCost  float64
	PoolCost           float64
	Type               db.BoundType
}

type SlotsPoolStatusResp struct {
	IsAutoCreated         bool
	GenerateDataSuccess   bool
	GenCombineDataSuccess bool
	// TotalDstCount         int
	RateStatus [][]RateDistributionStatus
}

func GetPoolStatus(_ mux.EmptyParams, ret *SlotsPoolStatusResp) (err error) {
	var statRet statRet
	statRet, err = stat()
	if err != nil {
		return
	}

	bounds := GBuckets.bounds

	tbl := make([][]RateDistributionStatus, bounds[len(bounds)-1].Group+1)

	selectedStat, err := stat_selected()
	if err != nil {
		return
	}

	for i, bound := range bounds {
		u := RateDistributionStatus{
			Index:        i,
			Flag:         strconv.Itoa(i) + ":  " + bound.name(),
			MinRate:      bound.Min,
			MaxRate:      bound.Max,
			CurrentCount: statRet.List[i].Count,
			Type:         bound.Type,
		}
		if s := selectedStat[i]; s != nil {
			u.DstCount = s.Count
			u.UseCombineCount = GCombineDataMng.getCombineUsedCount(i, 0)
			u.NotUseCombineCount = u.DstCount - u.UseCombineCount
			u.DstTimes = s.Times
			u.DstTimesMinusCost = u.DstTimes - float64(bound.PoolCost)*float64(s.Count)
			u.PoolCost = float64(bound.PoolCost)
		}

		tbl[bound.Group] = append(tbl[bound.Group], u)
	}

	ret.RateStatus = tbl
	return
}

type selectedBucket struct {
	BucketId int `bson:"_id"`
	Count    int
	Times    float64
}

func stat_selected() (m map[int]*selectedBucket, err error) {
	coll := db.Collection("rawSpinData")
	cursor, err := coll.Aggregate(context.TODO(), []bson.M{
		{
			"$match": bson.M{
				"selected": true,
			},
		},
		{
			"$group": bson.M{
				"_id":   "$bucketid",
				"count": bson.M{"$count": bson.M{}},
				"times": bson.M{"$sum": "$times"},
			},
		},
	})
	if err != nil {
		return
	}

	var ans []*selectedBucket
	err = cursor.All(context.TODO(), &ans)
	if err != nil {
		return
	}

	m = map[int]*selectedBucket{}

	for _, v := range ans {
		m[v.BucketId] = v
	}

	return
}
