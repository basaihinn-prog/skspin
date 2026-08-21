package jdbcomm

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"serve/comm/db"
	"time"
)

func CheckGameCnt(gameId string, needMap map[int]int64) {
	for ; ; time.Sleep(5 * time.Minute) {
		coll := db.Collection2(gameId, "simulate")
		isEnd := true
		for t, c := range needMap {
			cnt, _ := coll.CountDocuments(context.TODO(), bson.M{"type": t})
			if cnt < c {
				isEnd = false
				break
			}
		}
		if isEnd {
			coll = db.Collection2("robot", "game")
			coll.UpdateByID(context.TODO(), gameId, bson.M{"$set": bson.M{"isEnd": true, "status": 2, "endTime": time.Now().Unix()}})
			break
		}
	}
	//os.Exit(-1)
	fmt.Println("should be finished")
	//panic("should be finished")
}

func CheckGameCntHGC(gameId string, needMap map[int]int64) {
	for ; ; time.Sleep(5 * time.Minute) {
		coll := db.Collection2(gameId, "simulate")
		isEnd := true
		for t, c := range needMap {
			cnt, _ := coll.CountDocuments(context.TODO(), bson.M{"type": t})
			if cnt < c {
				isEnd = false
				break
			}
		}
		if isEnd {
			coll = db.Collection2("robot_hgc", "game")
			coll.UpdateByID(context.TODO(), gameId, bson.M{"$set": bson.M{"isEnd": true, "status": 2, "endTime": time.Now().Unix()}})
			break
		}
	}
	//os.Exit(-1)
	fmt.Println("should be finished")
	//panic("should be finished")
}
