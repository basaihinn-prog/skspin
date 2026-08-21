package gendata

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"serve/comm/redisx"
	"serve/servicejili/jili_226_witch/internal"
	"serve/servicepp/ppcomm"
	"sync"

	"serve/comm/db"
	"serve/servicejili/jili_226_witch/internal/models"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CombineDataMng struct {
	//原始
	combine      Combine
	combineDatas map[string]*CombineData
	mtx          sync.Mutex
	//修改
	buckets map[int][][]primitive.ObjectID
	//新建
	buyBucketMap     [][]primitive.ObjectID
	buyBuckets       []primitive.ObjectID
	combineDatas_buy map[string]*CombineDataBuy

	bucketExtras      map[int][][]primitive.ObjectID
	combineExtraDatas map[string]*CombineData
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
		Extra  [][]primitive.ObjectID
	}

	bucketConfig := map[int]*bucketGroups{
		ppcomm.BucketHeartBeat:      {makeSlice(), makeSlice()},
		ppcomm.BucketWave:           {makeSlice(), makeSlice()},
		ppcomm.BucketGov:            {makeSlice(), makeSlice()},
		ppcomm.BucketMix:            {makeSlice(), makeSlice()},
		ppcomm.BucketStable:         {makeSlice(), makeSlice()},
		ppcomm.BucketHighAward:      {makeSlice(), makeSlice()},
		ppcomm.BucketSuperHighAward: {makeSlice(), makeSlice()},
	}
	var (
		buyBucketMap = make([][]primitive.ObjectID, len(GBuckets.bounds))
		buyIds       []primitive.ObjectID
	)
	for _, doc := range docs {
		bound := GBuckets.GetBound(doc.BucketID)
		if bound.Type != doc.Type {
			return fmt.Errorf("type mismatch for doc %v", doc.ID)
		}
		if doc.Type == ppcomm.GameTypeNormal || doc.Type == internal.GameTypeExtra {
			handleBucket := func(bucketType int, config *bucketGroups, isOff bool) {
				if isOff {
					if doc.Type == ppcomm.GameTypeNormal {
						config.normal[doc.BucketID] = append(config.normal[doc.BucketID], doc.ID)
					}
					if doc.Type == internal.GameTypeExtra {
						config.Extra[doc.BucketID] = append(config.Extra[doc.BucketID], doc.ID)
					}
				}
			}
			handleBucket(ppcomm.BucketHeartBeat, bucketConfig[ppcomm.BucketHeartBeat], doc.BucketHeartBeat == ppcomm.BucketOff)
			handleBucket(ppcomm.BucketWave, bucketConfig[ppcomm.BucketWave], doc.BucketWave == ppcomm.BucketOff)
			handleBucket(ppcomm.BucketGov, bucketConfig[ppcomm.BucketGov], doc.BucketGov == ppcomm.BucketOff)
			handleBucket(ppcomm.BucketMix, bucketConfig[ppcomm.BucketMix], doc.BucketMix == ppcomm.BucketOff)
			handleBucket(ppcomm.BucketStable, bucketConfig[ppcomm.BucketStable], doc.BucketStable == ppcomm.BucketOff)
			handleBucket(ppcomm.BucketHighAward, bucketConfig[ppcomm.BucketHighAward], doc.BucketHighAward == ppcomm.BucketOff)
			handleBucket(ppcomm.BucketSuperHighAward, bucketConfig[ppcomm.BucketSuperHighAward], doc.BucketSuperHighAward == ppcomm.BucketOff)
		}
		if doc.Type == ppcomm.GameTypeGame {
			buyBucketMap[doc.BucketID] = append(buyBucketMap[doc.BucketID], doc.ID)
			buyIds = append(buyIds, doc.ID)
		}
	}

	buckets := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	bucketExtras := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	for k, v := range bucketConfig {
		buckets[k] = v.normal
		bucketExtras[k] = v.Extra
	}

	GCombineDataMng = &CombineDataMng{
		combine:      combine,
		combineDatas: map[string]*CombineData{},
		//新建
		buyBucketMap:      buyBucketMap,
		buckets:           buckets,
		buyBuckets:        buyIds,
		bucketExtras:      bucketExtras,
		combineExtraDatas: map[string]*CombineData{},
		combineDatas_buy:  map[string]*CombineDataBuy{},
	}

	return
}

func findFromSimulate(id primitive.ObjectID) (playResp *models.RawSpin, err error) {
	one := &models.RawSpin{}
	err = db.Collection("rawSpinData").FindOne(context.TODO(), db.ID(id)).Decode(one)
	if err != nil {
		return
	}
	playResp = one
	return
}

