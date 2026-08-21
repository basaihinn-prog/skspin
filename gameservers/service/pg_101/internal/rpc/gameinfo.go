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

	"serve/service/pg_101/internal/models"
)

// https://api.pg-demo.com/game-api/piggy-gold/v2/GameInfo/Get?traceId=CLJSNX12
func init() {
	mux.RegRpc("/game-api/101/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"bl":    10000.00,
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
	// ls := M{"si": M{"wp": nil, "lw": nil, "frl": D{6, 4, 5, 4, 0, 2, 6, 3, 4}, "pc": nil, "wm": 10, "tnbwm": nil, "gwt": -1, "fb": nil, "ctw": 0.0, "pmt": nil, "cwc": 0, "fstc": nil, "pcwc": 0, "rwsp": nil, "hashr": "0:4;0;3#MV#50#MT#1#MG#0#", "ml": 5, "cs": 10.0, "rl": D{4, 0, 3}, "sid": "176710195101281610240", "psid": "176710195101281610240", "st": 1, "nst": 1, "pf": 1, "aw": 0.00, "wid": 0, "wt": "C", "wk": "0_C", "wbn": nil, "wfg": nil, "blb": balance + 50, "blab": balance, "bl": balance, "tb": 50.00, "tbb": 50.00, "tw": 0.00, "np": -50.00, "ocr": nil, "mr": nil, "ge": D{1, 11}}}

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
		"wp":   nil,
		"bwp":  nil,
		"lw":   nil,
		"snww": nil,
		"ssaw": 0,
		"twbm": 0,
		"now":  1920,
		"nowpr": D{
			4,
			2,
			4,
			5,
			3,
			4,
		},
		"ptbr": nil,
		"orl": D{
			8,
			4,
			4,
			4,
			1,
			12,
			2,
			2,
			2,
			3,
			3,
			3,
			0,
			0,
			0,
			6,
			11,
			5,
			2,
			0,
			9,
			1,
			1,
			7,
			4,
			4,
			0,
			5,
			5,
			5,
			12,
			5,
			11,
			11,
			11,
			1,
		},
		"esb": M{
			"1": D{
				1,
				2,
				3,
			},
			"2": D{
				6,
				7,
				8,
			},
			"3": D{
				9,
				10,
				11,
			},
			"4": D{
				12,
				13,
				14,
			},
			"5": D{
				21,
				22,
			},
			"6": D{
				24,
				25,
			},
			"7": D{
				27,
				28,
				29,
			},
			"8": D{
				32,
				33,
				34,
			},
		},
		"ebb": M{
			"1": M{
				"fp": 1,
				"lp": 3,
				"bt": 2,
				"ls": 1,
			},
			"2": M{
				"fp": 6,
				"lp": 8,
				"bt": 1,
				"ls": 1,
			},
			"3": M{
				"fp": 9,
				"lp": 11,
				"bt": 1,
				"ls": 2,
			},
			"4": M{
				"fp": 12,
				"lp": 14,
				"bt": 2,
				"ls": 1,
			},
			"5": M{
				"fp": 21,
				"lp": 22,
				"bt": 2,
				"ls": 1,
			},
			"6": M{
				"fp": 24,
				"lp": 25,
				"bt": 1,
				"ls": 2,
			},
			"7": M{
				"fp": 27,
				"lp": 29,
				"bt": 1,
				"ls": 1,
			},
			"8": M{
				"fp": 32,
				"lp": 34,
				"bt": 2,
				"ls": 1,
			},
		},
		"es": M{
			"1": D{
				1,
				2,
				3,
			},
			"2": D{
				6,
				7,
				8,
			},
			"3": D{
				9,
				10,
				11,
			},
			"4": D{
				12,
				13,
				14,
			},
			"5": D{
				21,
				22,
			},
			"6": D{
				24,
				25,
			},
			"7": D{
				27,
				28,
				29,
			},
			"8": D{
				32,
				33,
				34,
			},
		},
		"eb": M{
			"1": M{
				"fp": 1,
				"lp": 3,
				"bt": 2,
				"ls": 1,
			},
			"2": M{
				"fp": 6,
				"lp": 8,
				"bt": 1,
				"ls": 1,
			},
			"3": M{
				"fp": 9,
				"lp": 11,
				"bt": 1,
				"ls": 2,
			},
			"4": M{
				"fp": 12,
				"lp": 14,
				"bt": 2,
				"ls": 1,
			},
			"5": M{
				"fp": 21,
				"lp": 22,
				"bt": 2,
				"ls": 1,
			},
			"6": M{
				"fp": 24,
				"lp": 25,
				"bt": 1,
				"ls": 2,
			},
			"7": M{
				"fp": 27,
				"lp": 29,
				"bt": 1,
				"ls": 1,
			},
			"8": M{
				"fp": 32,
				"lp": 34,
				"bt": 2,
				"ls": 1,
			},
		},
		"wbwp": nil,
		"fmbwf": M{
			"m":    1,
			"bwsk": nil,
			"bwsp": nil,
		},
		"rs":    nil,
		"fs":    nil,
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
			4,
			4,
			1,
			12,
			2,
			2,
			2,
			3,
			3,
			3,
			0,
			0,
			0,
			6,
			11,
			5,
			2,
			0,
			9,
			1,
			1,
			7,
			4,
			4,
			0,
			5,
			5,
			5,
			12,
			5,
			11,
			11,
			11,
			1,
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
