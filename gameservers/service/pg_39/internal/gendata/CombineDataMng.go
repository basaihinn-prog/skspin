package gendata

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math"
	"serve/comm/redisx"
	"sync"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/db"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	//Type            BoundType
}

type CombineDataMng struct {
	buckets      map[int][][]primitive.ObjectID
	buyBucketMap [][]primitive.ObjectID
	buyBuckets   []primitive.ObjectID
	combine      Combine

	combineDatas map[string]*CombineData

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

func LoadCombineData() (err error) {
	combine, err := getCombine()
	if err != nil {
		return
	}

	docs, err := db.GetSimulateByTableNames("pgSpinData")
	if err != nil {
		return
	}

	makeSlice := func() [][]primitive.ObjectID {
		return make([][]primitive.ObjectID, len(GBuckets.bounds))
	}

	type bucketGroups struct {
		normal [][]primitive.ObjectID
	}

	bucketConfig := map[int]*bucketGroups{
		BucketStable:         {makeSlice()},
		BucketHeartBeat:      {makeSlice()},
		BucketWave:           {makeSlice()},
		BucketGov:            {makeSlice()},
		BucketMix:            {makeSlice()},
		BucketHighAward:      {makeSlice()},
		BucketSuperHighAward: {makeSlice()},
	}

	for _, doc := range docs {
		handleBucket := func(bucketType int, config *bucketGroups, isOff bool) {
			if isOff {
				config.normal[doc.BucketID] = append(config.normal[doc.BucketID], doc.ID)
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
	for k, v := range bucketConfig {
		buckets[k] = v.normal
	}

	GCombineDataMng = &CombineDataMng{
		buckets:      buckets,
		buyBuckets:   []primitive.ObjectID{},
		combine:      combine,
		combineDatas: make(map[string]*CombineData),
	}

	return
}

func findFromSimulate(id primitive.ObjectID) (playResp *bson.M, err error) {
	one := &bson.M{}
	err = db.Collection("pgSpinData").FindOne(context.TODO(), db.ID(id)).Decode(one)
	if err != nil {
		slog.Error("findFromSimulate", "id", id.Hex(), "error", err)
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

func (mng *CombineDataMng) Next(bet string, WitchBucket int) (playResp *bson.M, err error) {
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
	data := mng.get(bet)

	id := data.next(WitchBucket)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleSimulate(bucketid, witchBucket int) (playResp *bson.M, err error) {
	data := mng.buckets[witchBucket][bucketid]

	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) getCombineUsedCount(bucketid int, witchBucket int) int {
	// data := mng.get(0)
	data := mng.get("")
	if len(data.combineUseCounts) == 0 {
		data.shuffle(witchBucket)
	}
	return data.combineUseCounts[bucketid]
}

func (mng *CombineDataMng) GetBigReward(witchBucket int) (playResp *bson.M, err error) {
	ids := []int{}
	ids = append(ids, GetBucketIds(10, 20, true)...)
	ids = append(ids, GetBucketIds(5, 10, false)...)
	ids = lo.Filter(ids, func(item int, index int) bool {
		return len(mng.buckets[witchBucket][item]) != 0
	})
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward2_5(witchBucket int) (playResp *bson.M, err error) {
	ids := []int{}
	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(2, 5, false)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
	}

	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(1, 6, false)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
	}

	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(0, 10, false)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
	}
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

// 控制下一局
func (mng *CombineDataMng) ControlNextData(nextData *redisx.NextMulti, witchBucket int) (playResp *bson.M, err error) {
	ids := []int{}

	if nextData != nil {
		filterBucket := func(min, max float64) (ids []int) {
			for _, b := range GBuckets.bounds {
				if min <= b.Min && b.Max <= max {
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

			tempIds := filterBucket(minMulti, maxMulti)
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
		log.Printf("警告: 无法找到符合条件的数据  witchBucket=%v, nextData=%+v",
			witchBucket, nextData)

		// 获取所有非空桶的ID作为备选
		for id, bucket := range mng.buckets[witchBucket] {
			if len(bucket) > 0 {
				ids = append(ids, id)
			}
		}

		// 如果真的一个都没有，返回错误
		if len(ids) == 0 {
			return nil, fmt.Errorf("没有可用的数据  witchBucket=%v", witchBucket)
		}
	}

	bucketid := lo.Sample(ids)
	data := []primitive.ObjectID{}
	data = mng.buckets[witchBucket][bucketid]

	id := lo.Sample(data)
	return findFromSimulate(id)
}
