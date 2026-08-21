package facaicomm

import (
	"fmt"
	"github.com/samber/lo"
	"log"
)

type BoundType int

const (
	GameTypeNormal     = iota // 普通模式
	GameTypeGame       = 10   // 购买游戏
	GameTypeSuperGame1 = 1    // 超级购买1
	GameTypeSuperGame2 = 2    // 超级购买2
	GameTypeSuperGame3        // 超级购买3
	GameTypeSuperGame4        // 超级购买4
	DoubleGame         = 10   // 双倍模式等等
	SpecialGame0       = 1000 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame1       = 1001 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame2       = 1002 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame3       = 1003 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame4       = 1004 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame5       = 1005 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame6       = 1006 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame7       = 1007 // N轮次小游戏，N根据不同游戏自我定义
	SpecialGame8       = 1008 // N轮次小游戏，N根据不同游戏自我定义

	BucketOff = 1
	BucketOn  = 0

	BucketHeartBeat      = 1
	BucketWave           = 2
	BucketGov            = 3
	BucketMix            = 4
	BucketStable         = 5
	BucketHighAward      = 6
	BucketSuperHighAward = 7
	HitBigAwardPercent   = 350
)

type Bound struct {
	ID       int
	Group    int
	Min, Max float64
	HasGame  bool

	// 需要额外扣除玩家个人奖池5Bet
	PoolCost int
	Type     BoundType
	God      int
}

var GBuckets = NewBuckets([]*Bound{
	{Group: 0, Min: -1, Max: 0, HasGame: false, Type: 0},
	{Group: 1, Min: 0, Max: 0.25, HasGame: false, Type: 0},
	{Group: 1, Min: 0.25, Max: 0.5, HasGame: false, Type: 0},
	{Group: 1, Min: 0.5, Max: 1, HasGame: false, Type: 0},
	{Group: 2, Min: 1, Max: 2, HasGame: false, Type: 0},
	{Group: 2, Min: 2, Max: 3, HasGame: false, Type: 0},
	{Group: 2, Min: 3, Max: 4, HasGame: false, Type: 0},
	{Group: 3, Min: 4, Max: 6, HasGame: false, Type: 0},
	{Group: 3, Min: 6, Max: 8, HasGame: false, Type: 0},
	{Group: 3, Min: 8, Max: 10, HasGame: false, Type: 0},
	{Group: 4, Min: 10, Max: 15, HasGame: false, Type: 0},
	{Group: 4, Min: 15, Max: 20, HasGame: false, Type: 0},
	{Group: 4, Min: 20, Max: 30, HasGame: false, Type: 0},
	{Group: 4, Min: 30, Max: 40, HasGame: false, Type: 0},
	{Group: 5, Min: 40, Max: 60, HasGame: false, Type: 0},
	{Group: 5, Min: 60, Max: 80, HasGame: false, Type: 0},
	{Group: 6, Min: 80, Max: 100, HasGame: false, Type: 0},
	{Group: 6, Min: 100, Max: 99999, HasGame: false, Type: 0},
	{Group: 7, Min: -1, Max: 1, HasGame: true, Type: 1},
	{Group: 7, Min: 1, Max: 30, HasGame: true, Type: 1},
	{Group: 7, Min: 30, Max: 40, HasGame: true, Type: 1},
	{Group: 7, Min: 40, Max: 99999, HasGame: true, Type: 1},
	{Group: 8, Min: -1, Max: 1, HasGame: true, Type: 2},
	{Group: 8, Min: 1, Max: 30, HasGame: true, Type: 2},
	{Group: 8, Min: 30, Max: 40, HasGame: true, Type: 2},
	{Group: 8, Min: 40, Max: 99999, HasGame: true, Type: 2},
	{Group: 9, Min: -1, Max: 1, HasGame: true, Type: 3},
	{Group: 9, Min: 1, Max: 30, HasGame: true, Type: 3},
	{Group: 9, Min: 30, Max: 40, HasGame: true, Type: 3},
	{Group: 9, Min: 40, Max: 99999, HasGame: true, Type: 3},
	{Group: 10, Min: -1, Max: 1, HasGame: true, Type: 4},
	{Group: 10, Min: 1, Max: 30, HasGame: true, Type: 4},
	{Group: 10, Min: 30, Max: 40, HasGame: true, Type: 4},
	{Group: 10, Min: 40, Max: 99999, HasGame: true, Type: 4},
	{Group: 11, Min: -1, Max: 1, HasGame: true, Type: 5},
	{Group: 11, Min: 1, Max: 30, HasGame: true, Type: 5},
	{Group: 11, Min: 30, Max: 40, HasGame: true, Type: 5},
	{Group: 11, Min: 40, Max: 99999, HasGame: true, Type: 5},
	{Group: 12, Min: -1, Max: 1, HasGame: true, Type: 6},
	{Group: 12, Min: 1, Max: 30, HasGame: true, Type: 6},
	{Group: 12, Min: 30, Max: 40, HasGame: true, Type: 6},
	{Group: 12, Min: 40, Max: 99999, HasGame: true, Type: 6},
	{Group: 13, Min: -1, Max: 1, HasGame: true, Type: 7},
	{Group: 13, Min: 1, Max: 30, HasGame: true, Type: 7},
	{Group: 13, Min: 30, Max: 40, HasGame: true, Type: 7},
	{Group: 13, Min: 40, Max: 99999, HasGame: true, Type: 7},
	{Group: 14, Min: -1, Max: 1, HasGame: true, Type: 8},
	{Group: 14, Min: 1, Max: 30, HasGame: true, Type: 8},
	{Group: 14, Min: 30, Max: 40, HasGame: true, Type: 8},
	{Group: 14, Min: 40, Max: 99999, HasGame: true, Type: 8},
	{Group: 15, Min: -1, Max: 1, HasGame: true, Type: 9},
	{Group: 15, Min: 1, Max: 30, HasGame: true, Type: 9},
	{Group: 15, Min: 30, Max: 40, HasGame: true, Type: 9},
	{Group: 15, Min: 40, Max: 99999, HasGame: true, Type: 9},
	{Group: 16, Min: -1, Max: 1, HasGame: true, Type: 10},
	{Group: 16, Min: 1, Max: 30, HasGame: true, Type: 10},
	{Group: 16, Min: 30, Max: 40, HasGame: true, Type: 10},
	{Group: 16, Min: 40, Max: 99999, HasGame: true, Type: 10},
})

