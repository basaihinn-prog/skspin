package gendata

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"serve/comm/redisx"
	"sync"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"serve/comm/db"
	"serve/servicejdb/jdbcomm"
)

type playitem struct {
	ID                   primitive.ObjectID `bson:"_id"`
	BucketHeartBeat      int                `bson:"BucketHeartBeat"`
	BucketWave           int                `bson:"BucketWave"`
	BucketGov            int                `bson:"BucketGov"`
	BucketMix            int                `bson:"BucketMix"`
	BucketStable         int                `bson:"BucketStable"`
	BucketHighAward      int                `bson:"BucketHighAward"`
	BucketSuperHighAward int                `bson:"BucketSuperHighAward"`
	BucketID             int
	Type                 jdbcomm.BoundType
}

type CombineDataMng struct {
	buckets           map[int][][]primitive.ObjectID
	buyBucketMap      [][]primitive.ObjectID
	buyBuckets        []primitive.ObjectID
	superBuyBucketMap [][]primitive.ObjectID
	superBuyBuckets   []primitive.ObjectID

	combine jdbcomm.Combine

	combineDatas           map[string]*jdbcomm.CombineData
	combineDatas_buy       map[string]*jdbcomm.CombineDataBuy
	combineDatas_super_buy map[string]*jdbcomm.CombineDataBuy

	mtx sync.Mutex
}

var GCombineDataMng *CombineDataMng
var ErrBucketCount = 0
var ErrTypeCount = 0

func LoadCombineData() (err error) {
	//combine, err := getCombine()
	//if err != nil {
	//	return
	//}

	//var docs []playitem
	docs, err := db.GetSimulate()
	if err != nil {
		return
	}

	makeSlice := func() [][]primitive.ObjectID {
		return make([][]primitive.ObjectID, len(jdbcomm.GBuckets.Bounds))
	}

	type bucketGroups struct {
		normal [][]primitive.ObjectID
	}

	bucketConfig := map[int]*bucketGroups{
		jdbcomm.BucketHeartBeat:      {makeSlice()},
		jdbcomm.BucketWave:           {makeSlice()},
		jdbcomm.BucketGov:            {makeSlice()},
		jdbcomm.BucketMix:            {makeSlice()},
		jdbcomm.BucketStable:         {makeSlice()},
		jdbcomm.BucketHighAward:      {makeSlice()},
		jdbcomm.BucketSuperHighAward: {makeSlice()},
	}
	var (
		buyBucketMap      = make([][]primitive.ObjectID, len(jdbcomm.GBuckets.Bounds))
		superBuyBucketMap = make([][]primitive.ObjectID, len(jdbcomm.GBuckets.Bounds))
		buyIds            []primitive.ObjectID
		superBuyIds       []primitive.ObjectID
	)
	for _, doc := range docs {
		if doc.Type == jdbcomm.GameTypeNormal {
			handleBucket := func(bucketType int, config *bucketGroups, isOff bool) {
				if isOff {
					if doc.Type == jdbcomm.GameTypeNormal {
						config.normal[doc.BucketID] = append(config.normal[doc.BucketID], doc.ID)
					}
				}
			}
			handleBucket(jdbcomm.BucketHeartBeat, bucketConfig[jdbcomm.BucketHeartBeat], doc.BucketHeartBeat == jdbcomm.BucketOff)
			handleBucket(jdbcomm.BucketWave, bucketConfig[jdbcomm.BucketWave], doc.BucketWave == jdbcomm.BucketOff)
			handleBucket(jdbcomm.BucketGov, bucketConfig[jdbcomm.BucketGov], doc.BucketGov == jdbcomm.BucketOff)
			handleBucket(jdbcomm.BucketMix, bucketConfig[jdbcomm.BucketMix], doc.BucketMix == jdbcomm.BucketOff)
			handleBucket(jdbcomm.BucketStable, bucketConfig[jdbcomm.BucketStable], doc.BucketStable == jdbcomm.BucketOff)
			handleBucket(jdbcomm.BucketHighAward, bucketConfig[jdbcomm.BucketHighAward], doc.BucketHighAward == jdbcomm.BucketOff)
			handleBucket(jdbcomm.BucketSuperHighAward, bucketConfig[jdbcomm.BucketSuperHighAward], doc.BucketSuperHighAward == jdbcomm.BucketOff)
		}
		if doc.Type == jdbcomm.GameTypeGame {
			buyBucketMap[doc.BucketID] = append(buyBucketMap[doc.BucketID], doc.ID)
			buyIds = append(buyIds, doc.ID)
		}
		if doc.Type == jdbcomm.GameTypeSuperGame1 {
			superBuyBucketMap[doc.BucketID] = append(superBuyBucketMap[doc.BucketID], doc.ID)
			superBuyIds = append(superBuyIds, doc.ID)
		}
	}

	buckets := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	for k, v := range bucketConfig {
		buckets[k] = v.normal
	}

	GCombineDataMng = &CombineDataMng{
		buckets:          buckets,
		buyBuckets:       buyIds,
		buyBucketMap:     buyBucketMap,
		combineDatas_buy: map[string]*jdbcomm.CombineDataBuy{},

		superBuyBuckets:        superBuyIds,
		superBuyBucketMap:      superBuyBucketMap,
		combineDatas_super_buy: map[string]*jdbcomm.CombineDataBuy{},

		//combine:      combine,
		combineDatas: map[string]*jdbcomm.CombineData{},
	}

	return
}

