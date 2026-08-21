package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 9
const BuyMul = 1

var Cs = []float32{0.1, 1, 5, 20}
var Ml = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
