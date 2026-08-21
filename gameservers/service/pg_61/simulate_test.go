package main

import (
	"fmt"
	"testing"

	"serve/comm/db"
	"serve/comm/ut"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
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
	str := `{
		"wp": {
		  "3": [
			2,
			5,
			8
		  ],
		  "5": [
			2,
			4,
			6
		  ],
		  "7": [
			2,
			5,
			7
		  ],
		  "11": [
			2,
			4,
			7
		  ],
		  "17": [
			2,
			4,
			8
		  ],
		  "19": [
			2,
			5,
			7
		  ]
		},
		"lw": {
		  "3": 0.25,
		  "5": 0.25,
		  "7": 0.25,
		  "11": 0.25,
		  "17": 0.25,
		  "19": 0.25
		},
		"fs": null,
		"orl": null,
		"gwt": -1,
		"fb": null,
		"ctw": 1.5,
		"pmt": null,
		"cwc": 1,
		"fstc": null,
		"pcwc": 1,
		"rwsp": {
		  "3": 5,
		  "5": 5,
		  "7": 5,
		  "11": 5,
		  "17": 5,
		  "19": 5
		},
		"hashr": "0:3;2;8;9;7#3;0;8;9;5#8;0;8;5;5#R#8#021222#MV#1.00#MT#1#R#8#021120#MV#1.00#MT#1#R#8#021221#MV#1.00#MT#1#R#8#021121#MV#1.00#MT#1#R#8#021122#MV#1.00#MT#1#R#8#021221#MV#1.00#MT#1#MG#1.50#",
		"ml": 1,
		"cs": 0.05,
		"rl": [
		  3,
		  3,
		  8,
		  2,
		  0,
		  0,
		  8,
		  8,
		  8,
		  9,
		  9,
		  5,
		  7,
		  5,
		  5
		],
		"sid": "1788115919684370432",
		"psid": "1788115919684370432",
		"st": 1,
		"nst": 1,
		"pf": 2,
		"aw": 1.5,
		"wid": 0,
		"wt": "C",
		"wk": "0_C",
		"wbn": null,
		"wfg": null,
		"blb": 31282.14,
		"blab": 31281.14,
		"bl": 31282.64,
		"tb": 1,
		"tbb": 1,
		"tw": 1.5,
		"np": 0.5,
		"ocr": null,
		"mr": null,
		"ge": [
		  1,
		  11
		]
	  }`
	bsondoc := lo.Must(db.Json2Bson([]byte(str)))

	num := getNoSelectPositionNum(bsondoc)
	fmt.Println("-1的个数是:", num)
}

func getNoSelectPositionNum(bsondoc bson.Raw) int {
	type SI_FS struct {
		BF []bson.M
	}

	type SI struct {
		FS *SI_FS
	}

	var si SI
	err := bson.Unmarshal(bsondoc, &si)
	if err != nil {
		fmt.Println("getNoSelectPositionNum Unmarshal error :", err)
	}
	if si.FS == nil {
		return 0
	}
	numOf_1 := 0
	for _, v := range si.FS.BF {
		fpValue := ut.GetInt(v["fp"])
		if fpValue == -1 {
			numOf_1++
		}
	}
	return numOf_1
}
