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
	"serve/service/pg_82/internal/models"
)

func init() {
	mux.RegRpc("/game-api/82/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type D = []any

var Cs = []float64{0.05, 0.5, 2.5, 10.0}
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
		"fb":    nil,
		"wt":    define.M{"mw": 5.0, "bw": 20.0, "mgw": 35.0, "smgw": 50.0},
		"maxwm": nil,
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
		"wp": nil,
		// "bwp": nil,
		"lw": nil,
		// "lwa": 0,
		// "orl": nil,
		// "now": 1920,
		// "nowpr": []int{
		// 	4,
		// 	2,
		// 	4,
		// 	5,
		// 	3,
		// 	4},
		"snww": nil,
		// "esb": map[string][]int{
		// 	"1": {
		// 		1,
		// 		2,
		// 		3,
		// 	},
		// 	"2": {
		// 		6,
		// 		7,
		// 		8,
		// 	},
		// 	"3": {
		// 		9,
		// 		10,
		// 		11,
		// 	},
		// 	"4": {
		// 		12,
		// 		13,
		// 		14,
		// 	},
		// 	"5": {
		// 		21,
		// 		22,
		// 	},
		// 	"6": {
		// 		24,
		// 		25,
		// 	},
		// 	"7": {
		// 		27,
		// 		28,
		// 		29,
		// 	},
		// 	"8": {
		// 		32,
		// 		33,
		// 		34,
		// 	},
		// },
		// "ebb": map[string]map[string]int{
		// 	"1": {
		// 		"fp": 1,
		// 		"lp": 3,
		// 		"bt": 2,
		// 		"ls": 2,
		// 	},
		// 	"2": {
		// 		"fp": 6,
		// 		"lp": 8,
		// 		"bt": 1,
		// 		"ls": 1,
		// 	},
		// 	"3": {
		// 		"fp": 9,
		// 		"lp": 11,
		// 		"bt": 1,
		// 		"ls": 2,
		// 	},
		// 	"4": {
		// 		"fp": 12,
		// 		"lp": 14,
		// 		"bt": 2,
		// 		"ls": 2,
		// 	},
		// 	"5": {
		// 		"fp": 21,
		// 		"lp": 22,
		// 		"bt": 2,
		// 		"ls": 2,
		// 	},
		// 	"6": {
		// 		"fp": 24,
		// 		"lp": 25,
		// 		"bt": 1,
		// 		"ls": 2,
		// 	},
		// 	"7": {
		// 		"fp": 27,
		// 		"lp": 29,
		// 		"bt": 1,
		// 		"ls": 1,
		// 	},
		// 	"8": {
		// 		"fp": 32,
		// 		"lp": 34,
		// 		"bt": 2,
		// 		"ls": 1,
		// 	},
		// },
		// "es": map[string][]int{
		// 	"1": {
		// 		1,
		// 		2,
		// 		3,
		// 	},
		// 	"2": {
		// 		6,
		// 		7,
		// 		8,
		// 	},
		// 	"3": {
		// 		9,
		// 		10,
		// 		11,
		// 	},
		// 	"4": {
		// 		12,
		// 		13,
		// 		14,
		// 	},
		// 	"5": {
		// 		21,
		// 		22,
		// 	},
		// 	"6": {
		// 		24,
		// 		25,
		// 	},
		// 	"7": {
		// 		27,
		// 		28,
		// 		29,
		// 	},
		// 	"8": {
		// 		32,
		// 		33,
		// 		34,
		// 	},
		// },
		// "eb": map[string]map[string]int{
		// 	"1": {
		// 		"fp": 1,
		// 		"lp": 3,
		// 		"bt": 2,
		// 		"ls": 2,
		// 	},
		// 	"2": {
		// 		"fp": 6,
		// 		"lp": 8,
		// 		"bt": 1,
		// 		"ls": 1,
		// 	},
		// 	"3": {
		// 		"fp": 9,
		// 		"lp": 11,
		// 		"bt": 1,
		// 		"ls": 2,
		// 	},
		// 	"4": {
		// 		"fp": 12,
		// 		"lp": 14,
		// 		"bt": 2,
		// 		"ls": 2,
		// 	},
		// 	"5": {
		// 		"fp": 21,
		// 		"lp": 22,
		// 		"bt": 2,
		// 		"ls": 2,
		// 	},
		// 	"6": {
		// 		"fp": 24,
		// 		"lp": 25,
		// 		"bt": 1,
		// 		"ls": 2,
		// 	},
		// 	"7": {
		// 		"fp": 27,
		// 		"lp": 29,
		// 		"bt": 1,
		// 		"ls": 1,
		// 	},
		// 	"8": {
		// 		"fp": 32,
		// 		"lp": 34,
		// 		"bt": 2,
		// 		"ls": 1,
		// 	},
		// },
		// "ssaw": 0,
		// "rs":   nil,
		"fs": nil,
		// "pr":   nil,
		"sc": 1,
		// "mi":   0,
		// "ma": []int{
		// 	1,
		// 	2,
		// 	3,
		// 	4,
		// 	5,
		// },
		"gwt":  -1,
		"fb":   nil,
		"ctw":  0,
		"pmt":  nil,
		"cwc":  0,
		"fstc": nil,
		// "pcwc":  0,
		"rwsp":  nil,
		"hashr": nil,
		"ml":    DefaultBetLevel,
		"cs":    DefaultCs,
		"rl": []int{
			5,
			5,
			4,
			9,
			5,
			5,
			8,
			3,
			3,
			5,
			7,
			9,
			3,
			1,
			9,
		},
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
		"tb":   1,
		"tbb":  1,
		"tw":   0,
		"np":   -1,
		"ocr":  nil,
		"mr":   nil,
		"ge":   []int{1, 11},
	}
}
