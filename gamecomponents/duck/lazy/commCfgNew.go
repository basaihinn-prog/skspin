package lazy

import (
	"sync"
)

type commCfgNew struct {
	IsDev         bool
	Alert         *alert
	CurrencyItems []*CurrencyItem `yaml:"currency"`
}

var (
	_commCfgNew = &commCfgNew{
		Alert: &alert{},
	}
	_commMtxNew sync.Mutex
	//plrStoreNew sync.Map
)

func CommCfgNew() *commCfgNew {
	_commMtxNew.Lock()
	defer _commMtxNew.Unlock()
	return _commCfgNew
}

func loadCommCfgNew(cfg *commCfgNew) error {
	_commMtxNew.Lock()
	_commCfgNew = cfg
	//ut.PrintJson(cfg)
	_commMtxNew.Unlock()
	return nil
}

func GetCurrencyItemNew(key string) *CurrencyItem {
	_commMtxNew.Lock()
	defer _commMtxNew.Unlock()
	if key == "" || len(_commCfgNew.CurrencyItems) == 0 {
		return &CurrencyItem{
			Key:    "THB",
			Symbol: "฿",
			Multi:  1,
		}
	}
	if len(_commCfgNew.CurrencyItems) == 1 {
		return _commCfgNew.CurrencyItems[0]
	}
	for i := range _commCfgNew.CurrencyItems {
		if _commCfgNew.CurrencyItems[i].Key == key {
			return _commCfgNew.CurrencyItems[i]
		}
	}
	return &CurrencyItem{
		Key:    "THB",
		Symbol: "฿",
		Multi:  1,
	}
}
