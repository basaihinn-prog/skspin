package gamedata

import (
	"fmt"

	"serve/comm/lazy"
	"serve/comm/ut"
	"serve/servicepp/pp_vs20piggybank/internal"
)

type Setting struct {
	RewardPercent        int   `yaml:"RewardPercent"`        // 千分比
	RewardPercentLess100 int   `yaml:"RewardPercentLess100"` // 千分比，转动次数小于100次
	NoAwardPercent       int   `yaml:"NoAwardPercent"`       // 个人奖池为负，不中奖千分比
	BuyMinAwardPercent   int   `yaml:"BuyMinAwardPercent"`   // 个人奖池为负数，购买游戏时，多少概率随机出0~10倍率
	HitBigAwardPercent   []int `yaml:"HitBigAwardPercent"`   // 进场前20场的爆奖池的千分比
}

var Settings *Setting

func Load() {
	cfg := lazy.GConfigManager
	file := fmt.Sprintf("%v_setting.yaml", internal.GameID)
	cfg.WatchUnmarshal(file, func(set *Setting) error {
		Settings = set
		ut.PrintJson(set)
		return nil
	})
}
