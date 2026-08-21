package rpc

import (
	"encoding/json"
	"log/slog"
	"serve/comm/redisx"
	"serve/service/pg_1879752/internal/gendata"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lazy"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_1879752/internal/models"
)

func init() {
	mux.RegRpc("/game-api/1879752/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type D = []any

var Cs = []float64{1.25, 12.5, 37.5, 125}
var Ml = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var Mxl = 20

func gameinfo(plr *models.Player, ps define.PGParams, ret *define.M) (err error) {
	gold, err := slotsmongo.GetBalance(ps.Pid)
	if err != nil {
		return
	}
	App, err := redisx.LoadAppIdCache(plr.AppID)
	if err != nil {
		return err
	}
	isEnd, _ := plr.IsEndO()
	c, err := redisx.GetPlayerCs(plr.AppID, plr.PID, isEnd)
	if err != nil {
		slog.Error("doSpin", "GetPlayerCs", err)
		return err
	}
	slotsmongo.UpdateEnterPlrCount(lazy.ServiceName, plr.AppID, plr.PID)
	plr.SpinCountOfThisEnter = 0
	balance := ut.Gold2Money(gold)
	curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
	s := define.M{
		"fb":    map[string]any{"is": true, "bm": 75, "t": 150 * curItem.Multi},
		"wt":    define.M{"mw": 15.0, "bw": 25.0, "mgw": 50.0, "smgw": 100.0},
		"maxwm": nil,
		"cs":    ut.FloatArrMul(c, curItem.Multi),
		"ml":    Ml,
		"mxl":   gendata.Line,
		"bl":    10000.00,
		"inwe":  false,
		"iuwe":  false,
		"cc":    curItem.Key,
		"gcs": map[string]any{"ab": map[string]any{
			"is":  true,
			"bm":  0,
			"bms": []float64{1.5},
			"t":   0.0,
		},
		},
		"abm": []float64{1.5},
	}
	s["bl"] = balance

	if len(plr.LS) != 0 {
		s["ls"] = json.RawMessage(plr.LS)
	} else {
		si := getInitSi(App.DefaultCs, App.DefaultBetLevel)
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
		"irfs":  false,
		"itr":   false,
		"itff":  false,
		"it":    false,
		"ifsw":  true,
		"gm":    0,
		"crtw":  0.0,
		"imw":   false,
		"orl":   []float64{2, 2, 2, 99, 0, 0, 0, 0, 3, 3, 3, 99},
		"rcs":   0,
		"gwt":   0,
		"ctw":   0.0,
		"pmt":   nil,
		"cwc":   0,
		"fstc":  nil,
		"pcwc":  0,
		"rwsp":  nil,
		"hashr": nil,
		"fb":    nil,
		"ab":    nil,
		"ml":    DefaultBetLevel,
		"cs":    DefaultCs,
		"rl":    []float64{2, 2, 2, 99, 0, 0, 0, 0, 3, 3, 3, 99},
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
		"bl":    100000.00,
		"tb":    0.00,
		"tbb":   0.00,
		"tw":    0.00,
		"np":    0.00,
		"ocr":   nil,
		"mr":    nil,
		"ge":    nil,
	}
}
