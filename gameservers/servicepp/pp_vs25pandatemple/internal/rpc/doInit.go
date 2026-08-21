package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs25pandatemple/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=6,7,4,3,8,4,3,5,6,7,8,5,7,3,4&balance=20,000.00&cfgs=4165&ver=2&index=1&balance_cash=20,000.00&def_sb=4,13,4,7,8&reel_set_size=3&def_sa=10,3,5,3,7&reel_set=0&balance_bonus=0.00&na=s&scatters=1~100,15,2,0,0~12,12,12,0,0~1,1,1,1,1&gmb=0,0,0&bg_i=3,0,2,1,1,2,2,3,1,4,3,5,2,6,1,7,2,8,1,9,5,10,15,20,25,50,75,100,150,200,250,500,1000,2500,4998,10,5,10,15,20,25,50,75,100,150,200,250,500,1000,2500,4998,11&rt=d&gameInfo={props:{max_rnd_sim:"1",max_rnd_hr:"1045296",max_rnd_win:"5000"}}&wl_i=tbm~5000&stime=1739167070397&sa=10,3,5,3,7&sb=4,13,4,7,8&sc=8.00,16.00,24.00,32.00,40.00,50.00,60.00,80.00,120.00,160.00,200.00,240.00,280.00,320.00,360.00,400.00,600.00,800.00,1200.00,1600.00,2000.00,2800.00,3200.00,3600.00,4000.00,6000.00,8000.00,12000.00,16000.00,20000.00,24000.00,28000.00,32000.00,36000.00,40000.00,48000.00&defc=40.00&sh=3&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=40.00&sver=5&bg_i_mask=pw,ic,pw,ic,pw,ic,pw,ic,pw,ic,pw,ic,pw,ic,pw,ic,pw,ic,pw,ic,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,ic,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,pw,ic&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;200,50,25,0,0;150,50,10,0,0;100,20,5,0,0;100,20,5,0,0;100,20,5,0,0;50,15,5,0,0;50,15,5,0,0;50,10,5,0,0;50,10,5,0,0;50,10,5,0,0;50,10,5,0,0&l=25&reel_set0=12,6,10,4,8,4,10,6,4,6,10,4,6,10,6,4,6,10,4~13,11,7,3,5,9,7,3,7,3,9,7,9,11,5,3,9,5,7,9,3,7,9,7,9,3,9,7,9,5,7,9,3,7,3~11,3,4,9,7,8,6,12,5,10,13,5,9,8,4,10,4,9,8,13,12,8,12,10,8,13,10,8,10,8,13,4~9,7,6,13,5,8,3,11,10,4,12,5,10,4~13,11,10,4,1,8,7,6,12,5,3,2,9,2,6,7,2,9,3,8,3,5,11,6,5,3&s=6,7,4,3,8,4,3,5,6,7,8,5,7,3,4&reel_set2=13,12,7,10,8,4,9,3,6,1,5,11,8,5,4,10,3~2,2,2,6,3,7,13,10,5,4,1,11,8,9,12,2,7,10,7,8,10,7,9,8,7,8,7,5,10,5,10~2,2,2,6,13,7,12,4,3,9,5,11,8,10,1,2,3,1,8,5,8,10,1~8,2,2,2,4,7,11,3,10,5,1,9,2,6,12,13,6,2,3,5,2,5,3,5,6,2,10,1,11,6,2,6,10,5,4,2,1,2,7,2,3,2,3,5,6,2~1,8,4,12,2,2,2,10,7,13,5,6,3,2,11,9,2,6,8,9,2,3,7,2,3,12,3,6,13,3,4,7,5,6,2,4,6,9,7,6,12,2,7,9,2,9,5,3,10,3,4,13,2,6,3,9,7,5,3,2,10,2,12,3,9,4,10,3,9,2,4,10,2,6&reel_set1=12,11,3,10,4,7,13,8,6,9,1,5,7,9,11,8,4,7,1,13~6,4,13,10,2,2,2,8,2,1,3,5,9,12,11,7,8,2,8,12,2,10,12,2,11,2,12,3,9,12,2,11,2,9,5~2,2,2,7,8,12,5,1,6,13,11,9,2,10,4,3,5,11,12~2,2,2,7,8,9,10,11,5,6,4,12,3,13,2,1,13,4,3,6,7,12,9,6,5,4~7,5,6,13,1,12,2,4,9,10,11,3,8,9,6,10,12,13,12,11,6,10,11,10,2,11,12,9,11,12,10,11,6,12,9,11`)
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
