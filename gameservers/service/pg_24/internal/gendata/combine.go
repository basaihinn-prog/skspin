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
		{ // 1倍下注 5次不中奖
			Name:  " 1倍下注 5次不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     5,
				},
			},
			Type: GameTypeNormal,
		}, { // 1倍下注 8次不中奖
			Name:  " 1倍下注 8次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     8,
				},
			},
			Type: GameTypeNormal,
		}, { // 1倍下注 12次不中奖
			Name:  " 1倍下注 12次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeNormal),
					Count:     12,
				},
			},
			Type: GameTypeNormal,
		}, { //  1倍下注 1次普通（10，20】+8次不中奖+4次剩余随机.
			Name:  " 1倍下注 普通（10，20】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(10, 20, GameTypeNormal),
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
		}, { //  1倍下注 1次普通（20，30】+8次不中奖+4次剩余随机
			Name:  "  1倍下注 1次普通（20，30】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 30, GameTypeNormal),
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
		}, { // 1倍下注 1次普通（30，50】+8次不中奖+4次剩余随机
			Name:  "1倍下注 1次普通（30，50】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 50, GameTypeNormal),
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
		}, { // 1倍下注 1次普通（50，80】+8次不中奖+4次剩余随机
			Name:  "1倍下注 1次普通（50，80】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(50, 80, GameTypeNormal),
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
		}, { // 1倍下注 1次普通（80，150】+8次不中奖+4次剩余随机
			Name:  "1倍下注 1次普通（80，150】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(80, 150, GameTypeNormal),
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
		}, { // 1倍下注 1次普通（150，400】+8次不中奖+4次剩余随机
			Name:  "1倍下注 1次普通（150，400】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(150, 400, GameTypeNormal),
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
		}, { // 1倍下注 1次普通（400，800】+8次不中奖+4次剩余随机
			Name:  "1倍下注 1次普通（400，800】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(400, 800, GameTypeNormal),
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
		},

		{ // 2倍下注 5次不中奖
			Name:  " 2倍下注 5次不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     5,
				},
			},
			Type: GameTypeGamex2,
		}, { // 2倍下注 8次不中奖
			Name:  " 2倍下注 8次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				},
			},
			Type: GameTypeGamex2,
		}, { // 2倍下注 12次不中奖
			Name:  " 2倍下注 12次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     12,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（0，20】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（0，20】+8次不中奖+4次剩余随机.",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(0, 20, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（20，40】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（20，40】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(20, 40, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（40，60】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（40，60】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(40, 60, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（60，100】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（60，100】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(60, 100, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（100，160】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（100，160】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(100, 160, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（160，300】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（160，300】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(160, 300, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（300，800】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（300，800】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(300, 800, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		}, { //  2倍下注 1次普通（800，1600】+8次不中奖+4次剩余随机
			Name:  " 2倍下注 1次普通（800，1600】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(800, 1600, GameTypeGamex2),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex2),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex2,
		},

		{ // 3倍下注 5次不中奖
			Name:  " 3倍下注 5次不中奖",
			Count: 10,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     5,
				},
			},
			Type: GameTypeGamex3,
		}, { // 3倍下注 8次不中奖
			Name:  " 3倍下注 8次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				},
			},
			Type: GameTypeGamex3,
		}, { // 3倍下注 12次不中奖
			Name:  " 3倍下注 12次不中奖",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     12,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（0，30】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（0，30】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(0, 30, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（30，60】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（30，60】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(30, 60, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（60，90】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（60，90】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(60, 90, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（90，150】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（90，150】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(90, 150, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（150，240】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（150，240】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(150, 240, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（240，450】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（240，450】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(240, 450, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（450，1200】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（450，1200】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(450, 1200, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		}, { //  3倍下注 1次普通（1200，2600】+8次不中奖+4次剩余随机
			Name:  " 3倍下注 1次普通（1200，2600】+8次不中奖+4次剩余随机",
			Count: 20,
			Meta: []CombineMeta{
				{
					BucketIds: GetBucketIds(1200, 2600, GameTypeGamex3),
					Count:     1,
				}, {
					BucketIds: getIdsNoReward(GameTypeGamex3),
					Count:     8,
				}, {
					BucketIds: getIdsFree(),
					Count:     4,
				},
			},
			Type: GameTypeGamex3,
		},
	}
	for i, v := range combine {
		v.ID = i
	}
	return combine
}
