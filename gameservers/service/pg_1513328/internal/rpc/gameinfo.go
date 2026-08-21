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
	"serve/service/pg_1513328/internal/models"
)

func init() {
	mux.RegRpc("/game-api/1513328/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"fb":    map[string]any{"is": true, "bm": 75, "t": 150 * curItem.Multi},
		"wt":    define.M{"mw": 5.0, "bw": 20.0, "mgw": 35.0, "smgw": 50.0},
		"maxwm": 8000,
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
		"lw":  nil,
		"bwp": nil,
		"orl": []int{
			8,
			4,
			7,
			8,
			8,
			12,
			9,
			2,
			2,
			2,
			3,
			3,
			1,
			1,
			11,
			11,
			0,
			5,
			1,
			1,
			11,
			11,
			0,
			5,
			9,
			2,
			2,
			2,
			3,
			3,
			8,
			4,
			7,
			8,
			8,
			12,
		},
		"now": 8100,
		"nowpr": []int{
			6,
			3,
			5,
			5,
			3,
			6,
		},
		"snww": nil,
		"esb": map[string][]int{
			"1": {
				3,
			},
			"2": {
				4,
			},
			"3": {
				7,
				8,
				9,
			},
			"4": {
				10,
				11,
			},
			"5": {
				12,
				13,
			},
			"6": {
				14,
			},
			"7": {
				17,
			},
			"8": {
				18,
				19,
			},
			"9": {
				20,
			},
			"10": {
				23,
			},
			"11": {
				25,
				26,
				27,
			},
			"12": {
				28,
				29,
			},
			"13": {
				33,
			},
			"14": {
				34,
			},
		},
		"ebb": map[string]map[string]int{
			"1": {
				"fp": 3,
				"lp": 3,
				"bt": 1,
				"ls": 1,
			},
			"2": {
				"fp": 4,
				"lp": 4,
				"bt": 1,
				"ls": 1,
			},
			"3": {
				"fp": 7,
				"lp": 9,
				"bt": 1,
				"ls": 1,
			},
			"4": {
				"fp": 10,
				"lp": 11,
				"bt": 2,
				"ls": 1,
			},
			"5": {
				"fp": 12,
				"lp": 13,
				"bt": 2,
				"ls": 1,
			},
			"6": {
				"fp": 14,
				"lp": 14,
				"bt": 1,
				"ls": 1,
			},
			"7": {
				"fp": 17,
				"lp": 17,
				"bt": 1,
				"ls": 1,
			},
			"8": {
				"fp": 18,
				"lp": 19,
				"bt": 2,
				"ls": 1,
			},
			"9": {
				"fp": 20,
				"lp": 20,
				"bt": 1,
				"ls": 1,
			},
			"10": {
				"fp": 23,
				"lp": 23,
				"bt": 1,
				"ls": 1,
			},
			"11": {
				"fp": 25,
				"lp": 27,
				"bt": 1,
				"ls": 1,
			},
			"12": {
				"fp": 28,
				"lp": 29,
				"bt": 2,
				"ls": 1,
			},
			"13": {
				"fp": 33,
				"lp": 33,
				"bt": 1,
				"ls": 1,
			},
			"14": {
				"fp": 34,
				"lp": 34,
				"bt": 1,
				"ls": 1,
			},
		},
		"es": map[string][]int{
			"1": {
				3,
			},
			"2": {
				4,
			},
			"3": {
				7,
				8,
				9,
			},
			"4": {
				10,
				11,
			},
			"5": {
				12,
				13,
			},
			"6": {
				14,
			},
			"7": {
				17,
			},
			"8": {
				18,
				19,
			},
			"9": {
				20,
			},
			"10": {
				23,
			},
			"11": {
				25,
				26,
				27,
			},
			"12": {
				28,
				29,
			},
			"13": {
				33,
			},
			"14": {
				34,
			},
		},
		"eb": map[string]map[string]int{
			"1": {
				"fp": 3,
				"lp": 3,
				"bt": 1,
				"ls": 1,
			},
			"2": {
				"fp": 4,
				"lp": 4,
				"bt": 1,
				"ls": 1,
			},
			"3": {
				"fp": 7,
				"lp": 9,
				"bt": 1,
				"ls": 1,
			},
			"4": {
				"fp": 10,
				"lp": 11,
				"bt": 2,
				"ls": 1,
			},
			"5": {
				"fp": 12,
				"lp": 13,
				"bt": 2,
				"ls": 1,
			},
			"6": {
				"fp": 14,
				"lp": 14,
				"bt": 1,
				"ls": 1,
			},
			"7": {
				"fp": 17,
				"lp": 17,
				"bt": 1,
				"ls": 1,
			},
			"8": {
				"fp": 18,
				"lp": 19,
				"bt": 2,
				"ls": 1,
			},
			"9": {
				"fp": 20,
				"lp": 20,
				"bt": 1,
				"ls": 1,
			},
			"10": {
				"fp": 23,
				"lp": 23,
				"bt": 1,
				"ls": 1,
			},
			"11": {
				"fp": 25,
				"lp": 27,
				"bt": 1,
				"ls": 1,
			},
			"12": {
				"fp": 28,
				"lp": 29,
				"bt": 2,
				"ls": 1,
			},
			"13": {
				"fp": 33,
				"lp": 33,
				"bt": 1,
				"ls": 1,
			},
			"14": {
				"fp": 34,
				"lp": 34,
				"bt": 1,
				"ls": 1,
			},
		},
		"ssaw":  0.00,
		"ptbr":  nil,
		"sc":    2,
		"gm":    0,
		"bmtw":  0.0,
		"ms":    0,
		"crtw":  0.0,
		"imw":   false,
		"rs":    nil,
		"fs":    nil,
		"gwt":   0,
		"fb":    nil,
		"ctw":   0.0,
		"pmt":   nil,
		"cwc":   0,
		"fstc":  nil,
		"pcwc":  0,
		"rwsp":  nil,
		"hashr": nil,
		"ml":    DefaultBetLevel,
		"cs":    DefaultCs,
		"rl": []int{
			8,
			4,
			7,
			8,
			8,
			12,
			9,
			2,
			2,
			2,
			3,
			3,
			1,
			1,
			11,
			11,
			0,
			5,
			1,
			1,
			11,
			11,
			0,
			5,
			9,
			2,
			2,
			2,
			3,
			3,
			8,
			4,
			7,
			8,
			8,
			12,
		},
		"sid":  "0",
		"psid": "0",
		"st":   1,
		"nst":  1,
		"pf":   0,
		"aw":   0.00,
		"wid":  0,
		"wt":   "C",
		"wk":   "0_C",
		"wbn":  nil,
		"wfg":  nil,
		"blb":  0.00,
		"blab": 0.00,
		"bl":   10000.00,
		"tb":   0.00,
		"tbb":  0.00,
		"tw":   0.00,
		"np":   0.00,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
}
