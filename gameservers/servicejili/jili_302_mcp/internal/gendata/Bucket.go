package gendata

import (
	"fmt"
	"log"
	"serve/comm/db"
)

const (
	GameTypeNormal = 0 // 普通模式
	GameTypeGame   = 1 // 购买小游戏
	GameTypeExtra  = 2 // 购买小游戏
)

type Bound struct {
	ID    int
	Group int

	Min, Max float64
	HasGame  bool
	//需要额外扣除玩家个人奖池5Bet
	PoolCost int

	Type db.BoundType
}

func (b *Bound) name() string {
	s := fmt.Sprintf("(%v, %v]", b.Min, b.Max)
	if b.HasGame {
		s += "小游戏"
	}
	return s
}

type Buckets struct {
	bounds []*Bound
}

var GBuckets = NewBuckets()

func (b *Buckets) GetPoolCost(bucketId int) int {
	return b.bounds[bucketId].PoolCost
}

func NewBuckets() *Buckets {
	bounds := []*Bound{
		// Type = 0
		{Group: 0, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 1, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 2, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 3, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 4, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 5, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 6, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 7, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 0},

		// Type = 1
		{Group: 8, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 1},
		{Group: 9, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 1},
		{Group: 10, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 1},
		{Group: 11, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 1},
		{Group: 12, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 1},
		{Group: 13, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 1},
		{Group: 14, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 1},
		{Group: 15, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 1},

		// Type = 2
		{Group: 16, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 17, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 18, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 19, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 20, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 21, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 22, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 23, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 2},

		// Type = 3
		{Group: 24, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 3},
		{Group: 25, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 3},
		{Group: 26, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 3},
		{Group: 27, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 3},
		{Group: 28, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 3},
		{Group: 29, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 3},
		{Group: 30, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 3},
		{Group: 31, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 3},

		// Type = 4
		{Group: 32, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 4},
		{Group: 33, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 4},
		{Group: 34, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 4},
		{Group: 35, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 4},
		{Group: 36, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 4},
		{Group: 37, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 4},
		{Group: 38, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 4},
		{Group: 39, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 4},

		// Type = 5
		{Group: 40, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 5},
		{Group: 41, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 5},
		{Group: 42, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 5},
		{Group: 43, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 5},
		{Group: 44, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 5},
		{Group: 45, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 5},
		{Group: 46, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 5},
		{Group: 47, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 5},

		// Type = 6
		{Group: 48, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 6},
		{Group: 49, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 6},
		{Group: 50, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 6},
		{Group: 51, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 6},
		{Group: 52, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 6},
		{Group: 53, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 6},
		{Group: 54, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 6},
		{Group: 55, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 6},

		// Type = 7
		{Group: 56, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 7},
		{Group: 57, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 7},
		{Group: 58, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 7},
		{Group: 59, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 7},
		{Group: 60, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 7},
		{Group: 61, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 7},
		{Group: 62, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 7},
		{Group: 63, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 7},

		// Type = 8
		{Group: 64, Min: 0, Max: 10, PoolCost: 0, HasGame: false, Type: 8},
		{Group: 65, Min: 10, Max: 20, PoolCost: 0, HasGame: false, Type: 8},
		{Group: 66, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 8},
		{Group: 67, Min: 30, Max: 60, PoolCost: 0, HasGame: false, Type: 8},
		{Group: 68, Min: 60, Max: 90, PoolCost: 0, HasGame: false, Type: 8},
		{Group: 69, Min: 90, Max: 100, PoolCost: 0, HasGame: false, Type: 8},
		{Group: 70, Min: 100, Max: 1000, PoolCost: 0, HasGame: false, Type: 8},
	}
	for i, v := range bounds {
		v.ID = i
	}

	return &Buckets{
		bounds: bounds,
	}
}

func (b *Buckets) GetBucket(multi float64, hasGame bool, gameType db.BoundType) int {
	for _, b := range b.bounds {
		if b.Min < multi && multi <= b.Max && b.HasGame == hasGame && b.Type == gameType {
			return b.ID
		}
	}
	return -1

}

func (b *Buckets) GetBound(i int) *Bound {
	return b.bounds[i]
}

func getBucketId(min, max float64, hasGame bool, ty int) int {
	for _, b := range GBuckets.bounds {
		if min <= b.Min && b.Max <= max && b.HasGame == hasGame && b.Type == db.BoundType(ty) {
			return b.ID
		}
	}

	log.Panicf("not found, min=%v, max=%v, hasGame=%v", min, max, hasGame)
	return -1
}

func getIdsNoReward(ty int) []int {
	ret := make([]int, 0, 2)
	ret = append(ret, getBucketId(-1, 0, false, ty))
	return ret
}

func getIdsFree() []int {
	return []int{}
}

func GetBucketIds(min, max float64, hasGame bool, ty int) (ids []int) {
	for _, b := range GBuckets.bounds {
		if min <= b.Min && b.Max <= max && b.HasGame == hasGame && b.Type == db.BoundType(ty) {
			ids = append(ids, b.ID)
		}
	}

	//lo.Must0(len(ids) != 0)
	return
}

func GetBuyMinBucketId() int {
	minRate, maxRate := float64(0), float64(10)
	ty := GameTypeGame
	for _, b := range GBuckets.bounds {
		if b.Min == minRate && b.Max == maxRate && b.Type == db.BoundType(ty) {
			return b.ID
		}
	}
	log.Panicf("GetBuyMinBucketId, min=%v, max=%v, hasGame=%v", minRate, maxRate)
	return -1
	//return getBucketId(0, 10, true)
}
