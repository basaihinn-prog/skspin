package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs10bbsplxmas/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`msi=12&def_s=6,7,4,2,8,9,8,5,6,7,8,6,7,3,9&msr=17&balance=1,000,000.00&cfgs=4469&ver=2&mo_s=11&index=1&balance_cash=1,000,000.00&reel_set_size=3&def_sb=10,11,9,6,8&mo_v=25,50,75,125,200,250,300,375,450,500,625,750,875&def_sa=7,5,4,4,3&reel_set=0&balance_bonus=0.00&na=s&scatters=1~0,0,1,0,0~0,0,8,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&stime=1733735201358&sa=7,5,4,4,3&sb=10,11,9,6,8&sc=0.04,0.08,0.20,0.40,1.00,2.00,4.00,12.00,20.00&defc=0.40&sh=3&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=0.40&sver=5&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;800,50,10,0,0;175,30,5,0,0;150,25,5,0,0;125,25,5,0,0;100,20,5,0,0;100,10,5,0,0;100,10,5,0,0;100,10,5,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0&l=25&rtp=96.50&reel_set0=5,7,6,4,10,9,4,5,7,6,9,7,11,11,11,11,10,8,4,7,6,9,5,10,6,8,3,4,10,5,7,9,6,10,4,9,5,8,9,6,10,3,7,6,10,8,9,10,6,9,5,3,10~7,2,2,2,2,2,4,8,5,9,1,10,6,9,4,7,8,11,11,11,11,10,6,9,4,8,3,9,5,10,4,9,1,7,5,9,3,8,5,4,9,8,3,7,4,8,5,9,6,10,1,8~3,2,2,2,2,2,8,5,9,1,10,4,8,7,3,8,6,5,8,4,7,1,8,5,9,6,8,9,2,2,2,8,10,6,8,5,7,1,8,6,9,11,11,11,11,8,6,3,9,5,8,10~10,2,2,2,2,5,6,4,7,10,3,7,11,11,11,11,9,5,7,3,10,1,7,6,9,4,7,8,9,3,7,4,10,6,7,1,8,7,3,9,5,7,8,6,7,9,6,7,4,5,8,3,10,4,9,1,7,3,5,8,6,9,7,4,10,3,9,1~7,2,2,6,4,3,9,5,8,6,4,7,6,10,7,6,12,9,4,5,7,6,10,5,7,4,10,6,7,4,9,6,8,3,6,9,5,3,10,6,7,5,4,8,9,6,10,4,7,6,8,4,9,12,10,6,7,4,10,9,5,6,10,8,4,7,8,10,4,8,5,9,6,10,9,4,7,8&s=6,7,4,2,8,9,8,5,6,7,8,6,7,3,9&reel_set2=18,18,18,18~18,18,18,18~18,18,18,18~18,18,18,18~7,2,2,2,4,8,6,10,3,9,5,12,3,4,7,6,10,9,5,10,7,6,9,4,3,7,6,5,10,7,12,10,6,7,4,9,6,8,7,6,9,5,3,10,6,7,5,3,8,9,12,10,4,7,3,9,4,8,9,4,10,6,5,3,10,9,5,6,9&reel_set1=5,7,6,3,10,9,4,5,7,6,9,7,11,11,11,11,11,9,5,10,6,8,3,4,10,5,7,9,6,10,4,9,5,8,9,6,10,3,7,6,8,10,9,6,11,11,11,11,11,7,4,10~7,2,2,2,2,2,2,2,2,2,3,8,5,9,1,10,6,9,4,11,11,11,11,11,5,10,6,9,4,8,3,9,5,10,4,9,1,7,5,9,11,11,11,11,11,8,3,7,4,8,5,9,6,10,1,8~8,2,2,2,2,2,2,2,2,2,2,2,5,9,1,10,4,8,7,9,8,6,5,8,4,7,1,8,5,9,11,11,11,11,11,4,10,8,6,5,8,7,1,8,6,9,11,11,11,11,11,6,8,9,5,8,3~10,2,2,2,2,2,2,2,2,2,2,2,2,2,2,10,3,7,11,11,11,11,11,5,7,3,10,1,7,6,9,4,7,8,9,3,7,4,10,6,7,1,8,7,3,9,5,7,8,6,11,11,11,11,11,11,5,8,10,4,9,1,7,3,5,8,6,9,7,4,10,3,9,1~7,2,2,8,6,10,3,9,5,12,3,4,7,6,10,7,6,12,9,4,8,7,6,12,5,7,4,10,6,7,12,9,6,8,3,6,9,5,4,10,6,7,12,4,8,9,6,10,12,7,6,8,4,9,12,10,6,7,3,10,9,12,6,10,4,8,7,12,10,4,8,5,9,6,12,10,4,7,8`)
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
