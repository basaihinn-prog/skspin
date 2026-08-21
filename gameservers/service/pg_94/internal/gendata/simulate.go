package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 10
const BuyMul = 75

var Cs = []float32{0.1, 1.0, 10.0, 50.0}
var Ml = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
