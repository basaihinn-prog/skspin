package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs25goldparty/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=3,4,5,6,7,8,9,10,3,4,5,6,7,8,9&balance=1,000,000.00&screenOrchInit={type:"mini_slots"}&cfgs=8018&ver=2&index=1&balance_cash=1,000,000.00&def_sb=4,5,6,7,8&reel_set_size=2&def_sa=4,5,6,7,8&reel_set=0&balance_bonus=0.00&na=s&scatters=1~0,0,0,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&gameInfo={rtps:{purchase:"96.50",regular:"96.50"},props:{max_rnd_sim:"1",max_rnd_hr:"20408163",jp1:"125000",max_rnd_win:"5163",jp3:"1250",jp2:"5000",jp4:"500"}}&wl_i=tbm~5163&stime=1733884287731&sa=4,5,6,7,8&sb=4,5,6,7,8&sc=0.04,0.08,0.12,0.16,0.20,0.32,0.40,0.80,1.20,1.60,2.00,3.00,4.00,8.00,12.00,16.00,20.00&defc=0.32&purInit_e=1&sh=3&wilds=2~0,0,0,0,0,0,0~1,1,1,1,1,1,1&bonuses=0&fsbonus=&st=rect&c=0.32&sw=5&sver=5&g={gp:{def_s:"3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3",reel_set:"1",sh:"6",st:"rect",sw:"10"}}&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;100,40,20,5,0;50,25,15,0,0;30,20,10,0,0;30,15,5,0,0;30,15,5,0,0;25,15,5,0,0;25,15,5,0,0;20,10,5,0,0;20,10,5,0,0;15,10,5,0,0;15,10,5,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0&l=25&total_bet_max=50,000.00&reel_set0=13,15,15,15,11,9,4,14,10,8,14,14,14,3,6,3,3,3,7,12,15,5,14,3,14,15,4,3,14,5,3,12,3,10,11,14,6,9,3,5,3,14,9,15,10,9,3,14,3,9,15,3,6,3,14,3,9,15,6,9,5,6,15,14,7,15,3,14,9,15,14,10,11,8,3,15,3,14,10,3,8,9,15,11,12,6,3,9,3,15,7,3,5,12,3,14,6,7,11,3,6,10,4,9,15,3,5,15,10,15,3,14,3,10,14,3~15,15,15,14,9,3,13,4,3,3,3,6,14,14,14,10,2,2,2,12,2,15,8,11,5,7,14,2,10,2,7,4,2,13,10,14,12,2,9,14,4,7,2,14,2,7,12,2,10,3,14,3,2,14,2,7,14,2,3,14,2,14,13,3,2,7,2,3,14,6,4,13,5,2,14,7,12,7,12,14,13,2,8,12,14,6,14,7,2,7,3,14,13,14,2~6,2,2,2,5,2,8,10,13,12,9,15,15,15,11,14,14,14,4,14,7,15,3,3,3,3,13,7,14,3,15,13,2,7,4,5,7,2,15,7,13,14,3,7,9,15,3,7,2,15,13,7,3,15,2,15,14,13,5,7,9,15,2,8,5,15,13,5,10,2,14,2,3,13,14,3,13,7,4,3,13,14,15,2,3,4,5,15,2,14,12,14,7,15,7,2,3,14,9,13,14,2,10,15,14,7,3,14,4,13,5,3,13,15,5,8,9,13,2,13,15,2,9,14,3,13,12,11,2,9,15,4,14~16,15,15,15,10,2,6,3,3,3,5,3,14,14,14,14,4,2,2,2,15,12,9,7,11,13,8,2,15,2,3,2,15,4,5,4,15,11,3,14,2,15,9,2,5,15,14,4,12,2,4,3,9~10,10,10,16,15,13,8,5,7,15,15,15,12,11,3,3,3,6,10,14,14,14,4,14,3,9,3,15,11,4,3,15,3,12,15,4,8,16,7,12,3,15,7,3,15,9,12,14,15,4,7,3,14,12,9,12,3,12,15,3,15,3,12,7,15,4,3,11,7,3,11,3,7,3,11,4,11,3,15,4,3,11,3,16,15,3,16,3,14,3,12,14,11,14,3,4,7,3,15,3,9,15,3,9,11,3,14&s=3,4,5,6,7,8,9,10,3,4,5,6,7,8,9&reel_set1=5,11,4,8,5,6,10,4,11,15,15,8,6,5,9,6,17,3,10,6,9,5,4,8,15,15,15,11,6,9,7,4,11,3,7,8,18,6,7,11,5,3,10,7,9,4,5,9,15,10,7,3,8&purInit=[{type:"d",bet:2500}]&total_bet_min=0.04`)
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
