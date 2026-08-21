package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs9aztecgemsdx/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=3,5,4,7,3,5,6,4,7&c_paytable=9~any~3,4~10,0,0~2&balance=16,618.86&cfgs=3243&reel1=2,3,3,3,3,7,5,6,6,6,8,4,4,4,6,4,7,7,7,5,5,5,4,3,5,3,7,4,3,4,6,5,6,4&ver=2&reel0=6,6,6,2,5,5,5,5,4,3,7,6,8,3,3,3,4,4,4,7,7,7,4,7,5,8,7&mo_s=8&index=1&balance_cash=16,618.86&def_sb=3,4,5&mo_v=5,8,18,28,58,68,88,108,128,288,888,900,2250&def_sa=3,4,5&reel2=4,2,3,6,5,8,4,4,4,7,6,6,6,5,5,5,7,7,7,5,7,5,6,7,6,5,7,5,6,5,6,5,7,3,6,2,5,8,5,2,3,5,3,5,7,5,6,5,6,2,7,2,5,6,5,3,5,6&bonusInit=[{bgid:2,bgt:42,bg_i:"18,28,58,88,108,128,188,288,388,100,250,500,1000",bg_i_mask:"w,w,w,w,w,w,w,w,w,w,w,w,w"},{bgid:3,bgt:42,bg_i:"2,3,5,8,10",bg_i_mask:"wlm,wlm,wlm,wlm,wlm"}]&mo_jp=900;2250;0&balance_bonus=0.00&na=s&scatters=1~0,0,0~0,0,0~1,1,1&gmb=0,0,0&rt=d&gameInfo={props:{jp1:"9000",jp-units:"coin",jp3:"2250",jp2:"4500",jp4:"900"}}&mo_jp_mask=jp4;jp3;jp1&stime=1737359462954&sa=3,4,5&sb=3,4,5&sc=0.10,0.20,0.30,0.50,1.00,2.00,3.00,4.00,6.67,13.33,25.00,50.00,100.00,200.00,250.00,300.00&defc=1.00&sh=3&wilds=2~250,0,0~1,1,1&bonuses=0&fsbonus=&c=1.00&sver=5&counter=2&ntp=0.00&paytable=0,0,0;0,0,0;0,0,0;88,0,0;58,0,0;28,0,0;18,0,0;8,0,0;0,0,0;0,0,0&l=9&s=3,5,4,7,3,5,6,4,7`)
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
