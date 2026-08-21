package gendata

import "serve/comm/db"

type CombineMeta struct {
	BucketIds []int
	Count     int
}

type CombineItem struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"-"`
	Count int
	Meta  []CombineMeta `bson:"-"`
	Type  db.BoundType
}

type Combine []*CombineItem

func newCombine() Combine {
	combine := Combine{
		{
			Name:  "5次不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     5,
				},
			},
			Type: GameTypeNormal,
		}, {
			Name:  "8次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     8,
				},
			},
			Type: GameTypeNormal,
		}, {
			Name:  "普通（0，5】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(0, 5, GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeNormal,
		}, {
			Name:  "普通（5，25】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(5, 25, GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeNormal,
		}, {
			Name:  "普通（25，50】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(25, 50, GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeNormal,
		}, {
			Name:  "普通（50，100】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(50, 100, GameTypeNormal),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeNormal,
		}}
	for i, v := range combine {
		v.ID = i
	}
	return combine
}
