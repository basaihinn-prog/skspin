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
	"serve/service/pg_1648578/internal/models"
)

type M = map[string]any
type D = []any

func init() {
	mux.RegRpc("/game-api/1648578/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"wp":   nil,
		"sw":   nil,
		"wsc":  nil,
		"wpl":  nil,
		"ssaw": 0.00,
		"orl": D{
			8,
			8,
			8,
			18,
			3,
			2,
			9,
			5,
			5,
			9,
			1,
			3,
			11,
			11,
			11,
			2,
			3,
			2,
			10,
			10,
			8,
			1,
			7,
			5,
			5,
			6,
			10,
			10,
			9,
			9,
		},
		"rns": D{
			D{},
			D{
				2,
			},
			D{},
			D{
				2,
			},
			D{
				8,
				1,
			},
			D{},
		},
		"crtw": 0,
		"imw":  false,
		"gm":   1,
		"ptbr": nil,
		// "pfmf": {},
		// "fmf":  {},
		"nfmf": nil,
		"mfmf": M{
			"3": M{
				"mv": 11,
				"s":  18,
			},
		},
		"cf": D{
			false,
			false,
			false,
			false,
			false,
			false,
		},
		"icp":  false,
		"ctbb": 0.8,
		"sc":   2,
		"imff": true,
		"fs":   nil,
		"gwt":  -1,
		"fb":   nil,
		"ctw":  0,
		"pmt":  nil,
		"cwc":  0,
		"fstc": M{
			"4": 5,
		},
		"pcwc":  0,
		"rwsp":  nil,
		"hashr": "5:8;2;1;2;8;6#8;9;3;3;1;10#8;5;11;2;7;10#18;5;11;10;5;9#3;9;11;10;5;9#MV#0#MT#1#MG#0.00#",
		"ml":    DefaultBetLevel,
		"cs":    DefaultCs,
		"rl": D{
			8,
			8,
			8,
			18,
			3,
			2,
			9,
			5,
			5,
			9,
			1,
			3,
			11,
			11,
			11,
			2,
			3,
			2,
			10,
			10,
			8,
			1,
			7,
			5,
			5,
			6,
			10,
			10,
			9,
			9,
		},
		"sid":  "0",
		"psid": "0",
		"st":   4,
		"nst":  1,
		"pf":   2,
		"aw":   0.00,
		"wid":  0,
		"wt":   "C",
		"wk":   "0_C",
		"wbn":  nil,
		"wfg":  nil,
		"blb":  0.0,
		"blab": 0.0,
		"bl":   0.0,
		"tb":   0,
		"tbb":  0.0,
		"tw":   0,
		"np":   0,
		"ocr":  nil,
		"mr":   nil,
		"ge": D{
			1,
			3,
			11,
		},
	}
}
