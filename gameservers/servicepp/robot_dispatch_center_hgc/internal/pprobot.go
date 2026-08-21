package internal

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"serve/comm/db"
	"serve/comm/robot"
	"serve/comm/ut"
	"serve/servicepp/ppcomm"
	"strconv"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PPRobot struct {
	*robot.Robot
	*Fetcher
	respArr  []ppcomm.Variables
	isFirst  bool
	RpGameId string

	arrC       map[float64]int
	insertC    float64
	coll       *mongo.Collection
	Buckets    *ppcomm.Buckets
	LastRid    string //官方试玩 使用记录 na /s /c
	LastStatus int    //官方试玩 使用记录 na /s /c
}

func GetMissions() []GameInfo {
	coll := db.Collection2("robot", "game")
	cur, err := coll.Find(context.Background(), bson.M{"status": 1})
	if err != nil {
		slog.Error(err.Error())
	}

	var robotsToUpdate []GameInfo
	for cur.Next(context.TODO()) {
		var robot GameInfo
		err := cur.Decode(&robot)
		if err != nil {
			slog.Error(err.Error())
		}
		robotsToUpdate = append(robotsToUpdate, robot)
	}
	return robotsToUpdate
}
func GetMissionsHGC() []GameInfo {
	coll := db.Collection2("robot_hgc", "game")
	cur, err := coll.Find(context.Background(), bson.M{"status": 1})
	if err != nil {
		slog.Error(err.Error())
	}

	var robotsToUpdate []GameInfo
	for cur.Next(context.TODO()) {
		var robot GameInfo
		err := cur.Decode(&robot)
		if err != nil {
			slog.Error(err.Error())
		}
		robotsToUpdate = append(robotsToUpdate, robot)
	}
	return robotsToUpdate
}

