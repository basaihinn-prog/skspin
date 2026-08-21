package main

import (
	"fmt"
	"sync"

	"serve/servicejili/jili_91_ifff/internal"

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
		go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("bk_123456_%d_%d", id, num), coll).run() {
					return
				}
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)

		time.Sleep(time.Second)
	}

	wg.Wait()
}
