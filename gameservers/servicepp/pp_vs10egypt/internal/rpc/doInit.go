package rpc

import (
	"maps"
	"serve/comm/redisx"
	"serve/servicepp/pp_vs10egypt/internal"
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
	sampleDoInit = ppcomm.ParseVariables(`wsc=1~bg~50,10,1,0,0~0,0,0,0,0~fs~50,10,1,0,0~10,10,10,0,0&def_s=5,8,7,9,8,8,7,3,4,4,11,6,8,11,10&balance=20,000.00&cfgs=2000&ver=2&index=1&balance_cash=20,000.00&reel_set_size=10&def_sb=8,8,7,5,9&def_sa=9,1,8,4,10&balance_bonus=0.00&na=s&scatters=&gmb=0,0,0&rt=d&stime=1739759580697&sa=9,1,8,4,10&sb=8,8,7,5,9&sc=20.00,40.00,60.00,80.00,100.00,120.00,150.00,200.00,300.00,400.00,500.00,600.00,700.00,800.00,900.00,1000.00,1500.00,2000.00,3000.00,4000.00,5000.00,6000.00,7000.00,8000.00,9000.00,10000.00,20000.00,30000.00,40000.00,50000.00,60000.00,70000.00,80000.00,90000.00,100000.00,120000.00&defc=100.00&sh=3&wilds=2~0,0,0,0,0~1,1,1,1,1&bonuses=0&fsbonus=&c=100.00&sver=5&n_reel_set=0&counter=2&ntp=0.00&paytable=0,0,0,0,0;0,0,0,0,0;0,0,0,0,0;5000,1000,100,10,0;2000,400,40,5,0;500,100,15,2,0;500,100,15,2,0;150,40,5,0,0;150,40,5,0,0;100,25,5,0,0;100,25,5,0,0;100,25,5,0,0&l=10&reel_set0=9,1,8,4,10,9,7,6,8,7,5,7,10,11,11,8,6,10,8,3~1,8,10,11,4,6,5,3,7,5,11,8,9,10,7,9~7,9,9,10,11,11,4,8,4,6,1,10,5,9,6,10,3,8~7,8,11,9,4,6,10,8,3,7,1,10,9,4,5,6,10,5~8,8,7,5,9,10,6,4,3,1,11,10,1,4,7,6,10,3,8,9&s=5,8,7,9,8,8,7,3,4,4,11,6,8,11,10&reel_set2=8,7,9,9,10,4,1,4,6,3,7,10,5,11,11~9,11,1,8,10,6,11,4,3,4,7,10,5~10,7,6,10,6,1,5,3,11,4,7,9,11,8~6,8,7,9,11,5,10,10,3,8,1,4,10,11,6,9,5~1,6,9,7,3,10,6,9,11,11,8,7,8,4,5,10&reel_set1=6,7,5,10,7,9,3,3,8,4,10,11,6,7,1~9,3,5,4,8,11,7,6,8,1,10,10~3,11,6,9,8,7,6,10,4,5,7,9,1~8,1,5,9,3,5,7,8,10,10,11,11,6,10,4~8,11,7,6,5,3,1,11,6,4,9,8,3,7,10,10&reel_set4=4,11,7,10,9,7,6,3,10,8,1,5,11,8,8,6~7,8,3,6,1,11,10,11,5,6,10,4,5,7,9~9,7,4,3,6,8,11,1,8,7,10,6,4,7,10,5~7,4,11,10,10,9,7,6,5,3,11,5,1,8~10,8,6,5,11,7,5,6,9,8,8,7,1,10,4,3,10&reel_set3=11,10,8,5,7,7,11,9,6,4,9,1,3~11,10,1,3,7,5,9,6,8,9,4,5~5,7,4,11,8,10,4,11,1,6,7,5,9,3,10~7,3,11,6,9,10,10,9,8,5,1,5,11,4~5,6,3,5,11,8,10,1,10,9,7,7,8,7,4&reel_set6=10,6,8,8,6,10,9,11,7,3,5,9,4,1~8,3,1,4,7,5,8,8,11,9,6,10,11~10,11,3,9,7,8,10,1,6,8,9,10,4,5,4,7~11,8,9,10,1,5,7,5,4,3,5,10,11,10,7,6~7,8,5,3,5,6,4,10,11,10,11,9,1,9,10,6,7&reel_set5=11,7,5,6,6,9,7,6,4,3,7,8,10,10,1,9~11,7,9,1,10,3,8,6,5,4,9,11,11~10,7,6,9,1,7,9,7,4,10,5,4,3,11,11,8~11,7,6,11,8,5,7,1,3,6,4,9,10,10,5,5~10,9,11,7,1,10,8,10,7,6,4,8,5,6,3,9,11,5&reel_set8=1,11,5,8,6,11,9,3,4,9,6,7,8,10,10,9~5,6,9,11,8,3,7,8,1,11,7,4,5,9,6,9,10~5,6,9,7,3,9,11,4,4,11,1,8,10,9~9,11,8,10,11,9,7,5,5,10,4,3,5,6,1,6,10~8,5,4,11,1,7,10,10,5,9,3,6,7&reel_set7=6,10,8,9,4,8,7,11,9,10,9,1,6,5,3,6~8,11,4,10,8,7,3,5,7,1,10,9,8,6,9~11,10,11,6,9,8,3,1,4,10,6,4,5,7,7,9~9,7,6,3,10,8,1,11,4,9,5,7,11,10,5,10,9,6~5,9,9,10,4,7,8,7,5,11,6,3,1,10,11,8&reel_set9=10,4,3,9,8,8,11,11,10,6,5,11,7,6,1~7,8,7,4,6,11,10,8,1,3,10,8,6,5,9,9,10~3,8,11,6,4,9,5,10,7,4,6,11,1~9,10,6,4,11,9,10,8,11,5,3,11,1,7,5,6~7,5,3,10,10,11,8,8,6,9,6,7,4,1,9,5`)
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
