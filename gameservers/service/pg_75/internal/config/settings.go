package config

import (
	"fmt"

	"serve/comm/lazy"
)

type Setting struct {
	RewardPercent        int `yaml:"RewardPercent"`        // 千分比
	RewardPercentLess100 int `yaml:"RewardPercentLess100"` // 千分比，转动次数小于100次
	NoAwardPercent       int `yaml:"NoAwardPercent"`       // 个人奖池为负，不中奖千分比
	// NoGoldPercent        int `yaml:"NoGoldPercent"`        // 强制不中奖的概率，千分比
	// PlatAwardPercent     int `yaml:"PlatAwardPercent"`     // 机台的爆奖千分比
	// PlatPoolBetPercent   int `yaml:"PlatPoolBetPercent"`   // 转动的金额的千分比计入机台奖池
	BuyMinAwardPercent int `yaml:"BuyMinAwardPercent"` // 个人奖池为负数，购买游戏时，多少概率随机出0~10倍率
	// 玩家转动的时候，有1%的概率命中不中奖
	HitBigAwardPercent     []int `yaml:"HitBigAwardPercent"`     // 进场前20场的爆奖池的千分比
	BuyGameBigAwardPercent []int `yaml:"BuyGameBigAwardPercent"` // 购买小游戏后接下来的20场的爆奖池的千分比
	// NoGold3BonusPercent    []int `yaml:"NoGold3BonusPercent"`    // 前面5把，不中奖且有3个小游戏块的概率
}

var Settings *Setting

var FinishFn func()

func LoadConfig(finishFn func()) {
	fmt.Println("start LoadConfig++++++++++++")
	FinishFn = finishFn
	cfg := lazy.GConfigManager
	cfg.WatchUnmarshal(fmt.Sprintf("%v_setting.yaml", lazy.ServiceName), func(set *Setting) error {
		Settings = set
		return nil
	})
	FinishFn()
}
