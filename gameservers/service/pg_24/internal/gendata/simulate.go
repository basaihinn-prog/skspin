package gendata

import (
	"serve/comm/slotsmongo"
)

type SimulateData = slotsmongo.SimulateData

const Line = 1

// const BuyMul = 75
const BuyMulx2 = 2
const BuyMulx3 = 3

var BuyMulArr = [...]int{BuyMulx2, BuyMulx3}
