package gendata

import (
	"fmt"
	"log"
	"serve/comm/db"

	"github.com/samber/lo"
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
		{Group: 0, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 1, Min: 0, Max: 0.25, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 1, Min: 0.25, Max: 0.5, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 1, Min: 0.5, Max: 1, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 2, Min: 1, Max: 2, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 2, Min: 2, Max: 3, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 2, Min: 3, Max: 4, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 3, Min: 4, Max: 6, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 3, Min: 6, Max: 8, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 3, Min: 8, Max: 10, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 4, Min: 10, Max: 15, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 4, Min: 15, Max: 20, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 4, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 4, Min: 30, Max: 40, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 5, Min: 40, Max: 60, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 5, Min: 60, Max: 80, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 6, Min: 80, Max: 100, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 6, Min: 100, Max: 99999, PoolCost: 0, HasGame: false, Type: 0},
		{Group: 7, Min: -1, Max: 0, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 8, Min: 0, Max: 10, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 8, Min: 10, Max: 20, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 8, Min: 20, Max: 30, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 8, Min: 30, Max: 50, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 9, Min: 50, Max: 80, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 9, Min: 80, Max: 100, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 9, Min: 100, Max: 150, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 10, Min: 150, Max: 200, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 10, Min: 200, Max: 250, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 10, Min: 250, Max: 99999, PoolCost: 0, HasGame: true, Type: 0},
		{Group: 12, Min: -1, Max: 0, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 13, Min: 0, Max: 0.25, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 13, Min: 0.25, Max: 0.5, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 13, Min: 0.5, Max: 1, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 14, Min: 1, Max: 2, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 14, Min: 2, Max: 3, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 14, Min: 3, Max: 4, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 15, Min: 4, Max: 6, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 15, Min: 6, Max: 8, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 15, Min: 8, Max: 10, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 15, Min: 10, Max: 15, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 16, Min: 15, Max: 20, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 16, Min: 20, Max: 30, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 16, Min: 30, Max: 40, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 16, Min: 40, Max: 60, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 17, Min: 60, Max: 80, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 17, Min: 80, Max: 100, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 17, Min: 100, Max: 99999, PoolCost: 0, HasGame: false, Type: 2},
		{Group: 18, Min: -1, Max: 0, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 19, Min: 0, Max: 10, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 19, Min: 10, Max: 20, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 19, Min: 20, Max: 30, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 20, Min: 30, Max: 50, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 20, Min: 50, Max: 80, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 20, Min: 80, Max: 100, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 21, Min: 100, Max: 150, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 21, Min: 150, Max: 200, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 21, Min: 200, Max: 250, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 21, Min: 250, Max: 99999, PoolCost: 0, HasGame: true, Type: 2},
		{Group: 22, Min: -1, Max: 0, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 23, Min: 0, Max: 10, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 23, Min: 10, Max: 20, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 23, Min: 20, Max: 30, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 24, Min: 30, Max: 50, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 24, Min: 50, Max: 80, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 24, Min: 80, Max: 100, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 25, Min: 100, Max: 150, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 25, Min: 150, Max: 200, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 25, Min: 200, Max: 250, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 26, Min: 250, Max: 300, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 26, Min: 300, Max: 350, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 26, Min: 350, Max: 500, PoolCost: 0, HasGame: true, Type: 1},
		{Group: 26, Min: 500, Max: 99999, PoolCost: 0, HasGame: true, Type: 1},
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

	lo.Must0(len(ids) != 0)
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
