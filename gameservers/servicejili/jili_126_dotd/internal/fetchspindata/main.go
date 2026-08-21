package main

import (
	"fmt"

	"serve/servicejili/jili_126_dotd/internal"
	"serve/servicejili/jiliut"

	"time"

	"serve/comm/db"
)

func main() {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 50; i++ {
		// go func(id int) {
		// 	var num int
		// 	for {
		// 		newFetcher(fmt.Sprintf("dotd_teemo_%d_%d", id, num), coll, false).run()
		// 		fmt.Println("拉取数据失败了，重试中。。。")
		// 		time.Sleep(5 * time.Second)
		// 		num++
		// 	}
		// }(i)
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("dotd_buy_teemo_%d_%d", id, num), coll, true).run()
				fmt.Println("拉取购买数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		time.Sleep(time.Second)
	}

	select {}
}
