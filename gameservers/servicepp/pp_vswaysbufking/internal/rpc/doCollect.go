package rpc

import (
	"maps"
	"strings"
	"time"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicepp/pp_vswaysbufking/internal"
	"serve/servicepp/ppcomm"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	ppcomm.RegRpc("doCollect", doCollect)
}

var (
	sampleCollect = ppcomm.ParseVariables(`balance_cash=1,009,701.65&balance_bonus=0.00&na=s&stime=1727072800469&sver=5&counter=4&ntp=-1.00`)
)

func doCollect(msg *nats.Msg) (ret []byte, err error) {
	// balance=999,967.95
	// index=5
	// balance_cash=999,967.95
	// balance_bonus=0.00
	// na=s
	// stime=1726280350909
	// sver=5
	// counter=10
	// ntp=-1.60

	ps := ppcomm.ParseVariables(string(msg.Data))
	//fmt.Println(ps)

	var pid int64
	pid, err = jwtutil.ParseToken(ps.Str("mgckey"))
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *ppcomm.Player) error {
		if plr.IsCollect {
			return nil
		}
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		plr.IsCollect = true
		info := maps.Clone(sampleCollect)
		plr.LastData["tw"] = strings.ReplaceAll(plr.LastData["tw"], ",", "")
		win := ut.Money2Gold(ut.GetFloat(plr.LastData.Float("tw")))
		balance, _ := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  win,
			RoundID: plr.LastBetId,
			Reason:  slotsmongo.ReasonWin,
			IsEnd:   true,
		})
		info.SetFloat("balance", ut.Gold2Money(balance))
		info.SetFloat("balance_cash", ut.Gold2Money(balance))
		info.SetInt("stime", int(time.Now().UnixMilli()))
		info.SetStr("index", ps.Str("index"))
		info.SetInt("counter", ps.Int("counter")+1)

		ppcomm.InsertBetHistoryEvery(plr, 0, false, false, ps, info, curItem.Key, curItem.Symbol, 0, primitive.ObjectID{}, sampleDoInit.Encode(), "")

		plr.LastData = info
		betLog := &slotsmongo.AddBetLogParams{
			ID:           plr.LastBetId + "End",
			Pid:          plr.PID,
			Bet:          0,
			Win:          win,
			Balance:      balance,
			RoundID:      plr.LastBetId,
			Completed:    true,
			TotalWinLoss: win - plr.CostGold,
			IsBuy:        plr.IsBuy,
			Grade:        int(plr.C * internal.Line),
			PGBetID:      plr.BetHistLastID,
			BigReward:    plr.BigReward,
			GameType:     define.GameType_Slot,
		}

		ppcomm.SetRoundID(0)
		ppcomm.SetTurnIndex(0)

		if msg.Header != nil && msg.Header.Get("jump_log") != "1" {
			slotsmongo.AddBetLog(betLog)
		}
		ret = info.Bytes()
		return nil
	})

	return
}
