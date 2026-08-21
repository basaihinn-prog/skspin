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
	"serve/service/pg_1635221/internal/models"
)

type M = map[string]any
type D = []any

func init() {
	mux.RegRpc("/game-api/1635221/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

var Cs = []float64{0.05, 0.5, 2.5, 10}
var Ml = []float64{1,
	2,
	3,
	4,
	5,
	6,
	7,
	8,
	9,
	10}
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
		"fb": M{
			"is": true,
			"bm": 75,
			"t":  100.0,
		},
		"wt": define.M{
			"mw":   5.0,
			"bw":   20.0,
			"mgw":  35.0,
			"smgw": 50.0,
		},
		"maxwm": 5000,
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
		"lw": nil,
		"sc": 2,
		"wc": 0,
		"orl": D{
			3,
			11,
			3,
			3,
			8,
			8,
			8,
			2,
			2,
			11,
			5,
			1,
			10,
			6,
			6,
			2,
			1,
			12,
			10,
			3,
			6,
			8,
		},
		"wm":   M{},
		"lml":  0,
		"rml":  0,
		"gm":   13,
		"snww": nil,
		"wpl":  nil,
		"ptbr": nil,
		"crtw": 0,
		"cc":   0,
		"rns":  nil,
		"twbm": 0,
		"ssaw": 0,
		"imw":  false,
		"fs": M{
			"s":    0,
			"ts":   15,
			"as":   0,
			"aw":   5.82,
			"ssaw": 0,
			"rc":   0,
		},
		"gwt": -1,
		"fb":  2,
		"ctw": 0,
		"pmt": nil,
		"cwc": 0,
		"fstc": M{
			"21": 15,
			"22": 9,
		},
		"pcwc":  0,
		"rwsp":  nil,
		"hashr": "24:3;8;11;6;10#11;8;5;6;3#3;8;1;2;6#3;2;10;1;8#2;12#MV#0#MT#13#MG#0#",
		"ml":    DefaultBetLevel,
		"cs":    DefaultCs,
		"rl": D{
			3,
			11,
			3,
			3,
			8,
			8,
			8,
			2,
			2,
			11,
			5,
			1,
			10,
			6,
			6,
			2,
			1,
			12,
			10,
			3,
			6,
			8,
		},
		"sid":  "0",
		"psid": "0",
		"st":   21,
		"nst":  1,
		"pf":   2,
		"aw":   0,
		"wid":  0,
		"wt":   "C",
		"wk":   "0_C",
		"wbn":  nil,
		"wfg":  nil,
		"blb":  0,
		"blab": 9965.86,
		"bl":   0.00,
		"tb":   0.00,
		"tbb":  0.00,
		"tw":   0.00,
		"np":   0.00,
		"ocr":  nil,
		"mr":   nil,
		"ge": D{
			1,
			3,
			11,
		},
	}
}
