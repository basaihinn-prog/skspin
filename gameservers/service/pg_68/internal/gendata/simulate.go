package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 5
const BuyMul = 1

var Cs = []float32{0.2, 2.0, 20.0, 100.0}
var Ml = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
