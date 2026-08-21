package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs10bxmasbnza/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=5,6,9,6,6,8,11,8,9,9,6,9,12,6,6&balance=100.00&cfgs=7839&ver=2&index=1&balance_cash=100.00&def_sb=10,7,7,8,11&reel_set_size=4&def_sa=11,10,7,10,4&reel_set=0&balance_bonus=0.00&na=s&scatters=1~0,0,0,0,0~20,15,10,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&gameInfo={rtps:{purchase:"96.71",regular:"96.71"},props:{max_rnd_sim:"1",max_rnd_hr:"3880481",max_rnd_win:"2100"}}&wl_i=tbm~2100&stime=1735616461576&sa=11,10,7,10,4&sb=10,7,7,8,11&sc=0.01,0.02,0.03,0.05,0.10,0.20,0.30,0.50,1.00,1.50,2.50,5.00,10.00,20.00,30.00,40.00&defc=0.10&purInit_e=1&sh=3&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=0.10&sver=5&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;2000,200,50,5,0;1000,150,30,0,0;500,100,20,0,0;500,100,20,0,0;200,50,10,0,0;100,25,5,0,0;100,25,5,0,0;100,25,5,0,0;100,25,5,0,0;100,25,5,0,0&l=10&total_bet_max=40,000.00&reel_set0=6,7,8,10,10,5,6,12,7,11,6,5,10,6,1,5,8,7,7,12,4,10,5,8,4,10,3,12,8,10,4,8,9,9,12,10,9,4,3,10,5,8,12,3,6,8,12,6,4,10,11,12,8,7,7,7,7,7~4,3,9,7,11,3,8,6,3,9,7,11,5,10,6,3,1,9,7,7,12,3,4,11,3,6,8,9,5,11,9,6,11,4,9,3,5,4,11,6,3,5,11,9,12,9,10,4,3,11,5,9,9,8,7,7,7,7,7~12,3,3,6,7,11,3,5,7,10,4,3,5,10,9,4,5,1,6,7,7,12,9,5,4,3,8,5,6,11,4,6,5,10,3,6,4,3,8,7,7,7,7,7~5,7,6,4,12,5,4,8,7,6,4,11,3,9,7,7,10,1,6,3,5,3,7,7,7,7,7~4,6,7,11,11,4,8,8,7,5,3,12,1,7,7,10,8,6,9,7,7,7,7,7&s=5,6,9,6,6,8,11,8,9,9,6,9,12,6,6&accInit=[{id:0,mask:"cp"},{id:1,mask:"cp;mp"}]&reel_set2=6,7,8,10,10,5,6,12,7,11,6,5,10,6,1,5,8,7,7,12,4,10,5,8,4,10,3,12,8,10,4,8,9,9,12,10,9,4,3,10,5,8,12,3,6,8,12,6,4,10,11,12,8,7,7,7,7,7~4,3,9,7,11,3,8,6,3,9,7,11,5,10,6,3,1,9,7,7,12,3,4,11,3,6,8,9,5,11,9,6,11,4,9,3,5,4,11,6,3,5,11,9,12,9,10,4,3,11,5,9,9,8,7,7,7,7,7~12,3,3,6,7,11,3,5,7,10,4,3,5,10,9,4,5,1,6,7,7,12,9,5,4,3,8,5,6,11,4,6,5,10,3,6,4,3,8,7,7,7,7,7~5,7,6,4,12,5,4,8,7,6,4,11,3,9,7,7,10,1,6,3,5,3,7,7,7,7,7~4,6,7,11,11,4,8,8,7,5,3,12,1,7,7,10,8,6,9,7,7,7,7,7&reel_set1=6,7,8,10,10,5,6,12,7,11,6,5,10,6,1,5,8,7,7,12,4,10,5,8,4,10,3,12,8,10,4,8,9,9,12,10,9,4,3,10,5,8,12,3,6,8,12,6,4,10,11,12,8,7,7,7,7,7~4,3,9,7,11,3,8,6,3,9,7,11,5,10,6,3,1,9,7,7,12,3,4,11,3,6,8,9,5,11,9,6,11,4,9,3,5,4,11,6,3,5,11,9,12,9,10,4,3,11,5,9,9,8,7,7,7,7,7~12,3,3,6,7,11,3,5,7,10,4,3,5,10,9,4,5,1,6,7,7,12,9,5,4,3,8,5,6,11,4,6,5,10,3,6,4,3,8,7,7,7,7,7~5,7,6,4,12,5,4,8,7,6,4,11,3,9,7,7,10,1,6,3,5,3,7,7,7,7,7~4,6,7,11,11,4,8,8,7,5,3,12,1,7,7,10,8,6,9,7,7,7,7,7&purInit=[{type:"fsbl",bet:1000,bet_level:0}]&reel_set3=6,7,8,10,2,5,6,12,7,11,6,5,10,6,5,8,7,7,12,4,10,5,8,4,10,3,12,8,10,4,8,9,2,12,10,9,4,3,10,5,8,12,3,6,8,12,6,4,10,11,12,8,7,7,7,7,7~4,3,9,7,11,3,8,6,3,9,7,11,5,10,6,3,9,7,7,12,3,4,11,3,6,2,9,5,11,9,6,11,4,9,3,5,4,11,6,3,5,11,2,12,9,10,4,3,11,5,2,9,8,7,7,7,7,7~12,2,3,6,7,11,3,5,7,10,4,3,5,2,9,4,5,6,7,7,12,9,5,4,3,8,5,6,11,4,6,2,10,3,6,4,3,8,7,7,7,7,7~5,7,6,2,12,5,2,8,7,6,4,11,3,9,7,7,10,6,3,5,2,7,7,7,7,7~4,6,7,11,2,4,8,2,7,5,3,12,7,7,10,2,6,9,7,7,7,7,7&total_bet_min=0.01`)
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
