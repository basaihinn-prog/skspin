package main

import (
	"fmt"
	"serve/comm/db"
	"serve/servicejili/jili_144_cai/internal"
	"time"
)

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 5; i++ {
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
		//		newFetcher(fmt.Sprintf("%d_extra_%d_%d", internal.GameNo, id, num), coll, gendata.GameTypeExtra).run()
		//		fmt.Println("normal 拉取数据失败了，重试中。。。")
		//		time.Sleep(5 * time.Second)
		//		num++
		//	}
		//}(i)
		//go func(id int) {
		//	var num int
		//	for {
		//		newFetcher(fmt.Sprintf("%d_extraplus_%d_%d", internal.GameNo, id, num), coll, gendata.GameTypeExtraPlus).run()
		//		fmt.Println("normal 拉取数据失败了，重试中。。。")
		//		time.Sleep(5 * time.Second)
		//		num++
		//	}
		//}(i)
	}
	select {}
}
