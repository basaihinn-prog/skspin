package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs1dragon8/internal"
	"time"

	"serve/comm/lazy"
	"serve/servicepp/ppcomm"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/slotsmongo"
	"serve/comm/ut"

	"github.com/nats-io/nats.go"
)

func init() {
	ppcomm.RegRpc("doInit", doInit)
}

var (
	sampleDoInit = ppcomm.ParseVariables(`def_s=3,6,5,6,3,6,4,6,4&balance=16,668.86&cfgs=2191&reel1=3,6,5,6,3,6,4,6,5,6,5,6,3,6,3,6,3,6&ver=2&reel0=5,6,3,6,3,6,4,6,4,6,4,6,5,6&index=1&balance_cash=16,668.86&def_sb=6,5,4&def_sa=6,3,3&reel2=4,6,4,6,3,6,4,6,5,6,5,6,5,6&balance_bonus=0.00&na=s&scatters=1~0,0,0~0,0,0~1,1,1&gmb=0,0,0&rt=d&stime=1737344180727&sa=6,3,3&sb=6,5,4&sc=1.00,2.00,3.00,5.00,10.00,20.00,30.00,45.00,75.00,150.00,300.00,450.00,750.00,1500.00,2250.00,3000.00&defc=10.00&sh=3&wilds=2~0,0,0~1,1,1&bonuses=0&fsbonus=&c=10.00&sver=5&counter=2&ntp=0.00&paytable=0,0,0;0,0,0;0,0,0;100,0,0;50,0,0;25,0,0;0,0,0&l=1&s=3,6,5,6,3,6,4,6,4`)
)

var cGrade = []float64{0.05, 0.10, 0.15, 0.25, 0.50, 1.00, 1.50, 2.00, 3.50, 7.50, 15.00, 25.00, 40.00, 75.00, 110.00, 150.00}

//var betGrade = []float64{0.05, 0.10}

func doInit(msg *nats.Msg) (ret []byte, err error) {
	// action=doInit&symbol=vs20olympx&cver=237859&index=1&counter=1&repeat=0&mgckey=AUTHTOKEN@7a4f4399cd1bd51a2eef10551a63ce1fdb69a05e0363abf9d5c47c3f5da92f13~stylename@hllgd_hollygod~SESSION@d5ec7d86-a060-4470-a302-0192a2d53f14~SN@54345f72

	ps := ppcomm.ParseVariables(string(msg.Data))
	//fmt.Println(ps)

	var pid int64
	pid, err = jwtutil.ParseToken(ps.Str("mgckey"))
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *ppcomm.Player) error {
		info := maps.Clone(sampleDoInit)
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)

		isEnd, _ := plr.IsEndO()
		c, err := redisx.GetPlayerCs(plr.AppID, plr.PID, isEnd)
		if err != nil {
			return err
		}
		minBet := c[0]
		maxBet := c[len(c)-1] * 100 * internal.Line

		info.SetFloat("c", minBet*curItem.Multi)
		for k, v := range plr.LastData {
			info.SetStr(k, v)
		}
		info.SetFloat("total_bet_max", maxBet*curItem.Multi)
		info.SetFloat("total_bet_min", minBet*curItem.Multi)
		info.SetInt("stime", int(time.Now().UnixMilli()))
		info.SetFloatArr("sc", ut.FloatArrMul(c, curItem.Multi))
		balance, err := slotsmongo.GetBalance(pid)
		if err != nil {
			return err
		}
		info.SetFloat("balance", ut.Gold2Money(balance))
		info.SetFloat("balance_cash", ut.Gold2Money(balance))

		ret = info.Bytes()
		plr.SpinCountOfThisEnter = 0
		return nil
	})

	return
}
