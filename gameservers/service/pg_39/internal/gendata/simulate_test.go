package gendata

import (
	"testing"

	"serve/comm/db"
	"serve/comm/ut"

	"github.com/stretchr/testify/assert"
)

/*
import (
	"game/comm/db"
	"game/comm/ut"
	"game/service/pg_39/internal/core"
	"testing"

	"github.com/zeebo/assert"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_39")

	ps := &core.SimulateParams{
		Count: 10000,
		Weight: core.SimulateParamsWeight{
			Wild:         1,
			Pig:          1,
			Gold:         1,
			Cabbage:      1,
			Firecracker3: 1,
			Firecracker2: 1,
			Firecracker1: 1,
			Nothing:      1,
		},
		Multi: [4]int{50, 30, 20, 10},
	}

	var ret SimulateResult
	err := simulate(ps, &ret)
	assert.Nil(t, err)
	ut.PrintJson(ret)
}

*/

func TestStat(t *testing.T) {

	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_39")

	statRet, err := stat()
	assert.Nil(t, err)

	ut.PrintJson(statRet)
}
