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
	"serve/service/pg_1432733/internal/models"
)

func init() {
	mux.RegRpc("/game-api/1432733/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"wp":   nil,
		"wpl":  nil,
		"lw":   nil,
		"bwp":  nil,
		"snww": nil,
		"now":  2700,
		"nowpr": []int{
			5,
			3,
			3,
			3,
			4,
			5,
		},
		"ssaw": 0,
		"esb": map[string][]int{
			"1": {
				6,
				7,
				8,
			},
			"2": {
				11,
				12,
				13,
			},
			"3": {
				16,
				17,
				18,
			},
			"4": {
				20,
				21,
			},
		},
		"ebb": map[string]map[string]int{
			"1": {
				"fp": 6,
				"lp": 8,
				"bt": 2,
				"ls": 1,
			},
			"2": {
				"fp": 11,
				"lp": 13,
				"bt": 2,
				"ls": 1,
			},
			"3": {
				"fp": 16,
				"lp": 18,
				"bt": 2,
				"ls": 1,
			},
			"4": {
				"fp": 20,
				"lp": 21,
				"bt": 2,
				"ls": 1,
			},
		},
		"es": map[string][]int{
			"1": {
				6,
				7,
				8,
			},
			"2": {
				11,
				12,
				13,
			},
			"3": {
				16,
				17,
				18,
			},
			"4": {
				20,
				21,
			},
		},
		"eb": map[string]map[string]int{
			"1": {
				"fp": 6,
				"lp": 8,
				"bt": 2,
				"ls": 1,
			},
			"2": {
				"fp": 11,
				"lp": 13,
				"bt": 2,
				"ls": 1,
			},
			"3": {
				"fp": 16,
				"lp": 18,
				"bt": 2,
				"ls": 1,
			},
			"4": {
				"fp": 20,
				"lp": 21,
				"bt": 2,
				"ls": 1,
			},
		},
		"orl": []int{
			8,
			9,
			10,
			9,
			8,
			5,
			2,
			2,
			2,
			8,
			6,
			3,
			3,
			3,
			9,
			7,
			4,
			4,
			4,
			10,
			10,
			10,
			0,
			1,
			9,
			8,
			9,
			10,
			9,
			8,
		},
		"sc": 1,
		"em": map[string]int{
			"2": 1,
			"3": 1,
			"4": 1,
		},
		"emb": nil,
		"eo": map[string]int{
			"2": 3,
			"3": 3,
			"4": 3,
		},
		"eob":   nil,
		"gm":    1,
		"imw":   false,
		"rs":    nil,
		"fs":    nil,
		"gwt":   0,
		"fb":    nil,
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
			8,
			9,
			10,
			9,
			8,
			5,
			2,
			2,
			2,
			8,
			6,
			3,
			3,
			3,
			9,
			7,
			4,
			4,
			4,
			10,
			10,
			10,
			0,
			1,
			9,
			8,
			9,
			10,
			9,
			8,
		},
		"sid":  "0",
		"psid": "0",
		"st":   1,
		"nst":  1,
		"pf":   0,
		"aw":   0,
		"wid":  0,
		"wt":   "C",
		"wk":   "0_C",
		"wbn":  nil,
		"wfg":  nil,
		"blb":  0,
		"blab": 0,
		"bl":   22771.74,
		"tb":   0,
		"tbb":  0,
		"tw":   0,
		"np":   0,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
}
