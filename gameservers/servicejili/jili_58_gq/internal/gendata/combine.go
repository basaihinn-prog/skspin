package gendata

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
		{
			Name:  "3次不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(),
					Count:     3,
				},
			},
		}, {
			Name:  "6次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(),
					Count:     6,
				},
			},
		}, {
			Name:  "10次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(),
					Count:     10,
				},
			},
		},
		{
			Name:  "普通（10，15】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 15, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "普通（15，20】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(15, 20, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "普通（20，30】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "普通（30，40】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 40, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "普通（40，60】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(40, 60, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "普通（60，80】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(60, 80, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "普通（80，100】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(80, 100, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "普通（100，99999】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(100, 99999, false),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（20，30】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30, true),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（30，50】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50, true),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（50，80】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(50, 80, true),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（80，100】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(80, 100, true),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（100，150】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(100, 150, true),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（150，200】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(150, 200, true),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
		}, {
			Name:  "小游戏（200，250】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(200, 250, true),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
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
