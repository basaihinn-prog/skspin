package gendata

import "serve/servicejili/jili_225_ctk/internal"

type CombineMeta struct {
	BucketIds []int
	Count     int
}

type CombineItem struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"-"`
	Count int
	Meta  []CombineMeta `bson:"-"`
	Type  int           `bson:"-"`
}

type Combine []*CombineItem

func newCombine() Combine {
	combine := Combine{
		// 普通
		{
			Name:  "普通3次不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     3,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "普通6次不中奖",
			Count: 0,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     0,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "普通10次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     10,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（10，15】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 15, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（15，20】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(15, 20, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（20，30】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（30，40】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 40, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（40，60】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(40, 60, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（60，80】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(60, 80, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（80，100】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(80, 100, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		}, {
			Name:  "1次普通（100，99999】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(100, 99999, false, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeNormal,
		},

		// 额外
		{
			Name:  "额外3次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     3,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "额外6次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "额外10次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     10,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（10，15】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 15, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（15，20】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(15, 20, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（20，30】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（30，40】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 40, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（40，60】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(40, 60, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（60，80】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(60, 80, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（80，100】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(80, 100, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		}, {
			Name:  "1次额外（100，99999】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(100, 99999, false, internal.GameTypeExtra),
					Count:     1,
				},
				{
					BucketIds: getIdsNoReward(internal.GameTypeExtra),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: internal.GameTypeExtra,
		},
	}
	for i, v := range combine {
		v.ID = i
	}
	return combine
}
