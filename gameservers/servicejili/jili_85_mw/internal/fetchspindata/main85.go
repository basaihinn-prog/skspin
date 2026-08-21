package main

import (
	"fmt"
	"serve/servicejili/jili_85_mw/internal"
	"time"

	"serve/comm/db"
)

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 5; i++ {
		go func(id int) {

			var num int
			for {
				newFetcher(fmt.Sprintf("%d_123456_%d_%d", internal.GameNo, id, num), coll, false).run()
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		//go func(id int) {
		//	var num int
		//	for {
		//		newFetcher(fmt.Sprintf("%d_buy_123456_%d_%d", internal.GameNo, id, num), coll, true).run()
		//		fmt.Println("拉取数据失败了，重试中。。。")
		//		time.Sleep(5 * time.Second)
		//		num++
		//	}
		//}(i)
	}

	select {}
	//lazy.SignalProc()
}
