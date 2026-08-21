package comm

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math"
	"serve/comm/db"
	"serve/servicepp/ppcomm"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------------由于Simulate表主键可能无序，该方法暂时启用
func ClearData(gname, mongoAddr, tableName string, filterRule bson.D, bsonParam string) {
	db.DialToMongo(mongoAddr, gname)
	var err error
	var wg sync.WaitGroup
	batchSize := int64(500)
	//从mongo中查询 桶子 == -1 且 times == 0的数据 删除掉
	coll := db.Collection(tableName)
	filter := db.D("selected", true, "bucketid", -1, "times", 0)
	many, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		slog.Error("coll.DeleteMany::err: ", "err", err)
		return
	}
	fmt.Println("delete num: " + strconv.Itoa(int(many.DeletedCount)))

	//查询需要处理的数据量 如果大于1w就开启协程
	//filterRule = db.D("selected", true, "bucketid", -1)
	totalCount, err := coll.CountDocuments(context.TODO(), filterRule) //10000
	if err != nil {
		slog.Error("coll.CountDocuments::err: ", "err", err)
		return
	}
	//开启协程
	// 计算总批次数
	totalBatches := totalCount / batchSize
	if totalCount%batchSize != 0 {
		totalBatches++
	}
	fmt.Println(fmt.Sprintf("一共有%v条数据, 需要分%v波完成", totalCount, totalBatches))
	if totalCount == 0 {
		fmt.Println("没有错，不玩了")
		return
	}

	//协程数量
	numGoroutines := 5
	// 创建一个通道来控制并发
	semaphore := make(chan int, numGoroutines)
	resChan := make(chan int, totalBatches)
	delChan := make(chan []primitive.ObjectID, totalBatches)
	defer close(semaphore)
	defer close(resChan)
	defer close(delChan)
	for i := int64(0); i < totalBatches; i++ {
		// 获取当前批次
		start := i * batchSize
		end := start + batchSize
		if end > totalCount {
			end = totalCount
		}
		// 启动协程处理批次
		semaphore <- int(i + 1) // 向管道发送数据
		wg.Add(1)
		fmt.Println(fmt.Sprintf("第%v波已经插入管道", int(i+1)))
		go gHandler(start, filterRule, coll, semaphore, &wg, resChan, int(i+1), totalBatches, bsonParam, delChan)
	}
	wg.Wait()
	cnt := 0
	delIds := make([]primitive.ObjectID, 0)
	for i := 0; i < int(totalBatches); i++ {
		cnt += <-resChan
		delIds = append(delIds, <-delChan...)
	}
	//删除错误数据
	filter = db.D("_id", db.D("$in", delIds))
	many, err = coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		slog.Error("coll.DeleteMany::err: ", "err", err)
		return
	}
	fmt.Println("所有数据处理完成!")
	fmt.Println(fmt.Sprintf("一共处理了%v个! 删除了%v个", cnt, len(delIds)))

}

func gHandler(start int64, filterRule bson.D, coll *mongo.Collection, semaphore chan int, wg *sync.WaitGroup,
	resChan chan int, num int, total int64, bsonParam string, delChan chan []primitive.ObjectID) {
	// 使用信号量控制协程数量
	//if len(semaphore) < cap(semaphore) {
	//	fmt.Println("管道长度是" + strconv.Itoa(len(semaphore)))
	defer wg.Done()
	var docs []ppcomm.SimulateData
	var err error
	findOptions := options.Find().
		SetSkip(start).
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(int64(500))
	cur, _ := coll.Find(context.TODO(), filterRule, findOptions)
	err = cur.All(context.TODO(), &docs)
	if err != nil {
		slog.Error("cur.All(context.TODO(), &docs)::err: ", "err", err)
		return
	}
	if len(docs) == 0 {
		fmt.Println(fmt.Sprintf("咋能是0呢"))
	}
	fmt.Println(fmt.Sprintf("一共%v波, 第%v波开始处理, 本次一共%v条数据", total, num, len(docs)))
	delIds, cnt := handleFG20250214(docs, coll, bsonParam, num)
	//time.Sleep(time.Second * 5)
	resChan <- cnt //放入处理了cnt个
	delChan <- delIds
	<-semaphore
	fmt.Println(fmt.Sprintf("一共%v波, 第%v波处理完成, 本次一共%v条数据, 处理了%v个, 标记了%v个错误数据, 移除管道", total, num, len(docs), cnt, len(delIds)))
	//} else {
	//	fmt.Println("管道已满，无法发送数据" + fmt.Sprintf("我是第%v波，我丢了", num))
	//}
}

