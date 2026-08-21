package main

import (
	"fmt"
	"sync"

	"serve/servicejili/jili_38_fs/internal"
	"serve/servicejili/jiliut"

	"time"

	"serve/comm/db"
)

var wg = sync.WaitGroup{}

func main() {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")
	var wgnum = 50
	for i := 0; i < wgnum; i++ {
		wg.Add(wgnum)
		go func(id int) {
			var num int
			for {
				if newFetcher(fmt.Sprintf("fs_teemo_%d_%d", id, num), coll, internal.GameTypeNormal).run() {
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
