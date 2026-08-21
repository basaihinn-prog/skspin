package main

import (
	"fmt"
	"serve/comm/db"
	"serve/servicejili/jili_109_fg/internal"
	"time"
)

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 20; i++ {
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("fg_123456_%d_%d", id, num), coll, internal.GameTypeNormal).run()
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		//go func(id int) {
		//	var num int
		//	for {
		//		newFetcher(fmt.Sprintf("fg_buy_123456_%d_%d", id, num), coll, internal.GameTypeExtra).run()
		//		fmt.Println("拉取购买数据失败了，重试中。。。")
		//		time.Sleep(5 * time.Second)
		//		num++
		//	}
		//}(i)
		time.Sleep(time.Second)
	}

	select {}
}
