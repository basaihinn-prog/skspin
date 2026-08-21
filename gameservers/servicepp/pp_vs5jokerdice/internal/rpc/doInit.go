package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs5jokerdice/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=8,7,4,9,8,6,7,4,9,8,3,7,7,6,6&balance=0.00&cfgs=6952&reel1=8,3,6,6,6,7,7,7,4,4,4,3,4,7,5,1,7,8,6,9,9,9,8,9,9,4,4&ver=2&reel0=8,8,8,7,7,7,4,6,6,6,5,7,6,1,4,3,6,8,9,8,7,9&index=1&balance_cash=0.00&def_sb=9,9,9,5,6&def_sa=8,8,8,7,7&reel3=6,6,6,9,9,9,7,6,6,8,8,8,4,8,3,5,3,9,9,9,6,6,5,4,8,1&reel2=4,9,9,9,9,8,8,8,3,6,8,9,5,4,7,8,5,6,1,7,6,8,7,8&reel4=9,9,9,5,6,6,6,3,9,6,9,3,5,6,7,8,8,8,6,3,1,8,6,8,4,4&balance_bonus=0.00&na=s&scatters=1~250,50,10,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&gameInfo={}&stime=1739154471302&sa=8,8,8,7,7&sb=9,9,9,5,6&sc=40.00,80.00,120.00,160.00,200.00,250.00,300.00,400.00,600.00,800.00,1000.00,1200.00,1400.00,1600.00,1800.00,2000.00,4000.00,5000.00,6000.00,8000.00,10000.00,12000.00,14000.00,16000.00,18000.00,20000.00,40000.00,60000.00,80000.00,100000.00,120000.00,140000.00,160000.00,180000.00,200000.00,240000.00&defc=200.00&sh=3&wilds=2~0,0,0~1,1,1&bonuses=0&fsbonus=&c=200.00&sver=5&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;5000,1000,100,0,0;1000,200,50,0,0;1000,200,50,0,0;200,50,20,0,0;200,50,20,0,0;200,40,20,0,0;200,40,20,5,0&l=5&s=8,7,4,9,8,6,7,4,9,8,3,7,7,6,6`)
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
