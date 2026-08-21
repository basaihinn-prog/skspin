package lazy

import (
	"context"
	"game/comm/db"
	"game/comm/slotspool"
	"game/comm/ut"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	IsDev         bool
	Alert         *alert
	CurrencyItems []*CurrencyItem `yaml:"currency"`
}

var (
	_commCfg = &commCfg{
		Alert: &alert{},
	}
	_commMtx sync.Mutex
	plrStore sync.Map
)

func load() {
	ConfigManager.WatchUnmarshal("comm_config.yaml", loadCommCfg)
	ConfigManager.WatchUnmarshal("slotsPool.yaml", slotspool.LoadSlotsPool)
	ConfigManager.WatchUnmarshal("gameSlotsPool.yaml", slotspool.LoadGameSlotsPool)
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

func GetCurrentItem(pid int64) *CurrencyItem {
	var item *CurrencyItem
	data, ok := plrStore.Load(pid)
	if ok { //存在
		item = data.(*CurrencyItem)
	} else {
		var doc struct {
			AppID string `bson:"AppID"` // 所属产品
		}
		// var appid string
		coll := db.Collection2("game", "Players")
		coll.FindOne(context.TODO(), db.ID(pid), options.FindOne().SetProjection(db.D("AppID", 1))).Decode(&doc)

		type currencyItem struct {
			CurrencyKey string `bson:"CurrencyKey"`
		}
		var tmp *currencyItem

		coll = db.Collection2("GameAdmin", "AdminOperator")
		coll.FindOne(context.TODO(), bson.M{"AppID": doc.AppID}, options.FindOne().SetProjection(db.D("CurrencyKey", 1))).Decode(&tmp)
		if tmp != nil {
			item = GetCurrencyItem(tmp.CurrencyKey)
		} else {
			item = GetCurrencyItem("")
		}
		plrStore.Store(pid, item)
	}
	return item
}
