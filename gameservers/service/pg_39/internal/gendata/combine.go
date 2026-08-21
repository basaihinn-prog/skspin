package gendata

import (
	"github.com/samber/lo"
)

type CombineMeta struct {
	BucketIds []int
	Count     int
}
type CombineItem struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"-"`
	Count int
	Meta  []CombineMeta `bson:"-" json:"-"`
}

type Combine []*CombineItem

func GetBucketIds(min, max float64, hasGame bool) (ids []int) {
	for _, b := range GBuckets.bounds {
		if min <= b.Min && b.Max <= max {
			ids = append(ids, b.ID)
		}
	}

	lo.Must0(len(ids) != 0)
	return
}

func getBucketIdsNoReward() (ids []int) {
	for _, b := range GBuckets.bounds {
		if b.Max == 0 {
			ids = append(ids, b.ID)
		}
	}

	lo.Must0(len(ids) != 0)
	return
}

func getBucketIdsFree() (ids []int) {
	return
}

func newCombine() Combine {
	combine := Combine{
		{
			Name:  "5次都不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getBucketIdsNoReward(),
					Count:     5,
				},
			},
		},
		{
			Name:  "8次都不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getBucketIdsNoReward(),
					Count:     8,
				},
			},
		},
		{
			Name:  "12次都不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getBucketIdsNoReward(),
					Count:     12,
				},
			},
		},
		{
			Name:  "1次普通中奖（10，15】+10次不中奖+4次剩余随机",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 15, false),
					Count:     1,
				},
				{
					BucketIds: getBucketIdsNoReward(),
					Count:     10,
				},
				{
					BucketIds: getBucketIdsFree(),
					Count:     4,
				},
			},
		},
		{
			Name:  "1次普通中奖（15，20】+10次不中奖+4次剩余随机",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(15, 20, false),
					Count:     1,
				},
				{
					BucketIds: getBucketIdsNoReward(),
					Count:     10,
				},
				{
					BucketIds: getBucketIdsFree(),
					Count:     4,
				},
			},
		},
		{
			Name:  "1次普通中奖（20，50】+10次不中奖+4次剩余随机",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 50, false),
					Count:     1,
				},
				{
					BucketIds: getBucketIdsNoReward(),
					Count:     10,
				},
				{
					BucketIds: getBucketIdsFree(),
					Count:     4,
				},
			},
		},
		{
			Name:  "1次普通中奖（50，100】+10次不中奖+4次剩余随机",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(50, 100, false),
					Count:     1,
				},
				{
					BucketIds: getBucketIdsNoReward(),
					Count:     10,
				},
				{
					BucketIds: getBucketIdsFree(),
					Count:     4,
				},
			},
		},
	}

	// 1次普通中奖（10，15】+10次不中奖+4次剩余随机
	// 1次普通中奖（15，20】+10次不中奖+4次剩余随机
	// 1次普通中奖（20，50】+10次不中奖+4次剩余随机
	// 1次普通中奖（50，100】+10次不中奖+4次剩余随机

	for i, v := range combine {
		v.ID = i
	}

	return combine
}
