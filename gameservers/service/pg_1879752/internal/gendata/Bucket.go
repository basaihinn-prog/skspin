package gendata

import (
	"fmt"
	"log"
	"serve/comm/db"

	"github.com/samber/lo"
)

const (
	GameTypeNormal = 0 // 普通模式
	GameTypeGame   = 1
	GameTypeDouble = 10 // 双滚轴
)

type Bound struct {
	ID    int
	Group int

	Min, Max float64

	//需要额外扣除玩家个人奖池5Bet
	PoolCost int

	Type db.BoundType
}

func (b *Bound) name() string {
	s := fmt.Sprintf("(%v, %v]", b.Min, b.Max)
	if b.Type == GameTypeDouble {
		s += "2倍购买"
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

		{Group: 0, Min: -1, Max: 0, Type: GameTypeNormal},
		{Group: 1, Min: 0, Max: 0.25, Type: GameTypeNormal},
		{Group: 1, Min: 0.25, Max: 0.5, Type: GameTypeNormal},
		{Group: 1, Min: 0.5, Max: 1, Type: GameTypeNormal},
		{Group: 2, Min: 1, Max: 2, Type: GameTypeNormal},
		{Group: 2, Min: 2, Max: 3, Type: GameTypeNormal},
		{Group: 2, Min: 3, Max: 4, Type: GameTypeNormal},
		{Group: 3, Min: 4, Max: 5, Type: GameTypeNormal},
		{Group: 3, Min: 5, Max: 8, Type: GameTypeNormal},
		{Group: 3, Min: 8, Max: 10, Type: GameTypeNormal},
		{Group: 4, Min: 10, Max: 15, Type: GameTypeNormal},
		{Group: 4, Min: 15, Max: 20, Type: GameTypeNormal},
		{Group: 4, Min: 20, Max: 25, Type: GameTypeNormal},
		{Group: 4, Min: 25, Max: 40, Type: GameTypeNormal},
		{Group: 5, Min: 40, Max: 50, Type: GameTypeNormal},
		{Group: 5, Min: 50, Max: 80, Type: GameTypeNormal},
		{Group: 5, Min: 80, Max: 100, Type: GameTypeNormal},
		{Group: 6, Min: 100, Max: 99999, Type: GameTypeNormal},
		{Group: 7, Min: -1, Max: 0, Type: GameTypeNormal},
		{Group: 7, Min: 0, Max: 10, Type: GameTypeNormal},
		{Group: 8, Min: 10, Max: 20, Type: GameTypeNormal},
		{Group: 8, Min: 20, Max: 30, Type: GameTypeNormal},
		{Group: 8, Min: 30, Max: 50, Type: GameTypeNormal},
		{Group: 8, Min: 50, Max: 80, Type: GameTypeNormal},
		{Group: 9, Min: 80, Max: 100, Type: GameTypeNormal},
		{Group: 9, Min: 100, Max: 150, Type: GameTypeNormal},
		{Group: 10, Min: 150, Max: 200, Type: GameTypeNormal},
		{Group: 10, Min: 200, Max: 250, Type: GameTypeNormal},
		{Group: 11, Min: 250, Max: 400, Type: GameTypeNormal},
		{Group: 11, Min: 400, Max: 99999, Type: GameTypeNormal},
		{Group: 12, Min: 0, Max: 20, Type: GameTypeDouble},
		{Group: 13, Min: 20, Max: 30, Type: GameTypeDouble},
		{Group: 14, Min: 30, Max: 40, Type: GameTypeDouble},
		{Group: 15, Min: 40, Max: 99999, Type: GameTypeDouble}}
	for i, v := range bounds {
		v.ID = i
	}

	return &Buckets{
		bounds: bounds,
	}
}

func (b *Buckets) GetBucket(multi float64, gameType db.BoundType) int {
	for _, b := range b.bounds {
		if b.Min < multi && multi <= b.Max && b.Type == gameType {
			return b.ID
		}
	}
	return -1
}

func (b *Buckets) GetBound(i int) *Bound {
	return b.bounds[i]
}

func getBucketId(min, max float64, gameType db.BoundType) int {
	for _, b := range GBuckets.bounds {
		if b.Min == min && b.Max == max && b.Type == gameType {
			return b.ID
		}
	}

	log.Panicf("not found, min=%v, max=%v", min, max)
	return -1
}

func getIdsNoReward(gameType db.BoundType) []int {
	ret := make([]int, 0, 2)
	// ret = append(ret, getBucketId(-1, 0, false))
	ret = append(ret, getBucketId(-1, 0, gameType))
	return ret
}

func getIdsFree() []int {
	return []int{}
}

func GetBucketIds(min, max float64, gameType db.BoundType) (ids []int) {
	for _, b := range GBuckets.bounds {
		if min <= b.Min && b.Max <= max && b.Type == gameType {
			ids = append(ids, b.ID)
		}
	}

	lo.Must0(len(ids) != 0)
	return
}

func GetBuyMinBucketId(t int) int {
	_, maxRate := float64(-1), float64(1)
	for _, b := range GBuckets.bounds {
		if b.Max <= maxRate && b.Type == db.BoundType(t) {
			return b.ID
		}
	}
	return -1
}
