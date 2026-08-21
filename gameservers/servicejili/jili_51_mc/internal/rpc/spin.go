package rpc

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"serve/comm/redisx"
	"strconv"
	"time"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lang"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_51_mc/internal"
	"serve/servicejili/jili_51_mc/internal/models"
	"serve/servicejili/jiliut/AckType"

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
	var spinReq serverOfficial.SpinReq
	err = proto.Unmarshal(data, &spinReq)
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {

		bet := spinReq.GetBet()
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)

		App, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}

		if !ut.FloatInArr(ut.FloatArrMul(App.Cs, curItem.Multi), bet) {
			plr.RewriteLastData()
			return lang.Error(plr.Language, "下注额非法")
		}

		gold := ut.Money2Gold(bet)
		costGold := ut.Money2Gold(bet)
		lo.Must0(bet > 0)

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
		ty := internal.BetMap[bet]
		var forcedKill, hitBigAward bool
		var doc *models.RawSpin
		selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(plr.PID)
		doc, hitBigAward, forcedKill, _, err = nextPlayResp(plr, costGold, selfPoolGold, ty, App)
		if err != nil {
			slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
				Pid:     plr.PID,
				Change:  costGold,
				RoundID: showIndex,
				Reason:  slotsmongo.ReasonRefund,
			})
			return err
		}
		//hexid, _ := primitive.ObjectIDFromHex(`67ee4ea962de760ae0cdfeb1`)
		//doc, err = gendata.FindFromSimulate(hexid)
		//if err != nil {
		//	slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
		//		Pid:     plr.PID,
		//		Change:  costGold,
		//		RoundID: showIndex,
		//		Reason:  slotsmongo.ReasonRefund,
		//	})
		//	return err
		//}
		fmt.Println("转动数据：", "pid:", plr.PID, "dataId:", doc.ID.Hex())

		//游戏内部数据
		olsData := doc.McData
		// 通用数据
		jiliUtData := doc.UtData
		//totalWinUtData := &(jiliUtData.TotalWin)
		//ut.FloatPtrMul(totalWinUtData, bet)
		totalWin := jiliUtData.GetTotalWin()

		balanceGold, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  ut.Money2Gold(totalWin),
			RoundID: showIndex,
			Reason:  slotsmongo.ReasonWin,
		})
		if err != nil {
			return err
		}
		jiliUtData.PostMoney = ut.Gold2Money(balanceGold)
		jiliUtData.Data, _ = proto.Marshal(olsData)

		var toSelfAwardPool, bigReward int64
		if hitBigAward {
			toSelfAwardPool += -ut.Money2Gold(totalWin)
			bigReward = ut.Money2Gold(totalWin)
		}

		plr.UpdatePool(costGold, selfPoolGold, toSelfAwardPool, forcedKill, false, App)
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
		slotsmongo.AddBetLog(betLog)
		if ps.Header.Get("isstat") == "1" {
			ps.Header.Set("stat_bet", strconv.Itoa(int(betLog.Bet)))
			ps.Header.Set("stat_win", strconv.Itoa(int(betLog.Win)))
		}

		jiliUtData.Data, _ = proto.Marshal(olsData)
		jiliUtData.Service = &serverOfficial.ServiceData{}
		//jiliUtData.TotalWin = totalWin
		//jiliUtData.PostMoney = 0
		jiliUtData.Hasspin = doc.HasGame
		jiliUtData.BaseBet = bet
		jiliUtData.RealBet = bet
		jiliUtData.RoundIndexV2 = roundIndexV2
		jiliUtData.SpinReq = &spinReq

		ret = jiliUtData

		doc.HistoryRecord.RoundIndex = showIndex
		his := doc.HistoryRecord
		costMoney := ut.Gold2Money(costGold)
		his.Bet = ut.Ftoa(costMoney)
		his.Win = ut.Ftoa(totalWin)
		his.NetValue = ut.Ftoa(totalWin - costMoney)
		his.PreMoney = ut.Ftoa(balance + costMoney)
		//his.PostMoney = ut.Ftoa(zeusData.GetNowMoney())
		nowMoney, _ := slotsmongo.GetBalance(plr.PID)
		nowBalance := ut.Gold2Money(nowMoney)
		his.PostMoney = ut.Ftoa(nowBalance)
		his.CreateTime = now.UnixMilli()

		//for _, sm := range doc.SingleRoundLogSummaries {
		//	for i, desc := range sm.Desc {
		//		// "Free game win: 0.4",
		//		// "Main game win: 1.6"
		//		rets := winReg.FindStringSubmatch(desc)
		//		if len(rets) == 3 {
		//			winstr := rets[2]
		//			ut.FloatStrPtrMul(&winstr, bet)
		//			ret := rets[1] + winstr
		//
		//			sm.Desc[i] = ret
		//		}
		//	}
		//}
		//
		//for _, pinfos := range doc.LogPlateInfos {
		//	for _, pinfo := range pinfos {
		//		for _, v := range pinfo.List {
		//			ut.FloatStrPtrMul(&v.W, bet)
		//		}
		//	}
		//}

		// doc.ID = betId

		hdoc := models.HistoryDoc{
			ID:                      logid,
			Pid:                     pid,
			HistoryRecord:           doc.HistoryRecord,
			SingleRoundLogSummaries: doc.SingleRoundLogSummaries,
			LogPlateInfos:           doc.LogPlateInfos,
			Tid:                     showIndex,
			OId:                     doc.ID,
		}

		_, err = db.Collection("BetHistory").InsertOne(context.TODO(), hdoc)
		if err != nil {
			slog.Error("insert bethistory fail", logid)
		}

		return nil
	})

	return
}
