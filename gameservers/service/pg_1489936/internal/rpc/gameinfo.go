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
	"serve/service/pg_1489936/internal/models"
)

func init() {
	mux.RegRpc("/game-api/1489936/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"fb":    define.M{"is": true, "bm": 75, "t": 200.0 * curItem.Multi},
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
		"wp":    nil,
		"lw":    nil,
		"bwp":   nil,
		"snww":  nil,
		"now":   4320,
		"nowpr": []any{6, 2, 4, 5, 3, 6},
		"es": define.M{
			"1": []any{6, 7, 8},
			"2": []any{9, 10, 11},
			"3": []any{12, 13, 14},
			"4": []any{21, 22},
			"5": []any{24, 25},
			"6": []any{27, 28, 29},
		},
		"eb": define.M{
			"1": define.M{"fp": 6, "lp": 8, "bt": 1, "ls": 1},
			"2": define.M{"fp": 9, "lp": 11, "bt": 1, "ls": 2},
			"3": define.M{"fp": 12, "lp": 14, "bt": 2, "ls": 1},
			"4": define.M{"fp": 21, "lp": 22, "bt": 2, "ls": 1},
			"5": define.M{"fp": 24, "lp": 25, "bt": 1, "ls": 2},
			"6": define.M{"fp": 27, "lp": 29, "bt": 1, "ls": 1},
		},
		"esb": define.M{
			"1": []any{6, 7, 8},
			"2": []any{9, 10, 11},
			"3": []any{12, 13, 14},
			"4": []any{21, 22},
			"5": []any{24, 25},
			"6": []any{27, 28, 29},
		},
		"ebb": define.M{
			"1": define.M{
				"fp": 6,
				"lp": 8,
				"bt": 1,
				"ls": 1,
			},
			"2": define.M{
				"fp": 9,
				"lp": 11,
				"bt": 1,
				"ls": 2,
			},
			"3": define.M{
				"fp": 12,
				"lp": 14,
				"bt": 2,
				"ls": 1,
			},
			"4": define.M{
				"fp": 21,
				"lp": 22,
				"bt": 2,
				"ls": 1,
			},
			"5": define.M{
				"fp": 24,
				"lp": 25,
				"bt": 1,
				"ls": 2,
			},
			"6": define.M{
				"fp": 27,
				"lp": 29,
				"bt": 1,
				"ls": 1,
			},
		},
		"ptbr":  nil,
		"orl":   []any{8, 4, 11, 11, 1, 12, 2, 2, 2, 3, 3, 3, 0, 0, 0, 6, 11, 5, 2, 0, 9, 1, 1, 7, 4, 4, 0, 5, 5, 5, 12, 5, 11, 9, 9, 1},
		"ssaw":  0.00,
		"sc":    3,
		"gm":    1,
		"omf":   []any{1, 2, 3, 5},
		"mf":    []any{1, 2, 3, 5},
		"mi":    0,
		"twbm":  0.0,
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
		"rl":    []any{8, 4, 11, 11, 1, 12, 2, 2, 2, 3, 3, 3, 0, 0, 0, 6, 11, 5, 2, 0, 9, 1, 1, 7, 4, 4, 0, 5, 5, 5, 12, 5, 11, 9, 9, 1},
		"sid":   "0",
		"psid":  "0",
		"st":    1,
		"nst":   1,
		"pf":    0,
		"aw":    0.00,
		"wid":   0,
		"wt":    "C",
		"wk":    "0_C",
		"wbn":   nil,
		"wfg":   nil,
		"blb":   0.00,
		"blab":  0.00,
		"bl":    0,
		"tb":    0.00,
		"tbb":   0.00,
		"tw":    0.00,
		"np":    0.00,
		"ocr":   nil,
		"mr":    nil,
		"ge":    nil,
	}
}
