package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs20sh/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`tw=0.00&def_s=4,6,4,3,8,5,8,3,5,6,4,7,5,4,7&rid=693124996000&balance=100.00&action=doSpin&cfgs=8819&ver=3&index=1&balance_cash=100.00&def_sb=3,7,7,5,8&reel_set_size=1&def_sa=3,6,5,4,7&reel_set=0&balance_bonus=0.00&na=s&scatters=1~500,15,5,0,0~0,0,0,0,0~1,1,1,1,1&rt=d&gameInfo={rtps:{regular:"96.33"},props:{max_rnd_sim:"1",max_rnd_hr:"2350058906",max_rnd_win:"2500"}}&stime=1735889925107&sa=6,5,5,5,7&sb=6,8,3,8,3&sc=0.01,0.02,0.03,0.05,0.10,0.15,0.20,0.25,0.50,0.75,1.25,2.50,5.00,10.00,15.00,20.00&defc=0.10&sh=3&wilds=2~2500,700,50,0,0~1,1,1,1,1&bonuses=0&st=rect&c=0.01&sw=5&sver=5&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;1200,250,15,0,0;230,50,15,0,0;230,50,15,0,0;110,25,10,0,0;110,25,10,0,0;110,25,10,0,0&l=20&reel_set0=8,6,8,7,2,7,8,6,7,4,8,7,4,4,4,4,4,3,8,8,6,7,8,7,4,8,5,8,8,7,7,7,7,6,4,1,5,2,4,4,7,6,7,7,3,3,6,6,6,6,6,3,6,4,5,8,3,6,3,1,6,5,6,5,5,5,5,3,6,3,7,3,6,7,4,7,6,8,4,8,8,8,8,8,7,7,8,6,8,1,3,6,4,4,8,6,4,3,3,3,3,8,6,5,5,7,7,1,3,7,3,4,6,7,2,2,2,7,2,6,4,4,8,4,6,8,7,3,3,4,4~3,2,6,1,5,3,3,6,7,5,5,8,6,6,6,8,7,6,4,7,8,7,5,8,8,7,5,8,5,5,5,5,8,7,5,1,7,7,5,4,6,7,6,8,6,8,8,8,8,4,1,7,5,8,8,3,5,4,5,7,3,2,7,7,7,8,4,5,6,7,8,1,8,6,4,5,6,2,4,4,4,6,5,7,7,8,6,7,8,6,6,4,5,7,2,2,2,7,8,5,1,5,6,6,8,3,4,6,5,3,3,3,3,5,3,7,6,7,6,8,7,5,8,5,8,5,7,3~5,3,3,5,3,7,8,8,8,8,5,8,6,6,8,3,5,5,4,4,4,4,6,8,4,8,7,3,3,1,3,3,3,6,8,4,1,5,5,2,5,6,6,6,6,8,4,7,6,7,3,5,5,5,5,6,3,2,3,6,4,4,7,7,7,5,1,5,4,6,6,3,8,2,2,2,3,6,7,4,4,5,7,4,6~7,6,8,6,3,4,4,7,8,4,7,7,7,7,4,2,6,6,5,4,1,8,7,8,4,5,8,8,8,8,5,5,3,4,4,8,3,4,8,8,5,4,4,4,4,7,4,1,2,7,3,8,8,5,5,8,8,6,6,6,6,4,7,8,3,5,4,3,5,7,7,8,3,3,3,3,7,5,6,6,5,8,7,6,1,8,4,5,5,5,5,4,5,7,4,8,6,7,5,5,7,7,4,2,2,2,2,5,7,5,6,7,2,4,7,2,7,2,4,8~3,7,5,6,8,8,2,7,6,4,4,8,7,8,7,7,7,7,7,6,6,1,7,6,1,4,7,7,3,8,4,7,6,3,6,5,5,5,5,5,6,2,6,4,6,5,3,8,3,4,8,2,3,3,7,6,6,6,6,6,8,3,8,4,7,6,3,1,7,7,4,3,6,8,6,7,3,3,3,3,8,3,6,3,3,7,3,4,7,8,3,4,4,7,7,2,8,8,8,8,4,3,8,6,6,7,1,7,6,8,3,7,4,7,8,4,4,4,4,3,6,8,2,4,2,7,2,6,7,6,5,7,8,2,4,2,2,2,2,8,5,7,4,5,4,8,5,3,4,5,3,8,4,4,6,6&s=6,5,5,7,8,6,5,1,7,6,6,1,3,7,3&w=0.00`)
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