func handleFG(docs []ppcomm.SimulateData, coll *mongo.Collection, bsonParam string, numbo int) ([]primitive.ObjectID, int) {
	var err error
	cnt := 0
	delIds := make([]primitive.ObjectID, 0)
	for _, doc := range docs {
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
		if (len(doc.DropPan) > 1 && doc.DropPan[len(doc.DropPan)-2].Get(bsonParam) == "3") && (doc.HasGame == false && doc.BucketId == -1 && doc.Type == 1) { //fg
			doc.HasGame = true
			//if doc.BucketId == -1 {
			doc.BucketId = GBuckets.GetBucket(doc.Times, doc.HasGame, 1)
			//}
			_, err = coll.UpdateByID(context.TODO(), doc.Id, bson.M{"$set": bson.M{"hasgame": doc.HasGame, "bucketid": doc.BucketId}})
			if err != nil {
				slog.Error("coll.UpdateByID::err: ", "err", err)
				continue
			}
			cnt++
			fmt.Println(fmt.Sprintf("这是第%v波, 一共%v个, 完成第%v个修改", numbo, len(docs), cnt))
		} else {
			delIds = append(delIds, doc.Id)
		}
	}
	return delIds, cnt
}

func handleFG20250214(docs []ppcomm.SimulateData, coll *mongo.Collection, bsonParam string, numbo int) ([]primitive.ObjectID, int) {
	var err error
	cnt := 0
	delIds := make([]primitive.ObjectID, 0)
	for _, doc := range docs {
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
		if true { //fg
			doc.HasGame = true
			//if doc.BucketId == -1 {
			doc.BucketId = GBuckets.GetBucket(doc.Times, doc.HasGame, 1)
			//}
			_, err = coll.UpdateByID(context.TODO(), doc.Id, bson.M{"$set": bson.M{"hasgame": doc.HasGame, "bucketid": doc.BucketId}})
			if err != nil {
				slog.Error("coll.UpdateByID::err: ", "err", err)
				continue
			}
			cnt++
			fmt.Println(fmt.Sprintf("这是第%v波, 一共%v个, 完成第%v个修改", numbo, len(docs), cnt))
		} else {
			delIds = append(delIds, doc.Id)
		}
	}
	return delIds, cnt
}