func dispatchRobots(gameId string, weights int64) error {
	// 1. 获取 game="aaa" 的机器人数量
	robotCur := db.Collection2("robot_hgc", "robot")
	count, err := robotCur.CountDocuments(context.TODO(), bson.M{"gameId": gameId})
	if err != nil {
		return err
	}

	fmt.Printf("gameId=\"%v\" 的机器人数量: %d\n", gameId, count)

	if count < weights {
		// 2. 计算缺少的数量
		missingCount := weights - count
		fmt.Printf("缺少的数量: %d\n", missingCount)

		opts := options.Find().SetLimit(missingCount) // 设置查询数量限制
		cur, err := robotCur.Find(context.TODO(), bson.M{"gameId": ""}, opts)
		if err != nil {
			return err
		}
		defer cur.Close(context.TODO())
		var robotsToUpdate []robot.Robot
		for cur.Next(context.TODO()) {
			var robot robot.Robot
			if err := cur.Decode(&robot); err != nil {
				return err
			}
			robotsToUpdate = append(robotsToUpdate, robot)
		}
		if err := cur.Err(); err != nil {
			return err
		}

		if len(robotsToUpdate) == 0 {
			fmt.Println("没有找到需要更新的机器人")
		} else {
			fmt.Println("将要更新的机器人:")
			for _, robot := range robotsToUpdate {
				fmt.Printf("ID: %v\n", robot.ID)
			}
			endTime := time.Now().Add(36 * time.Hour).Unix()
			// 4. 更新机器人game字段
			for _, robot := range robotsToUpdate {
				_, err := robotCur.UpdateOne(context.TODO(), bson.M{"_id": robot.ID}, bson.M{"$set": bson.M{"gameId": gameId, "endTime": endTime}})
				if err != nil {
					log.Printf("更新机器人失败: %v", err)
				} else {
					fmt.Printf("更新机器人成功: ID: %v\n", robot.ID)
				}
			}
		}
	} else {
		fmt.Println("不需要更新机器人")
	}
	return err
}
func dispatchRobotsHGC(gameId string, weights int64) error {
	// 1. 获取 game="aaa" 的机器人数量
	robotCur := db.Collection2("robot_hgc", "robot")
	count, err := robotCur.CountDocuments(context.TODO(), bson.M{"gameId": gameId})
	if err != nil {
		return err
	}

	fmt.Printf("gameId=\"%v\" 的机器人数量: %d\n", gameId, count)

	if count < weights {
		// 2. 计算缺少的数量
		missingCount := weights - count
		fmt.Printf("缺少的数量: %d\n", missingCount)

		opts := options.Find().SetLimit(missingCount) // 设置查询数量限制
		cur, err := robotCur.Find(context.TODO(), bson.M{"gameId": ""}, opts)
		if err != nil {
			return err
		}
		defer cur.Close(context.TODO())
		var robotsToUpdate []robot.Robot
		for cur.Next(context.TODO()) {
			var robot robot.Robot
			if err := cur.Decode(&robot); err != nil {
				return err
			}
			robotsToUpdate = append(robotsToUpdate, robot)
		}
		if err := cur.Err(); err != nil {
			return err
		}

		if len(robotsToUpdate) == 0 {
			fmt.Println("没有找到需要更新的机器人")
		} else {
			fmt.Println("将要更新的机器人:")
			for _, robot := range robotsToUpdate {
				fmt.Printf("ID: %v\n", robot.ID)
			}
			endTime := time.Now().Add(36 * time.Hour).Unix()
			// 4. 更新机器人game字段
			for _, robot := range robotsToUpdate {
				_, err := robotCur.UpdateOne(context.TODO(), bson.M{"_id": robot.ID}, bson.M{"$set": bson.M{"gameId": gameId, "endTime": endTime}})
				if err != nil {
					log.Printf("更新机器人失败: %v", err)
				} else {
					fmt.Printf("更新机器人成功: ID: %v\n", robot.ID)
				}
			}
		}
	} else {
		fmt.Println("不需要更新机器人")
	}
	return err
}
func NewPPRobot(r *robot.Robot, f *Fetcher, ps *FetchDataParam) *PPRobot {
	if len(ps.ArrC) == 0 {
		panic("invalid c arr")
	}
	if ps.L <= 0 {
		panic("invalid l")
	}
	f.l = ps.L
	f.ty = ps.Ty
	ppRobot := &PPRobot{
		Robot:    r,
		Fetcher:  f,
		RpGameId: ps.Game,
		respArr:  make([]ppcomm.Variables, 0),
		isFirst:  true,
		arrC:     ps.ArrC,
		insertC:  ps.InsertC,
		coll:     db.Collection2(ps.Game, "simulate"),
		Buckets:  ps.GBuckets,
	}
	ppRobot.RoundC()
	return ppRobot
}

func (ppr *PPRobot) Run() {
	defer func() {
		if r := recover(); r != nil {
			robotLock.Lock()
			delete(robotMap, ppr.ID.Hex())
			robotLock.Unlock()
			slog.Error("error:", r)
		}
	}()
	ppr.Fetcher.Do(ppcomm.Variables{
		"action": "doInit",
		"cver":   "271865",
		//"cver": "264484",
		//"cver": "265585",
	})
	ppr.RoundC()
	for {
		time.Sleep(500 * time.Millisecond)
		switch ppr.Fetcher.NextAction() {
		case "s":
			ppr.Fetcher.Do(ppcomm.Variables{
				"action": "doSpin",
				"bl":     "0",
			})
		case "c":
			ppr.Fetcher.Do(ppcomm.Variables{
				"action": "doCollect",
			})
		case "b":
			ppr.Fetcher.Do(ppcomm.Variables{
				"action": "doBonus",
			})
		case "cb": //vs243caishien
			ppr.Fetcher.Do(ppcomm.Variables{
				"action": "doCollectBonus",
			})
		default:
			//slog.Error("unkown status", "na", ppr.Fetcher.NextAction(), "user", ppr.Robot.ID.Hex(), "lastRespData", ppr.Fetcher.lastResp)
			panic(fmt.Sprintf("unkown status na:%s ,game:%s ,lastResp:%v", ppr.Fetcher.NextAction(), ppr.game, ppr.Fetcher.lastResp))
		}

		if ppr.DealData() {
			if rand.Int()%100 < 10 { //有10%的概率起拉取另外的数据
				ppr.RoundC()
			}
			//if ppr.EndTime <= time.Now().Unix() { // 机器人该结束了
			//	robotLock.Lock()
			//	delete(robotMap, ppr.ID.Hex())
			//	robotLock.Unlock()
			//	return
			//}
			//数据够了结束
			if robot.CheckGameEndHGC(ppr.RpGameId) || CheckGameEnoughHGC(ppr.RpGameId, ppr.ty) {
				robotLock.Lock()
				delete(robotMap, ppr.ID.Hex())
				robotLock.Unlock()
				break
			}
		}
	}
	fmt.Println("over:Run")
}

