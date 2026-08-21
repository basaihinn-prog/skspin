package gendata

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"serve/comm/redisx"
	"sync"

	"serve/comm/db"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CombineDataMng struct {
	buckets      map[int][][]primitive.ObjectID
	buckets2     map[int][][]primitive.ObjectID
	buyBucketMap [][]primitive.ObjectID
	buyBuckets   []primitive.ObjectID
	combine      Combine

	combineDatas  map[string]*CombineData
	combineDatas2 map[string]*CombineData

	mtx sync.Mutex
}

const (
	BucketOff = 1
	BucketOn  = 0

	BucketHeartBeat      = 1
	BucketWave           = 2
	BucketGov            = 3
	BucketMix            = 4
	BucketStable         = 5
	BucketHighAward      = 6
	BucketSuperHighAward = 7
	HitBigAwardPercent   = 350
)

var GCombineDataMng *CombineDataMng

func (mng *CombineDataMng) GetBigReward2_5(gameType db.BoundType, witchBucket int) (playResp *SimulateData, err error) {
	ids := []int{}
	if gameType == GameTypeNormal {
		ids = append(ids, GetBucketIds(0, 10, GameTypeNormal)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.buckets[witchBucket][bucketid]
		id := lo.Sample(data)
		return findFromSimulate(id)
	} else if gameType == GameTypeGame {
		ids = append(ids, GetBucketIds(0, 20, GameTypeGame)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets2[witchBucket][item]) != 0
		})
		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.buckets2[witchBucket][bucketid]
		id := lo.Sample(data)
		return findFromSimulate(id)
	}

	return nil, errors.New("not found")
}

func LoadCombineData() (err error) {
	combine, err := getCombine()
	if err != nil {
		return
	}

	docs, err := db.GetSimulate()
	if err != nil {
		return
	}

	makeSlice := func() [][]primitive.ObjectID {
		return make([][]primitive.ObjectID, len(GBuckets.bounds))
	}

	type bucketGroups struct {
		normal [][]primitive.ObjectID
		double [][]primitive.ObjectID
	}

	bucketConfig := map[int]*bucketGroups{
		BucketHeartBeat:      {makeSlice(), makeSlice()},
		BucketWave:           {makeSlice(), makeSlice()},
		BucketGov:            {makeSlice(), makeSlice()},
		BucketMix:            {makeSlice(), makeSlice()},
		BucketStable:         {makeSlice(), makeSlice()},
		BucketHighAward:      {makeSlice(), makeSlice()},
		BucketSuperHighAward: {makeSlice(), makeSlice()},
	}

	for _, doc := range docs {
		bound := GBuckets.GetBound(doc.BucketID)
		if bound.Type != doc.Type {
			return fmt.Errorf("type mismatch for doc %v", doc.ID)
		}

		handleBucket := func(bucketType int, config *bucketGroups, isOff bool) {
			if isOff {
				if doc.Type == GameTypeNormal {
					config.normal[doc.BucketID] = append(config.normal[doc.BucketID], doc.ID)
				}
				if doc.Type == GameTypeGame {
					config.double[doc.BucketID] = append(config.double[doc.BucketID], doc.ID)
				}
			}
		}

		handleBucket(BucketHeartBeat, bucketConfig[BucketHeartBeat], doc.BucketHeartBeat == BucketOff)
		handleBucket(BucketWave, bucketConfig[BucketWave], doc.BucketWave == BucketOff)
		handleBucket(BucketGov, bucketConfig[BucketGov], doc.BucketGov == BucketOff)
		handleBucket(BucketMix, bucketConfig[BucketMix], doc.BucketMix == BucketOff)
		handleBucket(BucketStable, bucketConfig[BucketStable], doc.BucketStable == BucketOff)
		handleBucket(BucketHighAward, bucketConfig[BucketHighAward], doc.BucketHighAward == BucketOff)
		handleBucket(BucketSuperHighAward, bucketConfig[BucketSuperHighAward], doc.BucketSuperHighAward == BucketOff)
	}

	buckets := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	bucketsDouble := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	for k, v := range bucketConfig {
		buckets[k] = v.normal
		bucketsDouble[k] = v.double
	}

	GCombineDataMng = &CombineDataMng{
		buckets:       buckets,
		buckets2:      bucketsDouble,
		buyBuckets:    []primitive.ObjectID{},
		combine:       combine,
		combineDatas:  make(map[string]*CombineData),
		combineDatas2: make(map[string]*CombineData),
	}

	return
}

func findFromSimulate(id primitive.ObjectID) (playResp *SimulateData, err error) {
	one := &SimulateData{}
	err = db.Collection("simulate").FindOne(context.TODO(), db.ID(id)).Decode(one)
	if err != nil {
		return
	}
	playResp = one
	return
}

func (mng *CombineDataMng) get(bet string) (data *CombineData) {
	mng.mtx.Lock()
	defer mng.mtx.Unlock()
	data = mng.combineDatas[bet]
	if data == nil {
		data = &CombineData{
			combine: mng.combine,
			buckets: mng.buckets,
		}
		mng.combineDatas[bet] = data
	}
	return
}

func (mng *CombineDataMng) get2(bet string) (data *CombineData) {
	mng.mtx.Lock()
	defer mng.mtx.Unlock()
	data = mng.combineDatas2[bet]
	if data == nil {
		data = &CombineData{
			combine: mng.combine,
			buckets: mng.buckets2,
		}
		mng.combineDatas2[bet] = data
	}
	return
}