var GBuckets = ppcomm.NewBuckets([]*ppcomm.Bound{
	// 普通模式
	{ // (-oo,0]
		Group: 0,
		Min:   -1,
		Max:   0,
	},

	{ // （0，0.25】
		Group: 1,
		Min:   0,
		Max:   0.25,
	}, { // （0.25，0.5】
		Group: 1,
		Min:   0.25,
		Max:   0.5,
	}, { // （0.5，1】
		Group: 1,
		Min:   0.5,
		Max:   1,
	},

	{ // (1,2]
		Group: 2,
		Min:   1,
		Max:   2,
	}, { // （2，3】
		Group: 2,
		Min:   2,
		Max:   3,
	}, { // （3，4】
		Group: 2,
		Min:   3,
		Max:   4,
	},

	{ // （4，6】
		Group: 3,
		Min:   4,
		Max:   6,
	}, { // （6，8】
		Group: 3,
		Min:   6,
		Max:   8,
	}, { // （8，10】
		Group: 3,
		Min:   8,
		Max:   10,
	},

	{ // （10，15】
		Group: 4,
		Min:   10,
		Max:   15,
	}, { // （15，20】
		Group: 4,
		Min:   15,
		Max:   20,
	}, { //（20，30】需要额外扣除玩家个人奖池5Bet
		Group:    4,
		Min:      20,
		Max:      30,
		PoolCost: 5,
	},

	{ // （30，40】需要额外扣除玩家个人奖池5Bet
		Group:    5,
		Min:      30,
		Max:      40,
		PoolCost: 5,
	}, { // （40，60】需要额外扣除玩家个人奖池6Bet
		Group:    5,
		Min:      40,
		Max:      60,
		PoolCost: 6,
	}, { // （60，80】需要额外扣除玩家个人奖池8Bet
		Group:    5,
		Min:      60,
		Max:      80,
		PoolCost: 8,
	},

	{ // (80, 100] 需要额外扣除玩家个人奖池10bet
		Group:    6,
		Min:      80,
		Max:      100,
		PoolCost: 10,
	}, { // (100, 150] 需要额外扣除玩家个人奖池15bet
		Group:    6,
		Min:      100,
		Max:      150,
		PoolCost: 15,
	}, { // (150, 200] 需要额外扣除玩家个人奖池20bet
		Group:    6,
		Min:      150,
		Max:      200,
		PoolCost: 20,
	}, { // (200, 250] 需要额外扣除玩家个人奖池20bet
		Group:    6,
		Min:      200,
		Max:      250,
		PoolCost: 20,
	},

	{ // (250, 300] 需要额外扣除玩家个人奖池25bet
		Group:    7,
		Min:      250,
		Max:      300,
		PoolCost: 25,
	}, { // (300, 350] 需要额外扣除玩家个人奖池30bet
		Group:    7,
		Min:      300,
		Max:      350,
		PoolCost: 30,
	}, { // (350, 400] 需要额外扣除玩家个人奖池35bet
		Group:    7,
		Min:      350,
		Max:      400,
		PoolCost: 35,
	},

	{ // (400, 450] 需要额外扣除玩家个人奖池40bet
		Group:    8,
		Min:      400,
		Max:      450,
		PoolCost: 40,
	}, { // (450, 500] 需要额外扣除玩家个人奖池45bet
		Group:    8,
		Min:      450,
		Max:      500,
		PoolCost: 45,
	}, { // (500, 550] 需要额外扣除玩家个人奖池50bet
		Group:    8,
		Min:      500,
		Max:      550,
		PoolCost: 50,
	},

	{ // (550, 600] 需要额外扣除玩家个人奖池55bet
		Group:    9,
		Min:      550,
		Max:      600,
		PoolCost: 55,
	}, { // (600, 650] 需要额外扣除玩家个人奖池60bet
		Group:    9,
		Min:      600,
		Max:      650,
		PoolCost: 60,
	}, { // (650, 700] 需要额外扣除玩家个人奖池60bet
		Group:    9,
		Min:      650,
		Max:      700,
		PoolCost: 60,
	},

	{ // (700, 750] 需要额外扣除玩家个人奖池60bet
		Group:    10,
		Min:      700,
		Max:      750,
		PoolCost: 60,
	}, { // (750, 800] 需要额外扣除玩家个人奖池60bet
		Group:    10,
		Min:      750,
		Max:      800,
		PoolCost: 60,
	}, { // (800, 850] 需要额外扣除玩家个人奖池60bet
		Group:    10,
		Min:      800,
		Max:      850,
		PoolCost: 60,
	},

	{ // (850, 900] 需要额外扣除玩家个人奖池60bet
		Group:    11,
		Min:      850,
		Max:      900,
		PoolCost: 60,
	}, { // (900, 950] 需要额外扣除玩家个人奖池60bet
		Group:    11,
		Min:      900,
		Max:      950,
		PoolCost: 60,
	},

	{ // (950, 1000] 需要额外扣除玩家个人奖池60bet
		Group:    12,
		Min:      950,
		Max:      1000,
		PoolCost: 60,
	}, { // (1000, 99999] 需要额外扣除玩家个人奖池60bet
		Group:    12,
		Min:      1000,
		Max:      99999,
		PoolCost: 60,
	},

	//小游戏
	{ // (-1，0】
		Group:   13,
		Min:     -1,
		Max:     0,
		HasGame: true,
	}, { // (0，10】
		Group:   13,
		Min:     0,
		Max:     10,
		HasGame: true,
	},

	{ // （10，20】
		Group:   14,
		Min:     10,
		Max:     20,
		HasGame: true,
	}, { // （20，30】
		Group:   14,
		Min:     20,
		Max:     30,
		HasGame: true,
	}, { // （30，50】需要额外扣除玩家个人奖池8Bet
		Group:    14,
		Min:      30,
		Max:      50,
		HasGame:  true,
		PoolCost: 8,
	},

	{ // (50, 80] 需要额外扣除玩家个人奖池10bet
		Group:    15,
		Min:      50,
		Max:      80,
		HasGame:  true,
		PoolCost: 10,
	}, { // (80, 100]小游戏 需要额外扣除玩家个人奖池10bet
		Group:    15,
		Min:      80,
		Max:      100,
		HasGame:  true,
		PoolCost: 10,
	}, { // (100, 150]小游戏 需要额外扣除玩家个人奖池20bet
		Group:    15,
		Min:      100,
		Max:      150,
		HasGame:  true,
		PoolCost: 20,
	},

	{ // (150, 200]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    16,
		Min:      150,
		Max:      200,
		HasGame:  true,
		PoolCost: 30,
	}, { // (200, 250]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    16,
		Min:      200,
		Max:      250,
		HasGame:  true,
		PoolCost: 30,
	}, { // (250, 300]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    16,
		Min:      250,
		Max:      300,
		HasGame:  true,
		PoolCost: 30,
	},

	{ // (300, 350]小游戏 需要额外扣除玩家个人奖池30bet
		Group:    17,
		Min:      300,
		Max:      350,
		HasGame:  true,
		PoolCost: 30,
	}, { // (350, 400]小游戏 需要额外扣除玩家个人奖池35bet
		Group:    17,
		Min:      350,
		Max:      400,
		HasGame:  true,
		PoolCost: 35,
	}, { // (400, 450]小游戏 需要额外扣除玩家个人奖池40bet
		Group:    17,
		Min:      400,
		Max:      450,
		HasGame:  true,
		PoolCost: 40,
	},

	{ // (450, 500]小游戏 需要额外扣除玩家个人奖池45bet
		Group:    18,
		Min:      450,
		Max:      500,
		HasGame:  true,
		PoolCost: 45,
	}, { // (500, 550]小游戏 需要额外扣除玩家个人奖池50bet
		Group:    18,
		Min:      500,
		Max:      550,
		HasGame:  true,
		PoolCost: 50,
	}, { // (550, 600]小游戏 需要额外扣除玩家个人奖池55bet
		Group:    18,
		Min:      550,
		Max:      600,
		HasGame:  true,
		PoolCost: 55,
	},

	{ // (600, 650]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    19,
		Min:      600,
		Max:      650,
		HasGame:  true,
		PoolCost: 60,
	}, { // (650, 700]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    19,
		Min:      650,
		Max:      700,
		HasGame:  true,
		PoolCost: 60,
	}, { // (700, 750]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    19,
		Min:      700,
		Max:      750,
		HasGame:  true,
		PoolCost: 60,
	},

	{ // (750, 800]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    20,
		Min:      750,
		Max:      800,
		HasGame:  true,
		PoolCost: 60,
	}, { // (800, 850]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    20,
		Min:      800,
		Max:      850,
		HasGame:  true,
		PoolCost: 60,
	}, { // (850, 900]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    20,
		Min:      850,
		Max:      900,
		HasGame:  true,
		PoolCost: 60,
	},

	{ // (900, 950]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    21,
		Min:      900,
		Max:      950,
		HasGame:  true,
		PoolCost: 60,
	}, { // (950, 1000]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    21,
		Min:      950,
		Max:      1000,
		HasGame:  true,
		PoolCost: 60,
	},

	{ // (1000, 99999]小游戏 需要额外扣除玩家个人奖池60bet
		Group:    22,
		Min:      1000,
		Max:      99999,
		HasGame:  true,
		PoolCost: 60,
	},

	//购买小游戏
	{ // (0，10】
		Group:   23,
		Min:     -1,
		Max:     0,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (0, 10]小游戏购买
		Group:   24,
		Min:     0,
		Max:     10,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (10, 20]小游戏购买
		Group:   24,
		Min:     10,
		Max:     20,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (20, 30]小游戏购买
		Group:   24,
		Min:     20,
		Max:     30,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (30, 50]小游戏购买
		Group:   25,
		Min:     30,
		Max:     50,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (50, 80]小游戏购买
		Group:   25,
		Min:     50,
		Max:     80,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (80, 100]小游戏购买
		Group:   25,
		Min:     80,
		Max:     100,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (100, 150]小游戏购买
		Group:   26,
		Min:     100,
		Max:     150,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (150, 200]小游戏购买
		Group:   26,
		Min:     150,
		Max:     200,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (200, 250]小游戏购买
		Group:   26,
		Min:     200,
		Max:     250,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (250, 300]小游戏购买
		Group:   27,
		Min:     250,
		Max:     300,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (300, 350]小游戏购买
		Group:   27,
		Min:     300,
		Max:     350,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (350, 500]小游戏购买
		Group:   27,
		Min:     350,
		Max:     500,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (500, 550]小游戏购买
		Group:   28,
		Min:     500,
		Max:     550,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (550, 600]小游戏购买
		Group:   28,
		Min:     550,
		Max:     600,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (600, 650]小游戏购买
		Group:   28,
		Min:     600,
		Max:     650,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (650, 700]小游戏购买
		Group:   29,
		Min:     650,
		Max:     700,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (700, 750]小游戏购买
		Group:   29,
		Min:     700,
		Max:     750,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (750, 800]小游戏购买
		Group:   29,
		Min:     750,
		Max:     800,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (800, 850]小游戏购买
		Group:   30,
		Min:     800,
		Max:     850,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (850, 900]小游戏购买
		Group:   30,
		Min:     850,
		Max:     900,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (900, 950]小游戏购买
		Group:   30,
		Min:     900,
		Max:     950,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},

	{ // (950, 1000]小游戏购买
		Group:   31,
		Min:     950,
		Max:     1000,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	}, { // (1000, 99999]小游戏购买
		Group:   31,
		Min:     1000,
		Max:     99999,
		HasGame: true,
		Type:    ppcomm.GameTypeGame,
	},
})

