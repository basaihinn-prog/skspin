package rpc

import (
	"context"
	"fmt"
	"regexp"
	"serve/comm/redisx"
	"strconv"
	"time"

	"serve/comm/define"
	"serve/comm/lazy"

	"serve/comm/db"
	"serve/comm/lang"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_208_phoenix/internal"
	"serve/servicejili/jili_208_phoenix/internal/gamedata"
	"serve/servicejili/jili_208_phoenix/internal/models"
	"serve/servicejili/jiliut/AckType"
	"serve/servicejili/jiliut/jiliUtMessage"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.Spin] = spin
}

var winReg = regexp.MustCompile(`^(\w+ game win: )([0-9.+-]+)$`)

func spin(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	var spinReq jiliUtMessage.Server_SpinReq
	err = proto.Unmarshal(data, &spinReq)
	if err != nil {
		return
	}

	//if slotsmongo.IsBanned(lazy.ServiceName, pid) {
	//	return nil, errors.New("Err:1302, Invalid player session")
	//}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		App, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		bet := spinReq.GetBet()
		if !ut.FloatInArr(ut.FloatArrMul(App.Cs, curItem.Multi), bet) {
			plr.RewriteLastData()
			return lang.Error(plr.Language, "下注额非法")
		}
		mallBet := spinReq.GetMallBet()
		gold := ut.Money2Gold(bet)
		costGold := ut.Money2Gold(bet)
		isBuy := mallBet != 0
		lo.Must0(bet > 0)

		now := time.Now()
		logid := primitive.NewObjectIDFromTimestamp(now)
		roundIndexV2 := now.UnixNano()/1000*1000 + internal.GameNo
		showIndex := fmt.Sprintf("%d", roundIndexV2)

		mul := int64(1)
		if isBuy {
			mul = int64(gamedata.GetSettings().BuyMulti)
			costGold *= mul
		}

		balanceGold, err := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  -costGold,
			RoundID: showIndex,
			Reason:  slotsmongo.ReasonBet,
		})
		if err != nil {
			return lang.Error(plr.Language, "金币不足")
		}
		balance := ut.Gold2Money(balanceGold)

		var forcedKill, hitBigAward, buyKill bool
		var doc *models.RawSpin
		selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(plr.PID)
		doc, hitBigAward, forcedKill, buyKill, err = nextPlayResp(plr, gold, selfPoolGold, isBuy, App)
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
		fmt.Println("转动数据：", "pid:", plr.PID, "dataId:", doc.ID.Hex())
		//var totalWin float64
		// 游戏内部数据
		phoenix_SpinAck := doc.Phoenix_SpinAck
		for _, roundQueue := range phoenix_SpinAck.GetRoundQueue() {
			ut.FloatPtrMul(roundQueue.RoundWin, bet)
			for _, awardDataVec := range roundQueue.AwardDataVec {
				ut.FloatPtrMul(awardDataVec.Win, bet)
			}

			for _, comboStageData := range roundQueue.ComboStageData {
				for _, tempAward := range comboStageData.TempAward {
					ut.FloatPtrMul(tempAward.Win, bet)
				}
			}
		}

		ut.FloatPtrMul(phoenix_SpinAck.TotalWin, bet)
		ut.FloatPtrMul(phoenix_SpinAck.FreeTotalWin, bet)
		// 通用数据
		jiliUtData := doc.UtData
		ut.FloatPtrMul(jiliUtData.TotalWin, bet)

		//lo.Must0(ut.FloatEQ(spinAck.GetTotalWin(), totalWin))
		totalWin := phoenix_SpinAck.GetTotalWin()

		//winMaxLimit := slotspool.GetMaxWin(plr.AppID, int64(float64(plr.GetWinLose())/10000), curItem.Multi)
		//if winMaxLimit != nil && totalWin > float64(winMaxLimit.MaxWin)*curItem.Multi {
		//	if isBuy {
		//		doc, err = gendata.GCombineDataMng.GetBuyMinData()
		//	} else {
		//		doc, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
		//		//doc, err = gendata.GCombineDataMng.SampleForceSimulate(ty,WitchBucket)
		//	}
		//	if err != nil {
		//		slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
		//			Pid:     plr.PID,
		//			Change:  costGold,
		//			RoundID: showIndex,
		//			Reason:  slotsmongo.ReasonRefund,
		//		})
		//		return err
		//	}
		//
		//	// 游戏内部数据
		//	phoenix_SpinAck = doc.Phoenix_SpinAck
		//	for _, roundQueue := range phoenix_SpinAck.GetRoundQueue() {
		//		ut.FloatPtrMul(roundQueue.RoundWin, bet)
		//		for _, awardDataVec := range roundQueue.AwardDataVec {
		//			ut.FloatPtrMul(awardDataVec.Win, bet)
		//		}
		//
		//		for _, comboStageData := range roundQueue.ComboStageData {
		//			for _, tempAward := range comboStageData.TempAward {
		//				ut.FloatPtrMul(tempAward.Win, bet)
		//			}
		//		}
		//	}
		//
		//	ut.FloatPtrMul(phoenix_SpinAck.TotalWin, bet)
		//	ut.FloatPtrMul(phoenix_SpinAck.FreeTotalWin, bet)
		//	// 通用数据
		//	jiliUtData = doc.UtData
		//	ut.FloatPtrMul(jiliUtData.TotalWin, bet)
		//
		//	realWin := phoenix_SpinAck.GetTotalWin()
		//	// 大奖处理
		//	slotsmongo.BigRewardDeal(&slotsmongo.BigRewardDealParams{
		//		Pid:         plr.PID,
		//		AppID:       plr.AppID,
		//		CurrencyKey: curItem.Key,
		//		ServiceName: lazy.ServiceName,
		//		OriginWin:   totalWin,
		//		RealWin:     realWin,
		//		WinLose:     plr.GetWinLose(),
		//		Desc:        winMaxLimit.Desc,
		//	})
		//	totalWin = realWin
		//	hitBigAward = false
		//	forcedKill = false
		//}

		balanceGold, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  ut.Money2Gold(totalWin),
			RoundID: showIndex,
			Reason:  slotsmongo.ReasonWin,
		})
		if err != nil {
			return err
		}
		phoenix_SpinAck.NowMoney = proto.Float64(ut.Gold2Money(balanceGold))
		jiliUtData.PostMoney = phoenix_SpinAck.NowMoney
		jiliUtData.Data, _ = proto.Marshal(phoenix_SpinAck)

		var toSelfAwardPool, bigReward int64
		var disC = false
		if isBuy {
			toSelfAwardPool += gold * int64(gamedata.GetSettings().BuyMulti-internal.BuyMul)
			disC = int64(gamedata.GetSettings().BuyMulti-internal.BuyMul) != 0
			if buyKill {
				toSelfAwardPool += costGold
			}
		} else {
			if hitBigAward {
				toSelfAwardPool += -ut.Money2Gold(totalWin)
				bigReward = ut.Money2Gold(totalWin)
			}
			//else {
			//				poolCost := gendata.GBuckets.GetPoolCost(doc.BucketId)
			//				if poolCost > 0 {
			//					toSelfAwardPool += -gold * int64(poolCost)
			//				}
			//			}
		}

		if !isBuy || buyKill || (isBuy && disC) {
			plr.UpdatePool(costGold, selfPoolGold, toSelfAwardPool, forcedKill, isBuy, App)
		}
		plr.OnSpinFinish(costGold, ut.Money2Gold(totalWin), isBuy, false, 1)
		betLog := &slotsmongo.AddBetLogParams{
			UserName:     "",
			CurrencyKey:  plr.CurrencyKey,
			ID:           showIndex,
			Pid:          plr.PID,
			Bet:          costGold,
			Win:          ut.Money2Gold(totalWin),
			Balance:      balanceGold,
			RoundID:      showIndex,
			Completed:    true,
			TotalWinLoss: ut.Money2Gold(totalWin) - costGold,
			IsBuy:        isBuy,
			Grade:        int(gold),
			PGBetID:      primitive.NilObjectID,
			BigReward:    bigReward,
			GameType:     define.GameType_Slot,
		}

		if ps.Header.Get("isstat") == "1" {
			ps.Header.Set("stat_bet", strconv.Itoa(int(betLog.Bet)))
			ps.Header.Set("stat_win", strconv.Itoa(int(betLog.Win)))
		}

		if ps.Header.Get("jump_log") != "1" {
			slotsmongo.AddBetLog(betLog)
		}

		jiliUtData.RoundIndexV2 = proto.Int64(roundIndexV2)
		ret = jiliUtData

		doc.HistoryRecord.RoundIndex = showIndex
		his := doc.HistoryRecord
		costMoney := ut.Gold2Money(costGold)
		his.Bet = ut.Ftoa(costMoney)
		his.Win = ut.Ftoa(totalWin)
		his.NetValue = ut.Ftoa(totalWin - costMoney)
		his.PreMoney = ut.Ftoa(balance + costMoney)
		his.PostMoney = ut.Ftoa(phoenix_SpinAck.GetNowMoney())
		his.CreateTime = now.UnixMilli()

		for _, sm := range doc.SingleRoundLogSummaries {
			for i, desc := range sm.Desc {
				// "Free game win: 0.4",
				// "Main game win: 1.6"

				rets := winReg.FindStringSubmatch(desc)
				if len(rets) == 3 {
					winstr := rets[2]
					ut.FloatStrPtrMul(&winstr, bet)

					ret := rets[1] + winstr

					sm.Desc[i] = ret
				}
			}
		}

		for _, pinfos := range doc.LogPlateInfos {
			for _, pinfo := range pinfos {
				for _, v := range pinfo.List {
					ut.FloatStrPtrMul(&v.W, bet)
				}
			}
		}

		// doc.ID = logid

		hdoc := models.HistoryDoc{
			ID:                      logid,
			Pid:                     pid,
			HistoryRecord:           doc.HistoryRecord,
			SingleRoundLogSummaries: doc.SingleRoundLogSummaries,
			LogPlateInfos:           doc.LogPlateInfos,
			Tid:                     showIndex,
			OId:                     doc.ID,
		}

		if ps.Header.Get("jump_log") != "1" {
			db.Collection("BetHistory").InsertOne(context.TODO(), hdoc)
		}
		return nil
	})

	return
}
