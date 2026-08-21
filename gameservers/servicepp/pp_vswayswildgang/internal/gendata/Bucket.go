package gendata

import (
	"log"

	"serve/comm/db"
	"serve/servicepp/ppcomm"

	"github.com/samber/lo"
)

var GBuckets = ppcomm.NewBuckets([]*ppcomm.Bound{
	// Group 0
	{Group: 0, Min: -1, Max: 0, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 1
	{Group: 1, Min: 0, Max: 0.3, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 1, Min: 0.3, Max: 0.5, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 1, Min: 0.5, Max: 1, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 2
	{Group: 2, Min: 1, Max: 2, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 2, Min: 2, Max: 3, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 2, Min: 3, Max: 4, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 3
	{Group: 3, Min: 4, Max: 6, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 3, Min: 6, Max: 8, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 3, Min: 8, Max: 10, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 4
	{Group: 4, Min: 10, Max: 15, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 4, Min: 15, Max: 20, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 4, Min: 20, Max: 30, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 4, Min: 30, Max: 40, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 5
	{Group: 5, Min: 40, Max: 60, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 5, Min: 60, Max: 80, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 5, Min: 80, Max: 100, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 6
	{Group: 6, Min: 100, Max: 150, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 6, Min: 150, Max: 200, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 6, Min: 200, Max: 250, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 7
	{Group: 7, Min: 250, Max: 300, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 7, Min: 300, Max: 350, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 7, Min: 350, Max: 400, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 8
	{Group: 8, Min: 400, Max: 450, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 8, Min: 450, Max: 500, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 8, Min: 500, Max: 550, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 9
	{Group: 9, Min: 550, Max: 600, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 9, Min: 600, Max: 650, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 9, Min: 650, Max: 700, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 10
	{Group: 10, Min: 700, Max: 750, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 10, Min: 750, Max: 800, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 10, Min: 800, Max: 850, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 11
	{Group: 11, Min: 850, Max: 900, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 11, Min: 900, Max: 950, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 11, Min: 950, Max: 1000, HasGame: false, Type: ppcomm.GameTypeNormal},
	{Group: 11, Min: 1000, Max: 99999, HasGame: false, Type: ppcomm.GameTypeNormal},

	// Group 12
	{Group: 12, Min: -1, Max: 0, HasGame: true, Type: ppcomm.GameTypeNormal},

	// Group 13
	{Group: 13, Min: 0, Max: 10, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 13, Min: 10, Max: 20, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 13, Min: 20, Max: 30, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 13, Min: 30, Max: 60, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 13, Min: 60, Max: 80, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 13, Min: 80, Max: 100, HasGame: true, Type: ppcomm.GameTypeNormal},

	// Group 14
	{Group: 14, Min: 100, Max: 150, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 14, Min: 150, Max: 200, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 14, Min: 200, Max: 250, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 14, Min: 250, Max: 300, HasGame: true, Type: ppcomm.GameTypeNormal},

	// Group 15
	{Group: 15, Min: 300, Max: 350, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 15, Min: 350, Max: 400, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 15, Min: 400, Max: 450, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 15, Min: 450, Max: 500, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 15, Min: 500, Max: 550, HasGame: true, Type: ppcomm.GameTypeNormal},

	// Group 16
	{Group: 16, Min: 550, Max: 600, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 16, Min: 600, Max: 650, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 16, Min: 650, Max: 700, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 16, Min: 700, Max: 750, HasGame: true, Type: ppcomm.GameTypeNormal},

	// Group 17
	{Group: 17, Min: 750, Max: 800, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 17, Min: 800, Max: 850, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 17, Min: 850, Max: 900, HasGame: true, Type: ppcomm.GameTypeNormal},

	// Group 18
	{Group: 18, Min: 900, Max: 950, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 18, Min: 950, Max: 1000, HasGame: true, Type: ppcomm.GameTypeNormal},
	{Group: 18, Min: 1000, Max: 99999, HasGame: true, Type: ppcomm.GameTypeNormal},

	// Group 19
	{Group: 19, Min: -1, Max: 0, HasGame: true, Type: ppcomm.GameTypeGame},

	// Group 20
	{Group: 20, Min: 0, Max: 10, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 20, Min: 10, Max: 20, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 20, Min: 20, Max: 30, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 20, Min: 30, Max: 50, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 20, Min: 50, Max: 80, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 20, Min: 80, Max: 100, HasGame: true, Type: ppcomm.GameTypeGame},

	// Group 21
	{Group: 21, Min: 100, Max: 150, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 21, Min: 150, Max: 200, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 21, Min: 200, Max: 250, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 21, Min: 250, Max: 300, HasGame: true, Type: ppcomm.GameTypeGame},

	// Group 22
	{Group: 22, Min: 300, Max: 350, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 22, Min: 350, Max: 500, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 22, Min: 500, Max: 550, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 22, Min: 550, Max: 600, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 22, Min: 600, Max: 650, HasGame: true, Type: ppcomm.GameTypeGame},

	// Group 23
	{Group: 23, Min: 650, Max: 700, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 23, Min: 700, Max: 750, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 23, Min: 750, Max: 800, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 23, Min: 800, Max: 850, HasGame: true, Type: ppcomm.GameTypeGame},

	// Group 24
	{Group: 24, Min: 850, Max: 900, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 24, Min: 900, Max: 950, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 24, Min: 950, Max: 1000, HasGame: true, Type: ppcomm.GameTypeGame},
	{Group: 24, Min: 1000, Max: 99999, HasGame: true, Type: ppcomm.GameTypeGame},

	// Group 25
	{Group: 25, Min: -1, Max: 0, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 25, Min: 0, Max: 10, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 25, Min: 10, Max: 20, HasGame: true, Type: ppcomm.GameTypeSuperGame1},

	// Group 26
	{Group: 26, Min: 20, Max: 30, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 26, Min: 30, Max: 50, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 26, Min: 50, Max: 80, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 26, Min: 80, Max: 100, HasGame: true, Type: ppcomm.GameTypeSuperGame1},

	// Group 27
	{Group: 27, Min: 100, Max: 150, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 27, Min: 150, Max: 200, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 27, Min: 200, Max: 250, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 27, Min: 250, Max: 300, HasGame: true, Type: ppcomm.GameTypeSuperGame1},

	// Group 28
	{Group: 28, Min: 300, Max: 350, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 28, Min: 350, Max: 500, HasGame: true, Type: ppcomm.GameTypeSuperGame1},

	// Group 29
	{Group: 29, Min: 500, Max: 550, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 29, Min: 550, Max: 600, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 29, Min: 600, Max: 650, HasGame: true, Type: ppcomm.GameTypeSuperGame1},

	// Group 30
	{Group: 30, Min: 650, Max: 700, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 30, Min: 700, Max: 750, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 30, Min: 750, Max: 800, HasGame: true, Type: ppcomm.GameTypeSuperGame1},

	// Group 31
	{Group: 31, Min: 800, Max: 850, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 31, Min: 850, Max: 900, HasGame: true, Type: ppcomm.GameTypeSuperGame1},

	// Group 32
	{Group: 32, Min: 900, Max: 950, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 32, Min: 950, Max: 1000, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
	{Group: 32, Min: 1000, Max: 99999, HasGame: true, Type: ppcomm.GameTypeSuperGame1},
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
	//if len(ids) != 0 {
	//	fmt.Println(min, max, hasGame)
	//}
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
func GetSuperBuyMinBucketId() int {
	//minRate, maxRate := float64(0), float64(10)
	minRate, maxRate := float64(20), float64(30)
	ty := ppcomm.GameTypeSuperGame1
	for _, b := range GBuckets.Bounds {
		if b.Min >= minRate && b.Max <= maxRate && b.Type == db.BoundType(ty) {
			return b.ID
		}
	}
	log.Panicf("GetSuperBuyMinBucketId, min=%v, max=%v, hasGame=%v", minRate, maxRate)
	return -1
}