func findFromSimulate(id primitive.ObjectID) (playResp *jdbcomm.SimulateData, err error) {
	one := &jdbcomm.SimulateData{}
	err = db.Collection("simulate").FindOne(context.TODO(), db.ID(id)).Decode(one)
	if err != nil {
		return
	}
	playResp = one
	return
}

func (mng *CombineDataMng) get(bet string) (data *jdbcomm.CombineData) {
	mng.mtx.Lock()
	defer mng.mtx.Unlock()
	data = mng.combineDatas[bet]
	if data == nil {
		data = &jdbcomm.CombineData{
			Combine: mng.combine,
			Buckets: mng.buckets,
		}
		mng.combineDatas[bet] = data
	}
	return
}

// 控制下一局
func (mng *CombineDataMng) ControlNextDataNormal(nextData *redisx.NextMulti, witchBucket int) (playResp *jdbcomm.SimulateData, err error) {
	gameType := jdbcomm.BoundType(jdbcomm.GameTypeNormal)
	ids := []int{}
	if nextData != nil {
		filterBucket := func(min, max float64, gameType jdbcomm.BoundType) (ids []int) {
			for _, b := range jdbcomm.GBuckets.Bounds {
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
			tempIds := filterBucket(minMulti, maxMulti, gameType)
			ids = lo.Filter(tempIds, func(item int, index int) bool {
				return len(mng.buckets[witchBucket][item]) != 0
			})
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
			for id, bucket := range mng.buckets[witchBucket] {
				if len(bucket) > 0 {
					ids = append(ids, id)
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
		for id, bucket := range mng.buckets[witchBucket] {
			if len(bucket) > 0 {
				ids = append(ids, id)
			}
		}
		// 如果真的一个都没有，返回错误
		if len(ids) == 0 {
			return nil, fmt.Errorf("没有可用的数据 gameType=%v, witchBucket=%v", gameType, witchBucket)
		}
	}
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) Next(bet string, WitchBucket int) (playResp *jdbcomm.SimulateData, err error) {
	data := mng.get(bet)

	id := data.Next(WitchBucket)
	//id, _ = primitive.ObjectIDFromHex("67d380cd882c2437cda34817")
	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleSimulate(bucketid, witchBucket int) (playResp *jdbcomm.SimulateData, err error) {
	data := mng.buckets[witchBucket][bucketid]

	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) getCombineUsedCount(bucketid int) int {
	data := mng.get("")
	if len(data.CombineUseCounts) == 0 {
		data.Shuffle(bucketid)
	}
	return data.CombineUseCounts[bucketid]
}

func (mng *CombineDataMng) NextBuy(bet string) (playResp *jdbcomm.SimulateData, err error) {
	mng.mtx.Lock()
	data := mng.combineDatas_buy[bet]
	if data == nil {
		data = &jdbcomm.CombineDataBuy{
			RawSeries: mng.buyBuckets,
		}
		mng.combineDatas_buy[bet] = data
	}
	mng.mtx.Unlock()

	id := data.Next()
	return findFromSimulate(id)
}

func (mng *CombineDataMng) NextSuperBuy(bet string) (playResp *jdbcomm.SimulateData, err error) {
	mng.mtx.Lock()
	data := mng.combineDatas_super_buy[bet]
	if data == nil {
		data = &jdbcomm.CombineDataBuy{
			RawSeries: mng.superBuyBuckets,
		}
		mng.combineDatas_super_buy[bet] = data
	}
	mng.mtx.Unlock()

	id := data.Next()
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward(witchBucket int) (playResp *jdbcomm.SimulateData, err error) {
	ids := []int{}
	ids = append(ids, jdbcomm.GetBucketIds(10, 20, true)...)
	ids = append(ids, jdbcomm.GetBucketIds(5, 10, false)...)
	ids = lo.Filter(ids, func(item int, index int) bool {
		return len(mng.buckets[witchBucket][item]) != 0
	})
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward2_5(witchBucket int) (playResp *jdbcomm.SimulateData, err error) {
	ids := []int{}

	if len(ids) == 0 {
		ids = append(ids, jdbcomm.GetBucketIds(2, 5, false)...)
	}

	if len(ids) == 0 {
		ids = append(ids, jdbcomm.GetBucketIds(1, 10, false)...)
	}

	if len(ids) == 0 {
		ids = append(ids, jdbcomm.GetBucketIds(0, 20, false)...)
	}

	ids = lo.Filter(ids, func(item int, index int) bool {
		return len(mng.buckets[witchBucket][item]) != 0
	})
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBuyMinData() (playResp *jdbcomm.SimulateData, err error) {
	bucketId := jdbcomm.GetBuyMinBucketId()
	data := mng.buyBucketMap[bucketId]
	id := lo.Sample(data)
	if id == primitive.NilObjectID {
		return nil, errors.New("not found")
	}
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetSuperBuyMinData() (playResp *jdbcomm.SimulateData, err error) {
	bucketId := jdbcomm.GetSuperBuyMinBucketId()
	data := mng.superBuyBucketMap[bucketId]
	id := lo.Sample(data)
	if id == primitive.NilObjectID {
		return nil, errors.New("not found")
	}
	return findFromSimulate(id)
}
