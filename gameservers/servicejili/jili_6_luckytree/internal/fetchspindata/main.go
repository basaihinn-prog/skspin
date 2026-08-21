package main

import (
	"fmt"

	"serve/servicejili/jili_6_luckytree/internal"

	"time"

	"serve/comm/db"
)

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	//go fetch(internal.GameID, coll)
	//go fetchBuy("buy102_123456", coll)

	for i := 0; i < 5; i++ {
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("luckytree_123456_%d_%d", id, num), coll, false).run()
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
	}
	select {}
}