func (mng *CombineDataMng) Next(bet string, ty, WitchBucket int) (playResp *models.RawSpin, err error) {
	var data *CombineData

	if ty == internal.GameTypeNormal {
		data = mng.get(bet)
	}
	if ty == internal.GameTypeExtra {
		data = mng.getExtra(bet)
	}

	id := data.next(ty, WitchBucket)
	//id, _ := primitive.ObjectIDFromHex("67f6181801f1999bfea5d255")

	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleSimulate(bucketid int, witchBucket int) (playResp *models.RawSpin, err error) {
	data := mng.buckets[witchBucket][bucketid]

	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) getCombineUsedCount(bucketid, ty int) int {
	var data *CombineData
	if ty == internal.GameTypeNormal {
		data = mng.get("")
		if len(data.combineUseCounts) == 0 {
			data.shuffle(ty, bucketid)
		}
	}

	if ty == internal.GameTypeExtra {
		data = mng.getExtra("")
		if len(data.combineUseCounts) == 0 {
			data.shuffle(ty, bucketid)
		}
	}
	return data.combineUseCounts[bucketid]
}

func (mng *CombineDataMng) NextBuy(bet string) (playResp *models.RawSpin, err error) {
	mng.mtx.Lock()
	data := mng.combineDatas_buy[bet]
	if data == nil {
		data = &CombineDataBuy{
			rawSeries: mng.buyBuckets,
		}
		mng.combineDatas_buy[bet] = data
	}
	mng.mtx.Unlock()

	id := data.next()
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBuyMinData() (playResp *models.RawSpin, err error) {
	bucketId := GetBuyMinBucketId()
	data := mng.buyBucketMap[bucketId]
	id := lo.Sample(data)
	if id == primitive.NilObjectID {
		return nil, errors.New("not found")
	}
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward(ty, witchBucket int) (playResp *models.RawSpin, err error) {
	ids := []int{}
	ids = append(ids, GetBucketIds(10, 20, true, ty)...) //todo 记得还原
	ids = append(ids, GetBucketIds(5, 10, false, ty)...)
	//ids = lo.Filter(ids, func(item int, index int) bool {
	//	return len(mng.buckets[witchBucket][item]) != 0
	//})
	//lo.Must0(len(ids) != 0)
	//bucketid := lo.Sample(ids)
	//data := mng.buckets[witchBucket][bucketid]
	//id := lo.Sample(data)

	var id primitive.ObjectID
	if ty == internal.GameTypeNormal {
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.buckets[witchBucket][bucketid]
		id = lo.Sample(data)
	}
	if ty == internal.GameTypeExtra {
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.bucketExtras[witchBucket][item]) != 0
		})
		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.bucketExtras[witchBucket][bucketid]
		id = lo.Sample(data)
	}

	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward2_5(ty, witchBucket int) (playResp *models.RawSpin, err error) {
	ids := []int{}

	var id primitive.ObjectID
	if ty == internal.GameTypeNormal {
		if len(ids) == 0 {
			ids = append(ids, GetBucketIds(2, 5, false, ty)...)
			ids = lo.Filter(ids, func(item int, index int) bool {
				return len(mng.buckets[witchBucket][item]) != 0
			})
		}
		if len(ids) == 0 {
			ids = append(ids, GetBucketIds(1, 6, false, ty)...)
			ids = lo.Filter(ids, func(item int, index int) bool {
				return len(mng.buckets[witchBucket][item]) != 0
			})
		}
		if len(ids) == 0 {
			ids = append(ids, GetBucketIds(0, 10, false, ty)...)
			ids = lo.Filter(ids, func(item int, index int) bool {
				return len(mng.buckets[witchBucket][item]) != 0
			})
		}

		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.buckets[witchBucket][bucketid]
		id = lo.Sample(data)
	}
	if ty == internal.GameTypeExtra {
		if len(ids) == 0 {
			ids = append(ids, GetBucketIds(2, 5, false, ty)...)
			ids = lo.Filter(ids, func(item int, index int) bool {
				return len(mng.buckets[witchBucket][item]) != 0
			})
		}
		if len(ids) == 0 {
			ids = append(ids, GetBucketIds(1, 6, false, ty)...)
			ids = lo.Filter(ids, func(item int, index int) bool {
				return len(mng.buckets[witchBucket][item]) != 0
			})
		}
		if len(ids) == 0 {
			ids = append(ids, GetBucketIds(0, 10, false, ty)...)
			ids = lo.Filter(ids, func(item int, index int) bool {
				return len(mng.buckets[witchBucket][item]) != 0
			})
		}
		lo.Must0(len(ids) != 0)
		bucketid := lo.Sample(ids)
		data := mng.bucketExtras[witchBucket][bucketid]
		id = lo.Sample(data)
	}

	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleForceSimulate(ty, witchBucket int) (playResp *models.RawSpin, err error) {
	var id primitive.ObjectID
	if ty == internal.GameTypeNormal {
		data := mng.buckets[witchBucket][0]
		id = lo.Sample(data)
	}
	if ty == internal.GameTypeExtra {
		data := mng.bucketExtras[witchBucket][29]
		id = lo.Sample(data)
	}
	return findFromSimulate(id)
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

func (mng *CombineDataMng) getExtra(bet string) (data *CombineData) {
	mng.mtx.Lock()
	defer mng.mtx.Unlock()
	data = mng.combineExtraDatas[bet]
	if data == nil {
		data = &CombineData{
			combine: mng.combine,
			buckets: mng.bucketExtras,
		}
		mng.combineExtraDatas[bet] = data
	}
	return
}

func (mng *CombineDataMng) ControlNextData(nextData *redisx.NextMulti, gameType int, witchBucket int) (playResp *models.RawSpin, err error) {
	ids := []int{}

	if nextData != nil {
		filterBucket := func(min, max float64, gameType int) (ids []int) {
			for _, b := range GBuckets.bounds {
				if min <= b.Min && b.Max <= max && b.Type == db.BoundType(gameType) {
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
			} else if gameType == GameTypeExtra {
				tempIds := filterBucket(minMulti, maxMulti, gameType)
				ids = lo.Filter(tempIds, func(item int, index int) bool {
					return len(mng.bucketExtras[witchBucket][item]) != 0
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
			} else if gameType == GameTypeExtra {
				for id, bucket := range mng.bucketExtras[witchBucket] {
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
		} else if gameType == GameTypeExtra {
			for id, bucket := range mng.bucketExtras[witchBucket] {
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
	} else if gameType == GameTypeExtra {
		data = mng.bucketExtras[witchBucket][bucketid]
	}
	id := lo.Sample(data)
	return findFromSimulate(id)
}
