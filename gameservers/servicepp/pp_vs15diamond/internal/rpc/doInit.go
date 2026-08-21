package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs15diamond/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=6,7,4,2,8,4,3,5,6,7,8,5,7,3,4&bgid=0&balance=16,617.86&cfgs=2188&ver=2&index=1&balance_cash=16,617.86&reel_set_size=2&def_sb=7,4,7,5,7&def_sa=3,6,6,8,8&bonusInit=[{bgid:0,bgt:18,bg_i:"1000,100,30,10",bg_i_mask:"pw,pw,pw,pw"}]&balance_bonus=0.00&na=s&scatters=1~0,0,0,0,0~0,0,8,0,0~1,1,1,1,1&gmb=0,0,0&bg_i=1000,100,30,10&rt=d&stime=1737367114607&bgt=18&sa=3,6,6,8,8&sb=7,4,7,5,7&sc=0.06,0.12,0.18,0.30,0.60,1.20,1.80,2.40,4.00,8.00,15.00,30.00,60.00,120.00,150.00,180.00&defc=0.60&sh=3&wilds=2~300,60,20,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=0.60&sver=5&n_reel_set=0&bg_i_mask=pw,pw,pw,pw&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;200,20,10,0,0;100,20,10,0,0;40,10,5,0,0;40,10,5,0,0;40,10,5,0,0;40,10,5,0,0&l=15&reel_set0=6,5,7,8,6,2,2,2,7,3,1,4,7,8,4,7,2,2,7,6,3,6,7,6,4,6~6,8,4,8,5,5,3,8,2,2,8,3,8,8,5,4,5,8,4,6,8,3,5,4,5,2,2,2,7~7,3,6,5,2,2,2,8,5,6,5,2,2,5,4,5,5,7,1~6,4,3,6,8,4,7,6,5,6,2,2,2,7,5,8,6,7,6,4,6,8,4,2~4,5,6,5,6,4,5,1,7,5,2,2,2,7,5,7,1,5,8,7,4,5,6,7,5,7,3,7,5,7,3,5,4,7,4,5,7&s=6,7,4,2,8,4,3,5,6,7,8,5,7,3,4&reel_set1=4,8,3,7,6,4,6,1,4,6,2,2,2,7,4,1,5,7,8,5,7,4,6,8,2,2,5,3~3,4,6,4,3,8,5,2,2,3,6,4,4,5,2,2,2,6,7,8,4,8,5~7,2,2,6,4,3,1,6,2,2,2,3,4,5,4,5,8,7,4,5,6,8~8,6,3,6,4,2,2,2,4,8,4,5,7,6,2,2,5,5,8,6,4,7~1,8,1,6,2,2,2,5,4,8,7,4,5,4,7,4,1,5,6,8,7,8,2,2,3,6`)
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
