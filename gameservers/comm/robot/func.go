package robot

import (
	"context"
	"fmt"
	"serve/comm/db"

	"go.mongodb.org/mongo-driver/bson"
)

func GetRobot(gameId string) []*Robot {
	//now := time.Now().Unix()
	coll := db.Collection2("robot", "robot")
	var ret []*Robot
	//cur, _ := coll.Find(context.TODO(), bson.M{"gameId": gameId, "endTime": bson.M{"$gt": now}})
	cur, _ := coll.Find(context.TODO(), bson.M{"gameId": gameId})
	err := cur.All(context.TODO(), &ret)
	if err != nil {
		fmt.Println(err)
	}
	return ret
}
func GetRobotHGC(gameId string) []*Robot {
	//now := time.Now().Unix()
	coll := db.Collection2("robot_hgc", "robot")
	var ret []*Robot
	//cur, _ := coll.Find(context.TODO(), bson.M{"gameId": gameId, "endTime": bson.M{"$gt": now}})
	cur, _ := coll.Find(context.TODO(), bson.M{"gameId": gameId})
	err := cur.All(context.TODO(), &ret)
	if err != nil {
		fmt.Println(err)
	}
	return ret
}

func CheckGameEnd(gameId string) bool {
	coll := db.Collection2("robot", "game")
	var doc struct {
		IsEnd bool `bson:"isEnd"`
	}

	err := coll.FindOne(context.TODO(), db.ID(gameId)).Decode(&doc)
	if err != nil {
		return true
	}
	return doc.IsEnd
}
func CheckGameEndHGC(gameId string) bool {
	coll := db.Collection2("robot_hgc", "game")
	var doc struct {
		IsEnd bool `bson:"isEnd"`
	}

	err := coll.FindOne(context.TODO(), db.ID(gameId)).Decode(&doc)
	if err != nil {
		return true
	}
	return doc.IsEnd
}
