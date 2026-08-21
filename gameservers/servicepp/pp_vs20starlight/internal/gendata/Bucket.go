package gendata

import (
	"log"
	"serve/comm/db"
	"serve/servicepp/ppcomm"

	"github.com/samber/lo"
)

var GBuckets = ppcomm.NewBuckets([]*ppcomm.Bound{
	// 普通模式
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
	},

	{ // （30，40】需要额外扣除玩家个人奖池5Bet
		Group:    5,
		Min:      30,
		Max:      40,
		PoolCost: 5,
	}, { // （40，60】需要额外扣除玩家个人奖池6Bet
		Group:    5,
		Min:      40,
		Max:      60,
		PoolCost: 6,
	}, { // （60，80】需要额外扣除玩家个人奖池8Bet
		Group:    5,
		Min:      60,
		Max:      80,
		PoolCost: 8,
	},

	{ // (80, 100] 需要额外扣除玩家个人奖池10bet
		Group:    6,
		Min:      80,
		Max:      100,
		PoolCost: 10,
	}, { // (100, 150] 需要额外扣除玩家个人奖池15bet
		Group:    6,
		Min:      100,
		Max:      150,
		PoolCost: 15,
	}, { // (150, 200] 需要额外扣除玩家个人奖池20bet
		Group:    6,
		Min:      150,
		Max:      200,
		PoolCost: 20,
	}, { // (200, 250] 需要额外扣除玩家个人奖池20bet
		Group:    6,
		Min:      200,
		Max:      250,
		PoolCost: 20,
	},

	{ // (250, 300] 需要额外扣除玩家个人奖池25bet
		Group:    7,
		Min:      250,
		Max:      300,
		PoolCost: 25,
	}, { // (300, 350] 需要额外扣除玩家个人奖池30bet
		Group:    7,
		Min:      300,
		Max:      350,
		PoolCost: 30,
	}, { // (350, 400] 需要额外扣除玩家个人奖池35bet
		Group:    7,
		Min:      350,
		Max:      400,
		PoolCost: 35,
	},

	{ // (400, 450] 需要额外扣除玩家个人奖池40bet
		Group:    8,
		Min:      400,
		Max:      450,
		PoolCost: 40,
	}, { // (450, 500] 需要额外扣除玩家个人奖池45bet
		Group:    8,
		Min:      450,
		Max:      500,
		PoolCost: 45,
	}, { // (500, 550] 需要额外扣除玩家个人奖池50bet
		Group:    8,
		Min:      500,
		Max:      550,
		PoolCost: 50,
	},

	{ // (550, 600] 需要额外扣除玩家个人奖池55bet
		Group:    9,
		Min:      550,
		Max:      600,
		PoolCost: 55,
	}, { // (600, 650] 需要额外扣除玩家个人奖池60bet
		Group:    9,
		Min:      600,
		Max:      650,
		PoolCost: 60,
	}, { // (650, 700] 需要额外扣除玩家个人奖池60bet
		Group:    9,
		Min:      650,
		Max:      700,
		PoolCost: 60,
	},

	{ // (700, 750] 需要额外扣除玩家个人奖池60bet
		Group:    10,
		Min:      700,
		Max:      750,
		PoolCost: 60,
	}, { // (750, 800] 需要额外扣除玩家个人奖池60bet
		Group:    10,
		Min:      750,
		Max:      800,
		PoolCost: 60,
	}, { // (800, 850] 需要额外扣除玩家个人奖池60bet
		Group:    10,
		Min:      800,
		Max:      850,
		PoolCost: 60,
	},

	{ // (850, 900] 需要额外扣除玩家个人奖池60bet
		Group:    11,
		Min:      850,
		Max:      900,
		PoolCost: 60,
	}, { // (900, 950] 需要额外扣除玩家个人奖池60bet
		Group:    11,
		Min:      900,
		Max:      950,
		PoolCost: 60,
	},

	{ // (950, 1000] 需要额外扣除玩家个人奖池60bet
		Group:    12,
		Min:      950,
		Max:      1000,
		PoolCost: 60,
	}, { // (1000, 99999] 需要额外扣除玩家个人奖池60bet
		Group:    12,
		Min:      1000,
		Max:      99999,
		PoolCost: 60,
	},

	//小游戏
	{ // (-1，0】
		Group:   13,
		Min:     -1,
		Max:     0,
		HasGame: true,
	}, { // (0，10】
		Group:   13,
		Min:     0,
		Max:     10,
		HasGame: true,
	},

	{ // （10，20】
		Group:   14,
		Min:     10,
		Max:     20,
		HasGame: true,
	}, { // （20，30】
		Group:   14,
		Min:     20,
		Max:     30,
		HasGame: true,
	}, { // （30，50】需要额外扣除玩家个人奖池8Bet
		Group:    14,
		Min:      30,
		Max:      50,
		HasGame:  true,
		PoolCost: 8,
	},

	{ // (50, 80] 需要额外扣除玩家个人奖池10bet
		Group:    15,
		Min:      50,
		Max:      80,
		HasGame:  true,
		PoolCost: 10,
	}, { // (80, 100]小游戏 需要额外扣除玩家个人奖池10bet
		Group:    15,
		Min:      80,
		Max:      100,
		HasGame:  true,
		PoolCost: 10,
	}, { // (100, 150]小游戏 需要额外扣除玩家个人奖池20bet
		Group:    15,
		Min:      100,
		Max:      150,
		HasGame:  true,
		PoolCost: 20,
	},

	{ // (150, 200]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    16,
		Min:      150,
		Max:      200,
		HasGame:  true,
		PoolCost: 30,
	}, { // (200, 250]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    16,
		Min:      200,
		Max:      250,
		HasGame:  true,
		PoolCost: 30,
	}, { // (250, 300]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    16,
		Min:      250,
		Max:      300,
		HasGame:  true,
		PoolCost: 30,
	},

	{ // (300, 350]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    17,
		Min:      300,
		Max:      350,
		HasGame:  true,
		PoolCost: 30,
	}, { // (350, 400]小游戏 需要额外扣除玩家个人奖池35bet
		Group:    17,
		Min:      350,
		Max:      400,
		HasGame:  true,
		PoolCost: 35,
	}, { // (400, 450]小游戏 需要额外扣除玩家个人奖池40bet
		Group:    17,
		Min:      400,
		Max:      450,
		HasGame:  true,
		PoolCost: 40,
	},

	{ // (450, 500]小游戏 需要额外扣除玩家个人奖池45bet
		Group:    18,
		Min:      450,
		Max:      500,
		HasGame:  true,
		PoolCost: 45,
	}, { // (500, 550]小游戏 需要额外扣除玩家个人奖池50bet
		Group:    18,
		Min:      500,
		Max:      550,
		HasGame:  true,
		PoolCost: 50,
	}, { // (550, 600]小游戏 需要额外扣除玩家个人奖池55bet
		Group:    18,
		Min:      550,
		Max:      600,
		HasGame:  true,
		PoolCost: 55,
	},

	{ // (600, 650]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    19,
		Min:      600,
		Max:      650,
		HasGame:  true,
		PoolCost: 60,
	}, { // (650, 700]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    19,
		Min:      650,
		Max:      700,
		HasGame:  true,
		PoolCost: 60,
	}, { // (700, 750]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    19,
		Min:      700,
		Max:      750,
		HasGame:  true,
		PoolCost: 60,
	},

	{ // (750, 800]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    20,
		Min:      750,
		Max:      800,
		HasGame:  true,
		PoolCost: 60,
	}, { // (800, 850]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    20,
		Min:      800,
		Max:      850,
		HasGame:  true,
		PoolCost: 60,
	}, { // (850, 900]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    20,
		Min:      850,
		Max:      900,
		HasGame:  true,
		PoolCost: 60,
	},

	{ // (900, 950]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    21,
		Min:      900,
		Max:      950,
		HasGame:  true,
		PoolCost: 60,
	}, { // (950, 1000]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    21,
		Min:      950,
		Max:      1000,
		HasGame:  true,
		PoolCost: 60,
	},

	{ // (1000, 99999]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    22,
		Min:      1000,
		Max:      99999,
		HasGame:  true,
		PoolCost: 60,
	},

	//购买小游戏
	{ // (0，10】
		Group:   23,
		Min:     -1,
		Max:     0,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (0, 10]小游戏购买
		Group:   24,
		Min:     0,
		Max:     10,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (10, 20]小游戏购买
		Group:   24,
		Min:     10,
		Max:     20,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (20, 30]小游戏购买
		Group:   24,
		Min:     20,
		Max:     30,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (30, 50]小游戏购买
		Group:   25,
		Min:     30,
		Max:     50,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (50, 80]小游戏购买
		Group:   25,
		Min:     50,
		Max:     80,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (80, 100]小游戏购买
		Group:   25,
		Min:     80,
		Max:     100,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (100, 150]小游戏购买
		Group:   26,
		Min:     100,
		Max:     150,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (150, 200]小游戏购买
		Group:   26,
		Min:     150,
		Max:     200,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (200, 250]小游戏购买
		Group:   26,
		Min:     200,
		Max:     250,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (250, 300]小游戏购买
		Group:   27,
		Min:     250,
		Max:     300,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (300, 350]小游戏购买
		Group:   27,
		Min:     300,
		Max:     350,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (350, 500]小游戏购买
		Group:   27,
		Min:     350,
		Max:     500,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (500, 550]小游戏购买
		Group:   28,
		Min:     500,
		Max:     550,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (550, 600]小游戏购买
		Group:   28,
		Min:     550,
		Max:     600,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (600, 650]小游戏购买
		Group:   28,
		Min:     600,
		Max:     650,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (650, 700]小游戏购买
		Group:   29,
		Min:     650,
		Max:     700,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (700, 750]小游戏购买
		Group:   29,
		Min:     700,
		Max:     750,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (750, 800]小游戏购买
		Group:   29,
		Min:     750,
		Max:     800,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (800, 850]小游戏购买
		Group:   30,
		Min:     800,
		Max:     850,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (850, 900]小游戏购买
		Group:   30,
		Min:     850,
		Max:     900,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (900, 950]小游戏购买
		Group:   30,
		Min:     900,
		Max:     950,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (950, 1000]小游戏购买
		Group:   31,
		Min:     950,
		Max:     1000,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (1000, 99999]小游戏购买
		Group:   31,
		Min:     1000,
		Max:     99999,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},
})

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

	lo.Must0(len(ids) != 0)
	return
}

func GetBuyMinBucketId() int {
	minRate, maxRate := float64(0), float64(10)
	ty := ppcomm.GameTypeGame
	for _, b := range GBuckets.Bounds {
		if b.Min == minRate && b.Max == maxRate && b.Type == db.BoundType(ty) {
			return b.ID
		}
	}
	log.Panicf("GetBuyMinBucketId, min=%v, max=%v, hasGame=%v", minRate, maxRate)
	return -1
}
