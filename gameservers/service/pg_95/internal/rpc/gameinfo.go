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
	"serve/service/pg_95/internal/models"
)

func init() {
	mux.RegRpc("/game-api/95/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type D = []any

var Cs = []float64{0.1, 1, 10, 50}
var Ml = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var Mxl = 10

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
		"wp": nil,
		// "bwp": nil,
		// "lw":  nil,
		// "lwa": 0,
		"orl": []int{
			7, 7, 8, 7, 6, 6, 6, 6, 4, 8, 3, 7, 3, 9, 4, 8, 7, 4, 2, 7, 8, 8, 6, 4, 8},
		"ogs": []int{
			0, 14, 17},
		// "now": 1920,
		// "nowpr": []int{
		// 	4,
		// 	2,
		// 	4,
		// 	5,
		// 	3,
		// 	4},
		// "snww": nil,

		"ssaw": 0,
		// "rs":   nil,
		"fs": nil,
		// "pr": nil,
		"sc": 0,
		// "mi": 0,
		// "ma": []int{
		// 	1,
		// 	2,
		// 	3,
		// 	4,
		// 	5,
		// },
		"gwt":   -1,
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
			7, 7, 8, 7, 6, 6, 6, 6, 4, 8, 3, 7, 3, 9, 4, 8, 7, 4, 2, 7, 8, 8, 6, 4, 8},
		"sid":  "0",
		"psid": "0",
		"st":   1,
		"nst":  1,
		"pf":   2,
		"aw":   0,
		"wid":  0,
		"wt":   "C",
		"wk":   "0_C",
		"wbn":  nil,
		"wfg":  nil,
		"blb":  0,
		"blab": 0,
		"bl":   0,
		"tb":   1,
		"tbb":  1,
		"tw":   0,
		"np":   -1,
		"ocr":  nil,
		"ogml": 1,
		"mr":   nil,
		"ge":   []int{1, 11},
		"gs":   []int{0, 14, 17},
		"gsd":  nil,
	}
}
