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
	"serve/service/pg_1543462/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	mux.RegRpc("/game-api/1543462/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"maxwm": 5000,
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
		s["ls"] = bson.M{"si": si}
	}

	*ret = s
	// *ret = json.RawMessage(jsonstr)
	return
}
func getInitSi(DefaultCs float64, DefaultBetLevel int64) map[string]any {
	si := map[string]any{
		"wp": nil,
		"lw": nil,
		"orl": []int{
			2,
			2,
			0,
			99,
			8,
			8,
			8,
			8,
			2,
			2,
			0,
			99,
		},
		"ift": false,
		"iff": false,
		"cpf": map[string]map[string]any{
			"1": {
				"p":  4,
				"bv": 5000,
				"m":  500,
			},
			"2": {
				"p":  5,
				"bv": 200,
				"m":  20,
			},
			"3": {
				"p":  6,
				"bv": 50,
				"m":  5,
			},
			"4": {
				"p":  7,
				"bv": 5,
				"m":  0.5,
			},
		},
		"cptw":  0,
		"crtw":  0,
		"imw":   false,
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
			2,
			2,
			0,
			99,
			8,
			8,
			8,
			8,
			2,
			2,
			0,
			99,
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
		"bl":   0,
		"tb":   0,
		"tbb":  0,
		"tw":   0,
		"np":   0,
		"ocr":  nil,
		"mr":   nil,
		"ge":   nil,
	}
	return si
}
