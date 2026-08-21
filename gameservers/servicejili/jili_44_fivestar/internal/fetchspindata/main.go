package main

import (
	"fmt"
	"serve/servicejili/jili_44_fivestar/internal"
	"sync"

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
				newFetcher(fmt.Sprintf("fivestar_123456_%d_%d", id, num), coll, internal.GameTypeNormal).run()
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)

		time.Sleep(time.Second)
	}

	wg.Wait()
}
