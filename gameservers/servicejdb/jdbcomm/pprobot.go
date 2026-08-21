package jdbcomm

//
//import (
//	"context"
//	"fmt"
//	"github.com/samber/lo"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"log/slog"
//	"math/rand/v2"
//	"os"
//	"serve/comm/db"
//	"serve/comm/robot"
//	"serve/comm/ut"
//	"sync"
//	"time"
//)
//
//type PPRobot struct {
//	*robot.Robot
//	*Fetcher
//	respArr  []Variables
//	isFirst  bool
//	RpGameId string
//
//	arrC    map[float64]int
//	insertC float64
//	coll    *mongo.Collection
//	Buckets *Buckets
//	LastRid string
//}
//
//func NewPPRobot(r *robot.Robot, f *Fetcher, ps *FetchDataParam) *PPRobot {
//	if len(ps.ArrC) == 0 {
//		panic("invalid c arr")
//	}
//	if ps.L <= 0 {
//		panic("invalid l")
//	}
//	f.l = ps.L
//	f.ty = ps.Ty
//	ppRobot := &PPRobot{
//		Robot:    r,
//		Fetcher:  f,
//		RpGameId: ps.Game,
//		respArr:  make([]Variables, 0),
//		isFirst:  true,
//		arrC:     ps.ArrC,
//		insertC:  ps.InsertC,
//		coll:     db.Collection2(ps.Game, "simulate"),
//		Buckets:  ps.GBuckets,
//	}
//	ppRobot.RoundC()
//	return ppRobot
//}
//
//func (ppr *PPRobot) Run() {
//	defer func() {
//		if r := recover(); r != nil {
//			robotLock.Lock()
//			delete(robotMap, ppr.ID.Hex())
//			robotLock.Unlock()
//			slog.Error("error:", r)
//		}
//	}()
//	ppr.Fetcher.Do(Variables{
//		"action": "doInit",
//		"cver":   "237859",
//	})
//	ppr.RoundC()
//	for {
//		time.Sleep(500 * time.Millisecond)
//		switch ppr.Fetcher.NextAction() {
//		case "s":
//			ppr.Fetcher.Do(Variables{
//				"action": "doSpin",
//				"bl":     "0",
//			})
//
//		case "c":
//			ppr.Fetcher.Do(Variables{
//				"action": "doCollect",
//			})
//		default:
//			slog.Error("unkown status", "na", ppr.Fetcher.NextAction(), "user", ppr.Robot.ID.Hex(), "lastRespData", ppr.Fetcher.lastResp)
//			panic(fmt.Sprintf("unkown status na:%s ,lastResp:%v", ppr.Fetcher.NextAction(), ppr.Fetcher.lastResp))
//		}
//
//		if ppr.DealData() {
//			if rand.Int()%100 < 10 { //有10%的概率起拉取另外的数据
//				ppr.RoundC()
//			}
//			//if ppr.EndTime <= time.Now().Unix() { // 机器人该结束了
//			//	robotLock.Lock()
//			//	delete(robotMap, ppr.ID.Hex())
//			//	robotLock.Unlock()
//			//	return
//			//}
//			//if robot.CheckGameEnd(ppr.RpGameId) {
//			//	robotLock.Lock()
//			//	delete(robotMap, ppr.ID.Hex())
//			//	robotLock.Unlock()
//			//	return
//			//}
//		}
//	}
//	fmt.Println("over:Run")
//}
//
//func (ppr *PPRobot) RoundC() {
//	allWeight := 0
//	for _, v := range ppr.arrC {
//		allWeight += v
//	}
//	randNum := rand.Int() % allWeight
//	for k, v := range ppr.arrC {
//		randNum -= v
//		if randNum < 0 {
//			ppr.c = k
//			break
//		}
//	}
//}
//
//func (ppr *PPRobot) DealData() bool {
//	//最新的一条数据
//	data := ppr.lastResp
//	//是否要入库
//	isAdd := false
//	//如果这轮数据没有tw 就跳过（推论一应该是验证数据用）
//	if _, ok := data["tw"]; !ok {
//		return false
//	}
//	// 判断是否新的数据是新的一轮，如果是新的一轮就把老数据那堆数组入库
//	if data.Str("rid") != ppr.LastRid && ppr.LastRid != "" {
//		isAdd = true
//	}
//	//上面判断是否是新的一轮了，这里直接修改最新值等待下次新一轮的数据
//	ppr.LastRid = data.Str("rid")
//	//被唤醒的首轮，直接把数据放入缓存，等下一轮do
//	if isAdd && ppr.isFirst {
//		ppr.isFirst = false
//		isAdd = false
//		ppr.respArr = []Variables{data}
//		return false
//	}
//	//确认入库了
//	if isAdd {
//		//投注额一直，并且缓存内有数据
//		if ppr.c == ppr.insertC && len(ppr.respArr) > 0 { // 只有固定的投注金额才会插入数据库
//			// TODO 插入数据
//			lastPan := ppr.respArr[len(ppr.respArr)-1]
//			win := lastPan.Float("tw")
//			objectId := primitive.NewObjectID()
//			hasGame := false
//			rid := ppr.respArr[0].Str("rid")
//			for j := 0; j < len(ppr.respArr); j++ {
//				if ppr.respArr[j].Str("rid") != rid {
//					fmt.Println("拉取数据错误，需要调整代码")
//					os.Exit(-1)
//				}
//				ppr.respArr[j]["gid"] = fmt.Sprintf("%s_%d_%d", objectId.Hex(), j+1, len(ppr.respArr))
//				if ppr.respArr[j].Int("fsmax") > 0 {
//					hasGame = true
//				}
//			}
//			times := ut.Round6(win / ppr.c / float64(ppr.l))
//			doc := SimulateData{
//				Id:       objectId,
//				DropPan:  ppr.respArr,
//				HasGame:  hasGame,
//				Times:    times,
//				BucketId: ppr.Buckets.GetBucket(times, hasGame, BoundType(ppr.ty)),
//				Type:     ppr.ty,
//				Selected: true,
//			}
//
//			lo.Must(ppr.coll.InsertOne(context.TODO(), doc))
//			fmt.Println("success-------->", objectId.Hex(), ppr.ID.Hex(), "type", ppr.ty)
//		}
//		//清空缓存
//		ppr.respArr = make([]Variables, 0)
//		// 如果继续使用，应该再清空一些数据
//	}
//	//加入缓存
//	ppr.respArr = append(ppr.respArr, data)
//	return isAdd
//}
//
//type FetchDataParam struct {
//	Game     string
//	ArrC     map[float64]int
//	GBuckets *Buckets
//	Ty       int
//	IsEnough bool //拉取够了设置为true
//	L        int
//	InsertC  float64
//}
//
//var robotMap map[string]bool
//var robotLock sync.Mutex
//
//func FetchData(ps *FetchDataParam) {
//	if ps == nil {
//		return
//	}
//	robotMap = make(map[string]bool)
//	for ; ; time.Sleep(5 * time.Second) {
//		robots := robot.GetRobot(ps.Game)
//
//		if len(robots) <= 0 {
//			slog.Info("get robot fail, robot busy, GameId:%s, type:%d\n", ps.Game, ps.Ty)
//		}
//		fmt.Println(len(robotMap))
//		for i := range robots {
//			if ps.IsEnough {
//				if rand.Int()%100 < 70 {
//					continue
//				}
//			}
//			robotLock.Lock()
//			if _, ok := robotMap[robots[i].ID.Hex()]; ok {
//				robotLock.Unlock()
//				continue
//			}
//			robotMap[robots[i].ID.Hex()] = true
//			robotLock.Unlock()
//
//			fetcher := NewFetcher(ps.Game, robots[i].ID.Hex())
//			if fetcher == nil {
//				robotLock.Lock()
//				delete(robotMap, robots[i].ID.Hex())
//				robotLock.Unlock()
//				continue
//			}
//			ppR := NewPPRobot(robots[i], fetcher, ps)
//			go ppR.Run()
//		}
//	}
//	slog.Info("over: FetchData")
//}
