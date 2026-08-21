package internal

import (
	"context"
	"fmt"
	"log"
	"serve/comm/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Game struct {
	ID string `bson:"_id"` // 使用 primitive.ObjectID 作为 _id 类型
}

func LoopUpdatePlayerSpinCountTo100() {
	ticker := time.NewTicker(time.Minute * 1)
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "game")

	gameCur, err := db.Collection2("game", "Games").Find(context.TODO(), bson.M{"Status": 0})
	if err != nil {

	}
	defer gameCur.Close(context.TODO())
	var games []Game
	for gameCur.Next(context.TODO()) {
		var game Game
		err = gameCur.Decode(&game)
		if err != nil {
			fmt.Println(err)
			return
		}
		games = append(games, game)
	}
	if err := gameCur.Err(); err != nil {
		fmt.Println(err)
		return
	}
	for {
		for _, game := range games {
			db.DialToMongo(mongoaddr, game.ID)
			coll := db.Collection2(game.ID, "players")
			update := bson.D{{"$inc", bson.D{{"spincount", -100}}}}
			// 更新所有文档，没有过滤器
			result, err := coll.UpdateMany(context.TODO(), bson.D{}, update)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("game:%v players fix data count: %v \n", game.ID, result.ModifiedCount)
		}
		<-ticker.C
	}
}
