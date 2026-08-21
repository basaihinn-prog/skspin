package main

import (
	"fmt"
	"serve/comm/db"
	"serve/servicejili/jili_300_fg3/internal"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("fg3_123456_%d_%d", id, num), coll, internal.GameTypeNormal).run() {
					return
				}
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("fg3_extra_123456_%d_%d", id, num), coll, internal.GameTypeExtra).run() {
					return
				}
				fmt.Println("拉取购买数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		time.Sleep(time.Second)
	}

	wg.Wait()
}
