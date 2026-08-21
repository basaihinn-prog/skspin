package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs10threestar/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`msi=10~11~12~13~14~15~16&def_s=3,4,5,6,7,8,9,8,7,6,5,4,3,4,5&balance=5,118,007.00&cfgs=5096&ver=2&index=1&balance_cash=5,118,007.00&reel_set_size=2&def_sb=3,4,5,6,7&def_sa=3,4,5,6,7&reel_set=0&balance_bonus=0.00&na=s&scatters=1~0,0,0,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&stime=1737429551928&sa=3,4,5,6,7&sb=3,4,5,6,7&sc=20.00,40.00,60.00,80.00,100.00,120.00,150.00,200.00,300.00,400.00,500.00,600.00,700.00,800.00,900.00,1000.00,1500.00,2000.00,3000.00,4000.00,5000.00,6000.00,7000.00,8000.00,9000.00,10000.00,20000.00,30000.00,40000.00,50000.00,60000.00,70000.00,80000.00,90000.00,100000.00,120000.00&defc=100.00&sh=3&wilds=2~0,0,0,0,0~1,1,1,1,1;17~0,0,0,0,0~1,1,1,1,1;18~0,0,0,0,0~1,1,1,1,1;19~0,0,0,0,0~1,1,1,1,1;20~0,0,0,0,0~1,1,1,1,1;21~0,0,0,0,0~1,1,1,1,1;22~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=100.00&sver=5&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;250,200,50,0,0;120,60,25,0,0;60,25,10,0,0;50,20,8,0,0;40,15,7,0,0;25,10,5,0,0;25,10,5,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0&l=10&rtp=96.27&reel_set0=9,9,9,9,9,6,6,6,6,6,7,7,9,9,5,5,5,7,6,6,7,7,7,6,6,6,5,5,8,9,9,9,4,4,6,6,6,3,3,3,6,6,6,9,9,7,7,7,3,6,8,8,6,6,4,4,4,7,7,7,9,9,9,8,8,8,8~3,3,3,5,5,4,4,4,8,8,8,6,6,6,6,8,8,5,5,5,4,6,6,7,7,7,8,8,8,17,5,5,8,6,6,8,8,8,8,3,9,9,9,5,5,5,9,6,6,7,7,6,6,6,7,8,8,8,5~8,8,8,9,9,4,4,4,4,4,8,8,8,8,5,5,5,9,6,9,9,9,4,9,7,7,7,7,7,4,4,4,7,7,7,5,7,3,18,6,6,5,5,5,3,3,3,3,8,8,7,7,7,7,8,8,8,6,6,6,8,7,7,9,9~5,5,4,4,4,5,5,5,3,3,4,4,8,4,9,8,3,3,5,5,3,3,3,3,3,4,4,4,8,8,8,19,3,3,3,5,5,9,9,6,6,8,8,7,7,7,7,5,9,9,9,8,8,8,4,4,4,5,5,7,8,8,8~6,6,9,9,7,7,6,7,7,6,6,6,9,9,9,9,7,7,7,4,9,7,7,9,9,7,4,4,6,6,6,5,7,7,8,8,8,5,5,5,7,7,9,9,9,6,6,6,7,7,7,3,3,3,8,8,8,6,6,9,9,9&s=3,4,5,6,7,8,9,8,7,6,5,4,3,4,5&reel_set1=10,10,10,10,10,10,10,10,10,10,15,10,10,15,10,10,10,10,10,15,10,15,10,10,15,10,10,15,10,10,15,10,10,10,15,10,10,10,15,10,10,10,10,10,10,10,10,10,15,10,10,10,10,10,15,10,10,10,10,15,10,10,10,10,15,10,10,15,15~15,11,11,15,15,11,11,11,15,11,11,15,11,11,15,11,11,15,11,11,11,11,15,11,11,11,11,11,11,11,11,11,15,11,11,11,11,11,17,11,11,11,11,11,15,11,15,11,11,15,11,11,15,11,11,15,11,11,11,11,11,11,11,11,11,11~12,12,15,12,15,12,12,15,12,12,18,12,12,12,12,12,12,12,12,12,12,15,12,12,12,12,15,12,12,15,12,12,15,12,12,12,12,12,15,12,12,12,12,12,12,12,12,12,15,15,12,12,12,15,12,12,12,15,12,12,15,12,12,15~13,13,16,13,19,13,13,16,13,13,13,16,13,13,13,13,13,13,13,16,13,13,13,13,13,16,13,13,16,13,13,13,13,13,13,16,16,13,13,13,13,16,13,13,16,13,13,13,13,13,13,13,13,13,16,13,13,13,16,16~14,14,14,16,14,16,16,14,14,16,16,14,14,16,14,14,14,16,14,14,14,14,14,14,16,14,14,14,14,14,14,14,14,14,16,16,14,16,14,14,14,14,14,14,14,16,14,14,14,16,14,14,14,14,14,14,14,14,14,14,16,14,14,14,14,16`)
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
