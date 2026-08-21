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
	"serve/service/pg_60/internal/models"
)

func init() {
	mux.RegRpc("/game-api/60/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type M = map[string]any
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
		"aw":   0,
		"bl":   0,
		"blab": 0,
		"blb":  0,
		"bwp":  nil,
		"cs":   DefaultCs,
		"ctw":  0,
		"cwc":  0,
		"eb": map[string]map[string]any{
			"1": {"fp": 6, "lp": 8, "bt": 1, "ls": 1},
			"2": {"fp": 9, "lp": 11, "bt": 1, "ls": 2},
			"3": {"fp": 24, "lp": 25, "bt": 1, "ls": 2},
			"4": {"fp": 27, "lp": 29, "bt": 1, "ls": 1},
		},
		"ebb": map[string]map[string]any{
			"1": {"fp": 6, "lp": 8, "bt": 1, "ls": 1},
			"2": {"fp": 9, "lp": 11, "bt": 1, "ls": 2},
			"3": {"fp": 24, "lp": 25, "bt": 1, "ls": 2},
			"4": {"fp": 27, "lp": 29, "bt": 1, "ls": 1},
		},
		"es": M{
			"1": D{6, 7, 8},
			"2": D{9, 10, 11},
			"3": D{24, 25},
			"4": D{27, 28, 29},
		},
		"esb": M{
			"1": D{6, 7, 8},
			"2": D{9, 10, 11},
			"3": D{24, 25},
			"4": D{27, 28, 29},
		},
		"fb":    nil,
		"fs":    nil,
		"fstc":  nil,
		"ge":    nil,
		"gwt":   0,
		"hashr": nil,
		"lw":    nil,
		"ml":    DefaultBetLevel,
		"mr":    nil,
		"now":   7776,
		"nowpr": D{6, 2, 6, 6, 3, 6},
		"np":    0,
		"nst":   1,
		"ocr":   nil,
		"orl":   D{8, 4, 7, 10, 1, 12, 2, 2, 2, 3, 3, 3, 0, 0, 0, 6, 11, 1, 2, 0, 9, 12, 7, 7, 4, 4, 0, 5, 5, 5, 12, 5, 11, 11, 6, 1},
		"pcwc":  0,
		"pf":    0,
		"pmt":   nil,
		"psid":  "0",
		"ptbr":  nil,
		"rl":    D{8, 4, 7, 10, 1, 12, 2, 2, 2, 3, 3, 3, 0, 0, 0, 6, 11, 1, 2, 0, 9, 12, 7, 7, 4, 4, 0, 5, 5, 5, 12, 5, 11, 11, 6, 1},
		"rs":    nil,
		"rwsp":  nil,
		"sid":   "0",
		"snww":  nil,
		"ssaw":  0,
		"st":    1,
		"tb":    0,
		"tbb":   0,
		"tw":    0,
		"wbn":   nil,
		"wfg":   nil,
		"wid":   0,
		"wk":    "0_C",
		"wp":    nil,
		"wt":    "C",
	}
}
