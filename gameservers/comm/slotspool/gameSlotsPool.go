package slotspool

import (
	"serve/comm/ut"
	"sync"
)

type GameSlotsPool struct {
	RewardPercent      int `yaml:"RewardPercent"`
	ExtraRewardPercent int `yaml:"ExtraRewardPercent"`
}

var gameSlotsPools map[string]*GameSlotsPool
var gameSlotsMutex sync.Mutex

func LoadGameSlotsPool(cfg map[string]*GameSlotsPool) error {
	gameSlotsMutex.Lock()
	defer gameSlotsMutex.Unlock()
	ut.PrintJson(cfg)
	if cfg["default"] == nil {
		cfg["default"] = &GameSlotsPool{
			RewardPercent:      30,
			ExtraRewardPercent: 30,
		}
	}
	gameSlotsPools = cfg
	return nil
}

func GetGameSlotsPool(gameId string) *GameSlotsPool {
	gameSlotsMutex.Lock()
	defer gameSlotsMutex.Unlock()

	if gameSlotsPools == nil {
		return &GameSlotsPool{
			RewardPercent: 30,
		}
	}
	ret, ok := gameSlotsPools[gameId]
	if !ok {
		ret, ok = gameSlotsPools["default"]
		if !ok {
			return &GameSlotsPool{
				RewardPercent:      30,
				ExtraRewardPercent: 30,
			}
		}
	}
	return ret
}
