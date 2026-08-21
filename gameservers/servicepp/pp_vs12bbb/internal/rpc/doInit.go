package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs12bbb/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`def_s=5,6,9,6,6,8,11,8,9,9,6,9,12,6,6,6,9,12,6,6&prg_m=cp,tp,lvl&balance=137,380.00&cfgs=7409&ver=2&prg=0,4,0&mo_s=7&index=1&balance_cash=137,380.00&def_sb=10,7,7,8,11&mo_v=24,60,120,180,240,300,600,48000&reel_set_size=4&def_sa=11,10,7,10,4&reel_set=0&prg_cfg_m=s&balance_bonus=0.00&na=s&scatters=1~0,0,0,0,0~20,15,10,0,0~1,1,1,1,1&gmb=0,0,0&rt=d&gameInfo={props:{max_rnd_sim:"1",max_rnd_hr:"981884",max_rnd_win:"4000"}}&wl_i=tbm~4000&stime=1733968483158&sa=11,10,7,10,4&sb=10,7,7,8,11&prg_cfg=2&sc=20.00,40.00,50.00,60.00,70.00,90.00,125.00,200.00,250.00,400.00,450.00,500.00,600.00,700.00,750.00,900.00,1700.00,2500.00,3500.00,4200.00,5000.00,6000.00,7000.00,8500.00,10000.00,12500.00,17000.00,25000.00,34000.00,42000.00,50000.00,60000.00,67000.00,75000.00,90000.00,100000.00&defc=70.00&sh=4&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=70.00&sver=5&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;2400,240,60,6,0;1200,180,36,0,0;600,120,24,0,0;600,120,24,0,0;0,60,12,0,0;120,30,6,0,0;120,30,6,0,0;120,30,6,0,0;120,30,6,0,0;120,30,6,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;0,0,0,0,0&l=12&rtp=96.71&reel_set0=9,1,5,8,7,10,9,8,11,9,7,3,7,7,7,7,11,7,12,6,11,7,5,9,6,6,4,5,8,4,11,7~12,7,8,4,6,11,5,1,7,7,7,12,7,9,7,10,3,10,6,8,5~5,4,7,6,11,10,9,7,7,7,8,9,7,12,3,10,6,7,1~4,7,7,6,9,11,1,7,7,7,12,5,9,6,7,5,3,8,8,10~4,5,10,9,6,9,10,5,6,1,7,6,7,7,7,3,8,7,12,9,11,8,5,7,8,7,4,12,11,7&s=5,6,9,6,6,8,11,8,9,9,6,9,12,6,6,6,9,12,6,6&accInit=[{id:0,mask:"cp"},{id:2,mask:"cp;mp"},{id:3,mask:"cp;mp"}]&reel_set2=10,3,7,4,3,7,4,4,9,6,7,10,12,6,6,7,6,7,7,7,12,7,10,3,12,9,6,7,10,7,12,7,12,10,7,6,6,12,10~11,11,7,10,8,8,7,1,5,4,7,8,7,11,7,7,7,1,7,1,1,4,5,7,10,4,11,7,11,5,7,1,5~5,6,9,3,6,3,11,6,12,5,9,9,3,6,3,8,11,12,1,8,1,5,11,12,12,11,1,9,9,5,9,1,11,1,8,6,1,8,6,6,8~8,7,10,5,8,7,9,12,12,4,6,8,10,9,6,9,8,5,7,7,7,11,7,3,7,6,9,12,6,7,5,6,11,4,3,7,11,7,10,4,6,12~3,12,7,11,8,10,6,6,7,11,10,7,12,8,7,7,7,4,5,6,12,6,9,4,8,5,11,10,7,9,7,9,6&reel_set1=10,7,11,5,8,7,1,8,7,7,7,7,11,4,5,7,1,11,7,1,8,4~12,6,7,10,9,10,12,7,4,10,10,12,3,6,7,6,7,7,7,7,12,9,7,3,12,6,6,4,6,4,7,7,10,10,4,6,7,7~12,9,9,11,9,6,3,9,5,5,8,6,12,6,8,8,3,3,11,6,11,8,6,5,9,3,12,8,5,6,9,11,11,6,12~6,1,3,12,6,7,4,5,11,10,6,1,9,11,9,1,12,4,12,6,6,8,7,7,7,3,8,8,7,7,9,10,6,9,8,10,1,5,7,11,1,1,4,7,5,7,6,7,12~3,6,5,11,4,10,9,12,8,7,6,4,6,9,9,12,7,7,7,7,5,5,7,10,10,6,11,6,11,9,6,4,8,7,12,8,7,7&reel_set3=12,4,12,10,7,8,6,3,8,11,4,12,10,7,7,7,7,10,12,5,8,9,10,7,12,11,9,8,10,8,6~11,10,4,6,11,9,5,11,3,11,7,7,11,12,10,9,5,12,3,8,9,9,8,7,9~4,10,12,3,6,6,7,10,4,7,9,12,4,8,5,11,11,7,4,5,9,7,8,6,6,5,5~7,6,5,4,5,3,6,3,8,10,7,11,5,9,12,6,7,4~6,11,3,7,9,8,4,12,10,7,5,7`)
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
