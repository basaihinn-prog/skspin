package rpc

import (
	"log/slog"
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs25btygold/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=3,10,2,3,10,6,7,8,4,7,4,10,4,5,11,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13&balance=3,482,007.00&cfgs=7669&accm=cp~tp~lvl~sc&ver=2&mo_s=11;12&acci=0&index=1&balance_cash=3,482,007.00&mo_v=25,50,75,100,125,150,175,200,250,375,500,625,1250,1875,2500,6250;250,375,500,625,1250,2500&reel_set_size=5&balance_bonus=0.00&na=s&accv=0~6~0~0&scatters=1~0,0,0,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&gameInfo={rtps:{regular:"96.50"},props:{max_rnd_sim:"1",max_rnd_hr:"826446",jp1:"25",max_rnd_win:"5000",jp3:"500",jp2:"50",jp4:"5000"}}&wl_i=tbm~5000&stime=1736481594038&sc=8.00,16.00,24.00,32.00,40.00,50.00,60.00,80.00,120.00,160.00,200.00,240.00,280.00,320.00,360.00,400.00,600.00,800.00,1200.00,1600.00,2000.00,2800.00,3200.00,3600.00,4000.00,6000.00,8000.00,12000.00,16000.00,20000.00,24000.00,28000.00,32000.00,36000.00,40000.00,48000.00&defc=40.00&sh=12&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&st=rect&c=40.00&sw=5&sver=5&g={base:{def_s:"3,10,2,3,10,6,7,8,4,7,4,10,4,5,11",def_sa:"4,7,5,2,11",def_sb:"6,7,2,4,7",reel_set:"0",s:"3,10,2,3,10,6,7,8,4,7,4,10,4,5,11",sa:"4,7,5,2,11",sb:"6,7,2,4,7",sh:"3",st:"rect",sw:"5"},matrix_2:{def_s:"13,13,13,13,13,13,13,13,13,13,13,13,13,13,13",def_sa:"13,13,13,13,13",def_sb:"13,13,13,13,13",reel_set:"2",s:"13,13,13,13,13,13,13,13,13,13,13,13,13,13,13",sa:"13,13,13,13,13",sb:"13,13,13,13,13",sh:"3",st:"rect",sw:"5"},matrix_3:{def_s:"13,13,13,13,13,13,13,13,13,13,13,13,13,13,13",def_sa:"13,13,13,13,13",def_sb:"13,13,13,13,13",reel_set:"3",s:"13,13,13,13,13,13,13,13,13,13,13,13,13,13,13",sa:"13,13,13,13,13",sb:"13,13,13,13,13",sh:"3",st:"rect",sw:"5"},matrix_4:{def_s:"13,13,13,13,13,13,13,13,13,13,13,13,13,13,13",def_sa:"13,13,13,13,13",def_sb:"13,13,13,13,13",reel_set:"4",s:"13,13,13,13,13,13,13,13,13,13,13,13,13,13,13",sa:"13,13,13,13,13",sb:"13,13,13,13,13",sh:"3",st:"rect",sw:"5"}}&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;200,50,10,0,0;150,30,5,0,0;125,25,5,0,0;75,25,5,0,0;50,10,5,0,0;50,10,5,0,0;50,10,5,0,0;50,10,5,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0&l=25&reel_set0=7,8,6,8,5,6,4,8,10,9,7,10,11,10,3,6,8,9,11,9,4,9,4,7,3,10,6,7,5,6,11,6,7,10,9,11,10,4,5,7,5,4,7,4,10,9,10,6,8,5,10,9,10,5,8,6~10,4,8,4,6,10,3,5,4,10,11,7,9,7,6,2,2,2,2,2,6,8,4,9,2,8,5,2,11,3,9,3,4,5,10,9,5,10,11,11,11,2,2,9,3,7,9,7,8,3,4,7,2,8,11,6,9,8,2,11~8,7,3,8,9,6,5,8,2,5,4,3,10,8,2,2,2,6,8,7,3,9,2,6,9,4,7,10,2,4,4,11,11,11,6,7,11,11,6,9,5,8,2,10,3,2,9,8,5~11,2,5,9,10,8,5,8,6,2,8,7,5,7,11,6,10,8,2,2,2,2,2,7,9,11,7,8,7,2,7,8,7,11,6,7,10,3,4,8,9,11,11,11,6,5,4,5,10,9,7,3,7,9,4,6,3,9,11,4,3,9,4,3~5,11,10,2,6,7,3,9,7,5,8,9,2,3,2,2,2,2,11,4,9,10,8,4,10,6,9,3,5,3,6,10,7,10,11,11,11,7,6,2,3,4,8,6,5,2,7,4,10,5,8,4,6&s=3,10,2,3,10,6,7,8,4,7,4,10,4,5,11,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13,13&accInit=[{id:0,mask:"cp;tp;lvl;sc;cl"}]&reel_set2=11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13&reel_set1=11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13&reel_set4=11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13&reel_set3=11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13~11,13,11,13,13,13,11,11,13,11,13,13`)
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
