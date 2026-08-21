package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs5aztecgems/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=7,4,9,7,4,9,7,4,9&balance=1,000,000.00&cfgs=3677&reel1=4,4,4,9,9,9,7,7,7,7,7,7,2,2,2,8,8,8,9,7,7,4,7,7,5,5,5,5,5,6,6,6,9,2,3,3,3,9,9,9,8,5,8,9,5,3,6,4&ver=2&reel0=8,8,8,8,7,7,7,6,6,6,6,6,8,8,4,4,4,4,8,6,5,5,5,7,8,8,6,9,9,9,9,9,8,6,8,4,3,3,3,3,9,8,9,4,3,2,2,2,9,5,2,2,6,6,7&index=1&balance_cash=1,000,000.00&def_sb=4,4,4&def_sa=8,8,8&reel2=4,4,4,4,8,8,8,9,9,9,9,7,7,7,7,5,5,5,3,3,3,7,7,3,8,9,5,3,8,4,9,7,9,5,5,9,9,6,6,6,6,7,2,2,2,9,2&balance_bonus=0.00&na=s&aw=3&scatters=1~0,0,0~0,0,0~1,1,1&gmb=0,0,0&rt=d&base_aw=m~1;m~2;m~3;m~5;m~10;m~15&stime=1733822014888&sa=8,8,8&sb=4,4,4&sc=0.04,0.08,0.20,0.40,1.00,2.00,4.00,12.00,20.00&defc=0.40&def_aw=3&sh=3&wilds=2~25,0,0~1,1,1&bonuses=0&fsbonus=&c=0.40&sver=5&counter=2&ntp=0.00&paytable=0,0,0;0,0,0;0,0,0;20,0,0;15,0,0;12,0,0;10,0,0;8,0,0;5,0,0;2,0,0&l=5&rtp=96.52&s=7,4,9,7,4,9,7,4,9&awt=6rl`)
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
