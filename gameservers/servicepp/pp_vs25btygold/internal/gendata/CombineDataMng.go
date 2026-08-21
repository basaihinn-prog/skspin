package gendata

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/servicepp/ppcomm"
	"strconv"
	"sync"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	buckets      map[int][][]primitive.ObjectID
	buyBucketMap [][]primitive.ObjectID
	buyBuckets   []primitive.ObjectID
	combine      ppcomm.Combine

	combineDatas     map[string]*ppcomm.CombineData
	combineDatas_buy map[string]*ppcomm.CombineDataBuy

	mtx sync.Mutex
}

var GCombineDataMng *CombineDataMng
var ErrBucketCount = 0
var ErrTypeCount = 0

func (mng *CombineDataMng) GetBigReward2_5(witchBucket int) (playResp *ppcomm.SimulateData, err error) {
	ids := []int{}

	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(2, 5, false, ppcomm.GameTypeNormal)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
	}
	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(1, 6, false, ppcomm.GameTypeNormal)...)
		ids = lo.Filter(ids, func(item int, index int) bool {
			return len(mng.buckets[witchBucket][item]) != 0
		})
	}
	if len(ids) == 0 {
		ids = append(ids, GetBucketIds(0, 10, false, ppcomm.GameTypeNormal)...)
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
		normal [][]primitive.ObjectID
	}

	bucketConfig := map[int]*bucketGroups{
		ppcomm.BucketHeartBeat:      {makeSlice()},
		ppcomm.BucketWave:           {makeSlice()},
		ppcomm.BucketGov:            {makeSlice()},
		ppcomm.BucketMix:            {makeSlice()},
		ppcomm.BucketStable:         {makeSlice()},
		ppcomm.BucketHighAward:      {makeSlice()},
		ppcomm.BucketSuperHighAward: {makeSlice()},
	}
	var (
		buyBucketMap      = make([][]primitive.ObjectID, len(GBuckets.Bounds))
		superBuyBucketMap = make([][]primitive.ObjectID, len(GBuckets.Bounds))
		buyIds            []primitive.ObjectID
		superBuyIds       []primitive.ObjectID
	)
	for _, doc := range docs {
		bound := GBuckets.GetBound(doc.BucketID)
		if bound.Type != doc.Type {
			return fmt.Errorf("type mismatch for doc %v", doc.ID)
		}
		if doc.Type == ppcomm.GameTypeNormal {
			handleBucket := func(bucketType int, config *bucketGroups, isOff bool) {
				if isOff {
					if doc.Type == ppcomm.GameTypeNormal {
						config.normal[doc.BucketID] = append(config.normal[doc.BucketID], doc.ID)
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
		if doc.Type == ppcomm.GameTypeSuperGame1 {
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
		combineDatas_buy: map[string]*ppcomm.CombineDataBuy{},

		combine:      combine,
		combineDatas: map[string]*ppcomm.CombineData{},
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

func (mng *CombineDataMng) get(bet string) (data *ppcomm.CombineData) {
	mng.mtx.Lock()
	defer mng.mtx.Unlock()
	data = mng.combineDatas[bet]
	if data == nil {
		data = &ppcomm.CombineData{
			Combine: mng.combine,
			Buckets: mng.buckets,
		}
		mng.combineDatas[bet] = data
	}
	return
}

// 控制下一局
func (mng *CombineDataMng) ControlNextDataNormal(nextData *redisx.NextMulti, witchBucket int) (playResp *ppcomm.SimulateData, err error) {
	gameType := db.BoundType(ppcomm.GameTypeNormal)
	ids := []int{}
	if nextData != nil {
		filterBucket := func(min, max float64, gameType db.BoundType) (ids []int) {
			for _, b := range GBuckets.Bounds {
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

func (mng *CombineDataMng) Next(bet string, WitchBucket int) (playResp *ppcomm.SimulateData, err error) {
	data := mng.get(bet)

	id := data.Next(WitchBucket)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) SampleSimulate(bucketid, witchBucket int) (playResp *ppcomm.SimulateData, err error) {
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

func (mng *CombineDataMng) NextBuy(bet string) (playResp *ppcomm.SimulateData, err error) {
	mng.mtx.Lock()
	data := mng.combineDatas_buy[bet]
	if data == nil {
		data = &ppcomm.CombineDataBuy{
			RawSeries: mng.buyBuckets,
		}
		mng.combineDatas_buy[bet] = data
	}
	mng.mtx.Unlock()

	id := data.Next()
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBigReward(witchBucket int) (playResp *ppcomm.SimulateData, err error) {
	ids := []int{}
	ids = append(ids, GetBucketIds(30, 50, true, ppcomm.GameTypeNormal)...)
	ids = append(ids, GetBucketIds(30, 40, false, ppcomm.GameTypeNormal)...)
	ids = lo.Filter(ids, func(item int, index int) bool {
		return len(mng.buckets[witchBucket][item]) != 0
	})
	lo.Must0(len(ids) != 0)
	bucketid := lo.Sample(ids)
	data := mng.buckets[witchBucket][bucketid]
	id := lo.Sample(data)
	return findFromSimulate(id)
}

func (mng *CombineDataMng) GetBuyMinData() (playResp *ppcomm.SimulateData, err error) {
	bucketId := GetBuyMinBucketId()
	data := mng.buyBucketMap[bucketId]
	id := lo.Sample(data)
	if id == primitive.NilObjectID {
		return nil, errors.New("not found")
	}
	return findFromSimulate(id)
}

func handler() {
	db.DialToMongo("mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin", "pp_vs25btygold")
	var err error
	var docs []ppcomm.SimulateData
	//从mongo中查询桶子-1且times是0的数据
	coll := db.Collection("simulate")
	filter := db.D("selected", true, "bucketid", -1, "times", 0)
	many, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		slog.Error("coll.DeleteMany::err: ", "err", err)
		return
	}
	fmt.Println("delete num: " + strconv.Itoa(int(many.DeletedCount)))
	//从mongo中查询桶子-1的数据
	filter = db.D("selected", true)
	cur, _ := coll.Find(context.TODO(), filter, options.Find())
	err = cur.All(context.TODO(), &docs)
	if err != nil {
		slog.Error("cur.All(context.TODO(), &docs)::err: ", "err", err)
		return
	}
	//将查询出的数据hasgame设置为true，然后使用getBucketId方法赋值bucketid
	num := 0
	cnt := 0
	for _, doc := range docs {
		num++
		//doc := ppcomm.SimulateData{
		//	Id:              objectId,
		//	DropPan:         ppr.respArr,
		//	HasGame:         hasGame,
		//	Times:           times,
		//	BucketId:        ppr.Buckets.GetBucket(times, hasGame, db.BoundType(ppr.ty)),
		//	Type:            ppr.ty,
		//	Selected:        true,
		//	BucketMix:       1,
		//	BucketStable:    1,
		//	BucketHeartBeat: 1,
		//}
		if len(doc.DropPan) > 1 && doc.DropPan[len(doc.DropPan)-2].Get("rs_m") == "3" { //fg
			doc.HasGame = true
			if doc.BucketId == -1 {
				doc.BucketId = GBuckets.GetBucket(doc.Times, doc.HasGame, 1)
			}
			_, err = coll.UpdateByID(context.TODO(), doc.Id, bson.M{"$set": bson.M{"hasgame": doc.HasGame, "bucketid": doc.BucketId}})
			if err != nil {
				slog.Error("coll.UpdateByID::err: ", "err", err)
				continue
			}
			fmt.Println("完成一个修改")
			cnt++
		}
	}
	fmt.Println("----------------------------我滴任务结束辣----------------------------")
	fmt.Println("一共跑了" + strconv.Itoa(num) + "条数据")
	fmt.Println("修了" + strconv.Itoa(cnt) + "条数据")
}
