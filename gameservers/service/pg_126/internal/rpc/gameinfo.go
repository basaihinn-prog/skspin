package rpc

import (
	"encoding/json"
	"log/slog"
	"serve/comm/lazy"
	"serve/comm/redisx"

	"serve/comm/db"

	"serve/comm/define"

	"serve/comm/mux"

	"serve/comm/slotsmongo"

	"serve/comm/ut"

	"serve/service/pg_126/internal/models"
)

// https://api.pg-demo.com/game-api/piggy-gold/v2/GameInfo/Get?traceId=CLJSNX12
func init() {
	mux.RegRpc("/game-api/126/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type M = map[string]any
type D = []any

var Cs = []float64{0.2, 2.0, 20.0, 100.0}
var Ml = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var Mxl = 5

func gameinfo(plr *models.Player, ps define.PGParams, ret *M) (err error) {
	//s := M{"fb": nil, "wt": M{"mw": 5.0, "bw": 10.0, "mgw": 20.0, "smgw": 35.0}, "maxwm": nil, "cs": D{1.0, 10.0, 100.0, 500.0}, "ml": D{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "mxl": 1, "bl": 549.00, "inwe": false, "iuwe": false, "cc": "PGC"}

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
	curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
	s := M{
		"fb": nil,
		"wt": M{
			"mw":   5.0,
			"bw":   20.0,
			"mgw":  35.0,
			"smgw": 50.0,
		},
		"maxwm": nil,
		"cs":    ut.FloatArrMul(c, curItem.Multi),
		"ml":    Ml,
		"mxl":   Mxl,
		"bl":    0.00,
		"inwe":  false,
		"iuwe":  false,
		"cc":    curItem.Key,
	}

	gold, err := slotsmongo.GetBalance(ps.Pid)
	if err != nil {
		return
	}
	balance := ut.Gold2Money(gold)
	// s["ls"] = nil
	s["bl"] = balance
	// ls := M{"si": M{"wp": nil, "lw": nil, "frl": D{6, 4, 5, 4, 0, 2, 6, 3, 4}, "pc": nil, "wm": 10, "tnbwm": nil, "gwt": -1, "fb": nil, "ctw": 0.0, "pmt": nil, "cwc": 0, "fstc": nil, "pcwc": 0, "rwsp": nil, "hashr": "0:4;0;3#MV#50#MT#1#MG#0#", "ml": 5, "cs": 10.0, "rl": D{4, 0, 3}, "sid": "1767749575281610240", "psid": "1767749575281610240", "st": 1, "nst": 1, "pf": 1, "aw": 0.00, "wid": 0, "wt": "C", "wk": "0_C", "wbn": nil, "wfg": nil, "blb": balance + 50, "blab": balance, "bl": balance, "tb": 50.00, "tbb": 50.00, "tw": 0.00, "np": -50.00, "ocr": nil, "mr": nil, "ge": D{1, 11}}}

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
		"wc":    3,
		"ist":   false,
		"itw":   false,
		"fws":   0,
		"wp":    nil,
		"orl":   nil,
		"lw":    nil,
		"irs":   false,
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
		"rl": D{
			2,
			2,
			5,
			0,
			0,
			0,
			6,
			3,
			3,
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
		"tb":   0.00,
		"tbb":  0.00,
		"tw":   0.00,
		"np":   0.00,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
}
