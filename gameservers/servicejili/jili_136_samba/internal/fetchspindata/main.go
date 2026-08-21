package main

import (
	"fmt"
	"serve/servicejili/jili_136_samba/internal"
	"time"

	"serve/comm/db"
)

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 10; i++ {
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("%d_123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeNormal).run()
				fmt.Println("normal 拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		//go func(id int) {
		//	var num int
		//	for {
		//		newFetcher(fmt.Sprintf("%d_buy_123456_%d_%d", internal.GameNo, id, num), coll, gendata.GameTypeGame).run()
		//		fmt.Println("buy 拉取数据失败了，重试中。。。")
		//		time.Sleep(5 * time.Second)
		//		num++
		//	}
		//}(i)
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("%d_extra_123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeExtra).run()
				fmt.Println("buy 拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
	}

	select {}
}