// -------------------------直接取所有桶子 == -1的数据到内存处理
func ClearData2(gname, mongoAddr, tableName string, filterRule bson.D, bsonParam string) {
	db.DialToMongo(mongoAddr, gname)
	var err error
	var wg sync.WaitGroup
	var docs []ppcomm.SimulateData
	batchSize := 500
	//从mongo中查询 桶子 == -1 且 times == 0的数据 删除掉
	coll := db.Collection(tableName)
	filter := db.D("selected", true, "bucketid", -1, "times", 0)
	many, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		slog.Error("coll.DeleteMany::err: ", "err", err)
		return
	}
	fmt.Println("delete num: " + strconv.Itoa(int(many.DeletedCount)))

	//查询需要处理的数据量 如果大于1w就开启协程
	//filterRule = db.D("selected", true, "bucketid", -1)
	cur, _ := coll.Find(context.TODO(), filterRule)
	err = cur.All(context.TODO(), &docs)
	if err != nil {
		slog.Error("cur.All(context.TODO(), &docs)::err: ", "err", err)
		return
	}
	//开启协程
	// 计算总批次数
	totalCount := len(docs)
	totalBatches := totalCount / batchSize
	if totalCount%batchSize != 0 {
		totalBatches++
	}
	fmt.Println(fmt.Sprintf("一共有%v条数据, 需要分%v波完成", totalCount, totalBatches))
	if totalCount == 0 {
		fmt.Println("没有错，不玩了")
		return
	}

	//协程数量
	numGoroutines := 5
	// 创建一个通道来控制并发
	semaphore := make(chan int, numGoroutines)
	resChan := make(chan int, totalBatches)
	delChan := make(chan []primitive.ObjectID, totalBatches)
	defer close(semaphore)
	defer close(resChan)
	defer close(delChan)
	for i := 0; i < totalBatches; i++ {
		// 获取当前批次
		start := i * batchSize
		end := start + batchSize
		if end > totalCount {
			end = totalCount
		}
		batch := docs[start:end]

		// 启动协程处理批次
		semaphore <- int(i + 1) // 向管道发送数据
		wg.Add(1)
		fmt.Println(fmt.Sprintf("第%v波已经插入管道", int(i+1)))
		go gHandler2(batch, coll, semaphore, &wg, resChan, int(i+1), totalBatches, bsonParam, delChan)
	}
	wg.Wait()
	cnt := 0
	delIds := make([]primitive.ObjectID, 0)
	for i := 0; i < int(totalBatches); i++ {
		cnt += <-resChan
		delIds = append(delIds, <-delChan...)
	}
	//删除错误数据
	filter = db.D("_id", db.D("$in", delIds))
	many, err = coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		slog.Error("coll.DeleteMany::err: ", "err", err)
		return
	}
	fmt.Println("所有数据处理完成!")
	fmt.Println(fmt.Sprintf("一共处理了%v个! 删除了%v个", cnt, len(delIds)))

}

