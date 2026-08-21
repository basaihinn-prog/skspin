package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs10noodles/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=6,12,7,6,10,5,8,5,4,9,6,10,6,6,10&balance=137,380.00&cfgs=14526&ver=3&index=1&balance_cash=137,380.00&def_sb=3,8,11,7,9&reel_set_size=6&def_sa=5,12,11,7,10&reel_set=0&balance_bonus=0.00&na=s&scatters=1~0,0,0,0,0~0,0,0,0,0~1,1,1,1,1&rt=d&gameInfo={rtps:{ante:"96.58",purchase:"96.58",regular:"96.58"},props:{max_rnd_sim:"1",max_rnd_hr:"531872",max_rnd_win:"5100",max_rnd_win_a:"3400",max_rnd_hr_a:"305159"}}&wl_i=tbm~5100;tbm_a~3400&bl=0&stime=1733968350863&sa=5,12,11,7,10&sb=3,8,11,7,9&sc=20.00,40.00,60.00,80.00,100.00,120.00,150.00,200.00,300.00,400.00,500.00,600.00,700.00,800.00,900.00,1000.00,1500.00,2000.00,3000.00,4000.00,5000.00,6000.00,7000.00,8000.00,9000.00,10000.00,20000.00,30000.00,40000.00,50000.00,60000.00,70000.00,80000.00,90000.00,100000.00,120000.00&defc=100.00&purInit_e=1&sh=3&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&st=rect&c=100.00&sw=5&sver=5&bls=10,15&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;2000,200,50,5,0;1000,150,30,0,0;500,100,20,0,0;500,100,20,0,0;200,50,10,0,0;100,25,5,0,0;100,25,5,0,0;100,25,5,0,0;100,25,5,0,0;100,25,5,0,0;0,0,0,0,0&l=10&total_bet_max=120,000,000.00&reel_set0=6,12,1,3,11,12,5,7,8,3,4,9,5,6,8,3,9,10,7,5~6,8,11,3,4,9,3,5,6,10,1,12,10,7,4~10,6,9,7,8,4,11,6,4,1,12,5,3,12,11,12,3,5~4,12,8,9,4,12,3,6,5,11,5,7,10,7,4,12,1,11,7~11,12,1,3,4,5,6,7,8,9,10&s=6,12,7,6,10,5,8,5,4,9,6,10,6,6,10&reel_set2=12,3,5,6,12,8,5,1,3,9,7,8,5,9,3,4,11,7,4,10,5~11,12,3,4,5,6,7,8,9,10~11,5,12,6,4,6,10,9,10,7,5,12,4,12,3,8,4,3~12,11,5,4,5,10,11,3,12,7,11,7,8,1,5,12,4,9,4,6,7~11,12,3,4,5,6,7,8,9,10&reel_set1=6,5,11,3,7,5,6,12,3,9,4,1,7,5,10,8,9,12,5,3,8~11,12,3,4,5,6,7,8,9,10~10,8,12,3,4,10,5,4,11,9,5,6,1,6,3,7,10,12~11,12,3,4,5,6,7,8,9,10~11,12,3,4,5,6,7,8,9,10&reel_set4=11,12,3,4,5,6,7,8,9,10~11,12,1,3,4,5,6,7,8,9,10~11,12,3,4,5,6,7,8,9,10~11,12,1,3,4,5,6,7,8,9,10~11,12,3,4,5,6,7,8,9,10&purInit=[{bet:1000,type:"fs"}]&reel_set3=11,12,3,4,5,6,7,8,9,10~1,10,4,8,12,4,3,7,6,9,5,10,11,3,7,6~11,12,1,3,4,5,6,7,8,9,10~11,12,3,4,5,6,7,8,9,10~11,12,3,4,5,6,7,8,9,10&reel_set5=5,4,10,3,9,3,5,11,4,3,7,12,5,7,8,9,12,8,6~4,8,3,6,4,7,11,3,10,6,10,9,5,7,12~11,10,3,8,10,4,9,3,6,12,5,11,12,4,12,7,5,6~9,4,5,11,7,12,5,8,12,10,6,3,4,11,7,4,7,5~11,12,3,4,5,6,7,8,9,10&total_bet_min=20.00`)
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
