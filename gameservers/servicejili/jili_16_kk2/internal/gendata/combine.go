package gendata

import "serve/servicejili/jili_16_kk2/internal"

type CombineMeta struct {
	BucketIds []int
	Count     int
}

type CombineItem struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"-"`
	Count int
	Meta  []CombineMeta `bson:"-"`
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
		}, {
			Name:  "普通6次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     6,
				},
			},
		}, {
			Name:  "普通10次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     10,
				},
			},
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
		},

		//小游戏
		{
			Name:  "1次小游戏(20，30】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "1次小游戏(30，50】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "1次小游戏(50，80】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(50, 80, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "1次小游戏(80，100】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(80, 100, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "1次小游戏(100，150】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(100, 150, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "1次小游戏(150，200】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(150, 200, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "1次小游戏(200，250】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(200, 250, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		},
	}
	for i, v := range combine {
		v.ID = i
	}
	return combine
}
