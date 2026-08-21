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
	"serve/service/pg_24/internal/models"
)

func init() {
	mux.RegRpc("/game-api/24/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type D = []any

var Cs = []float64{1, 5, 20, 100, 500}
var Ml = []float64{1}
var Mxl = 1

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
		"wt":    define.M{"mw": 0.0, "bw": 60.0, "mgw": 150.0, "smgw": 400.0},
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
		"aw":    0,
		"bl":    10000,
		"blab":  0,
		"blb":   0,
		"bn":    1,
		"cs":    DefaultCs,
		"ctw":   0,
		"cwc":   0,
		"fb":    nil,
		"frl":   D{4, 5, 6, 1, 2, 3, 6, 5, 4},
		"fstc":  nil,
		"ge":    nil,
		"gwt":   0,
		"hashr": nil,
		"lw":    nil,
		"ml":    DefaultBetLevel,
		"mp":    nil,
		"mr":    nil,
		"mv":    nil,
		"np":    0,
		"nst":   1,
		"ocr":   nil,
		"pcwc":  0,
		"pf":    0,
		"pmt":   nil,
		"pp":    nil,
		"pprl":  nil,
		"psid":  "0",
		"rl":    D{1, 2, 3},
		"rwsp":  nil,
		"sid":   "0",
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
