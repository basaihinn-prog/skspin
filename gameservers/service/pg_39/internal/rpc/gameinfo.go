package rpc

import (
	"encoding/json"
	"serve/comm/lazy"
	"serve/comm/redisx"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_39/internal/models"
)

// https://api.pg-demo.com/game-api/piggy-gold/v2/GameInfo/Get?traceId=CLJSNX12
func init() {
	mux.RegRpc("/game-api/39/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type M = map[string]any
type D = []any

var Cs = []float64{1.0, 10.0, 100.0, 500.0}
var Ml = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var Mxl = 1

func gameinfo(plr *models.Player, ps define.PGParams, ret *M) (err error) {
	gold, err := slotsmongo.GetBalance(ps.Pid)
	if err != nil {
		return
	}
	info, err := redisx.LoadAppIdCache(plr.AppID)
	if err != nil {
		return err
	}
	//修改end
	slotsmongo.UpdateEnterPlrCount(lazy.ServiceName, plr.AppID, plr.PID)
	plr.SpinCountOfThisEnter = 0
	balance := ut.Gold2Money(gold)
	curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
	//1.0, 10.0, 100.0, 500.0
	s := M{
		"fb":    nil,
		"wt":    M{"mw": 5.0, "bw": 10.0, "mgw": 20.0, "smgw": 35.0},
		"maxwm": nil,
		"cs":    ut.FloatArrMul(info.Cs, curItem.Multi),
		"ml":    Ml,
		"mxl":   Mxl,
		"bl":    549.00,
		"inwe":  false,
		"iuwe":  false,
		"cc":    curItem.Key,
		"gc":    info.ShowNameAndTimeOff,
		"ign":   info.ShowNameAndTimeOff,
		"asc":   info.StopLoss,
	}
	// s["ls"] = nil
	s["bl"] = balance
	// ls := M{"si": M{"wp": nil, "lw": nil, "frl": D{6, 4, 5, 4, 0, 2, 6, 3, 4}, "pc": nil, "wm": 10, "tnbwm": nil, "gwt": -1, "fb": nil, "ctw": 0.0, "pmt": nil, "cwc": 0, "fstc": nil, "pcwc": 0, "rwsp": nil, "hashr": "0:4;0;3#MV#50#MT#1#MG#0#", "ml": 5, "cs": 10.0, "rl": D{4, 0, 3}, "sid": "1767749575281610240", "psid": "1767749575281610240", "st": 1, "nst": 1, "pf": 1, "aw": 0.00, "wid": 0, "wt": "C", "wk": "0_C", "wbn": nil, "wfg": nil, "blb": balance + 50, "blab": balance, "bl": balance, "tb": 50.00, "tbb": 50.00, "tw": 0.00, "np": -50.00, "ocr": nil, "mr": nil, "ge": D{1, 11}}}

	if len(plr.LS) != 0 {
		s["ls"] = json.RawMessage(plr.LS)
	} else {

		si := getInitSi(info.DefaultCs, info.DefaultBetLevel)
		si["bl"] = balance

		s["ls"] = M{"si": si}
	}

	*ret = s
	// *ret = json.RawMessage(jsonstr)
	return
}

func getInitSi(DefaultCs float64, DefaultBetLevel int64) M {
	si := M{
		"wp":    nil,
		"lw":    nil,
		"frl":   D{3, 1, 4, 4, 2, 2, 3, 1, 4},
		"pc":    nil,
		"wm":    nil,
		"tnbwm": nil,
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
		"rl":    D{1, 2, 1},
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
	return si
}