func (mng *CombineDataMng) Next(bet string, gameType db.BoundType, WitchBucket int) (playResp *SimulateData, err error) {
	// mng.mtx.Lock()
	// data := mng.combineDatas[bet]
	// if data == nil {
	// 	data = &CombineData{
	// 		combine: mng.combine,
	// 		buckets: mng.buckets,
	// 	}
	// 	mng.combineDatas[bet] = data
	// }
	// mng.mtx.Unlock()
	id := primitive.NilObjectID
	var data *CombineData

	if gameType == GameTypeNormal {
		data = mng.get(bet)
	} else if gameType == GameTypeGame {
		data = mng.get2(bet)
	} else {
		return nil, errors.New("error type")
	}
	id = data.next(gameType, WitchBucket)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleSimulate(bucketid, witchBucket int, boundType db.BoundType) (playResp *SimulateData, err error) {
	if boundType == GameTypeNormal {
		data := mng.buckets[witchBucket][bucketid]
		id := lo.Sample(data)
		return findFromSimulate(id)
	} else if boundType == GameTypeGame {
		data := mng.buckets2[witchBucket][bucketid]
		id := lo.Sample(data)
		return findFromSimulate(id)
	}
	return nil, errors.New("error Type")
}

func (mng *CombineDataMng) getCombineUsedCount(bucketid int, ty db.BoundType, witchBucket int) int {
	// data := mng.get(0)
	data := mng.get("")
	if ty == GameTypeGame {
		data = mng.get2("")
	}
	if len(data.combineUseCounts) == 0 {
		data.shuffle(ty, witchBucket)
	}
	return data.combineUseCounts[bucketid]
}

func (mng *CombineDataMng) GetBigReward(gameType db.BoundType, witchBucket int) (playResp *SimulateData, err error) {
	ids := []int{}
	if gameType == GameTypeNormal {
		ids = append(ids, GetBucketIds(5, 25, GameTypeNormal)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.buckets[witchBucket][bucketid]
		id := lo.Sample(data)
		return findFromSimulate(id)
	} else if gameType == GameTypeGame {
		ids = append(ids, GetBucketIds(0, 20, GameTypeGame)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets2[witchBucket][item]) != 0
		})
		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.buckets2[witchBucket][bucketid]
		id := lo.Sample(data)
		return findFromSimulate(id)
	}

	return nil, errors.New("not found")
}

// 控制下一局
func (mng *CombineDataMng) ControlNextData(nextData *redisx.NextMulti, gameType db.BoundType, witchBucket int) (playResp *SimulateData, err error) {
	ids := []int{}

	if nextData != nil {
		filterBucket := func(min, max float64, gameType db.BoundType) (ids []int) {
			for _, b := range GBuckets.bounds {
				if min <= b.Min && b.Max <= max && b.Type == gameType {
					ids = append(ids, b.ID)
				}
			}
			return
		}
		// 初始最小和最大倍数
		minMulti := nextData.MinMulti
		maxMulti := nextData.MaxMulti

		// 尝试最多5次扩大范围
		for attempt := 0; attempt < 20; attempt++ {
			if gameType == GameTypeNormal {
				tempIds := filterBucket(minMulti, maxMulti, gameType)
				ids = lo.Filter(tempIds, func(item int, index int) bool {
					return len(mng.buckets[witchBucket][item]) != 0
				})
			} else if gameType == GameTypeGame {
				tempIds := filterBucket(minMulti, maxMulti, gameType)
				ids = lo.Filter(tempIds, func(item int, index int) bool {
					return len(mng.buckets2[witchBucket][item]) != 0
				})
			}

			// 如果找到了符合条件的数据，跳出循环
			if len(ids) > 0 {
				break
			}

			// 扩大搜索范围
			minMulti = math.Max(1, minMulti-5) // 最小不低于1
			maxMulti = maxMulti + 5
		}

		// 如果还是没有找到，使用所有可用的桶
		if len(ids) == 0 {
			if gameType == GameTypeNormal {
				for id, bucket := range mng.buckets[witchBucket] {
					if len(bucket) > 0 {
						ids = append(ids, id)
					}
				}
			} else if gameType == GameTypeGame {
				for id, bucket := range mng.buckets2[witchBucket] {
					if len(bucket) > 0 {
						ids = append(ids, id)
					}
				}
			}
		}
	}

	// 确保ids不为空
	if len(ids) == 0 {
		// 最后的保底措施：记录错误日志，并使用默认值
		log.Printf("警告: 无法找到符合条件的数据 gameType=%v, witchBucket=%v, nextData=%+v",
			gameType, witchBucket, nextData)

		// 获取所有非空桶的ID作为备选
		if gameType == GameTypeNormal {
			for id, bucket := range mng.buckets[witchBucket] {
				if len(bucket) > 0 {
					ids = append(ids, id)
				}
			}
		} else if gameType == GameTypeGame {
			for id, bucket := range mng.buckets2[witchBucket] {
				if len(bucket) > 0 {
					ids = append(ids, id)
				}
			}
		}

		// 如果真的一个都没有，返回错误
		if len(ids) == 0 {
			return nil, fmt.Errorf("没有可用的数据 gameType=%v, witchBucket=%v", gameType, witchBucket)
		}
	}

	bucketid := lo.Sample(ids)
	data := []primitive.ObjectID{}
	if gameType == GameTypeNormal {
		data = mng.buckets[witchBucket][bucketid]
	} else if gameType == GameTypeGame {
		data = mng.buckets2[witchBucket][bucketid]
	}
	id := lo.Sample(data)
	return findFromSimulate(id)
}
