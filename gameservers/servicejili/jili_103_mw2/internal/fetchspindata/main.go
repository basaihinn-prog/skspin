package main

import (
	"fmt"
	"serve/servicejili/jili_103_mw2/internal"
	"serve/servicejili/jiliut"

	"time"

	"serve/comm/db"
)

func main() {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 20; i++ {
		//go func(id int) {
		//	var num int
		//	for {
		//		newFetcher(fmt.Sprintf("mw2_123456_%d_%d", id, num), coll, false).run()
		//		fmt.Println("拉取数据失败了，重试中。。。")
		//		time.Sleep(5 * time.Second)
		//		num++
		//	}
		//}(i)
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("mw2_buy_123456_%d_%d", id, num), coll, true).run()
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		time.Sleep(time.Second)
	}

	select {}
}
