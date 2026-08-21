package main

import (
	"fmt"
	"sync"

	"serve/servicejili/jiliut"

	"serve/servicejili/jili_26_sweetheart3/internal"

	"time"

	"serve/comm/db"
)

var wg = sync.WaitGroup{}

func main() {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("%d_sweetheart3_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeNormal).run()
				fmt.Println("normal 拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		time.Sleep(time.Second)
	}

	wg.Wait()
}
