package slotspool

import (
	"game/comm/ut"
	"sync"
)

type SlotsPool struct {
	Value int //千分比
}

var slotsPools map[string]*SlotsPool
var slotsMutex sync.Mutex

func LoadSlotsPool(cfg map[string]*SlotsPool) error {
	slotsMutex.Lock()
	defer slotsMutex.Unlock()
	ut.PrintJson(cfg)
	slotsPools = cfg
	return nil
}

func GetSlotsPool(appId string) *SlotsPool {
	if slotsPools == nil {
		return &SlotsPool{
			Value: 0,
		}
	}
	ret, ok := slotsPools[appId]
	if !ok {
		ret, _ = slotsPools["default"]
	}
	if ret == nil {
		return &SlotsPool{
			Value: 0,
		}
	}
	return ret
}
