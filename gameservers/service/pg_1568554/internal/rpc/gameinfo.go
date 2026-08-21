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
	"serve/service/pg_1568554/internal/models"
)

func init() {
	mux.RegRpc("/game-api/1568554/v2/GameInfo/Get", "gameinfo", "game-api", db.WrapRpcPlayer(gameinfo), nil)
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
		"fb":    map[string]any{"is": true, "bm": 75, "t": 250 * curItem.Multi},
		"wt":    define.M{"mw": 5.0, "bw": 20.0, "mgw": 35.0, "smgw": 50.0},
		"maxwm": 2500,
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

		"trl": []int{
			6,
			6,
			12,
			2,
			4,
			0,
			0,
			0,
			0,
			8,
			1,
			5,
			13,
			10,
			0,
		},
		"twbm":   0,
		"rtwbm":  0,
		"trtwbm": 0,
		"rtw":    0,
		"trtw":   0,
		"wp":     nil,
		"trwp":   nil,
		"wpl":    nil,
		"trwpl":  nil,
		"lw":     nil,
		"trlw":   nil,
		"snww":   nil,
		"trsnww": nil,
		"pswp":   0,
		"swp":    99,
		"ptrswp": 0,
		"trswp":  99,
		"ptbr":   nil,
		"trptbr": nil,
		"wfrp":   nil,
		"wftrp":  nil,
		"fswp":   nil,
		"nrswp":  nil,
		"rsmw":   false,
		"orl": []int{
			6,
			6,
			12,
			2,
			4,
			0,
			0,
			0,
			0,
			8,
			1,
			5,
			13,
			10,
			0,
		},
		"otrl": []int{
			6,
			6,
			12,
			2,
			4,
			0,
			0,
			0,
			0,
			8,
			1,
			5,
			13,
			10,
			0,
		},
		"rns":  nil,
		"trns": nil,
		"ssaw": 0,
		"sc":   0,
		"gm":   1,
		"gml": []int{
			1,
			2,
			3,
			5,
		},
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
			6,
			6,
			12,
			2,
			4,
			0,
			0,
			0,
			0,
			8,
			1,
			5,
			13,
			10,
			0,
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
}
