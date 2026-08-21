package rpc

import (
	"context"
	"fmt"
	"regexp"
	"serve/comm/redisx"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"strconv"
	"time"

	"serve/comm/define"
	"serve/comm/lazy"

	"serve/comm/db"
	"serve/comm/lang"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_33_pp/internal"
	"serve/servicejili/jili_33_pp/internal/models"
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
var winReg2 = regexp.MustCompile(`^(\w+ Win: )([0-9.+-]+)$`)

//var myVar = int64(1000) //该游戏特殊处理 用来对一些数据进行处理

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
		//bet := float64(spinReq.GetBet()) / float64(myVar)
		bet := spinReq.GetBet()
		if !ut.FloatInArr(ut.FloatArrMul(App.Cs, curItem.Multi), bet) {
			plr.RewriteLastData()
			return lang.Error(plr.Language, "下注额非法")
		}
		//mallBet := spinReq.GetMallBet()
		mallBet := 0.0
		if spinReq.Mall != nil && spinReq.Mall.Bet > 0 {
			mallBet = spinReq.Mall.Bet
		}
		isBuy := mallBet != 0
		//mul := int64(1)
		gold := ut.Money2Gold(bet)
		costGold := ut.Money2Gold(bet)
		//isBuy := mallBet != 0
		lo.Must0(bet > 0)
		if isBuy {
			//ty = internal.GameTypeGame
			mul := 29.5
			temp := float64(costGold)
			temp *= mul
			costGold = int64(temp)
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

		//id, _ := primitive.ObjectIDFromHex(`67f62dee1460e571746b96db`)
		//doc, _ = gendata.FindFromSimulate(id)

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
		ppData := doc.PpData
		jiliUtData := doc.UtData
		betMul := bet / jiliUtData.BaseBet
		for _, roundInfo := range ppData.GetPlateInfo() {
			for _, AwardData := range roundInfo.AwardDataList {
				ut.FloatPtrMul(&AwardData.AwardMoney, betMul)
			}
			if roundInfo.GetPlateWin() != 0 {
				ut.FloatPtrMul(&roundInfo.PlateWin, betMul)
			}
			if roundInfo.GetLineWin() != 0 {
				ut.FloatPtrMul(&roundInfo.LineWin, betMul)
			}
			if roundInfo.GetJPWin() != 0 {
				ut.FloatPtrMul(&roundInfo.JPWin, betMul)
			}
		}

		//lo.Must0(ut.FloatEQ(spinAck.GetTotalWin(), totalWin))
		ut.FloatPtrMul(&ppData.TotalWin, betMul)

		ut.FloatPtrMul(&jiliUtData.TotalWin, betMul)
		totalWin := jiliUtData.TotalWin

		balanceGold, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  ut.Money2Gold(float64(totalWin)),
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
			//toSelfAwardPool += gold * int64(gamedata.Settings.BuyMulti-internal.BuyMul)
			//disC = int64(gamedata.Settings.BuyMulti-internal.BuyMul) != 0
			if buyKill {
				toSelfAwardPool += costGold
			}
		} else {
			if hitBigAward {
				toSelfAwardPool += -ut.Money2Gold(float64(totalWin))
				bigReward = ut.Money2Gold(float64(totalWin))
			}
		}

		if !isBuy || buyKill || (isBuy && disC) {
			plr.UpdatePool(costGold, selfPoolGold, toSelfAwardPool, forcedKill, isBuy, App)
		}
		plr.OnSpinFinish(costGold, ut.Money2Gold(float64(totalWin)), isBuy, false, 1)
		betLog := &slotsmongo.AddBetLogParams{
			UserName:     "",
			CurrencyKey:  plr.CurrencyKey,
			ID:           showIndex,
			Pid:          plr.PID,
			Bet:          costGold,
			Win:          ut.Money2Gold(float64(totalWin)),
			Balance:      balanceGold,
			RoundID:      showIndex,
			Completed:    true,
			TotalWinLoss: ut.Money2Gold(float64(totalWin)) - costGold,
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
		slotsmongo.AddBetLog(betLog)

		jiliUtData.TotalWin = totalWin
		jiliUtData.PostMoney = postMoney
		jiliUtData.BaseBet = bet
		jiliUtData.RealBet = costMoney
		jiliUtData.RoundIndex = roundIndexV2
		jiliUtData.SpinReq = &spinReq
		jiliUtData.Data, _ = proto.Marshal(ppData)
		ret = jiliUtData

		doc.HistoryRecord.RoundIndex = showIndex
		his := doc.HistoryRecord
		his.Bet = ut.Ftoa(costMoney)
		his.Win = ut.Ftoa(float64(totalWin))
		his.NetValue = ut.Ftoa(float64(totalWin) - costMoney)
		his.PreMoney = ut.Ftoa(balance + costMoney)
		his.CreateTime = now.UnixMilli()
		for _, sm := range doc.SingleRoundLogSummaries {
			for i, desc := range sm.Desc {
				rets := winReg.FindStringSubmatch(desc)
				if len(rets) == 3 {
					winstr := rets[2]
					ut.FloatStrPtrMul(&winstr, betMul)
					ret := rets[1] + winstr
					sm.Desc[i] = ret
				}
				rets2 := winReg2.FindStringSubmatch(desc)
				if len(rets2) == 3 {
					winstr := rets2[2]
					ut.FloatStrPtrMul(&winstr, betMul)
					ret := rets2[1] + winstr
					sm.Desc[i] = ret
				}
			}
		}
		for _, pinfos := range doc.LogPlateInfos {
			for _, pinfo := range pinfos {
				for _, v := range pinfo.List {
					ut.FloatStrPtrMul(&v.W, betMul)
				}
			}
		}
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
