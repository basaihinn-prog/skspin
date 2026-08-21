package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs20tweethouse/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=9,3,2,3,9,10,4,1,4,10,9,3,2,3,9&balance=27,050.00&cfgs=7709&ver=2&index=1&balance_cash=27,050.00&reel_set_size=2&def_sb=8,10,8,11,7&def_sa=8,9,10,11,12&balance_bonus=0.00&na=s&scatters=1~0,0,5,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&mbri=1,2,3&rt=d&gameInfo={rtps:{regular:"96.51"}}&stime=1737624853192&sa=8,9,10,11,12&sb=8,10,8,11,7&sc=10.00,20.00,30.00,40.00,50.00,60.00,80.00,100.00,150.00,200.00,250.00,300.00,350.00,400.00,450.00,500.00,750.00,1000.00,1500.00,2000.00,2500.00,3000.00,4000.00,4500.00,5000.00,7500.00,10000.00,15000.00,20000.00,25000.00,30000.00,35000.00,40000.00,45000.00,50000.00,60000.00&defc=50.00&sh=3&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=50.00&sver=5&n_reel_set=0&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;750,150,50,0,0;500,100,35,0,0;300,60,25,0,0;200,40,20,0,0;150,25,12,0,0;100,20,8,0,0;50,10,5,0,0;50,10,5,0,0;25,5,2,0,0;25,5,2,0,0;25,5,2,0,0&l=20&reel_set0=9,8,12,8,10,7,5,11,4,1,3,7,10,13,1,6,9,13,6,11,12~3,6,8,13,7,10,9,11,10,9,6,5,12,2,4,8,11,12,13,7~4,9,13,12,13,7,8,12,6,1,2,10,11,7,5,11,3,10,8,9,6~2,6,10,7,11,13,12,5,9,3,6,7,12,9,13,8,10,11,4,8~8,11,7,6,13,9,10,5,1,12,6,3,8,4,7,10,13,12,11,9&s=9,3,2,3,9,10,4,1,4,10,9,3,2,3,9&reel_set1=12,5,11,9,13,8,13,3,3,3,10,12,11,10,13,11,8,8,9,6,9,10,12,6,3,7,4,7,5~13,11,7,9,4,12,7,3,10,9,8,13,11,10,13,5,6,9,2,7,6,10,12,8,11~6,12,10,13,7,12,5,10,8,7,2,13,3,6,9,8,11,8,5,12,9,4,11,10,9,13~13,9,5,7,13,6,12,11,6,10,13,12,9,7,8,10,4,2,8,7,5,9,11,3,12,8,6,10,11~13,12,11,7,10,11,7,13,4,9,12,6,10,3,3,3,8,6,11,8,9,13,7,9,5,8,12&mbr=1,1,1`)
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