func (b *Bound) Name() string {
	s := fmt.Sprintf("(%v, %v]", b.Min, b.Max)
	if b.HasGame {
		s += "小游戏"
	}
	if b.Type == GameTypeGame {
		s += "购买"
	}
	return s
}

type Buckets struct {
	Bounds []*Bound
}

func (bs *Buckets) GetPoolCost(bucketId int) int {
	return bs.Bounds[bucketId].PoolCost
}

func NewBuckets(bounds []*Bound) *Buckets {
	for i, v := range bounds {
		v.ID = i
	}

	return &Buckets{
		Bounds: bounds,
	}
}

func (bs *Buckets) GetBucket(multi float64, hasGame bool, gameType BoundType) int {
	for _, b := range bs.Bounds {
		if b.Min < multi && multi <= b.Max && b.HasGame == hasGame && b.Type == gameType {
			return b.ID
		}
	}
	return -1

}

func (bs *Buckets) GetBound(i int) *Bound {
	return bs.Bounds[i]
}

func (bs *Buckets) getBucketId(min, max float64, hasGame bool) int {
	for _, b := range bs.Bounds {
		if b.Min == min && b.Max == max && b.HasGame == hasGame {
			return b.ID
		}
	}

	log.Panicf("not found, min=%v, max=%v, hasGame=%v", min, max, hasGame)
	return -1
}

func (bs *Buckets) getIdsNoReward() []int {
	ret := make([]int, 0, 2)
	ret = append(ret, bs.getBucketId(-1, 0, false))
	return ret
}

func (bs *Buckets) getIdsFree() []int {
	return []int{}
}

func (bs *Buckets) GetBucketIds(min, max float64, hasGame bool) (ids []int) {
	for _, b := range bs.Bounds {
		if min <= b.Min && b.Max <= max && b.HasGame == hasGame {
			ids = append(ids, b.ID)
		}
	}

	lo.Must0(len(ids) != 0)
	return
}

func (bs *Buckets) GetBuyMinBucketId() int {
	minRate, maxRate := float64(0), float64(10)
	ty := GameTypeGame
	for _, b := range bs.Bounds {
		if b.Min == minRate && b.Max == maxRate && b.Type == BoundType(ty) {
			return b.ID
		}
	}
	log.Panicf("GetBuyMinBucketId, min=%v, max=%v", minRate, maxRate)
	return -1
}

func getBucketId(min, max float64, hasGame bool) int {
	for _, b := range GBuckets.Bounds {
		if b.Min == min && b.Max == max && b.HasGame == hasGame {
			return b.ID
		}
	}

	log.Panicf("not found, min=%v, max=%v, hasGame=%v", min, max, hasGame)
	return -1
}

func getIdsNoReward() []int {
	ret := make([]int, 0, 2)
	ret = append(ret, getBucketId(-1, 0, false))
	return ret
}

func getIdsFree() []int {
	return []int{}
}

func GetBucketIds(min, max float64, hasGame bool) (ids []int) {
	for _, b := range GBuckets.Bounds {
		if min <= b.Min && b.Max <= max && b.HasGame == hasGame {
			ids = append(ids, b.ID)
		}
	}

	//lo.Must0(len(ids) != 0)
	return
}

func GetBuyMinBucketIdByTy(t int) int {
	_, maxRate := float64(-1), float64(1)
	for _, b := range GBuckets.Bounds {
		if b.Max <= maxRate && b.Type == BoundType(t) {
			return b.ID
		}
	}
	return -1
}

func GetBuyMinBucketId() int {
	//minRate, maxRate := float64(0), float64(10)
	minRate, maxRate := float64(0), float64(10)
	ty := GameTypeGame
	for _, b := range GBuckets.Bounds {
		if b.Min >= minRate && b.Max <= maxRate && b.Type == BoundType(ty) {
			return b.ID
		}
	}
	log.Panicf("GetBuyMinBucketId, min=%v, max=%v, hasGame=%v", minRate, maxRate)
	return -1
}

func GetSuperBuyMinBucketId() int {
	//minRate, maxRate := float64(0), float64(10)
	minRate, maxRate := float64(0), float64(10)
	ty := GameTypeSuperGame1
	for _, b := range GBuckets.Bounds {
		if b.Min >= minRate && b.Max <= maxRate && b.Type == BoundType(ty) {
			return b.ID
		}
	}
	log.Panicf("GetSuperBuyMinBucketId, min=%v, max=%v, hasGame=%v", minRate, maxRate)
	return -1
}
