package gendata

import (
	"fmt"
	"log"
	"serve/comm/db"

	"github.com/samber/lo"
)

const (
	GameTypeNormal = 0          // 普通模式
	GameTypeGame   = BuyMul - 1 // 双滚轴
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
	if b.Type == GameTypeGame {
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

		// 单盘面
		{ // (-oo,0]
			Group: 0,
			Min:   -1,
			Max:   0,
			Type:  GameTypeNormal,
		}, { // （0，5】
			Group: 0,
			Min:   0,
			Max:   5,
			Type:  GameTypeNormal,
		},

		{ // （5，25】
			Group:    1,
			Min:      5,
			Max:      25,
			Type:     GameTypeNormal,
			PoolCost: 3,
		}, { // （25，50】
			Group:    1,
			Min:      25,
			Max:      50,
			Type:     GameTypeNormal,
			PoolCost: 6,
		}, { // （50，100】
			Group:    1,
			Min:      50,
			Max:      100,
			Type:     GameTypeNormal,
			PoolCost: 10,
		},

		// 双盘面  -----------------------
		{ // (-oo,0]
			Group: 2,
			Min:   -1,
			Max:   0,
			Type:  GameTypeGame,
		}, { // （0，20】
			Group: 2,
			Min:   0,
			Max:   20,
			Type:  GameTypeGame,
		},

		{ // （20，100】
			Group:    3,
			Min:      20,
			Max:      100,
			Type:     GameTypeGame,
			PoolCost: 8,
		}, { // （100，200】
			Group:    3,
			Min:      100,
			Max:      200,
			Type:     GameTypeGame,
			PoolCost: 20,
		}, { // （200，400】
			Group:    3,
			Min:      200,
			Max:      400,
			Type:     GameTypeGame,
			PoolCost: 50,
		},
	}
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
