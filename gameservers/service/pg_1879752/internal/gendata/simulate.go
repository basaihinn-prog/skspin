package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 10
const BuyMul = 1.5

var Cs = []float32{1, 10, 100, 500}
var Ml = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
