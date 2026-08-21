package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs9hotroll/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=8,4,8,3,8,3,8,4,8&c_paytable=9~explicit~5,2,5~5000;10~explicit~5,3,5~3000;11~explicit~5,4,5~2000;12~explicit~5,5,5~1000;13~any_with_added_wild_multiplier~6~100,0,0;14~any_with_added_wild_multiplier~7~75,0,0;15~any_with_added_wild_multiplier~6,7~50,0,0;16~any_with_added_wild_multiplier~8~30,0,0;17~any_with_added_wild_multiplier~9~20,0,0;18~any_with_added_wild_multiplier~10~10,0,0;19~any_with_added_wild_multiplier~8,9,10~5,0,0;20~anywhere_on_line~11~10,5,2&balance=20,000.00&cfgs=2457&reel1=5,6,12,9,12,7,12,11,12,8,12,10,12,8,12,8,12,6,12,9,12,7,12,9,12,10,12,4,12,9,12,8,12,9,12,10,12,7,12,9,12,6,2,6,12,9,12,8,12,9,12,9,12,8,12,10,12,9,12,3,12,9,12,7,12,8,12,6,12&ver=2&reel0=6,12,8,12,10,12,5,12,6,12,10,12,9,12,7,12,7,12,10,12,8,12,11,12,10,12,7,12,7,12,10,12,7,12,8,12,10,12,6,5,6,12,10,12,9,12,7,12,10,12,9,12,10,12,8,12&index=1&balance_cash=20,000.00&def_sb=8,9,10&def_sa=5,6,7&reel2=10,12,9,12,10,12,8,12,8,12,5,12,10,12,9,12,10,12,10,12,7,12,9,12,8,12,11,12,10,12,8,12,6,5,6,12,10,12,7,12,8,12,10,12,9,12,8,12,9,12,8,12,10,12,6,12,11,12,8,12,7,12,10,12,8,12,10,12,8,12,9,12&balance_bonus=0.00&na=s&scatters=1~0,0,0~0,0,0~1,1,1&gmb=0,0,0&rt=d&stime=1739759417458&sa=5,6,7&sb=8,9,10&sc=25.00,45.00,60.00,70.00,90.00,120.00,180.00,250.00,350.00,450.00,600.00,700.00,800.00,900.00,1000.00,1200.00,1700.00,2300.00,3400.00,4500.00,5600.00,6700.00,7800.00,8900.00,10000.00,11500.00,23000.00,34000.00,45000.00,56000.00,68000.00,79000.00,90000.00,100000.00,112000.00,140000.00&defc=90.00&sh=3&wilds=2~0,0,0~1,1,5;3~0,0,0~1,1,4;4~0,0,0~1,1,3;5~0,0,0~1,1,2&bonuses=0&fsbonus=&c=90.00&sver=5&counter=2&ntp=0.00&paytable=0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0;0,0,0&l=9&s=8,4,8,3,8,3,8,4,8`)
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