func gHandler2(docs []ppcomm.SimulateData, coll *mongo.Collection, semaphore chan int, wg *sync.WaitGroup,
	resChan chan int, num int, total int, bsonParam string, delChan chan []primitive.ObjectID) {
	// 使用信号量控制协程数量
	//if len(semaphore) < cap(semaphore) {
	//	fmt.Println("管道长度是" + strconv.Itoa(len(semaphore)))
	defer wg.Done()
	if len(docs) == 0 {
		fmt.Println(fmt.Sprintf("咋能是0呢"))
	}
	fmt.Println(fmt.Sprintf("一共%v波, 第%v波开始处理, 本次一共%v条数据", total, num, len(docs)))
	delIds, cnt := handleFG20250214(docs, coll, bsonParam, num)
	//time.Sleep(time.Second * 5)
	resChan <- cnt //放入处理了cnt个
	delChan <- delIds
	<-semaphore
	fmt.Println(fmt.Sprintf("一共%v波, 第%v波处理完成, 本次一共%v条数据, 处理了%v个, 标记了%v个错误数据, 移除管道", total, num, len(docs), cnt, len(delIds)))
	//} else {
	//	fmt.Println("管道已满，无法发送数据" + fmt.Sprintf("我是第%v波，我丢了", num))
	//}
}

