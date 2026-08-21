package gendata

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"serve/comm/redisx"
	"serve/servicefacai/facai_22043/internal"
	"serve/servicefacai/facaicomm"
	"sync"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"serve/comm/db"
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
	Type                 facaicomm.BoundType
}

type CombineDataMng struct {
	buckets map[int][][]primitive.ObjectID

	buyBucketModMap  map[int][][]primitive.ObjectID //桶子分布
	buyBucketsModMap map[int][]primitive.ObjectID   //购买序列

	combine facaicomm.Combine

	combineDatas     map[string]*facaicomm.CombineData
	combineDatas_buy map[int]map[string]*facaicomm.CombineDataBuy

	mtx sync.Mutex
}

var GCombineDataMng *CombineDataMng

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
		return make([][]primitive.ObjectID, len(facaicomm.GBuckets.Bounds))
	}

	type bucketGroups struct {
		normal [][]primitive.ObjectID
	}

	bucketConfig := map[int]*bucketGroups{
		facaicomm.BucketHeartBeat:      {makeSlice()},
		facaicomm.BucketWave:           {makeSlice()},
		facaicomm.BucketGov:            {makeSlice()},
		facaicomm.BucketMix:            {makeSlice()},
		facaicomm.BucketStable:         {makeSlice()},
		facaicomm.BucketHighAward:      {makeSlice()},
		facaicomm.BucketSuperHighAward: {makeSlice()},
	}
	var (
		buyBucketMap = make(map[int][][]primitive.ObjectID, internal.ModCount)
		buyIdsMap    = map[int][]primitive.ObjectID{}
	)
	for _, doc := range docs {
		if doc.Type == facaicomm.GameTypeNormal {
			handleBucket := func(bucketType int, config *bucketGroups, isOff bool) {
				if isOff {
					if doc.Type == facaicomm.GameTypeNormal {
						config.normal[doc.BucketID] = append(config.normal[doc.BucketID], doc.ID)
					}
				}
			}
			handleBucket(facaicomm.BucketHeartBeat, bucketConfig[facaicomm.BucketHeartBeat], doc.BucketHeartBeat == facaicomm.BucketOff)
			handleBucket(facaicomm.BucketWave, bucketConfig[facaicomm.BucketWave], doc.BucketWave == facaicomm.BucketOff)
			handleBucket(facaicomm.BucketGov, bucketConfig[facaicomm.BucketGov], doc.BucketGov == facaicomm.BucketOff)
			handleBucket(facaicomm.BucketMix, bucketConfig[facaicomm.BucketMix], doc.BucketMix == facaicomm.BucketOff)
			handleBucket(facaicomm.BucketStable, bucketConfig[facaicomm.BucketStable], doc.BucketStable == facaicomm.BucketOff)
			handleBucket(facaicomm.BucketHighAward, bucketConfig[facaicomm.BucketHighAward], doc.BucketHighAward == facaicomm.BucketOff)
			handleBucket(facaicomm.BucketSuperHighAward, bucketConfig[facaicomm.BucketSuperHighAward], doc.BucketSuperHighAward == facaicomm.BucketOff)
		} else {
			mod := int(doc.Type)
			if buyBucketMap[mod] == nil {
				buyBucketMap[mod] = makeSlice()
			}
			if buyIdsMap[mod] == nil {
				buyIdsMap[mod] = []primitive.ObjectID{}
			}
			buyBucketMap[mod][doc.BucketID] = append(buyBucketMap[mod][doc.BucketID], doc.ID)
			buyIdsMap[mod] = append(buyIdsMap[mod], doc.ID)
		}

	}

	buckets := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	for k, v := range bucketConfig {
		buckets[k] = v.normal
	}

	GCombineDataMng = &CombineDataMng{
		buckets:          buckets,
		buyBucketModMap:  buyBucketMap,
		buyBucketsModMap: buyIdsMap,
		combineDatas_buy: map[int]map[string]*facaicomm.CombineDataBuy{},

		//combine:      combine,
		combineDatas: map[string]*facaicomm.CombineData{},
	}

	return
}

func findFromSimulate(id primitive.ObjectID) (playResp *facaicomm.SimulateData, err error) {
	one := &facaicomm.SimulateData{}
	err = db.Collection("simulate").FindOne(context.TODO(), db.ID(id)).Decode(one)
	if err != nil {
		return
	}
	playResp = one
	return
}

func (mng *CombineDataMng) get(bet string) (data *facaicomm.CombineData) {
	mng.mtx.Lock()
	defer mng.mtx.Unlock()
	data = mng.combineDatas[bet]
	if data == nil {
		data = &facaicomm.CombineData{
			Combine: mng.combine,
			Buckets: mng.buckets,
		}
		mng.combineDatas[bet] = data
	}
	return
}

// 控制下一局
func (mng *CombineDataMng) ControlNextDataNormal(nextData *redisx.NextMulti, witchBucket int) (playResp *facaicomm.SimulateData, err error) {
	gameType := facaicomm.BoundType(facaicomm.GameTypeNormal)
	ids := []int{}
	if nextData != nil {
		filterBucket := func(min, max float64, gameType facaicomm.BoundType) (ids []int) {
			for _, b := range facaicomm.GBuckets.Bounds {
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

func (mng *CombineDataMng) Next(bet string, WitchBucket int) (playResp *facaicomm.SimulateData, err error) {
	data := mng.get(bet)

	id := data.Next(WitchBucket)
	//id, _ = primitive.ObjectIDFromHex("67ebf18db99b9f8895b91cec")
	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleSimulate(bucketid, witchBucket int) (playResp *facaicomm.SimulateData, err error) {
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

func (mng *CombineDataMng) NextBuy(mod int, bet string) (playResp *facaicomm.SimulateData, err error) {
	mng.mtx.Lock()
	if mng.combineDatas_buy[mod] == nil {
		mng.combineDatas_buy[mod] = make(map[string]*facaicomm.CombineDataBuy)
	}
	data := mng.combineDatas_buy[mod][bet]
	if data == nil {
		data = &facaicomm.CombineDataBuy{
			RawSeries: mng.buyBucketsModMap[mod],
		}
		mng.combineDatas_buy[mod][bet] = data
	}
	mng.mtx.Unlock()

	id := data.Next()
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward(witchBucket int) (playResp *facaicomm.SimulateData, err error) {
	ids := []int{}
	ids = append(ids, facaicomm.GetBucketIds(5, 10, false)...)
	ids = lo.Filter(ids, func(item int, index int) bool {
		return len(mng.buckets[witchBucket][item]) != 0
	})
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward2_5(witchBucket int) (playResp *facaicomm.SimulateData, err error) {
	ids := []int{}
	ids = append(ids, facaicomm.GetBucketIds(2, 5, false)...)
	ids = lo.Filter(ids, func(item int, index int) bool {
		return len(mng.buckets[witchBucket][item]) != 0
	})
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBuyMinData(mod int) (playResp *facaicomm.SimulateData, err error) {
	//bucketId := facaicomm.GetBuyMinBucketId(mod)
	var bucketId int
	for i := range mng.buyBucketModMap[mod] {
		if mng.buyBucketModMap[mod][i] != nil {
			bucketId = i
			break
		}
	}
	data := mng.buyBucketModMap[mod][bucketId]
	id := lo.Sample(data)
	if id == primitive.NilObjectID {
		return nil, errors.New("not found")
	}
	return findFromSimulate(id)
}
