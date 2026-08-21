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
	"serve/service/pg_1671262/internal/models"
)

func init() {
	mux.RegRpc("/game-api/1671262/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type D = []any

var Cs = []float64{0.05, 0.5, 2.5, 10.0}
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
		"fb":    define.M{"is": true, "bm": 75, "t": 200.0 * curItem.Multi},
		"wt":    define.M{"mw": 5.0, "bw": 20.0, "mgw": 35.0, "smgw": 50.0},
		"maxwm": 8000,
		"cs":    ut.FloatArrMul(c, curItem.Multi),
		"ml":    Ml,
		"mxl":   Mxl,
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
		"ssaw": 0,
		"orl": []int{
			9,
			9,
			9,
			6,
			6,
			2,
			2,
			2,
			8,
			8,
			10,
			1,
			10,
			5,
			5,
			10,
			1,
			10,
			5,
			5,
			2,
			2,
			2,
			8,
			8,
			9,
			9,
			9,
			6,
			6,
		},
		"cgm":  0,
		"gm":   1,
		"rns":  nil,
		"twbm": 0,
		"crtw": 0,
		"ir":   false,
		"imw":  false,
		"sc":   0,
		"mf": map[string]map[string]int{
			"v": {
				"0": 0,
				"1": 50,
				"2": 0,
				"3": 10,
				"4": 0,
				"5": 2,
			},
			"t": {
				"0": -1,
				"1": 2,
				"2": -1,
				"3": 1,
				"4": -1,
				"5": 0,
			},
			"l": {
				"0": -1,
				"1": 3,
				"2": -1,
				"3": 3,
				"4": -1,
				"5": 3,
			},
		},
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
		"rl": []int{
			9,
			9,
			9,
			6,
			6,
			2,
			2,
			2,
			8,
			8,
			10,
			1,
			10,
			5,
			5,
			10,
			1,
			10,
			5,
			5,
			2,
			2,
			2,
			8,
			8,
			9,
			9,
			9,
			6,
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