// RoundC 随机取权重，下注
func (ppr *PPRobot) RoundC() {
	allWeight := 0
	for _, v := range ppr.arrC {
		allWeight += v
	}
	randNum := rand.Int() % allWeight
	for k, v := range ppr.arrC {
		randNum -= v
		if randNum < 0 {
			ppr.c = k
			break
		}
	}
}

func (ppr *PPRobot) DealData() bool {
	//最新的一条数据
	data := ppr.lastResp
	//是否要入库
	isAdd := false
	//如果这轮数据没有tw 就跳过（推论一应该是验证数据用）
	if _, ok := data["tw"]; !ok {
		return false
	}
	// 判断是否新的数据是新的一轮，如果是新的一轮就把老数据那堆数组入库
	if data.Str("rid") != ppr.LastRid && ppr.LastRid != "" {
		isAdd = true
	}
	//上面判断是否是新的一轮了，这里直接修改最新值等待下次新一轮的数据
	ppr.LastRid = data.Str("rid")
	//被唤醒的首轮，直接把数据放入缓存，等下一轮do
	if isAdd && ppr.isFirst {
		ppr.isFirst = false
		isAdd = false
		ppr.respArr = []ppcomm.Variables{data}
		return false
	}
	//确认入库了
	if isAdd {
		//投注额一直，并且缓存内有数据
		//if ppr.c == ppr.insertC && len(ppr.respArr) > 0 { // 只有固定的投注金额才会插入数据库
		if len(ppr.respArr) > 0 { // 只有固定的投注金额才会插入数据库
			// TODO 插入数据
			lastPan := ppr.respArr[len(ppr.respArr)-1]
			win := lastPan.Float("tw")
			objectId := primitive.NewObjectID()
			hasGame := false
			//新试玩站需要自己制作rid
			sf := ut.NewSnowflake()
			newRid := strconv.FormatInt(sf.NextID(), 10)
			rid := ppr.respArr[0].Str("rid")
			for j := 0; j < len(ppr.respArr); j++ {
				if ppr.respArr[j].Str("rid") != rid {
					fmt.Println("")
					panic(fmt.Sprintf("拉取数据错误，需要调整代码rid:%s,newRid:%s ,game:%s ,respArr:%v", rid, newRid, ppr.game, ppr.respArr))
				}
				ppr.respArr[j]["gid"] = fmt.Sprintf("%s_%d_%d", objectId.Hex(), j+1, len(ppr.respArr))
				ppr.respArr[j]["rid"] = newRid
				if ppr.respArr[j].Int("fsmax") > 0 {
					hasGame = true
				}
			}
			times := ut.Round6(win / ppr.c / float64(ppr.l))
			doc := ppcomm.SimulateData{
				Id:              objectId,
				DropPan:         ppr.respArr,
				HasGame:         hasGame,
				Times:           times,
				BucketId:        ppr.Buckets.GetBucket(times, hasGame, db.BoundType(ppr.ty)),
				Type:            ppr.ty,
				Selected:        true,
				BucketMix:       1,
				BucketStable:    1,
				BucketHeartBeat: 1,
			}

			lo.Must(ppr.coll.InsertOne(context.TODO(), doc))
			fmt.Println("success-------->", ppr.game, objectId.Hex(), ppr.ID.Hex(), "type", ppr.ty)
		}
		//清空缓存
		ppr.respArr = make([]ppcomm.Variables, 0)
		// 如果继续使用，应该再清空一些数据
	}
	//加入缓存
	ppr.respArr = append(ppr.respArr, data)
	return isAdd
	////最新的一条数据
	//data := ppr.lastResp
	////是否要入库
	//isAdd := false
	////fmt.Println(ppr.lastResp.GetGameSt())
	//// 判断是否新的数据是新的一轮，如果是新的一轮就把老数据那堆数组入库（s or c）
	//if ppr.LastRid == "c" || ppr.lastResp.GetGameSt() == 1 && ppr.LastRid != "c" && ppr.LastStatus != 0 && ppr.LastRid != "" || ppr.lastResp.GetGameSt() == 2 && ppr.LastRid == data.Str("na") && ppr.LastStatus == 1 && ppr.LastRid != "" {
	//	isAdd = true
	//}
	////上面判断是否是新的一轮了，这里直接修改最新值等待下次新一轮的数据
	//ppr.LastRid = data.Str("na")
	//ppr.LastStatus = data.GetGameSt()
	////如果这轮数据没有tw 就跳过（推论一应该是验证数据用） doCollect 跳过
	////if _, ok := data["tw"]; !ok {
	////	return false
	////}
	////被唤醒的首轮，直接把数据放入缓存，等下一轮do
	////if isAdd && ppr.isFirst {
	////	ppr.isFirst = false
	////	isAdd = false
	////	ppr.respArr = []ppcomm.Variables{data}
	////	return false
	////}
	////确认入库了
	//if isAdd {
	//	//投注额一直，并且缓存内有数据
	//	if ppr.c == ppr.insertC && len(ppr.respArr) > 0 { // 只有固定的投注金额才会插入数据库
	//		//if len(ppr.respArr) > 0 { // 只有固定的投注金额才会插入数据库
	//		// TODO 插入数据
	//		//加入缓存新试玩站，本结束时本局应该是c doCollect 所以要加入队列
	//		//ppr.respArr = append(ppr.respArr, data)
	//
	//		lastPan := ppr.respArr[len(ppr.respArr)-1]
	//		win := lastPan.Float("tw")
	//		objectId := primitive.NewObjectID()
	//		hasGame := false
	//		//新试玩站需要自己制作rid
	//		sf := ut.NewSnowflake()
	//		rid := strconv.FormatInt(sf.NextID(), 10)
	//		for j := 0; j < len(ppr.respArr); j++ {
	//			//无法验证连续性 试玩站没round id
	//			//if ppr.respArr[j].Str("rid") != rid {
	//			//	fmt.Println("拉取数据错误，需要调整代码")
	//			//	os.Exit(-1)
	//			//}
	//			ppr.respArr[j]["rid"] = rid
	//			ppr.respArr[j]["gid"] = fmt.Sprintf("%s_%d_%d", objectId.Hex(), j+1, len(ppr.respArr))
	//			if ppr.respArr[j].Int("fsmax") > 0 {
	//				hasGame = true
	//			}
	//		}
	//		times := ut.Round6(win / ppr.c / float64(ppr.l))
	//		doc := ppcomm.SimulateData{
	//			Id:       objectId,
	//			DropPan:  ppr.respArr,
	//			HasGame:  hasGame,
	//			Times:    times,
	//			BucketId: ppr.Buckets.GetBucket(times, hasGame, db.BoundType(ppr.ty)),
	//			Type:     ppr.ty,
	//			Selected: true,
	//		}
	//
	//		lo.Must(ppr.coll.InsertOne(context.TODO(), doc))
	//		fmt.Println("success-------->", objectId.Hex(), ppr.ID.Hex(), "type", ppr.ty)
	//	}
	//	//清空缓存
	//	ppr.respArr = make([]ppcomm.Variables, 0)
	//	// 如果继续使用，应该再清空一些数据
	//	if _, ok := data["tw"]; !ok {
	//		return false
	//	}
	//}
	////加入缓存
	//ppr.respArr = append(ppr.respArr, data)
	//return isAdd
}
