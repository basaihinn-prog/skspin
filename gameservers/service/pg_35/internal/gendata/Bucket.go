package gendata

import (
	"fmt"
	"log"
	"serve/comm/db"

	"github.com/samber/lo"
)

const (
	GameTypeNormal = 0          // 普通模式
	GameTypeGame   = BuyMul - 1 // 购买小游戏
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
	if b.Type == GameTypeGame {
		s += "购买"
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
		{ // (-oo,0]
			Group: 0,
			Min:   -1,
			Max:   0,
		},

		{ // （0，0.25】
			Group: 1,
			Min:   0,
			Max:   0.25,
		}, { // （0.25，0.5】
			Group: 1,
			Min:   0.25,
			Max:   0.5,
		}, { // （0.5，1】
			Group: 1,
			Min:   0.5,
			Max:   1,
		},

		{ // (1,2]
			Group: 2,
			Min:   1,
			Max:   2,
		}, { // （2，3】
			Group: 2,
			Min:   2,
			Max:   3,
		}, { // （3，4】
			Group: 2,
			Min:   3,
			Max:   4,
		},

		{ // （4，6】
			Group: 3,
			Min:   4,
			Max:   6,
		}, { // （6，8】
			Group: 3,
			Min:   6,
			Max:   8,
		}, { // （8，10】
			Group: 3,
			Min:   8,
			Max:   10,
		},

		{ // （10，15】
			Group: 4,
			Min:   10,
			Max:   15,
		}, { // （15，20】
			Group: 4,
			Min:   15,
			Max:   20,
		}, { //（20，30】需要额外扣除玩家个人奖池5Bet
			Group:    4,
			Min:      20,
			Max:      30,
			PoolCost: 5,
		}, { // （30，40】需要额外扣除玩家个人奖池5Bet
			Group:    4,
			Min:      30,
			Max:      40,
			PoolCost: 5,
		},

		{ // （40，60】需要额外扣除玩家个人奖池6Bet
			Group:    5,
			Min:      40,
			Max:      60,
			PoolCost: 6,
		}, { // （60，80】需要额外扣除玩家个人奖池8Bet
			Group:    5,
			Min:      60,
			Max:      80,
			PoolCost: 8,
		}, { // （80，100】需要额外扣除玩家个人奖池10et
			Group:    5,
			Min:      80,
			Max:      100,
			PoolCost: 10,
		},

		{ // （100，99999】
			Group:    6,
			Min:      100,
			Max:      99999,
			PoolCost: 50,
		},

		//小游戏
		{ // (0，10】
			Group:   7,
			Min:     -1,
			Max:     0,
			HasGame: true,
		},
		{ // (0，10】
			Group:   7,
			Min:     0,
			Max:     10,
			HasGame: true,
		},

		{ // （10，20】
			Group:   8,
			Min:     10,
			Max:     20,
			HasGame: true,
		}, { // （20，30】
			Group:   8,
			Min:     20,
			Max:     30,
			HasGame: true,
		}, { // （30，50】需要额外扣除玩家个人奖池8Bet
			Group:    8,
			Min:      30,
			Max:      50,
			HasGame:  true,
			PoolCost: 8,
		}, { // （50，80】需要额外扣除玩家个人奖池10Bet
			Group:    8,
			Min:      50,
			Max:      80,
			HasGame:  true,
			PoolCost: 10,
		},

		{ // （80，100】需要额外扣除玩家个人奖池10Bet
			Group:    9,
			Min:      80,
			Max:      100,
			HasGame:  true,
			PoolCost: 10,
		}, { // （100，150】需要额外扣除玩家个人奖池20Bet
			Group:    9,
			Min:      100,
			Max:      150,
			HasGame:  true,
			PoolCost: 20,
		},

		{ // （150，200】需要额外扣除玩家个人奖池30Bet
			Group:    10,
			Min:      150,
			Max:      200,
			HasGame:  true,
			PoolCost: 30,
		}, { // （200，250】需要额外扣除玩家个人奖池30Bet
			Group:    10,
			Min:      200,
			Max:      250,
			HasGame:  true,
			PoolCost: 30,
		},

		{ // （250，99999】
			Group:    11,
			Min:      250,
			Max:      99999,
			HasGame:  true,
			PoolCost: 50,
		},

		//购买小游戏
		{ // 【0，10】
			Group:   12,
			Min:     -1,
			Max:     0,
			HasGame: true,
			Type:    GameTypeGame,
		},
		{ // 【0，10】
			Group:   12,
			Min:     0,
			Max:     10,
			HasGame: true,
			Type:    GameTypeGame,
		},

		{ // （10，20】
			Group:   13,
			Min:     10,
			Max:     20,
			HasGame: true,
			Type:    GameTypeGame,
		}, { // （20，30】
			Group:   13,
			Min:     20,
			Max:     30,
			HasGame: true,
			Type:    GameTypeGame,
		}, { // （30，50】
			Group:   13,
			Min:     30,
			Max:     50,
			HasGame: true,
			Type:    GameTypeGame,
		}, { // （50，80】
			Group:   13,
			Min:     50,
			Max:     80,
			HasGame: true,
			Type:    GameTypeGame,
		},

		{ // （80，100】
			Group:   14,
			Min:     80,
			Max:     100,
			HasGame: true,
			Type:    GameTypeGame,
		}, { // （100，150】
			Group:   14,
			Min:     100,
			Max:     150,
			HasGame: true,
			Type:    GameTypeGame,
		},

		{ // （150，200】
			Group:   15,
			Min:     150,
			Max:     200,
			HasGame: true,
			Type:    GameTypeGame,
		}, { // （200，250】
			Group:   15,
			Min:     200,
			Max:     250,
			HasGame: true,
			Type:    GameTypeGame,
		},

		{ // （250，300】
			Group:   16,
			Min:     250,
			Max:     300,
			HasGame: true,
			Type:    GameTypeGame,
		}, { // （300，350】
			Group:   16,
			Min:     300,
			Max:     350,
			HasGame: true,
			Type:    GameTypeGame,
		}, { // （350，500】
			Group:   16,
			Min:     350,
			Max:     500,
			HasGame: true,
			Type:    GameTypeGame,
		},

		{ // （500，99999】
			Group:   17,
			Min:     500,
			Max:     99999,
			HasGame: true,
			Type:    GameTypeGame,
		},
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

func getBucketId(min, max float64, hasGame bool) int {
	for _, b := range GBuckets.bounds {
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
	for _, b := range GBuckets.bounds {
		if min <= b.Min && b.Max <= max && b.HasGame == hasGame {
			ids = append(ids, b.ID)
		}
	}

	lo.Must0(len(ids) != 0)
	return
}
