package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 25
const BuyMul = 75

var Cs = []float32{0.04, 0.4, 2, 8}
var Ml = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