func handleFG2(docs []ppcomm.SimulateData, coll *mongo.Collection, bsonParam string, numbo int) ([]primitive.ObjectID, int) {
	var err error
	cnt := 0
	delIds := make([]primitive.ObjectID, 0)
	for _, doc := range docs {
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
		if (len(doc.DropPan) > 1 && doc.DropPan[len(doc.DropPan)-2].Get(bsonParam) == "3") && (doc.HasGame == false && doc.BucketId == -1 && doc.Type == 1) { //fg
			doc.HasGame = true
			//if doc.BucketId == -1 {
			doc.BucketId = GBuckets.GetBucket(doc.Times, doc.HasGame, 1)
			//}
			_, err = coll.UpdateByID(context.TODO(), doc.Id, bson.M{"$set": bson.M{"hasgame": doc.HasGame, "bucketid": doc.BucketId}})
			if err != nil {
				slog.Error("coll.UpdateByID::err: ", "err", err)
				continue
			}
			cnt++
			fmt.Println(fmt.Sprintf("这是第%v波, 一共%v个, 完成第%v个修改", numbo, len(docs), cnt))
		} else {
			delIds = append(delIds, doc.Id)
		}
	}
	return delIds, cnt
}

//------------------------------------夺的方法----------------------------------------

// 刷times,去除tw中的,
func test1() {
	uri := "mongodb://myAdmin:myAdminPassword1@18.61.185.51:27017/?authSource=admin"
	// 创建一个新的客户端
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	// 尝试连接到MongoDB服务器
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	var dbtable = []string{
		"pp_vs10bxmasbnza",
		"pp_vs20chickdrop",
		"pp_vs20emptybank",
		"pp_vs20fparty2",
		"pp_vs20starlightx",
		"pp_vswaysmfreya",
	}
	for _, v := range dbtable {
		fmt.Println(v)
		coll := client.Database(v).Collection("simulate")

		projection := bson.D{
			{"_id", 1},
			{"times", 1},
			{"droppan", 1},
		}
		opts := options.Find().SetProjection(projection)

		filter := bson.M{
			"times": 0,
		}
		cursor, err := coll.Find(context.TODO(), filter, opts)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var result struct {
				Id      string              `bson:"_id,omitempty"`
				Times   float64             `bson:"times"`
				Droppan []map[string]string `bson:"droppan"`
			}
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			data := result.Droppan[len(result.Droppan)-1]
			tw, _ := strconv.ParseFloat(strings.ReplaceAll(data["tw"], ",", ""), 64)
			l, _ := strconv.ParseFloat(data["l"], 64)
			c, _ := strconv.ParseFloat(data["c"], 64)
			times := tw / c / l
			times = math.Floor(times*100+0.5) / 100

			if result.Times != times {
				fmt.Println(result.Times, times, v, result.Id)

				objectID, err := primitive.ObjectIDFromHex(result.Id)
				if err != nil {
					log.Fatal(err)
				}
				updateFilter := bson.M{"_id": objectID}
				update := bson.M{"$set": bson.M{"times": times}}
				_, err = coll.UpdateOne(context.TODO(), updateFilter, update)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		}
		test3(v, coll)
	}
}

// 根据times分桶
func test3(dbName string, coll *mongo.Collection) {
	//coll := db.Collection2(dbName, "simulate")
	for k, v := range GBuckets.Bounds {
		fmt.Println(k, v)
		updateFilter := bson.M{
			"$and": []bson.M{
				{"type": v.Type},
				{"selected": true},
				{"hasgame": v.HasGame},
				{"times": bson.M{"$gt": v.Min}},
				{"times": bson.M{"$lte": v.Max}},
			},
		}
		updateData := bson.M{"$set": bson.M{"bucketid": k}}
		_, err := coll.UpdateMany(context.TODO(), updateFilter, updateData)
		fmt.Println(err)
	}

	return
}
