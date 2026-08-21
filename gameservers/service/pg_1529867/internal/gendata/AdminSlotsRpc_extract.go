package gendata

import (
	"serve/comm"
	"serve/comm/mux"
)

func init() {
	mux.RegRpc("/pg_1529867/AdminSlotsRpc/ExtractAction", "提取数据", "gendata", comm.ExtractAction, []comm.ExtractActionPsItem{
		{
			Index: 5,
			Count: 1000,
		},
		{
			Index: 6,
			Count: 2000,
		},
	}).SetOnlyDev()
}
