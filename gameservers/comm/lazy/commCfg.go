package lazy

import (
	"serve/comm/slotspool"
	"sync"

	"serve/comm/ut"
)

type alert struct {
	TelegramNotify          string
	AlertGameThreshold      float64
	AlertSingleWinThreshold float64
}

type CurrencyItem struct {
	Key    string  `yaml:"key"`
	Symbol string  `yaml:"symbol"`
	Multi  float64 `yaml:"multi"`
}

type commCfg struct {
	IsDev bool `yaml:"isDev"`
	Alert *alert
	// BackupHost string
	CurrencyItems []*CurrencyItem `yaml:"currency"`
}

var (
	_commCfg = &commCfg{
		Alert: &alert{},
	}
	_commMtx sync.Mutex
)

func load() {
	GConfigManager.WatchUnmarshal("comm_config.yaml", loadCommCfg)
	GConfigManager.WatchUnmarshal("slotsPool.yaml", slotspool.LoadSlotsPool)
}

func loadCommCfg(cfg *commCfg) error {
	_commMtx.Lock()
	_commCfg = cfg
	ut.PrintJson(cfg)
	_commMtx.Unlock()
	return nil
}

func CommCfg() *commCfg {
	_commMtx.Lock()
	defer _commMtx.Unlock()
	return _commCfg
}

func GetCurrencyItem(key string) *CurrencyItem {
	_commMtx.Lock()
	defer _commMtx.Unlock()
	if key == "" || len(_commCfg.CurrencyItems) == 0 {
		return &CurrencyItem{
			Key:    "THB",
			Symbol: "฿",
			Multi:  1,
		}
	}
	if len(_commCfg.CurrencyItems) == 1 {
		return _commCfg.CurrencyItems[0]
	}
	for i := range _commCfg.CurrencyItems {
		if _commCfg.CurrencyItems[i].Key == key {
			return _commCfg.CurrencyItems[i]
		}
	}
	return &CurrencyItem{
		Key:    "THB",
		Symbol: "฿",
		Multi:  1,
	}
}
