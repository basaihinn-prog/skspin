package rpc

import (
	"log/slog"
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs40bigjuan/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=12,5,7,11,3,12,5,2,11,10,11,5,5,11,10,10,5,10,11,7&balance=3,613,507.00&cfgs=7609&accm=cp~tp~lvl~sc;cp~tp~lvl~sc;cp~tp~lvl~sc;cp~tp~lvl~sc&ver=2&mo_s=15&acci=1;2;3;4&index=1&balance_cash=3,613,507.00&mo_v=20,40,80,120,200,320,400,600,800,1000,1600,2000,4000,5000,8000,10000&def_sb=4,10,10,7,8&reel_set_size=4&def_sa=8,8,9,6,3&reel_set=0&balance_bonus=0.00&na=s&accv=0~5~0~0;0~5~0~0;0~4~0~0;0~3~0~0&scatters=1~0,0,0,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&gameInfo={rtps:{purchase:"96.53",regular:"96.70"},props:{max_rnd_sim:"1",max_rnd_hr:"2896871",jp1:"100000",max_rnd_win:"2600",jp3:"2000",jp2:"10000",jp4:"500"}}&wl_i=tbm~2600&stime=1736394453440&sa=8,8,9,6,3&sb=4,10,10,7,8&sc=5.00,10.00,15.00,20.00,25.00,30.00,40.00,50.00,75.00,100.00,125.00,150.00,175.00,200.00,225.00,250.00,375.00,500.00,750.00,1000.00,1250.00,1500.00,1750.00,2000.00,2500.00,3750.00,5000.00,7500.00,10000.00,12500.00,15000.00,17500.00,20000.00,22500.00,25000.00,30000.00&defc=25.00&purInit_e=1&sh=4&wilds=2~500,150,50,0,0~1,1,1,1,1&bonuses=0&fsbonus=&st=rect&c=25.00&sw=5&sver=5&g={fs_collect:{def_s:"22,14,14",def_sa:"23",def_sb:"22",reel_set:"2",s:"22,14,14",sa:"23",sb:"22",sh:"3",st:"rect",sw:"1"},fs_main:{def_s:"15,15,18,14,16,14,15,14,20",def_sa:"14,14,14",def_sb:"14,14,14",reel_set:"1",s:"15,15,18,14,16,14,15,14,20",sa:"14,14,14",sb:"14,14,14",sh:"3",st:"rect",sw:"3"}}&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;250,100,25,0,0;200,80,20,0,0;150,40,15,0,0;100,20,10,0,0;100,20,10,0,0;40,10,5,0,0;40,10,5,0,0;40,10,5,0,0;40,10,5,0,0;40,10,5,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0&l=40&total_bet_max=120,000,000.00&reel_set0=8,7,7,13,13,11,13,7,1,12,11,9,5,4,3,13,11,12,3,6,3,3,3,3,12,11,12,11,8,8,5,6,8,12,5,3,5,2,11,2,3,12,8,13,10,8,7~12,3,13,6,8,8,11,5,13,9,9,9,9,10,13,8,10,6,13,10,5,11,10,4,13,13,13,13,9,13,8,7,10,6,6,7,9,3,13,8,8,8,8,2,11,11,10,4,4,6,12,8,11,11,11,11,9,1,5,2,11,9,13,10,6,1,4,10,10,10,10,5,8,8,11,11,2,4,9,9,11,5,6,6,6,6,13,5,10,9,10,2,10,8,4,4~3,10,9,12,7,2,10,7,12,7,13,9,5,13,13,13,13,3,12,9,12,1,13,5,11,8,6,5,12,13,4,2,9,9,9,9,10,5,3,11,12,9,7,13,12,13,2,10,10,9,12,12,12,8,9,9,10,9,1,13,3,9,12,4,10,6,9,10,3~6,7,5,7,2,13,7,7,11,4,13,13,13,13,8,13,10,6,1,8,6,2,13,9,7,11,11,11,11,7,12,11,1,3,10,8,11,6,4,13,6,6,6,11,4,4,8,8,4,13,8,10,11,13,8,8,8,8,11,7,10,3,8,10,6,11,13,10,5,7,7,7,7,6,9,6,4,12,11,6,7,8,4,2,10,10~12,3,9,11,5,11,3,11,8,3,9,13,6,2,5,2,12,13,9,8,6,13,8,9,12,10,7,11,11,13,5,6,3,6,7,13,2,10,7,6,6,3,3,3,3,11,9,5,6,11,12,10,12,5,3,9,5,10,3,3,13,9,12,13,5,11,9,12,13,8,5,7,2,3,7,3,4,5,13,12,11,6,7,4,13,5,9,12,13,13,13,2,7,3,11,1,3,6,7,11,6,9,3,12,4,12,5,1,10,7,2,3,7,12,8,13,8,13,12,10,12,11,7,12,6,9,5,12,6,7,5,13,11,1,6,13&s=12,5,7,11,3,12,5,2,11,10,11,5,5,11,10,10,5,10,11,7&accInit=[{id:1,mask:"cp;tp;lvl;sc;cl"},{id:2,mask:"cp;tp;lvl;sc;cl"},{id:3,mask:"cp;tp;lvl;sc;cl"},{id:4,mask:"cp;tp;lvl;sc;cl"}]&reel_set2=23,22,14,23,14,23,23,23,22,22,14,14,23,14,23,22,22,22,14,22,22,14,23,14,14&reel_set1=14,14,14,14,14~14,14,14,14,14~14,14,14,14,14&purInit=[{type:"fsbl",bet:4000,bet_level:0}]&reel_set3=15,24,15,15,21,24,14,14,15,14,24,14,21,15,14,15,15,15,15,14,15,15,14,15,24,15,15,21,15,14,14,15,24,14,24,24,24,15,24,14,15,14,24,24,21,15,21,14,15,14,15,14,21,14&total_bet_min=5.00`)
)

// var cGrade = []float64{0.01, 150000.0}
var cGrade = []float64{0.05, 0.10, 0.15, 0.25, 0.50, 1.00, 1.50, 2.00, 3.50, 7.50, 15.00, 25.00, 40.00, 75.00, 110.00, 150.00}

//var betGrade = []float64{0.05, 0.10}

func doInit(msg *nats.Msg) (ret []byte, err error) {
	// action=doInit&symbol=vs20olympx&cver=237859&index=1&counter=1&repeat=0&mgckey=AUTHTOKEN@7a4f4399cd1bd51a2eef10551a63ce1fdb69a05e0363abf9d5c47c3f5da92f13~stylename@hllgd_hollygod~SESSION@d5ec7d86-a060-4470-a302-0192a2d53f14~SN@54345f72

	ps := ppcomm.ParseVariables(string(msg.Data))
	//fmt.Println(ps)

	var pid int64
	pid, err = jwtutil.ParseToken(ps.Str("mgckey"))
	if err != nil {
		slog.Error("mgckey err")
		return nil, err
	}

	err = db.CallWithPlayer(pid, func(plr *ppcomm.Player) error {
		info := maps.Clone(sampleDoInit)
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)

		// ev : "FR0~1.00,20,10,0,0,1729739148,1,,"
		// fra : "0.00"
		// frn : "10"
		// frt : "N"

		// "FR1~0.50,20,33.00,,;FR0~1.00,20,10,0,0,1729740027,1,,"

		ppcomm.FRBOnInit(plr, info)
		isEnd, _ := plr.IsEndO()
		c, err := redisx.GetPlayerCs(plr.AppID, plr.PID, isEnd)
		if err != nil {
			return err
		}
		minBet := c[0]
		maxBet := c[len(c)-1] * 100 * internal.Line

		info.SetFloat("c", minBet*curItem.Multi)
		for k, v := range plr.LastData {
			if !ppcomm.IsFRBField(k) {
				info.SetStr(k, v)
			}
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
