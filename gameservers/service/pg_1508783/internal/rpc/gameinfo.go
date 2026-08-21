package rpc

import (
	"encoding/json"
	"log/slog"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lazy"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_1508783/internal/models"
)

func init() {
	mux.RegRpc("/game-api/1508783/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type D = []any

var Cs = []float64{0.05, 0.5, 2.5, 10}
var Ml = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var Mxl = 20

func gameinfo(plr *models.Player, ps define.PGParams, ret *define.M) (err error) {
	gold, err := slotsmongo.GetBalance(ps.Pid)
	if err != nil {
		return
	}
	info, err := redisx.LoadAppIdCache(plr.AppID)
	if err != nil {
		return err
	}
	isEnd, _ := plr.IsEndO()
	c, err := redisx.GetPlayerCs(plr.AppID, plr.PID, isEnd)
	if err != nil {
		slog.Error("doSpin", "GetPlayerCs", err)
		return err
	}
	//修改end
	slotsmongo.UpdateEnterPlrCount(lazy.ServiceName, plr.AppID, plr.PID)
	plr.SpinCountOfThisEnter = 0
	balance := ut.Gold2Money(gold)
	curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
	s := define.M{
		"fb":    map[string]any{"is": true, "bm": 75, "t": 100 * curItem.Multi},
		"wt":    define.M{"mw": 5.0, "bw": 20.0, "mgw": 35.0, "smgw": 50.0},
		"maxwm": 10000,
		"cs":    ut.FloatArrMul(c, curItem.Multi),
		"ml":    Ml,
		"mxl":   Mxl,
		"bl":    10000.00,
		"inwe":  false,
		"iuwe":  false,
		"cc":    curItem.Key,
		"gc":    info.ShowNameAndTimeOff,
		"ign":   info.ShowNameAndTimeOff,
		"asc":   info.StopLoss,
	}
	s["bl"] = balance

	if len(plr.LS) != 0 {
		s["ls"] = json.RawMessage(plr.LS)
	} else {
		si := getInitSi(info.DefaultCs, info.DefaultBetLevel)
		si["bl"] = balance
		s["ls"] = map[string]any{"si": si}
	}

	*ret = s
	// *ret = json.RawMessage(jsonstr)
	return
}

func getInitSi(DefaultCs float64, DefaultBetLevel int64) map[string]any {
	return map[string]any{
		"wp":  nil,
		"bwp": nil,
		"lw":  nil,
		// "lwa":  0,
		"orl": []int{
			12, 12, 5, 5, 7, 7, 11, 11, 11, 11, 11, 2, 11, 11, 99, 99, 0, 4, 6, 8, 99, 99, 12, 8, 11, 12, 12, 12, 7, 7, 8, 8, 8, 9, 9, 10},
		// "now":  1920,
		"snww": nil,
		"esb": map[string][]int{
			"1": {
				2,
				3,
			},
			"2": {
				4,
				5,
			},
			"3": {
				6,
				7,
				8,
			},
			"4": {
				9,
				10,
			},
			"5": {
				25,
				26,
				27,
			},
			"6": {
				28,
				29,
			},
			"7": {
				30,
				31,
				32,
			},
			"8": {
				33,
				34,
			},
			"9": {
				14,
				15,
			},
			"10": {
				20,
				21,
			},
			"11": {
				11,
			},
			"12": {
				22,
			},
		},
		"ebb": map[string]map[string]int{
			"1": {
				"fp": 2,
				"lp": 3,
				"bt": 2,
				"ls": 1,
			},
			"2": {
				"fp": 4,
				"lp": 5,
				"bt": 2,
				"ls": 1,
			},
			"3": {
				"fp": 6,
				"lp": 8,
				"bt": 3,
				"ls": 1,
			},
			"4": {
				"fp": 9,
				"lp": 10,
				"bt": 2,
				"ls": 1,
			},
			"5": {
				"fp": 25,
				"lp": 27,
				"bt": 2,
				"ls": 1,
			},
			"6": {
				"fp": 28,
				"lp": 29,
				"bt": 3,
				"ls": 1,
			},
			"7": {
				"fp": 30,
				"lp": 32,
				"bt": 2,
				"ls": 1,
			},
			"8": {
				"fp": 33,
				"lp": 34,
				"bt": 2,
				"ls": 1,
			},
			"9": {
				"fp": 14,
				"lp": 15,
				"bt": 2,
				"ls": 1,
			},
			"10": {
				"fp": 20,
				"lp": 21,
				"bt": 2,
				"ls": 1,
			},
			"11": {
				"fp": 11,
				"lp": 11,
				"bt": 3,
				"ls": 1,
			},
			"12": {
				"fp": 22,
				"lp": 22,
				"bt": 3,
				"ls": 1,
			},
		},
		"es": map[string][]int{
			"1": {
				2,
				3,
			},
			"2": {
				4,
				5,
			},
			"3": {
				6,
				7,
				8,
			},
			"4": {
				9,
				10,
			},
			"5": {
				25,
				26,
				27,
			},
			"6": {
				28,
				29,
			},
			"7": {
				30,
				31,
				32,
			},
			"8": {
				33,
				34,
			},
			"9": {
				14,
				15,
			},
			"10": {
				20,
				21,
			},
			"11": {
				11,
			},
			"12": {
				22,
			},
		},
		"eb": map[string]map[string]int{
			"1": {
				"fp": 2,
				"lp": 3,
				"bt": 2,
				"ls": 1,
			},
			"2": {
				"fp": 4,
				"lp": 5,
				"bt": 2,
				"ls": 1,
			},
			"3": {
				"fp": 6,
				"lp": 8,
				"bt": 3,
				"ls": 1,
			},
			"4": {
				"fp": 9,
				"lp": 10,
				"bt": 2,
				"ls": 1,
			},
			"5": {
				"fp": 25,
				"lp": 27,
				"bt": 2,
				"ls": 1,
			},
			"6": {
				"fp": 28,
				"lp": 29,
				"bt": 3,
				"ls": 1,
			},
			"7": {
				"fp": 30,
				"lp": 32,
				"bt": 2,
				"ls": 1,
			},
			"8": {
				"fp": 33,
				"lp": 34,
				"bt": 2,
				"ls": 1,
			},
			"9": {
				"fp": 14,
				"lp": 15,
				"bt": 2,
				"ls": 1,
			},
			"10": {
				"fp": 20,
				"lp": 21,
				"bt": 2,
				"ls": 1,
			},
			"11": {
				"fp": 11,
				"lp": 11,
				"bt": 3,
				"ls": 1,
			},
			"12": {
				"fp": 22,
				"lp": 22,
				"bt": 3,
				"ls": 1,
			},
		},
		"ssaw":  0,
		"rs":    nil,
		"fs":    nil,
		"pr":    nil,
		"sc":    0,
		"mi":    0,
		"crtw":  0,
		"gwt":   -1,
		"gm":    1,
		"fb":    nil,
		"imw":   false,
		"ctw":   0,
		"pmt":   nil,
		"cwc":   0,
		"fstc":  nil,
		"pcwc":  0,
		"rwsp":  nil,
		"hashr": nil,
		"ml":    DefaultBetLevel,
		"cs":    DefaultCs,
		"rl": []int{
			12, 12, 5, 5, 7, 7, 11, 11, 11, 11, 11, 2, 11, 11, 99, 99, 0, 4, 6, 8, 99, 99, 12, 8, 11, 12, 12, 12, 7, 7, 8, 8, 8, 9, 9, 10},
		"sid":  "0",
		"psid": "0",
		"st":   1,
		"nst":  1,
		"pf":   2,
		"aw":   0,
		"wid":  0,
		"wt":   "C",
		"wk":   "0_C",
		"wbn":  nil,
		"wfg":  nil,
		"blb":  0,
		"blab": 0,
		"bl":   0,
		"tb":   0,
		"tbb":  0,
		"tw":   0,
		"np":   0,
		"ocr":  nil,
		"mr":   nil,
		"ge":   []int{1, 11},
	}
}
