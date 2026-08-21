package gamedata

import (
	"fmt"
	"sync"

	"serve/comm/lazy"
	"serve/comm/ut"
)

type Setting struct {
	RewardPercent        int `yaml:"RewardPercent"`        // 千分比
	RewardPercentLess100 int `yaml:"RewardPercentLess100"` // 千分比，转动次数小于100次
	NoAwardPercent       int `yaml:"NoAwardPercent"`       // 个人奖池为负，不中奖千分比
	// 玩家转动的时候，有1%的概率命中不中奖
	HitBigAwardPercent []int `yaml:"HitBigAwardPercent"` // 进场前20场的爆奖池的千分比

	Extra5RewardPercent        int `yaml:"Extra5RewardPercent"`        // 千分比
	Extra5RewardPercentLess100 int `yaml:"Extra5RewardPercentLess100"` // 千分比，转动次数小于100次
	Extra5NoAwardPercent       int `yaml:"Extra5NoAwardPercent"`       // 个人奖池为负，不中奖千分比
	// 玩家转动的时候，有1%的概率命中不中奖
	Extra5HitBigAwardPercent []int `yaml:"Extra5HitBigAwardPercent"` // 进场前20场的爆奖池的千分比

	Extra50RewardPercent        int `yaml:"Extra50RewardPercent"`        // 千分比
	Extra50RewardPercentLess100 int `yaml:"Extra50RewardPercentLess100"` // 千分比，转动次数小于100次
	Extra50NoAwardPercent       int `yaml:"Extra50NoAwardPercent"`       // 个人奖池为负，不中奖千分比
	// 玩家转动的时候，有1%的概率命中不中奖
	Extra50HitBigAwardPercent []int `yaml:"Extra50HitBigAwardPercent"` // 进场前20场的爆奖池的千分比
}

var settings *Setting
var settingsLock sync.Mutex

func Load() {
	cfg := lazy.GConfigManager

	file := fmt.Sprintf("%v_setting.yaml", lazy.ServiceName)
	cfg.WatchUnmarshal(file, func(set *Setting) error {
		settingsLock.Lock()
		defer settingsLock.Unlock()
		settings = set
		ut.PrintJson(set)
		return nil
	})
}

func GetSettings() *Setting {
	settingsLock.Lock()
	defer settingsLock.Unlock()
	return settings
}
