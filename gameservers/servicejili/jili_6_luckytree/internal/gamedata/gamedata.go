package gamedata

import (
	"cmp"
	"fmt"

	"serve/comm/lazy"
	"serve/comm/ut"
)

type Setting struct {
	RewardPercent        int `yaml:"RewardPercent"`        // 千分比
	RewardPercentLess100 int `yaml:"RewardPercentLess100"` // 千分比，转动次数小于100次
	NoAwardPercent       int `yaml:"NoAwardPercent"`       // 个人奖池为负，不中奖千分比
	BuyMinAwardPercent   int `yaml:"BuyMinAwardPercent"`   // 个人奖池为负数，购买游戏时，多少概率随机出0~10倍率
	// 玩家转动的时候，有1%的概率命中不中奖
	HitBigAwardPercent []int   `yaml:"HitBigAwardPercent"` // 进场前20场的爆奖池的千分比
	PriceOdd           float64 `yaml:"PriceOdd"`           // 购买freegame 倍数
	BuyMulti           int     `yaml:"BuyMulti"`           // 购买实际倍率
}

var Settings *Setting

func Load() {
	cfg := lazy.GConfigManager

	file := fmt.Sprintf("%v_setting.yaml", lazy.ServiceName)
	cfg.WatchUnmarshal(file, func(set *Setting) error {
		set.PriceOdd = cmp.Or(set.PriceOdd, 50)
		Settings = set
		ut.PrintJson(set)
		return nil
	})
}
