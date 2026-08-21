package rpc

import (
	"encoding/json"
	"log/slog"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"

	"serve/comm/define"

	"serve/comm/mux"

	"serve/comm/slotsmongo"

	"serve/comm/ut"

	"serve/service/pg_79/internal/models"
)

// https://api.pg-demo.com/game-api/piggy-gold/v2/GameInfo/Get?traceId=CLJSNX12
func init() {
	mux.RegRpc("/game-api/79/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type M = map[string]any
type D = []any

var Cs = []float64{0.05, 0.5, 2.5, 10.0}
var Ml = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var Mxl = 20

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
		"fb": M{
			"is": true,
			"bm": 75,
			"t":  150.0 * curItem.Multi,
		},
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
		"bl":    10254.45,
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
	// ls := M{"si": M{"wp": nil, "lw": nil, "frl": D{6, 4, 5, 4, 0, 2, 6, 3, 4}, "pc": nil, "wm": 10, "tnbwm": nil, "gwt": -1, "fb": nil, "ctw": 0.0, "pmt": nil, "cwc": 0, "fstc": nil, "pcwc": 0, "rwsp": nil, "hashr": "0:4;0;3#MV#50#MT#1#MG#0#", "ml": 5, "cs": 10.0, "rl": D{4, 0, 3}, "sid": "1767799579281610240", "psid": "1767799579281610240", "st": 1, "nst": 1, "pf": 1, "aw": 0.00, "wid": 0, "wt": "C", "wk": "0_C", "wbn": nil, "wfg": nil, "blb": balance + 50, "blab": balance, "bl": balance, "tb": 50.00, "tbb": 50.00, "tw": 0.00, "np": -50.00, "ocr": nil, "mr": nil, "ge": D{1, 11}}}

	if len(plr.LS) != 0 {
		s["ls"] = json.RawMessage(plr.LS)
	} else {
		si := getInitSi(info.DefaultCs, info.DefaultBetLevel)
		s["ls"] = M{"si": si}
	}
	*ret = s
	//*ret = json.RawMessage(jsonstr)
	return
}

func getInitSi(DefaultCs float64, DefaultBetLevel int64) map[string]any {
	return M{
		"sc":   0,
		"wbwp": nil,
		"wp":   nil,
		"twp":  nil,
		"lw":   nil,
		"trl": D{
			3,
			6,
			2,
			4,
		},
		"torl": nil,
		"bwp":  nil,
		"now":  7200,
		"nowpr": D{
			5,
			3,
			6,
			6,
			4,
			5,
		},
		"snww": nil,
		"esb": M{
			"1": D{
				5,
				6,
				7,
			},
			"2": D{
				8,
				9,
			},
			"3": D{
				10,
				11,
				12,
			},
			"4": D{
				20,
				21,
			},
			"5": D{
				23,
				24,
			},
		},
		"ebb": M{
			"1": M{
				"fp": 5,
				"lp": 7,
				"bt": 1,
				"ls": 1,
			},
			"2": M{
				"fp": 8,
				"lp": 9,
				"bt": 1,
				"ls": 2,
			},
			"3": M{
				"fp": 10,
				"lp": 12,
				"bt": 1,
				"ls": 3,
			},
			"4": M{
				"fp": 20,
				"lp": 21,
				"bt": 1,
				"ls": 2,
			},
			"5": M{
				"fp": 23,
				"lp": 24,
				"bt": 1,
				"ls": 1,
			},
		},
		"es": M{
			"1": D{
				5,
				6,
				7,
			},
			"2": D{
				8,
				9,
			},
			"3": D{
				10,
				11,
				12,
			},
			"4": D{
				20,
				21,
			},
			"5": D{
				23,
				24,
			},
		},
		"eb": M{
			"1": M{
				"fp": 5,
				"lp": 7,
				"bt": 1,
				"ls": 1,
			},
			"2": M{
				"fp": 8,
				"lp": 9,
				"bt": 1,
				"ls": 2,
			},
			"3": M{
				"fp": 10,
				"lp": 12,
				"bt": 1,
				"ls": 3,
			},
			"4": M{
				"fp": 20,
				"lp": 21,
				"bt": 1,
				"ls": 2,
			},
			"5": M{
				"fp": 23,
				"lp": 24,
				"bt": 1,
				"ls": 1,
			},
		},
		"fs":   nil,
		"rs":   nil,
		"ssaw": 0,
		"pr":   nil,
		"tpr":  nil,
		"orl": D{
			8,
			4,
			1,
			10,
			7,
			2,
			2,
			2,
			3,
			3,
			0,
			0,
			0,
			11,
			10,
			0,
			9,
			12,
			7,
			7,
			4,
			4,
			0,
			5,
			5,
			12,
			5,
			1,
			11,
			6,
		},
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
		"rl": D{
			8,
			4,
			1,
			10,
			7,
			2,
			2,
			2,
			3,
			3,
			0,
			0,
			0,
			11,
			10,
			0,
			9,
			12,
			7,
			7,
			4,
			4,
			0,
			5,
			5,
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
		"bl":   10000,
		"tb":   0,
		"tbb":  0,
		"tw":   0,
		"np":   0,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
}
