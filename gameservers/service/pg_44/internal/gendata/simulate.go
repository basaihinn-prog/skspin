package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 30
const BuyMul = 75

var Cs = []float32{0.03, 0.3, 1, 6}
var Ml = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
