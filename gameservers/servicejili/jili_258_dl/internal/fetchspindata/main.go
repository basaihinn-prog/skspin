package main

import (
	"fmt"
	"sync"
	"time"

	"serve/servicejili/jili_258_dl/internal"

	"serve/comm/db"
)

var wg = sync.WaitGroup{}

func main() {
	mongoaddr := "mongodb://myUserAdmin:aajohsujie9vecieSohqu4weivai7oxayei9rie5Yoh4vuojohwaothee0waethi@127.0.0.1:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	// go func() {
	// 	var num int
	// 	for {
	// 		if newFetcher(fmt.Sprintf("dl_%d", num), coll, false).run() {
	// 			return
	// 		}
	// 		fmt.Println("拉取数据失败了，重试中。。。")
	// 		time.Sleep(5 * time.Second)
	// 		num++
	// 	}
	// }()
	// select {}

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("dl_%d_%d", id, num), coll, false).run() {
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
				if newFetcher(fmt.Sprintf("dl_buy_%d_%d", id, num), coll, true).run() {
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
