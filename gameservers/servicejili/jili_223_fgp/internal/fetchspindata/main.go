package main

import (
	"fmt"
	"sync"

	"serve/servicejili/jili_223_fgp/internal"
	"serve/servicejili/jiliut"

	"time"

	"serve/comm/db"
)

var wg = sync.WaitGroup{}

func main() {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("%d_xx123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeNormal).run() {
					return
				}
				fmt.Println("normal 拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("%d_extra_xx123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeExtra).run() {
					return
				}
				fmt.Println(" 拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
		time.Sleep(time.Second)
	}
	wg.Wait()
}
