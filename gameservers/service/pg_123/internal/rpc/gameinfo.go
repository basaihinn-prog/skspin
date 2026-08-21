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
	"serve/service/pg_123/internal/models"
)

func init() {
	mux.RegRpc("/game-api/123/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"wp":   nil,
		"lw":   nil,
		"otw":  0.0,
		"ssaw": 0.0,
		"snww": nil,
		"wpl":  nil,
		"ptbr": nil,
		"rns":  nil,
		"orl":  nil,
		"sc":   0,
		"fs":   nil,
		"swlb": map[string]map[string]any{
			"1": {
				"p":  6,
				"ls": 3,
			},
			"2": {
				"p":  21,
				"ls": 3,
			},
		},
		"swl": map[string]map[string]any{
			"1": {
				"p":  6,
				"ls": 3,
			},
			"2": {
				"p":  21,
				"ls": 3,
			},
		},
		"gm":    1,
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
			5,
			2,
			13,
			5,
			2,
			0,
			2,
			12,
			7,
			1,
			6,
			4,
			11,
			7,
			1,
			6,
			4,
			11,
			5,
			3,
			0,
			3,
			12,
			8,
			5,
			3,
			13,
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
		"bl":   60000.00,
		"tb":   0.00,
		"tbb":  0.00,
		"tw":   0.00,
		"np":   0.00,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
}
