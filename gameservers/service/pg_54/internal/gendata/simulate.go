package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 20

var Cs = []float32{0.05, 0.5, 2.5, 10}
var Ml = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
