package internal

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"serve/comm/db"
	"serve/servicepp/ppcomm"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MaxConcurrencyTaskNum = 2

var gameMap map[string]bool
var gameLock sync.Mutex

func RobotsDispatch() {
	gameMap = make(map[string]bool)
	for ; ; time.Sleep(5 * time.Minute) {
		for _, mission := range GetMissions() {
			if mission.Id == "" || len(mission.Id) <= 0 {
				fmt.Println("Mission id empty.")
				continue
			}
			{
				gameLock.Lock()
				if _, ok := gameMap[mission.Id]; ok {
					gameLock.Unlock()
					continue
				}
				gameMap[mission.Id] = true
				gameLock.Unlock()
			}
			var (
				game      = mission.Id          // 游戏ID
				line      = mission.Line        // 游戏line, 参照游戏发送doSpin请求里面的‘l’值
				insertC   = mission.InsertC     // 插入数据库时的投注，参照游戏发送doSpin请求里面的‘c’值
				normalCnt = mission.NormalCount // 需要拉取普通盘的次数
				gameCnt   = mission.BuyCount    // 需要拉取游戏的次数
				// 使用最低3个挡位
				arrC        = ToMapFI(mission.ArrC)
				gameWeight  = mission.Weight // 如果需要更多机器人拉数据，可以把这个值改大点，建议不修改
				ApiVersion  = mission.APIVersion
				HuiDuSerial = mission.HuiDuSerial
			)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						gameLock.Lock()
						delete(gameMap, game)
						gameLock.Unlock()
						slog.Error("dispatch error:", r)
					}
				}()
				//fmt.Println(line, insertC, normalCnt, gameCnt, gameWeight, arrC, ApiVersion, HuiDuSerial)
				dispatchRobots(game, gameWeight)
				//if os.Getenv("GO_ENV") != "" {
				//fmt.Println("Running in development mode.")
				// 开发环境特有的代码
				if normalCnt > 0 {
					go FetchData(&FetchDataParam{
						Game:        game,
						ArrC:        arrC,
						GBuckets:    GBuckets,
						InsertC:     insertC,
						L:           line,
						Ty:          ppcomm.GameTypeNormal,
						ApiVersion:  ApiVersion,
						HuiDuSerial: HuiDuSerial,
					})
				}
				if gameCnt > 0 {
					go FetchData(&FetchDataParam{
						Game:        game,
						ArrC:        arrC,
						GBuckets:    GBuckets,
						InsertC:     insertC,
						L:           line,
						Ty:          ppcomm.GameTypeGame,
						bonusTy:     0,
						ApiVersion:  ApiVersion,
						HuiDuSerial: HuiDuSerial,
					})
				}
				if mission.SuperBuy > 0 {
					go FetchData(&FetchDataParam{
						Game:        game,
						ArrC:        arrC,
						GBuckets:    GBuckets,
						InsertC:     insertC,
						L:           line,
						Ty:          ppcomm.GameTypeSuperGame1,
						bonusTy:     1,
						ApiVersion:  ApiVersion,
						HuiDuSerial: HuiDuSerial,
					})
				}
				if mission.SuperBuy2 > 0 {
					go FetchData(&FetchDataParam{
						Game:        game,
						ArrC:        arrC,
						GBuckets:    GBuckets,
						InsertC:     insertC,
						L:           line,
						Ty:          ppcomm.GameTypeSuperGame2,
						bonusTy:     2,
						ApiVersion:  ApiVersion,
						HuiDuSerial: HuiDuSerial,
					})
				}
				//}
				needMap := map[int]int64{}
				if normalCnt > 0 {
					needMap[ppcomm.GameTypeNormal] = int64(normalCnt)
				}
				if gameCnt > 0 {
					needMap[ppcomm.GameTypeGame] = int64(gameCnt)
				}
				if mission.SuperBuy > 0 {
					needMap[ppcomm.GameTypeSuperGame1] = int64(mission.SuperBuy)
				}
				if mission.SuperBuy2 > 0 {
					needMap[ppcomm.GameTypeSuperGame2] = int64(mission.SuperBuy2)
				}
				go ppcomm.CheckGameCnt(game, needMap)
				select {}
			}()

		}
	}

}

func TasksDispatch() {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		<-ticker.C

		taskCur := db.Collection2("robot", "game")
		count, err := taskCur.CountDocuments(context.TODO(), bson.M{"status": 1, "isEnd": false})
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("当前并发任务数量: %d\n", count)

		if count < MaxConcurrencyTaskNum {
			// 2. 计算缺少的数量
			missingCount := MaxConcurrencyTaskNum - count
			fmt.Printf("缺少的数量: %d\n", missingCount)

			opts := options.Find().SetLimit(missingCount) // 设置查询数量限制
			cur, err := taskCur.Find(context.TODO(), bson.M{"status": 2, "isEnd": false}, opts)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer cur.Close(context.TODO())
			var tasksToUpdate []GameInfo
			for cur.Next(context.TODO()) {
				var task GameInfo
				if err := cur.Decode(&task); err != nil {
					fmt.Println(err)
					return
				}
				tasksToUpdate = append(tasksToUpdate, task)
			}
			if err := cur.Err(); err != nil {
				fmt.Println(err)
				return
			}

			if len(tasksToUpdate) == 0 {
				fmt.Println("没有找到需要更新的任务")
			} else {
				fmt.Println("将要更新的任务:")
				for _, task := range tasksToUpdate {
					fmt.Printf("ID: %v\n", task.Id)
				}
				// 4. 更新机器人game字段
				for _, task := range tasksToUpdate {
					_, err := taskCur.UpdateOne(context.TODO(), bson.M{"_id": task.Id}, bson.M{"$set": bson.M{"status": 1, "startTime": time.Now().Unix()}})
					if err != nil {
						log.Printf("更新任务失败: %v", err)
					} else {
						fmt.Printf("更新任务成功: ID: %v\n", task.Id)
					}
				}
			}
		} else {
			fmt.Println("不需要更新任务")
		}
	}
}
