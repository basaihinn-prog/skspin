package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs20eking/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=3,4,5,6,7,3,4,5,6,7,3,4,5,6,7&balance=1,425,958.39&cfgs=5496&accm=cp~pp;cp~pp&ver=2&acci=0;1&index=1&balance_cash=1,425,958.39&def_sb=5,7,7,8,2&reel_set_size=3&def_sa=8,3,4,3,3&reel_set=0&balance_bonus=0.00&na=s&accv=1~1;0~0&scatters=1~0,0,0,0,0~0,0,0,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&wl_i=tbm~20000&cpri=2&stime=1739782055334&sa=8,3,4,3,3&sb=5,7,7,8,2&sc=10.00,20.00,30.00,40.00,50.00,60.00,80.00,100.00,150.00,200.00,250.00,300.00,350.00,400.00,450.00,500.00,750.00,1000.00,1500.00,2000.00,2500.00,3000.00,4000.00,4500.00,5000.00,7500.00,10000.00,15000.00,20000.00,25000.00,30000.00,35000.00,40000.00,45000.00,50000.00,60000.00&defc=50.00&sh=3&wilds=2~500,100,50,0,0~1,1,1,1,1&bonuses=0&fsbonus=&st=rect&c=50.00&sw=5&sver=5&g={comm:{reel_set:"2",screenOrchInit:"{type:\"mini_slots\",layout_h:3,layout_w:5}"},ms00:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms01:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms02:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms03:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms04:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms05:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms06:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms07:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms08:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms09:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms10:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms11:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms12:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms13:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"},ms14:{def_s:"13,13,13",sh:"1",st:"rect",sw:"3"}}&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;500,100,50,0,0;100,50,30,0,0;50,30,20,0,0;30,20,10,0,0;30,20,10,0,0;20,10,5,0,0;20,10,5,0,0;400,0,0;300,0,0;200,0,0;0,0,0&l=20&rtp=96.51&reel_set0=3,6,8,9,7,7,9,9,9,7,9,9,9,5,5,7,9,4,6,5,5,5,7,8,7,8,7,5,7,7,7,3,9,9,9,9,2,5,6,9,9,9,7,7,7,7~3,8,8,8,9,9,4,4,4,8,8,8,6,6,6,8,8,6,6,6,4,4,8,5,8,8,6,6,8,8,8,8,6,6,6,2,6,6,8,4,7,3,3~3,9,9,9,7,9,7,9,7,7,7,4,4,4,9,9,9,5,9,9,8,8,8,6,6,6,5,5,5,9,7,2,9,3,3~3,7,8,6,6,8,8,8,6,5,5,5,8,8,8,6,5,9,7,7,7,7,4,4,4,2,8,9,9,9,4,8,3,3,3,6,6,6~3,9,5,7,7,7,4,4,4,9,9,7,7,7,9,9,9,8,8,8,7,5,5,5,5,9,9,5,5,5,6,9,9,9,6,6,6,2,5,2,3,3,3,5&s=3,4,5,6,7,3,4,5,6,7,3,4,5,6,7&accInit=[{id:0,mask:"cp;pp;mp"},{id:1,mask:"cp;pp;mp"}]&reel_set2=10,12,12,12,11,11,10,12,12,10,10,11,13,13,13,10,10,10,12,12,11,11,11,11,13~10,12,12,12,10,10,10,11,13,13,11,11,10,10,10,11,11,11,12,12,12,12,13~10,12,12,12,13,13,10,10,13,13,13,11,11,10,10,10,11&reel_set1=3,3,3,6,6,2,4,4,4,9,9,2,2,2,5,5,5,7,2,3,2,2,2,9,8,8,8,8,7,9,9,9,9,7,7,7,2,2,9,7,9~3,3,3,8,8,2,9,6,4,4,4,8,8,8,9,9,9,6,6,6,2,2,4,5,5,5,8,8,7,7,7,6,2,2,2,3~3,3,3,6,6,6,2,7,3,8,7,7,7,5,5,9,8,8,8,2,2,2,4,4,4,2,6,6,5,5,5,2,2,9,9,9,9,7~3,3,3,7,8,9,9,9,7,8,2,2,8,5,8,8,8,6,2,2,2,7,7,7,2,2,2,8,6,6,6,6,4,4,7,4,4,4,8~3,3,3,9,5,7,7,4,4,4,6,2,2,2,8,7,9,9,8,8,8,5,3,9,2,6,5,2,2,2,9,9,9,7,7,7,2,9,9&cpri_mask=tbw`)
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
