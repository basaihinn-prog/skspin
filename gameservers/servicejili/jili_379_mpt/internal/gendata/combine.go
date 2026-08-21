package gendata

import "serve/servicejili/jili_379_mpt/internal"

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
		{
			Name:  "5次不中奖",
			Count: 10,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     5,
				},
			},
		}, {
			Name:  "8次不中奖",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     8,
				},
			},
		}, {
			Name:  "12次不中奖",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     12,
				},
			},
		}, {
			Name:  "普通（15，20】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(15, 20, false, internal.GameTypeNormal),
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
			Name:  "普通（20，30】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30, false, internal.GameTypeNormal),
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
			Name:  "普通（30，40】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 40, false, internal.GameTypeNormal),
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
			Name:  "普通（40，60】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(40, 60, false, internal.GameTypeNormal),
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
			Name:  "普通（60，80】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(60, 80, false, internal.GameTypeNormal),
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
			Name:  "普通（80，100】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(80, 100, false, internal.GameTypeNormal),
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
			Name:  "普通（100，99999】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(100, 99999, false, internal.GameTypeNormal),
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
			Name:  "小游戏（0，10】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(0, 10, true, internal.GameTypeNormal),
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
			Name:  "小游戏（10，20】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 20, true, internal.GameTypeNormal),
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
			Name:  "小游戏（20，30】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
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
			Name:  "小游戏（30，50】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
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
			Name:  "小游戏（50，80】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
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
			Name:  "小游戏（80，100】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
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
			Name:  "小游戏（100，150】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
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
			Name:  "小游戏（150，200】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
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
			Name:  "小游戏（200，250】+8次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
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
		}, {
			Name:  "小游戏（30，50】+小游戏（10，20】+12次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: GetBucketIds(10, 20, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     12,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（30，50】+小游戏（30，50】+12次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50, true, internal.GameTypeNormal),
					Count:     2,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     12,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（30，50】+小游戏（50，80】+12次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: GetBucketIds(50, 80, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     12,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（30，50】+小游戏（80，100】+12次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: GetBucketIds(80, 100, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     12,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（10，20】+小游戏（20，30】+12次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 20, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: GetBucketIds(20, 30, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     12,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（10，20】+小游戏（150，200】+12次不中奖+4次剩余随机",
			Count: 20,
			Type:  internal.GameTypeNormal,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 20, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: GetBucketIds(150, 200, true, internal.GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(internal.GameTypeNormal),
					Count:     12,
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
