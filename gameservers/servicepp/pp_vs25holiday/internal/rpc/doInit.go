package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs25holiday/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=7,9,13,6,6,9,10,6,12,12,11,5,4,9,10&balance=20,000.00&cfgs=7803&ver=2&index=1&balance_cash=20,000.00&def_sb=5,7,10,3,3&reel_set_size=4&def_sa=8,8,9,10,11&reel_set=0&balance_bonus=0.00&na=s&scatters=1~0,0,3,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&gameInfo={props:{max_rnd_sim:"1",max_rnd_hr:"1000000",max_rnd_win:"5000",max_rnd_win_a:"5000"}}&wl_i=tbm~5000&bl=0&stime=1739166641487&sa=8,8,9,10,11&sb=5,7,10,3,3&sc=8.00,16.00,24.00,32.00,40.00,50.00,60.00,80.00,120.00,160.00,200.00,240.00,280.00,320.00,360.00,400.00,600.00,800.00,1200.00,1600.00,2000.00,2800.00,3200.00,3600.00,4000.00,6000.00,8000.00,12000.00,16000.00,20000.00,24000.00,28000.00,32000.00,36000.00,40000.00,48000.00&defc=40.00&purInit_e=1&sh=3&wilds=2~1000,250,50,0,0~1,1,1,1,1,1;14~1000,250,50,0,0~1,1,1,1,1,1;15~1000,250,50,0,0~1,1,1,1,1,1;16~1000,250,50,0,0~1,1,1,1,1,1&bonuses=0&fsbonus=&c=40.00&sver=5&bls=25,35&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;750,150,50,0,0;500,100,35,0,0;300,60,25,0,0;200,40,20,0,0;150,25,12,0,0;100,20,8,0,0;50,10,5,0,0;50,10,5,0,0;25,5,2,0,0;25,5,2,0,0;25,5,2,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0&l=25&total_bet_max=120,000,000.00&reel_set0=8,7,12,9,3,5,6,10,13,11,4,2,1,6,11,9,1,9,3,9,7,9,4,5,6,1,9,4,13,6,12,7,13,4,11,7,5,9,7,5,13,4,9,13,4,9,12,4,12,11,6,7,1,7,6,12,11,9,12,3,11,4,7,9~3,13,8,7,6,2,10,9,11,5,4,12,9,6,9~1,3,13,12,2,11,7,10,9,5,4,6,8,13,8,6,5,8,5,6,8,13,10,5,10,13,6,10,7,10,8,5,13,2,5,10,3,10,8,7,13,2,7,8,11,10,5,10,12,5~6,3,12,7,13,9,11,5,8,2,4,10,12,8,12,8,13,4,8,11,4,9,12,4,12~10,5,12,9,6,7,1,2,8,11,13,3,4,11,1,11,7,1,11,5,1,7,4,1,7,11,8,11,6,7,1,11,1,6,7,9,6,13,6,7,1,5,11,7,11,1,9,12,11,8,9,13,3,13,11,7,5,7,11,1,9,8&s=7,9,13,6,6,9,10,6,12,12,11,5,4,9,10&reel_set2=6,8,13,9,11,5,3,5,9,3,11,9,5,8,11,8,11,9,11,5,9,13,11,5,11,8,11,5,1~12,7,10,4,10~3,13,9,5,6,8,11,6,8,13,1~4,10,12,7,10~8,3,6,13,9,5,11,5,1&reel_set1=3,3,3,8,3,13,6,12,11,5,7,9,4,10,13,10,5,8,12,7,12,9,4,13,11,9,11,13,9,11,4,12,7,13,5,9,11,7,13,8,12,4,11~11,13,7,8,3,12,9,5,6,4,10,5,10,7,10,4,3,10,8,10,4,12,10,12,9,10,12,5,10,4,10,4,3,13,10,7,10~5,8,9,3,12,6,13,4,10,7,11,8,12,6,7,6,11,12~8,9,12,11,3,13,10,6,7,5,4,3,7,5,3,9,3,4,5,3,13,6,9,3,4,5,6,9,5,13,9,7~12,6,13,8,5,3,7,10,3,3,3,9,11,4,8,4,3,8,13,3,4,11,13,3,13,11,5,13,5,11,10,3,10,7,8,6,5,11,3,5,11,5,3,5,8,5,3,13,7,5,6,11,10,5,3,4,13,11,4,8,7,10&purInit=[{type:"fsbl",bet:2500,bet_level:0}]&reel_set3=6,13,9,7,11,5,3,4,1,12,2,10,8,10,2,7,4,10,4,7,5,10,4,2,4,8,2,11,10,4,3,4,8,1,8,12,10,11,10,7,2,4,7,9,8,7,10,1,4,1,3,8,4,1,8,2,7,10,8,2,4,10,2,4,2,4~9,8,11,3,2,12,6,7,13,10,4,5,11,5,10,11,7,8,6,2,7,5,4,7,11,4,10,4,6~5,6,11,2,8,7,13,9,10,1,4,3,12,11,6,10,8,7,3,4,8,11,6,7,10,13,3,9,1,8,6,7,6,2,13,9,11,6,7,13,6,13,6,7,6,1,10,6,7,10,13,6,13,6,2,7,3,7,3,6,13,7,1,6,2,8,7,10,8,13,10,13,9~3,2,6,11,7,12,10,8,9,13,4,5,12,9,6,10,9,11,4,7,11,9,13,4,7,11,12,13,10,11,13,12,11,7,13,6,12,10,9,12,9,12,13,4,11,12,11,7,4,6,9,12,7,12,10,12,8,10,12,9,12,13,4,7,12,11,12,2,4,12~11,13,10,7,9,3,12,1,2,6,5,8,4,6&total_bet_min=8.00`)
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
