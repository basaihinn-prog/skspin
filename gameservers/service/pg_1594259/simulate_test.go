package main

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

/*
import (
	"encoding/json"
	"fmt"
	"game/duck/lazy"
	"game/duck/logger"
	"game/service/pg_75/internal/config"
	"os"
	"testing"
	"time"

	_ "gopkg.in/yaml.v3"
)

func TestMain(m *testing.M) {
	fmt.Println("sss_ssss:", "____begin")
	lazy.Init(GameName)
	logger.SetDaily("log", "pg_75")
	config.LoadConfig(func() {})
	c := m.Run()
	os.Exit(c)
}

type Player struct {
	Id   int
	Name string
	Play bool
}

type IndexCount struct {
	Index int
	Count int
}

func TestExtract(t *testing.T) {
	list := make([]*IndexCount, 40)
	for i := 0; i < 40; i++ {
		list[i] = &IndexCount{Index: i, Count: 10000}
	}
	listBytes, _ := json.Marshal(list)
	fmt.Println("listBytes is:", string(listBytes))
}

func getMonthFirstAndLastDay(date string) (day1st, daylast string) {
	year, _ := time.Parse("200601", date)
	month := year.Month()

	// 获取该月的第一天
	firstDay := time.Date(year.Year(), month, 1, 0, 0, 0, 0, time.UTC)

	// 获取该月的下一个月的第一天，然后减去一天得到该月的最后一天
	nextMonth := firstDay.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	// fmt.Println("该月的第一天是：", firstDay)
	// fmt.Println("该月的最后一天是：", lastDay)

	return firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02")
}

func TestDate(t *testing.T) {
	start, end := getMonthFirstAndLastDay("202404")
	fmt.Println("start:", start)
	fmt.Println("end:", end)

	// start, _ := time.Parse("20060102", "20240401")
	// end, _ := time.Parse("20060102", "20240430")

	// fmt.Println("start:", start.Format("2006-01-02"))
	// fmt.Println("end:", end.Format("2006-01-02"))

	//start.Format("2006-01-02")
}

*/

func TestLo(t *testing.T) {
	buckets := [][]int{
		{1, 2, 3, 8, 8},
		{4, 6, 8},
		{9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9},
	}

	oldcountmap := lo.Map(buckets, func(arr []int, _ int) int {
		return len(arr)
	})

	fmt.Println("oldcountmap is:", oldcountmap)

	freeIds := lo.Range(len(buckets))

	fmt.Println("freeIds is:", freeIds)

	freeBucket := lo.Flatten(buckets)
	freePending := lo.Map(freeBucket, func(item int, _ int) []int {
		return []int{item}
	})

	fmt.Println("freeBucket is:", freeBucket)
	fmt.Println("freePending is:", freePending)
}
