package main

import (
	"fmt"
	"math"
	"time"

	"serve/servicejili/jili_115_aa/internal"
	"serve/servicejili/jiliut"

	"serve/comm/db"
)

func main() {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 5; i++ {
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("teemo%s_%d_%d", internal.GameShortName, id, num), coll, false).run()
				fmt.Println("拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
				if num > math.MaxInt16 {
					num = 0
				}
			}
		}(i)
	}
	select {}

}
