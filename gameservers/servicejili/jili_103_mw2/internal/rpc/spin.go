package rpc

import (
	"context"
	"fmt"
	"serve/comm/redisx"
	"serve/servicejili/jili_103_mw2/internal/message"
	"strconv"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lang"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_103_mw2/internal"
	"serve/servicejili/jili_103_mw2/internal/gamedata"
	"serve/servicejili/jili_103_mw2/internal/models"

	"github.com/nats-io/nats.go"

	"regexp"
	"time"

	"serve/servicejili/jiliut/AckType"
	"serve/servicejili/jiliut/jiliUtMessage"

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
		gold := ut.Money2Gold(bet)
		costGold := ut.Money2Gold(bet)
		isBuy := spinReq.GetMallBet() > 0
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
		mw2Data := doc.Mw2Data
		for _, roundInfo := range mw2Data.GetRoundQueue() {
			ut.FloatPtrMul(roundInfo.RoundWin, bet)
			for _, comboStageData := range roundInfo.ComboStageData {
				ut.FloatPtrMul(comboStageData.ComboStageWin, bet)
				for _, awardData := range comboStageData.AwardDataVec {
					ut.FloatPtrMul(awardData.Win, bet)
				}
			}
		}
		ut.FloatPtrMul(mw2Data.TotalWin, bet)
		ut.FloatPtrMul(mw2Data.FreeTotalWin, bet)
		// 通用数据
		jiliUtData := doc.UtData
		ut.FloatPtrMul(jiliUtData.TotalWin, bet)

		//lo.Must0(ut.FloatEQ(spinAck.GetTotalWin(), totalWin))
		totalWin := mw2Data.GetTotalWin()

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
		//	mw2Data = doc.Mw2Data
		//	for _, roundInfo := range mw2Data.GetRoundQueue() {
		//		ut.FloatPtrMul(roundInfo.RoundWin, bet)
		//		for _, comboStageData := range roundInfo.ComboStageData {
		//			ut.FloatPtrMul(comboStageData.ComboStageWin, bet)
		//			for _, awardData := range comboStageData.AwardDataVec {
		//				ut.FloatPtrMul(awardData.Win, bet)
		//			}
		//		}
		//	}
		//	ut.FloatPtrMul(mw2Data.TotalWin, bet)
		//	ut.FloatPtrMul(mw2Data.FreeTotalWin, bet)
		//	// 通用数据
		//	jiliUtData = doc.UtData
		//	ut.FloatPtrMul(jiliUtData.TotalWin, bet)
		//
		//	realWin := mw2Data.GetTotalWin()
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
		mw2Data.NowMoney = proto.Float64(ut.Gold2Money(balanceGold))
		jiliUtData.PostMoney = mw2Data.NowMoney
		jiliUtData.Data, _ = proto.Marshal(mw2Data)

		var toSelfAwardPool, bigReward int64
		var disC = false

		if isBuy {
			toSelfAwardPool += gold * int64(mul-internal.BuyMul)
			disC = int64(mul-internal.BuyMul) != 0
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
		plr.OnSpinFinish(costGold, ut.Money2Gold(totalWin), false, false, 1)
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
			IsBuy:        false,
			Grade:        int(gold),
			PGBetID:      primitive.NilObjectID,
			BigReward:    bigReward,
			GameType:     define.GameType_Slot,
		}
		if ps.Header.Get("isstat") == "1" {
			ps.Header.Set("stat_bet", strconv.Itoa(int(betLog.Bet)))
			ps.Header.Set("stat_win", strconv.Itoa(int(betLog.Win)))
		}
		slotsmongo.AddBetLog(betLog)

		jiliUtData.RoundIndexV2 = proto.Int64(roundIndexV2)

		ret = jiliUtData
		doc.HistoryRecord.RoundIndex = showIndex
		his := doc.HistoryRecord
		costMoney := ut.Gold2Money(costGold)
		his.Bet = ut.Ftoa(costMoney)
		his.Win = ut.Ftoa(totalWin)
		his.NetValue = ut.Ftoa(totalWin - costMoney)
		his.PreMoney = ut.Ftoa(balance + costMoney)
		his.PostMoney = ut.Ftoa(mw2Data.GetNowMoney())
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

		for i, sm := range doc.SingleRoundLogSummaries {
			roundInfos := mw2Data.RoundQueue[i]

			plateInfos := hdoc.LogPlateInfos[sm.LogIndex]
			for j, csd := range roundInfos.ComboStageData {
				symbolLength := getSymbolLength(csd)
				plateInfos[j].SymbolLength = symbolLength
				plateInfos[j].Mult = fmt.Sprintf("%v", *csd.ComboStageMult)
			}
		}

		db.Collection("BetHistory").InsertOne(context.TODO(), hdoc)

		return nil
	})

	return
}

func getSymbolLength(stageInfo *message.Mw2_ComboStageInfo) (symbolLength string) {
	for _, columns := range stageInfo.ComboStageSymbol {
		for _, symbol := range columns.MColumn {
			if symbol.Length != nil {
				symbolLength = fmt.Sprintf("%v%v", symbolLength, *symbol.Length)
			}
		}
	}
	return
}
