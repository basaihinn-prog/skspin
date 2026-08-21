package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"maps"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lang"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicehacksaw/hacksaw_1117/internal"
	"serve/servicehacksaw/hacksaw_1117/internal/gendata"
	"serve/servicehacksaw/hacksawcomm"
	"strconv"
)

func init() {
	hacksawcomm.RegRpc("doSpin", doSpin)
}

func doSpin(msg *nats.Msg) (ret []byte, err error) {
	pid, sessionData, err := hacksawcomm.ParseHackSawReq(msg.Data)
	if err != nil {
		return nil, err
	}
	err = db.CallWithPlayer(pid, func(plr *hacksawcomm.Player) error {
		//curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		App, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		bets := sessionData.Bets[0]
		//betRequest, _ := entity.GetSFSObject("betRequest")
		//hacksaw用的是分
		playerBet := bets.BetAmount
		// 参数处理
		bet, err := strconv.ParseFloat(playerBet, 64)
		if err != nil {
			return err
		}
		//缩放倍数
		multiple := bet / internal.FetchBet
		//if !ut.FloatInArr(ut.FloatArrMul(App.Cs, curItem.Multi*internal.Line), bet) {
		//	//当触发这个错误的时候直接一刀切，把该用户的上一局历史删除
		//	plr.RewriteLastData()
		//	return lang.Error(plr.Language, "下注额非法")
		//}
		//c := ps.Float("c") //下注金额
		isBuy := bets.BuyBonus != ""
		sf := ut.NewSnowflake()
		showIndex := strconv.FormatInt(sf.NextID(), 10)

		//传入450 实际为4.5 游戏以分为单位 需要转换
		realBet := ut.JdbBet2Money(bet)
		gold := ut.Money2Gold(realBet)
		costGold := ut.Money2Gold(realBet)
		mul := int64(1)
		mod := internal.FindBuy(bets.BuyBonus)
		if isBuy {
			mul = int64(internal.GetFreeMultiply()[mod])
			costGold *= mul
		}
		//balanceGold, err := slotsmongo.ModifyGold(plr.PID, -gold*mul, "下注")
		balanceGold, err := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  -costGold,
			RoundID: showIndex,
			Reason:  slotsmongo.ReasonBet,
		})
		if err != nil {
			return lang.Error(plr.Language, "金币不足")
		}
		var forcedKill, hitBigAward, buyKill bool
		var doc *hacksawcomm.SimulateData
		selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(plr.PID)
		doc, hitBigAward, forcedKill, buyKill, err = hacksawcomm.NextPlayResp(&hacksawcomm.NextPlayRespParam{
			Player:                 plr,
			Bet:                    gold,
			AppStore:               App,
			SelfPoolGold:           selfPoolGold,
			IsBuy:                  isBuy,
			BuyMinAwardPercent:     App.BuyMinAwardPercent,
			Mod:                    mod,
			Combine:                hacksawcomm.NewCombine(),
			NextFunc:               gendata.GCombineDataMng.Next,           //正常轮次下一轮获取
			SimulateByBucketIdFunc: gendata.GCombineDataMng.SampleSimulate, //根据既定bucketId获取
			BigRewardFunc:          gendata.GCombineDataMng.GetBigReward,   //大奖，从整体奖池里找即可
			BigRewardFunc2_5:       gendata.GCombineDataMng.GetBigReward2_5,
			MinBuyFunc:             gendata.GCombineDataMng.GetBuyMinData,         //杀！
			NextBuyFunc:            gendata.GCombineDataMng.NextBuy,               //购买游戏轮的次获取，包括首轮和后续轮次
			ControlNextData:        gendata.GCombineDataMng.ControlNextDataNormal, //控制游戏的下一次func
		})
		if err != nil {
			//slotsmongo.ModifyGold(plr.PID, gold*mul, "下注-退回, err:"+err.Error())
			slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
				Pid:     plr.PID,
				Change:  costGold,
				RoundID: showIndex,
				Reason:  slotsmongo.ReasonRefund,
			})
			return err
		}
		allWin := ut.JdbBet2Money(ut.GetFloat(doc.Times) * bet)
		fmt.Println("转动数据：", "pid:", plr.PID, "dataId:", doc.Id.Hex())
		//var totalWin float64
		//处理数据
		data := doc.Deal2(multiple, ut.HackGold2Money(balanceGold), showIndex)
		plr.LastData = maps.Clone(data)
		//data.SetCurrency("balance", ut.Gold2Money(balanceGold)) //设置余额
		//data.SetStr("gameSeq", showIndex)                       //设置回合
		//他要的余额需要在doCollect中添加
		balanceGold, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  ut.Money2Gold(allWin),
			RoundID: showIndex,
			Reason:  slotsmongo.ReasonWin,
			IsEnd:   true,
		})
		if err != nil {
			return err
		}
		//装载数据
		var toSelfAwardPool, bigReward int64

		if isBuy {
			if buyKill {
				//toSelfAwardPool += gold * mul
			}
		} else {
			if hitBigAward {
				toSelfAwardPool += -ut.Money2Gold(allWin)
				plr.BigReward = ut.Money2Gold(allWin)
			}

		}

		if !isBuy || buyKill {
			plr.UpdatePool(&hacksawcomm.UpdatePoolParam{
				Bet:                  gold,
				App:                  App,
				SelfPoolGold:         selfPoolGold,
				ToSelfAwardPool:      toSelfAwardPool,
				ForcedKill:           forcedKill,
				IsBuy:                isBuy,
				RewardPercentLess100: 0,
			})
		}
		plr.OnSpinFinish(costGold, ut.Money2Gold(allWin), isBuy, false, 1)

		betLog := &slotsmongo.AddBetLogParams{
			UserName:     "",
			CurrencyKey:  plr.CurrencyKey,
			ID:           showIndex,
			Pid:          plr.PID,
			Bet:          costGold,
			Win:          ut.Money2Gold(allWin),
			Balance:      balanceGold,
			RoundID:      showIndex,
			Completed:    true,
			TotalWinLoss: ut.Money2Gold(allWin) - costGold,
			IsBuy:        isBuy,
			Grade:        int(gold),
			PGBetID:      primitive.NilObjectID,
			BigReward:    bigReward,
			GameType:     define.GameType_Slot,
			AppID:        plr.AppID,
			Uid:          plr.Uid,
		}

		if msg.Header.Get("isstat") == "1" {
			msg.Header.Set("stat_bet", strconv.Itoa(int(betLog.Bet)))
			msg.Header.Set("stat_win", strconv.Itoa(int(betLog.Win)))
		}

		if msg.Header.Get("jump_log") != "1" {
			slotsmongo.AddBetLog(betLog)
		}
		//返回数据
		ret, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return nil
	})
	return
}
