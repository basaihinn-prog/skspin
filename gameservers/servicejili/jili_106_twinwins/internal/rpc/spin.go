package rpc

import (
	"context"
	"fmt"
	"regexp"
	"serve/comm/redisx"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"serve/comm/define"
	"serve/comm/lazy"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"serve/comm/db"
	"serve/comm/lang"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_106_twinwins/internal"
	"serve/servicejili/jili_106_twinwins/internal/models"
	"serve/servicejili/jiliut/AckType"
)

func init() {
	reqMux[AckType.Spin] = spin
}

var winReg = regexp.MustCompile(`^(\w+ game win: )([0-9.+-]+)$`)
var winReg2 = regexp.MustCompile(`^(Bonus win : )([0-9.+-]+)$`)

func spin(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	var spinReq serverOfficial.SpinReq
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
		extra := spinReq.GetExtraSpin()
		costGold := ut.Money2Gold(bet)
		ty := internal.GameTypeNormal
		lo.Must0(bet > 0)
		if extra != nil {
			ty = internal.GameTypeExtra
			costGold = int64(float64(costGold) * 1.5)
		}

		mallBet := 0.0
		if spinReq.Mall != nil {
			mallBet = spinReq.Mall.Bet
		}
		isBuy := mallBet != 0
		mul := int64(1)
		if isBuy {
			ty = internal.GameTypeGame
			mul = 46
			costGold *= mul
		}

		now := time.Now()
		logid := primitive.NewObjectIDFromTimestamp(now)
		roundIndexV2 := now.UnixNano()/1000*1000 + internal.GameNo
		showIndex := fmt.Sprintf("%d", roundIndexV2)

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
		doc, hitBigAward, forcedKill, buyKill, err = nextPlayResp(plr, costGold, selfPoolGold, ty, App)
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
		twinwinsData := doc.TwinwinsData
		// 通用数据
		jiliUtData := doc.UtData
		betMul := bet / jiliUtData.BaseBet
		ut.FloatPtrMul(&jiliUtData.TotalWin, bet)
		ut.FloatPtrDiv(&jiliUtData.TotalWin, internal.BaseBet)
		totalWin := jiliUtData.TotalWin

		for _, roundInfo := range twinwinsData.GetPlate() {
			ut.FloatPtrMul(&roundInfo.Win, betMul)
		}

		//lo.Must0(ut.FloatEQ(spinAck.GetTotalWin(), totalWin))

		//winMaxLimit := slotspool.GetMaxWin(plr.AppID, int64(float64(plr.GetWinLose())/10000), curItem.Multi)
		//if winMaxLimit != nil && totalWin > float64(winMaxLimit.MaxWin)*curItem.Multi {
		//	if isBuy {
		//		doc, err = gendata.GCombineDataMng.GetBuyMinData()
		//	} else {
		//		doc, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
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
		//	mw3Data = doc.Mw3Data
		//
		//	for _, roundInfo := range mw3Data.GetRoundQueue() {
		//		ut.FloatPtrMul(roundInfo.RoundWin, bet)
		//		ut.FloatPtrMul(roundInfo.WaysWin, bet)
		//		ut.FloatPtrMul(roundInfo.LampWin, bet)
		//
		//		for _, awardInfo := range roundInfo.GetAwardDataVec() {
		//			ut.FloatPtrMul(awardInfo.Win, bet)
		//		}
		//	}
		//	ut.FloatPtrMul(mw3Data.FreeTotalWin, bet)
		//	ut.FloatPtrMul(mw3Data.TotalWin, bet)
		//
		//	// 通用数据
		//	jiliUtData = doc.UtData
		//	for _, serviceInfo := range jiliUtData.GetService() {
		//		ut.FloatPtrMul(serviceInfo.JpWin, bet)
		//		ut.FloatPtrMul(serviceInfo.FullJpWin, bet)
		//	}
		//	ut.FloatPtrMul(jiliUtData.TotalWin, bet)
		//
		//	realWin := mw3Data.GetTotalWin()
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

		costMoney := ut.Gold2Money(costGold)
		postMoney := ut.Gold2Money(balanceGold)

		var toSelfAwardPool, bigReward int64
		var disC = false
		if isBuy {
			if buyKill {
				toSelfAwardPool += costGold
			}
		} else {
			if hitBigAward {
				toSelfAwardPool += -ut.Money2Gold(totalWin)
				bigReward = ut.Money2Gold(totalWin)
			}
		}
		if !isBuy || buyKill || (isBuy && disC) {
			plr.UpdatePool(costGold, selfPoolGold, toSelfAwardPool, forcedKill, ty, App)
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
			IsBuy:        false,
			Grade:        int(gold),
			PGBetID:      primitive.NilObjectID,
			BigReward:    bigReward,
			GameType:     define.GameType_Slot,
			Extra:        extra != nil,
		}
		if ps.Header.Get("isstat") == "1" {
			ps.Header.Set("stat_bet", strconv.Itoa(int(betLog.Bet)))
			ps.Header.Set("stat_win", strconv.Itoa(int(betLog.Win)))
		}
		if ps.Header["jump_log"] == nil {
			slotsmongo.AddBetLog(betLog)
		}

		jiliUtData.Data, _ = proto.Marshal(twinwinsData)
		jiliUtData.Service = &serverOfficial.ServiceData{}
		jiliUtData.TotalWin = totalWin
		jiliUtData.PostMoney = postMoney
		jiliUtData.Hasspin = doc.HasGame
		jiliUtData.BaseBet = bet
		jiliUtData.RealBet = bet
		jiliUtData.RoundIndexV2 = roundIndexV2
		jiliUtData.SpinReq = &spinReq

		ret = jiliUtData

		doc.HistoryRecord.RoundIndex = showIndex
		his := doc.HistoryRecord
		his.Bet = ut.Ftoa(costMoney)
		his.Win = ut.Ftoa(totalWin)
		his.NetValue = ut.Ftoa(totalWin - costMoney)
		his.PreMoney = ut.Ftoa(balance + costMoney)
		his.PostMoney = ut.Ftoa(balance)
		his.CreateTime = now.UnixMilli()

		for _, sm := range doc.SingleRoundLogSummaries {
			for i, desc := range sm.Desc {
				// "Free game win: 0.4",
				// "Main game win: 1.6"

				rets := winReg.FindStringSubmatch(desc)
				if len(rets) == 3 {
					winstr := rets[2]
					ut.FloatStrPtrMul(&winstr, betMul)

					ret := rets[1] + winstr

					sm.Desc[i] = ret
				} else {
					rets = winReg2.FindStringSubmatch(desc)
					if len(rets) == 3 {
						winstr := rets[2]
						ut.FloatStrPtrMul(&winstr, betMul)

						ret := rets[1] + winstr

						sm.Desc[i] = ret
					}
				}
			}
		}

		for _, pinfos := range doc.LogPlateInfos {
			for _, pinfo := range pinfos {
				for _, v := range pinfo.List {
					ut.FloatStrPtrMul(&v.W, betMul)
				}
				splits := strings.Split(pinfo.LampTag, ",")
				for i := range splits {
					ut.FloatStrPtrMul1(&splits[i], betMul)
				}
				pinfo.LampTag = strings.Join(splits, ",")
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

		db.Collection("BetHistory").InsertOne(context.TODO(), hdoc)

		return nil
	})

	return
}
