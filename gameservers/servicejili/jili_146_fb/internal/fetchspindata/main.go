package main

import (
	"fmt"
	"sync"
	"time"

	"serve/comm/db"
	"serve/servicejili/jili_146_fb/internal"
	"serve/servicejili/jiliut"
	// "serve/servicejili/jili_146_fb/internal/gendata"
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
				if newFetcher(fmt.Sprintf("%d_123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeNormal).run() {
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
				if newFetcher(fmt.Sprintf("%d_extra_123456_%d_%d", internal.GameNo, id, num), coll, internal.GameTypeExtra).run() {
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
