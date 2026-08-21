package rpc

import (
	"encoding/json"
	"log/slog"
	"serve/comm/lazy"
	"serve/comm/redisx"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_25/internal/models"
)

func init() {
	mux.RegRpc("/game-api/25/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
}

type D = []any

var Cs = []float64{0.03, 0.3, 1, 6}
var Ml = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var Mxl = 30

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
		"wt":    define.M{"mw": 2.5, "bw": 5, "mgw": 15, "smgw": 35},
		"maxwm": nil,
		"cs":    ut.FloatArrMul(c, curItem.Multi),
		"ml":    Ml,
		"mxl":   Mxl,
		"inwe":  false,
		"iuwe":  false,
		"cc":    curItem.Key,
		"bl":    10000.00,
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
		s["ls"] = bson.M{"si": si}
	}
	*ret = s
	return
}

func getInitSi(DefaultCs float64, DefaultBetLevel int64) map[string]any {
	si := map[string]any{
		"wp":    nil,
		"lw":    nil,
		"tlwa":  0.00,
		"rs":    nil,
		"cf":    nil,
		"bns":   nil,
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
			3,
			6,
			4,
			5,
			7,
			12,
			8,
			9,
			11,
			5,
			7,
			10,
			3,
			6,
			4,
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
		"bl":   0.00,
		"tb":   0.00,
		"tbb":  0.00,
		"tw":   0.00,
		"np":   0.00,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
	return si
}
