package main

import (
	"fmt"
	"sync"

	"serve/servicejili/jili_87_bog/internal"

	"time"

	"serve/comm/db"
)

var wg = sync.WaitGroup{}

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")
	for i := 0; i < 20; i++ {
		wg.Add(1)
		/*go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("%d_mw3_123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeNormal).run() {
					return
				}
				fmt.Println("normal 拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)*/
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("%d_buy_mw3_123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeGame).run()
				fmt.Println("buy拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		/*go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("%d_extra_mw3_123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeExtra).run() {
					return
				}
				fmt.Println("extra拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)*/
		time.Sleep(time.Second)
	}

	wg.Wait()
}
