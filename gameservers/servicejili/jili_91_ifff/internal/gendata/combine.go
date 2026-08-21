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
		}, {
			Name:  "1次普通（10，20】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 20),
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
			Name:  "1次普通（20，30】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30),
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
			Name:  "1次普通（30，50】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50),
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
			Name:  "1次普通（50，90】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(50, 90),
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
			Name:  "1次普通（90，150】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(90, 150),
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
			Name:  "1次普通（150，300】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(150, 300),
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
			Name:  "1次普通（300，500】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(300, 500),
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
			Name:  "1次普通（500，1000】+6次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(500, 1000),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(),
					Count:     6,
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
