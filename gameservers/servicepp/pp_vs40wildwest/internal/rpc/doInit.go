package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs40wildwest/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=11,5,7,7,5,1,6,9,9,6,12,11,9,9,11,12,11,5,5,11&apvi=10&balance=1,000,000.00&cfgs=5005&ver=2&index=1&balance_cash=1,000,000.00&reel_set_size=2&def_sb=8,2,6,6,1&def_sa=11,9,5,3,9&reel_set=0&balance_bonus=0.00&na=s&scatters=1~0,0,0,0,0~0,0,8,0,0~1,1,1,1,1;14~0,0,0,0,0~0,0,8,0,0~1,1,1,1,1&cls_s=-1&gmb=0,0,0&mbri=1,2,3&rt=d&wl_i=tbm~10000&apti=bet_mul&stime=1733885460604&sa=11,9,5,3,9&sb=8,2,6,6,1&sc=0.04,0.08,0.12,0.16,0.20,0.40,0.80,1.20,1.60,2.00,3.00,4.00,8.00,12.00,16.00,20.00&defc=0.40&purInit_e=1&sh=4&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=0.40&sver=5&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;400,100,30,0,0;250,75,25,0,0;150,40,15,0,0;100,25,10,0,0;75,15,7,0,0;50,10,5,0,0;30,6,3,0,0;30,6,3,0,0;20,5,2,0,0;20,5,2,0,0;20,5,2,0,0;0,0,0,0,0&l=20&rtp=96.51&total_bet_max=40,000.00&reel_set0=7,11,11,1,12,12,6,8,4,10,10,5,11,11,9,9,3,13,13,5,8,12,12,1,13,13,6,10,10~7,11,11,2,12,12,6,8,4,9,9,5,13,13,3,11,11,5,8,12,12,2,13,13,6,10,10~9,7,11,11,2,13,13,6,8,4,9,9,5,10,10,1,6,8,3,11,11,5,8,12,12,2,13,13,6~7,10,10,2,12,12,6,8,11,11,4,9,9,5,6,7,3,11,11,5,6,12,12,7,13,13,6~7,10,10,1,12,12,6,8,4,9,9,5,6,7,3,11,11,5,6,13,13,7,13,13,6,10,10&s=11,5,7,7,5,1,6,9,9,6,12,11,9,9,11,12,11,5,5,11&reel_set1=10,5,9,9,7,10,10,8,12,12,6,13,13,8,9,9,4,9,9,5,6,8,3,3,3,3,11,11~7,10,10,2,12,12,6,8,4,9,9,5,6,3,11,11,5,6,12,12,7,13,13,2,10,10,7,4~7,10,10,2,12,12,6,8,4,9,9,5,6,3,11,11,5,6,12,12,8,13,13,2,10,10,7,4~7,10,10,2,12,12,6,8,4,9,9,5,6,7,3,11,11,5,6,12,12,7,13,13,6,10,10,7~10,10,6,12,12,8,4,9,9,5,6,7,3,3,3,3,11,11,6,12,12,7,13,13,6,10,10,7&purInit=[{type:"fs",bet:2000,fs_count:8}]&mbr=1,1,1&total_bet_min=0.80`)
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
