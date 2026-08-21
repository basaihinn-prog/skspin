package gendata

import (
	"fmt"
	"log"
	"serve/comm/db"

	"github.com/samber/lo"
)

const (
	GameTypeNormal db.BoundType = 0 // (0) 普通模式
	// GameTypeGame   = BuyMul - 1 // 购买小游戏

	GameTypeGamex2 db.BoundType = BuyMulx2 - 1 // (1) 双倍
	GameTypeGamex3 db.BoundType = BuyMulx3 - 1 // (2) 三倍
)

type Bound struct {
	ID    int
	Group int

	Min, Max float64
	//HasGame  bool

	//需要额外扣除玩家个人奖池5Bet
	PoolCost int

	Type db.BoundType
}

func (b *Bound) name() string {
	s := fmt.Sprintf("(%v, %v]", b.Min, b.Max)
	// if b.Type == GameTypeGame {
	// 	s += "购买"
	// }

	if b.Type == GameTypeGamex2 {
		s += "2倍购买"
	}
	if b.Type == GameTypeGamex3 {
		s += "3倍购买"
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
		// 1 倍购买
		{ // (-oo,0]
			Group: 0,
			Min:   -1,
			Max:   0,
			Type:  GameTypeNormal,
		}, { // （0，10】
			Group: 0,
			Min:   0,
			Max:   10,
			Type:  GameTypeNormal,
		},

		{ // （10，20】
			Group: 1,
			Min:   10,
			Max:   20,
			Type:  GameTypeNormal,
		}, { //（20，30】需要额外扣除玩家个人奖池5Bet
			Group:    1,
			Min:      20,
			Max:      30,
			Type:     GameTypeNormal,
			PoolCost: 5,
		}, { //（30，50】需要额外扣除玩家个人奖池6Bet
			Group:    1,
			Min:      30,
			Max:      50,
			Type:     GameTypeNormal,
			PoolCost: 6,
		}, { //（50，80】需要额外扣除玩家个人奖池8Bet
			Group:    1,
			Min:      50,
			Max:      80,
			Type:     GameTypeNormal,
			PoolCost: 8,
		}, { //（80，150】需要额外扣除玩家个人奖池10Bet
			Group:    1,
			Min:      80,
			Max:      150,
			Type:     GameTypeNormal,
			PoolCost: 10,
		}, { //（150，400】需要额外扣除玩家个人奖池20Bet
			Group:    1,
			Min:      150,
			Max:      400,
			Type:     GameTypeNormal,
			PoolCost: 20,
		}, { //（400，800】需要额外扣除玩家个人奖池30Bet
			Group:    1,
			Min:      400,
			Max:      800,
			Type:     GameTypeNormal,
			PoolCost: 30,
		},

		// 2 倍购买 ===========================
		{ // (-oo,0]
			Group: 2,
			Min:   -1,
			Max:   0,
			Type:  GameTypeGamex2,
		},

		{ // （0，20】
			Group: 3,
			Min:   0,
			Max:   20,
			Type:  GameTypeGamex2,
		}, { // （20，40】需要额外扣除玩家个人奖池5Bet
			Group:    3,
			Min:      20,
			Max:      40,
			Type:     GameTypeGamex2,
			PoolCost: 5,
		}, { //（40，60】需要额外扣除玩家个人奖池6Bet
			Group:    3,
			Min:      40,
			Max:      60,
			Type:     GameTypeGamex2,
			PoolCost: 6,
		}, { //（60，100】需要额外扣除玩家个人奖池8Bet
			Group:    3,
			Min:      60,
			Max:      100,
			Type:     GameTypeGamex2,
			PoolCost: 8,
		}, { //（100，160】需要额外扣除玩家个人奖池10Bet
			Group:    3,
			Min:      100,
			Max:      160,
			Type:     GameTypeGamex2,
			PoolCost: 10,
		}, { //（160，300】需要额外扣除玩家个人奖池20Bet
			Group:    3,
			Min:      160,
			Max:      300,
			Type:     GameTypeGamex2,
			PoolCost: 20,
		},

		{ //（300，800】需要额外扣除玩家个人奖池30Bet
			Group:    4,
			Min:      300,
			Max:      800,
			Type:     GameTypeGamex2,
			PoolCost: 30,
		}, { //（800，1600】需要额外扣除玩家个人奖池50Bet
			Group:    4,
			Min:      800,
			Max:      1600,
			Type:     GameTypeGamex2,
			PoolCost: 50,
		},

		// 3 倍购买 ===========================
		{ //【-1，0】
			Group: 5,
			Min:   -1,
			Max:   0,
			Type:  GameTypeGamex3,
		},

		{ //（0，30】需要额外扣除玩家个人奖池5Bet
			Group:    6,
			Min:      0,
			Max:      30,
			Type:     GameTypeGamex3,
			PoolCost: 5,
		}, { //（30，60】需要额外扣除玩家个人奖池6Bet
			Group:    6,
			Min:      30,
			Max:      60,
			Type:     GameTypeGamex3,
			PoolCost: 6,
		}, { //（60，90】需要额外扣除玩家个人奖池8Bet
			Group:    6,
			Min:      60,
			Max:      90,
			Type:     GameTypeGamex3,
			PoolCost: 8,
		}, { //（90，150】需要额外扣除玩家个人奖池10Bet
			Group:    6,
			Min:      90,
			Max:      150,
			Type:     GameTypeGamex3,
			PoolCost: 10,
		}, { //（150，240】需要额外扣除玩家个人奖池20Bet
			Group:    6,
			Min:      150,
			Max:      240,
			Type:     GameTypeGamex3,
			PoolCost: 20,
		}, { //（240，450】需要额外扣除玩家个人奖池25Bet
			Group:    6,
			Min:      240,
			Max:      450,
			Type:     GameTypeGamex3,
			PoolCost: 25,
		},

		{ //（450，1200】需要额外扣除玩家个人奖池35Bet
			Group:    7,
			Min:      450,
			Max:      1200,
			Type:     GameTypeGamex3,
			PoolCost: 35,
		}, { //（1200，2600】需要额外扣除玩家个人奖池50Bet
			Group:    7,
			Min:      1200,
			Max:      2600,
			Type:     GameTypeGamex3,
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
