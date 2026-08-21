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
	"serve/service/pg_120/internal/models"
)

func init() {
	mux.RegRpc("/game-api/120/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"wp":  nil,
		"twp": nil,
		"bwp": nil,
		"lw":  nil,
		"lwa": 0,
		"orl": []int{
			8,
			4,
			1,
			10,
			7,
			0,
			0,
			0,
			11,
			10,
			2,
			2,
			2,
			3,
			3,
			2,
			2,
			2,
			3,
			3,
			0,
			9,
			12,
			7,
			7,
			12,
			5,
			1,
			11,
			6,
		},
		"trl": []int{
			6,
			3,
			3,
			2,
		},
		"now":   0,
		"nowpr": nil,
		"snww":  nil,
		"esb": map[string][]int{
			"1": {
				5,
				6,
				7,
			},
			"2": {
				10,
				11,
				12,
			},
			"3": {
				13,
				14,
			},
			"4": {
				15,
				16,
				17,
			},
			"5": {
				18,
				19,
			},
		},
		"ebb": map[string]map[string]int{
			"1": {
				"fp": 5,
				"lp": 7,
				"bt": 2,
				"ls": 1,
			},
			"2": {
				"fp": 10,
				"lp": 12,
				"bt": 1,
				"ls": 1,
			},
			"3": {
				"fp": 13,
				"lp": 14,
				"bt": 1,
				"ls": 2,
			},
			"4": {
				"fp": 15,
				"lp": 17,
				"bt": 1,
				"ls": 1,
			},
			"5": {
				"fp": 18,
				"lp": 19,
				"bt": 1,
				"ls": 2,
			},
		},
		"esbor": nil,
		"ebbor": nil,
		"es": map[string][]int{
			"1": {
				5,
				6,
				7,
			},
			"2": {
				10,
				11,
				12,
			},
			"3": {
				13,
				14,
			},
			"4": {
				15,
				16,
				17,
			},
			"5": {
				18,
				19,
			},
		},
		"eb": map[string]map[string]int{
			"1": {
				"fp": 5,
				"lp": 7,
				"bt": 2,
				"ls": 1,
			},
			"2": {
				"fp": 10,
				"lp": 12,
				"bt": 1,
				"ls": 1,
			},
			"3": {
				"fp": 13,
				"lp": 14,
				"bt": 1,
				"ls": 2,
			},
			"4": {
				"fp": 15,
				"lp": 17,
				"bt": 1,
				"ls": 1,
			},
			"5": {
				"fp": 18,
				"lp": 19,
				"bt": 1,
				"ls": 2,
			},
		},
		"ssaw":  0,
		"pr":    nil,
		"tptbr": nil,
		"sc":    0,
		"rs":    nil,
		"fs":    nil,
		"csr":   nil,
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
			4,
			1,
			10,
			7,
			0,
			0,
			0,
			11,
			10,
			2,
			2,
			2,
			3,
			3,
			2,
			2,
			2,
			3,
			3,
			0,
			9,
			12,
			7,
			7,
			12,
			5,
			1,
			11,
			6,
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
		"bl":   22934.19,
		"tb":   0,
		"tbb":  0,
		"tw":   0,
		"np":   0,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
}
