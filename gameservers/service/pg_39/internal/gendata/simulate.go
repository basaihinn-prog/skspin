package gendata

/*
import (
	"context"
	"fmt"
	"game/comm/db"
	"game/duck/logger"
	"game/service/pg_39/internal/core"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func simulateOne(ps *core.SimulateParams) *core.PlayResp {
	// spinResp := spin(ps)

	playResp := core.Play(ps)
	// playResp.ResumeByDi(1)

	bucketId := GBuckets.GetBucket(playResp.Times)
	if bucketId == -1 {
		return nil
	}

	playResp.BucketId = bucketId

	return playResp
}

type SimulateResult struct {
	AvgMulti float64
}

func simulate(ps *core.SimulateParams, ret *SimulateResult) (err error) {
	statRet, err := stat()
	if err != nil {
		return
	}

	coll := db.Collection("simulate")

	var sum float64
	var count int

	opts := options.InsertMany().SetOrdered(false)
	arr := make([]interface{}, 0, 1000)
	for i := 0; i < ps.Count; i++ {
		one := simulateOne(ps)
		if one == nil {
			continue
		}

		statU := statRet.List[one.BucketId]
		if statU.Need <= statU.Count {
			continue
		}

		arr = append(arr, one)
		sum += one.Times
		count++
		statU.Count++

		if len(arr) >= 1000 {
			_, inserterr := coll.InsertMany(context.TODO(), arr, opts)
			if inserterr != nil {
				logger.Err(inserterr)
			}

			arr = arr[:0]
			fmt.Printf("进度 %d/%d\n", i, ps.Count)
		}
	}

	if len(arr) != 0 {
		_, inserterr := coll.InsertMany(context.TODO(), arr, opts)
		if inserterr != nil {
			logger.Err(inserterr)
		}

		arr = arr[:0]
	}

	if count != 0 {
		ret.AvgMulti = float64(sum) / float64(count)
	}
	return
}

*/
