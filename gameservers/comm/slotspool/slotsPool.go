package slotspool

import (
	"serve/comm/ut"
	"sync"
)

type SlotsPool struct {
	Value        int            //千分比
	WinMaxLimits []*WinMaxLimit `yaml:"winMaxLimits"` // 运营商最大赢分限制
}

type WinMaxLimit struct {
	Min    int    `yaml:"min"`    // 最小净输赢
	Max    int    `yaml:"max"`    // 最大净输赢
	MaxWin int    `yaml:"maxWin"` // 最大赢分
	Desc   string `yaml:"desc"`   // 描述
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

func GetMaxWin(appId string, winLose int64, multi float64) *WinMaxLimit {
	slotsPool := GetSlotsPool(appId)
	for i := range slotsPool.WinMaxLimits {
		if int64(float64(slotsPool.WinMaxLimits[i].Min)*multi) <= winLose &&
			int64(float64(slotsPool.WinMaxLimits[i].Max)*multi) > winLose {
			return slotsPool.WinMaxLimits[i]
		}
	}
	slotsPool = GetSlotsPool("default")
	for i := range slotsPool.WinMaxLimits {
		if int64(float64(slotsPool.WinMaxLimits[i].Min)*multi) <= winLose &&
			int64(float64(slotsPool.WinMaxLimits[i].Max)*multi) > winLose {
			return slotsPool.WinMaxLimits[i]
		}
	}
	return nil
}
