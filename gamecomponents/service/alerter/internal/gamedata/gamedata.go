package gamedata

import (
	"fmt"
	"game/comm/mq"
	"game/comm/ut"
	"game/duck/lazy"
	"log/slog"
	"math"
	"slices"
)

type SingleCond struct {
	Single float64
	Total  float64
	RTP    float64
}

type CurrencySetting struct { //这是币种的配置
	Currency    string  //币种
	CurrencyMax float64 //币种最大配置
	CurrencyMin float64 //币种最小配置
}

type Setting struct {
	SingleMin            float64
	SystemInternalSec    int
	SystemMin, SystemMax float64
	GameInternalSec      int
	GameMin, GameMax     float64
	CurrencySettings     []CurrencySetting //币种的配置
}

// Localization 表示多语言的结构体
type Localization struct {
	Zh  Language `json:"zh"`
	En  Language `json:"en"`
	It  Language `json:"it"`
	Es  Language `json:"es"`
	Th  Language `json:"th"`
	Idr Language `json:"idr"`
}

// Language 表示单种语言的字段
type Language struct {
	Profit              string `json:"Profit"`
	BetID               string `json:"BetID"`
	TotalProfitAndLoss  string `json:"TotalProfitAndLoss"`
	TotalBet            string `json:"TotalBet"`
	TotalRateOfReturn   string `json:"TotalRateOfReturn"`
	Currency            string `json:"Currency"`
	CommercialOwner     string `json:"CommercialOwner"`
	LargeAmountTransfer string `json:"LargeAmountTransfer"`
	PlayerID            string `json:"PlayerID"`
	MerchantPlayerID    string `json:"MerchantPlayerID"`
	Merchant            string `json:"Merchant"`
	OrderNumber         string `json:"OrderNumber"`
	TransFerAmount      string `json:"TransFerAmount"`
	LatestBalance       string `json:"LatestBalance"`
	BalanceAlert        string `json:"BalanceAlert"`
	CooperationType     string `json:"CooperationType"`
	WalletType          string `json:"WalletType"`
	MerchantBalance     string `json:"MerchantBalance"`
	TransferWallet      string `json:"TransferWallet"`
	SingleWallet        string `json:"SingleWallet"`
	RevenueShare        string `json:"RevenueShare"`
	CashFlowShare       string `json:"CashFlowShare"`
	BalanceAlertInfo    string `json:"BalanceAlertInfo"`
}

// GetSettingCurrency 获得币种配置的具体值
func (set *Setting) GetSettingCurrency(currency string) CurrencySetting {
	currencySettings := set.CurrencySettings
	fmt.Println("配置的币种集合:", currencySettings)
	for _, currencySetting := range currencySettings {
		fmt.Println("配置的币种:", currencySetting)
		if currencySetting.Currency == currency {
			return currencySetting
		}
	}
	return CurrencySetting{}
}

func (set *Setting) MatchPlayer(single, bet, win float64, cond []*SingleCond) bool {
	if bet <= 0 {
		return false
	}
	totalWin := win - bet
	rtp := 100 * win / bet
	return slices.ContainsFunc(cond, func(c *SingleCond) bool {
		return single > c.Single && totalWin > c.Total && rtp > c.RTP
	})
}

func (set *Setting) GetSingleCond(appId string) []*SingleCond {
	var cond struct {
		List []*SingleCond
	}
	// 调用接口根据appid获取对应的告警配置，获取到配置需要给单笔盈利最小值设置值
	err := mq.Invoke("/AdminInfo/Interior/PushRickRuleReturnRate", map[string]any{
		"AppID": appId,
	}, &cond)
	if len(cond.List) != 0 {
		set.SingleMin = slices.MinFunc(cond.List, func(x, y *SingleCond) int {
			return int(x.Single) - int(y.Single)
		}).Single
	}
	if err != nil {
		slog.Info("获取告警配置失败", "err", err)
		return cond.List
	}
	return cond.List
}

func (set *Setting) GetTransferOutMoney(appId string) float64 {
	var transferOut struct {
		TransferOutMoney float64
	}
	transferOut.TransferOutMoney = math.MaxFloat32
	// 调用接口获取最小的转出预警阈值
	err := mq.Invoke("/AdminInfo/Interior/PushRickRuleTransferOut", map[string]any{
		"AppID": appId,
	}, &transferOut)
	if err != nil {
		slog.Info("获取转出预警阈值失败", "err", err)
		return transferOut.TransferOutMoney
	}
	if transferOut.TransferOutMoney == -1 {
		transferOut.TransferOutMoney = math.MaxFloat32
	}
	return transferOut.TransferOutMoney
}

func GetLanguage(appId string) string {
	var language struct {
		Lang string
	}
	// 调用接口获取最小的转入预警阈值
	err := mq.Invoke("/AdminInfo/Interior/operatorInfo", map[string]any{
		"AppID": appId,
	}, &language)
	if err != nil {
		slog.Info("获取语言失败", "err", err)
	}
	if language.Lang == "" {
		language.Lang = "en"
	}
	return language.Lang
}

// type

var Settings *Setting

func Load() {
	cfg := lazy.ConfigManager
	file := fmt.Sprintf("%v_setting.yaml", lazy.ServiceName)
	cfg.WatchUnmarshal(file, func(set *Setting) error {
		ut.PrintJson(set)
		set.SingleMin = math.MaxFloat32
		Settings = set
		return nil
	})
}

func LoadLocalization() {
	lazy.ConfigManager.WatchUnmarshal("alert_notice.json", loadLanguage)
}

var MapLanguage = map[string]*Language{}

func loadLanguage(settint_ *Localization) error {
	//map init
	MapLanguage = map[string]*Language{}
	MapLanguage["zh"] = &settint_.Zh
	MapLanguage["en"] = &settint_.En
	MapLanguage["it"] = &settint_.It
	MapLanguage["es"] = &settint_.Es
	MapLanguage["th"] = &settint_.Th
	MapLanguage["idr"] = &settint_.Idr
	return nil
}
