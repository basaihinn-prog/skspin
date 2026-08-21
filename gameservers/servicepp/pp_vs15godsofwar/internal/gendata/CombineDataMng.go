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
	"serve/servicepp/pp_vs15godsofwar/internal"
	"serve/servicepp/ppcomm"

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
	Type                 db.BoundType
}

type CombineDataMng struct {
	combine ppcomm.Combine

	buckets      []map[int][][]primitive.ObjectID
	combineDatas []map[string]*ppcomm.CombineData

	buyBuckets       [][]primitive.ObjectID
	buyBucketMap     [][][]primitive.ObjectID
	combineDatas_buy []map[string]*ppcomm.CombineDataBuy

	SuperBuyBuckets        [][]primitive.ObjectID
	superBuyBucketMap      [][][]primitive.ObjectID
	combineDatas_super_buy []map[string]*ppcomm.CombineDataBuy

	mtx sync.Mutex
}

var GCombineDataMng *CombineDataMng
var ErrBucketCount = 0
var ErrTypeCount = 0

func (mng *CombineDataMng) GetBigReward2_5(witchBucket, WitchGod int) (playResp *ppcomm.SimulateData, err error) {
	ids := []int{}
	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(2, 5, false)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[WitchGod][witchBucket][item]) != 0
		})
	}
	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(1, 6, false)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[WitchGod][witchBucket][item]) != 0
		})
	}
	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(0, 10, false)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[WitchGod][witchBucket][item]) != 0
		})
	}

	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[WitchGod][witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
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
		return make([][]primitive.ObjectID, len(GBuckets.Bounds))
	}

	type bucketGroups struct {
		hbucket [][]primitive.ObjectID
		zbucket [][]primitive.ObjectID
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
		buckets           []map[int][][]primitive.ObjectID
		buyBucketMap      = make([][][]primitive.ObjectID, len(internal.GetWitchGod()))
		buySuperBucketMap = make([][][]primitive.ObjectID, len(internal.GetWitchGod()))
		buyIds            [][]primitive.ObjectID
		superBuyIds       [][]primitive.ObjectID
	)
	for i := range buyBucketMap {
		buyBucketMap[i] = make([][]primitive.ObjectID, len(GBuckets.Bounds))
		buySuperBucketMap[i] = make([][]primitive.ObjectID, len(GBuckets.Bounds))
	}
	gods := internal.GetWitchGod()
	buckets = make([]map[int][][]primitive.ObjectID, len(gods))
	buyIds = make([][]primitive.ObjectID, len(gods))
	superBuyIds = make([][]primitive.ObjectID, len(gods))
	var hbuyId []primitive.ObjectID
	var hsuperBuyId []primitive.ObjectID

	var zbuyId []primitive.ObjectID
	var zsuperBuyId []primitive.ObjectID
	for _, doc := range docs {
		//bound := GBuckets.GetBound(doc.BucketID)
		//if bound.Type != doc.Type {
		//	return fmt.Errorf("type mismatch for doc %v", doc.ID)
		//}
		handleBucket := func(bucketType int, config *bucketGroups, isOff bool) {
			if isOff {
				switch doc.Type {
				case internal.GodFindBucket(gods[internal.HadesInd], ppcomm.GameTypeNormal):
					config.hbucket[doc.BucketID] = append(config.hbucket[doc.BucketID], doc.ID)
				case internal.GodFindBucket(gods[internal.HadesInd], ppcomm.GameTypeGame):
					buyBucketMap[internal.HadesInd][doc.BucketID] = append(buyBucketMap[internal.HadesInd][doc.BucketID], doc.ID)
					hbuyId = append(hbuyId, doc.ID)
				case internal.GodFindBucket(gods[internal.HadesInd], ppcomm.GameTypeSuperGame1):
					buySuperBucketMap[internal.HadesInd][doc.BucketID] = append(buySuperBucketMap[internal.HadesInd][doc.BucketID], doc.ID)
					hsuperBuyId = append(hsuperBuyId, doc.ID)
				case internal.GodFindBucket(gods[internal.ZuseInd], ppcomm.GameTypeNormal):
					config.zbucket[doc.BucketID] = append(config.zbucket[doc.BucketID], doc.ID)
				case internal.GodFindBucket(gods[internal.ZuseInd], ppcomm.GameTypeGame):
					buyBucketMap[internal.ZuseInd][doc.BucketID] = append(buyBucketMap[internal.ZuseInd][doc.BucketID], doc.ID)
					zbuyId = append(zbuyId, doc.ID)
				case internal.GodFindBucket(gods[internal.ZuseInd], ppcomm.GameTypeSuperGame1):
					buySuperBucketMap[internal.ZuseInd][doc.BucketID] = append(buySuperBucketMap[internal.ZuseInd][doc.BucketID], doc.ID)
					zsuperBuyId = append(zsuperBuyId, doc.ID)
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

	hbucket := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	zbucket := make(map[int][][]primitive.ObjectID, len(bucketConfig))
	for k, v := range bucketConfig {
		hbucket[k] = v.hbucket
		zbucket[k] = v.zbucket
	}
	buckets[internal.HadesInd] = hbucket
	buckets[internal.ZuseInd] = zbucket
	buyIds[internal.HadesInd] = hbuyId
	buyIds[internal.ZuseInd] = zbuyId
	superBuyIds[internal.HadesInd] = hsuperBuyId
	superBuyIds[internal.ZuseInd] = zsuperBuyId
	GCombineDataMng = &CombineDataMng{
		buckets:                buckets,
		buyBucketMap:           buyBucketMap,
		superBuyBucketMap:      buySuperBucketMap,
		buyBuckets:             buyIds,
		SuperBuyBuckets:        superBuyIds,
		combine:                combine,
		combineDatas:           make([]map[string]*ppcomm.CombineData, len(internal.GetWitchGod())),
		combineDatas_buy:       make([]map[string]*ppcomm.CombineDataBuy, len(internal.GetWitchGod())),
		combineDatas_super_buy: make([]map[string]*ppcomm.CombineDataBuy, len(internal.GetWitchGod())),
	}

	return
}

func findFromSimulate(id primitive.ObjectID) (playResp *ppcomm.SimulateData, err error) {
	one := &ppcomm.SimulateData{}
	err = db.Collection("simulate").FindOne(context.TODO(), db.ID(id)).Decode(one)
	if err != nil {
		return
	}
	playResp = one
	return
}

func (mng *CombineDataMng) get(bet string, WitchGod int) (data *ppcomm.CombineData) {
	mng.mtx.Lock()
	defer mng.mtx.Unlock()
	data = mng.combineDatas[WitchGod][bet]
	if data == nil {
		if mng.combineDatas[WitchGod] == nil {
			mng.combineDatas[WitchGod] = make(map[string]*ppcomm.CombineData)
		}
		data = &ppcomm.CombineData{
			Combine: mng.combine,
			Buckets: mng.buckets[WitchGod],
		}
		mng.combineDatas[WitchGod][bet] = data
	}
	return
}

// 控制下一局
func (mng *CombineDataMng) ControlNextDataNormal(nextData *redisx.NextMulti, witchBucket, god int) (playResp *ppcomm.SimulateData, err error) {
	gameType := db.BoundType(ppcomm.GameTypeNormal)
	ids := []int{}
	if nextData != nil {
		filterBucket := func(min, max float64, gameType db.BoundType, god int) (ids []int) {
			for _, b := range GBuckets.Bounds {
				if min <= b.Min && b.Max <= max && b.Type == gameType && b.God == god {
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
			tempIds := filterBucket(minMulti, maxMulti, gameType, god)
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
	data := mng.buckets[god][witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) Next(bet string, WitchBucket, WitchGod int) (playResp *ppcomm.SimulateData, err error) {
	data := mng.get(bet, WitchGod)

	id := data.Next(WitchBucket)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleSimulate(bucketid, witchBucket, WitchGod int) (playResp *ppcomm.SimulateData, err error) {
	if WitchGod == internal.ZuseInd {
		bucketid = 110
	}
	data := mng.buckets[WitchGod][witchBucket][bucketid]

	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) getCombineUsedCount(bucketid, WitchGod int) int {
	data := mng.get("", WitchGod)
	if len(data.CombineUseCounts) == 0 {
		data.Shuffle(bucketid)
	}
	return data.CombineUseCounts[bucketid]
}

func (mng *CombineDataMng) GetBigReward(witchBucket, WitchGod int) (playResp *ppcomm.SimulateData, err error) {
	ids := []int{}
	ids = append(ids, GetBucketIds(10, 20, true)...)
	ids = append(ids, GetBucketIds(5, 10, false)...)
	ids = lo.Filter(ids, func(item int, index int) bool {
		return len(mng.buckets[WitchGod][witchBucket][item]) != 0
	})
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[WitchGod][witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) NextBuy(bet string, witchGod, ty int) (playResp *ppcomm.SimulateData, err error) {
	mng.mtx.Lock()
	//这里没跟神走
	var data *ppcomm.CombineDataBuy
	if ty == internal.PurHadesBuy || ty == internal.PurZuseBuy {
		data = mng.combineDatas_buy[witchGod][bet]
		if data == nil {
			if mng.combineDatas_buy[witchGod] == nil {
				mng.combineDatas_buy[witchGod] = make(map[string]*ppcomm.CombineDataBuy)
			}
			data = &ppcomm.CombineDataBuy{
				RawSeries: mng.buyBuckets[witchGod],
			}
			mng.combineDatas_buy[witchGod][bet] = data
		}
	}
	if ty == internal.PurHadesSuperBuy || ty == internal.PurZuseSuperBuy {
		data = mng.combineDatas_super_buy[witchGod][bet]
		if data == nil {
			if mng.combineDatas_super_buy[witchGod] == nil {
				mng.combineDatas_super_buy[witchGod] = make(map[string]*ppcomm.CombineDataBuy)
			}
			data = &ppcomm.CombineDataBuy{
				RawSeries: mng.SuperBuyBuckets[witchGod],
			}
			mng.combineDatas_super_buy[witchGod][bet] = data
		}
	}
	mng.mtx.Unlock()
	if data == nil {
		return nil, errors.New("not find data")
	}
	id := data.Next()
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBuyMinData(witchGod, ty int) (playResp *ppcomm.SimulateData, err error) {

	var id primitive.ObjectID
	if ty == internal.PurHadesBuy || ty == internal.PurZuseBuy {
		bucketId := GetBuyMinBucketId(witchGod, ppcomm.GameTypeGame)
		data := mng.buyBucketMap[witchGod][bucketId]
		id = lo.Sample(data)
	}
	if ty == internal.PurHadesSuperBuy || ty == internal.PurZuseSuperBuy {
		bucketId := GetBuyMinBucketId(witchGod, ppcomm.GameTypeSuperGame1)
		data := mng.superBuyBucketMap[witchGod][bucketId]
		id = lo.Sample(data)
	}
	if id == primitive.NilObjectID {
		return nil, errors.New("not found")
	}
	return findFromSimulate(id)
}
