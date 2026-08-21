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
	"serve/service/pg_20/internal/models"
)

func init() {
	mux.RegRpc("/game-api/20/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"aw":    0,
		"bl":    10000.0,
		"blab":  0,
		"blb":   0,
		"cs":    DefaultCs,
		"ctw":   0,
		"cwc":   0,
		"fb":    nil,
		"fs":    nil,
		"fstc":  nil,
		"ge":    nil,
		"gm":    nil,
		"gtf":   nil,
		"gwt":   0,
		"hashr": nil,
		"lw":    nil,
		"lwbm":  nil,
		"ml":    DefaultBetLevel,
		"mr":    nil,
		"np":    0,
		"nst":   1,
		"ocr":   nil,
		"orl":   D{2, 2, 2, 5, 1, 8, 3, 3, 3, 7, 0, 6, 4, 4, 4},
		"pcwc":  0,
		"pf":    0,
		"pmt":   nil,
		"psid":  "0",
		"rl":    D{2, 2, 2, 5, 1, 8, 3, 3, 3, 7, 0, 6, 4, 4, 4},
		"rwsp":  nil,
		"sc":    0,
		"sid":   "0",
		"sp":    nil,
		"st":    1,
		"tb":    0,
		"tbb":   0,
		"tw":    0,
		"wa":    0,
		"wbn":   nil,
		"wfg":   nil,
		"wid":   0,
		"wk":    "0_C",
		"wp":    nil,
		"wt":    "C",
	}
}
