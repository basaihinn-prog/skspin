package rpc

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"serve/comm/redisx"
	"strconv"
	"time"

	"serve/comm/define"
	"serve/comm/lazy"

	"serve/servicejili/jili_142_bh/internal/gamedata"

	"serve/comm/lang"
	"serve/servicejili/jili_142_bh/internal"
	"serve/servicejili/jili_142_bh/internal/message"
	"serve/servicejili/jili_142_bh/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/slotsmongo"
	"serve/comm/ut"

	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

func init() {
	// mux.RegRpc("/csh/game/spin", "spin", "game-api", spin, nil)
	jiliut.RegRpc(fmt.Sprintf("/%s/game/spin", internal.GameShortName), spin)
}

var winReg = regexp.MustCompile(`^(\w+ game win: )([0-9.+-]+)$`)

func spin(ps *nats.Msg) (ret []byte, err error) {
	var gameReq message.Bh_GameReqData
	err = proto.Unmarshal(ps.Data, &gameReq)
	if err != nil {
		return
	}
	token := gameReq.GetToken()

	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	// gold, err := slotsmongo.GetBalance(pid)
	// if err != nil {
	// 	return err
	// }
	// balance := ut.Gold2Money(gold)
	//if slotsmongo.IsBanned(lazy.ServiceName, pid) {
	//	return nil, errors.New("Err:1302, Invalid player session")
	//}

	var spinReq message.Bh_SpinReq
	err = proto.Unmarshal(gameReq.Encode, &spinReq)
	if err != nil {
		return
	}

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
		showIndex := fmt.Sprintf("%d%d", now.UnixNano(), rand.Int()%10)

		mul := int64(1)
		if isBuy {
			mul = int64(gamedata.GetSettings().BuyMulti)
			costGold *= mul
		}

		//balanceGold, err := slotsmongo.ModifyGold(plr.PID, -costGold, "下注")
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

		// Teemo, replace the previous line for test
		// doc, hitBigAward, forcedKill, buyKill, err = findDoc("6683be0f503aea7a1c174aaa")

		if err != nil {
			//slotsmongo.ModifyGold(plr.PID, -gold*mul, "下注-退回, err:"+err.Error())
			slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
				Pid:     plr.PID,
				Change:  costGold,
				RoundID: showIndex,
				Reason:  slotsmongo.ReasonRefund,
			})
			return err
		}
		fmt.Println("转动数据：", "pid:", plr.PID, "costGold:", costGold, "dataId:", doc.ID.Hex(), "HBA:", hitBigAward)
		spinAck := doc.Data
		// 处理数据
		for _, round := range spinAck.MainGame {
			for _, platedata := range round.PlateQueue {
				ut.FloatPtrMul(platedata.SingleWin, bet)
				ut.FloatPtrDiv(platedata.SingleWin, internal.BaseBet)
				for _, awarddatavec := range platedata.AwardDataVec {
					ut.FloatPtrMul(awarddatavec.Win, bet)
					ut.FloatPtrDiv(awarddatavec.Win, internal.BaseBet)
				}
			}
			ut.FloatPtrMul(round.RoundWin, bet)
			ut.FloatPtrDiv(round.RoundWin, internal.BaseBet)
		}
		ut.FloatPtrMul(spinAck.TotalWin, bet)
		ut.FloatPtrDiv(spinAck.TotalWin, internal.BaseBet)
		totalWin := spinAck.GetTotalWin()
		for _, freeGameInfo := range spinAck.FreeGame {
			ut.FloatPtrMul(freeGameInfo.FreeTotalWin, bet)
			ut.FloatPtrDiv(freeGameInfo.FreeTotalWin, internal.BaseBet)
			for _, queuedata := range freeGameInfo.FreeQueue {
				ut.FloatPtrMul(queuedata.RoundWin, bet)
				ut.FloatPtrDiv(queuedata.RoundWin, internal.BaseBet)
				for _, platedata := range queuedata.PlateQueue {
					ut.FloatPtrMul(platedata.SingleWin, bet)
					ut.FloatPtrDiv(platedata.SingleWin, internal.BaseBet)
					for _, awarddata := range platedata.AwardDataVec {
						ut.FloatPtrMul(awarddata.Win, bet)
						ut.FloatPtrDiv(awarddata.Win, internal.BaseBet)
					}
				}
			}
		}

		//if totalWin > 0 {
		//	balance, err := slotsmongo.ModifyGold(plr.PID, ut.Money2Gold(totalWin), "赢分")
		//	if err != nil {
		//		return err
		//	}
		//	spinAck.NowMoney = proto.Float64(ut.Gold2Money(balance))
		//}

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
		//	spinAck = doc.Data
		//	// 处理数据
		//	for _, round := range spinAck.MainGame {
		//		for _, platedata := range round.PlateQueue {
		//			ut.FloatPtrMul(platedata.SingleWin, bet)
		//			ut.FloatPtrDiv(platedata.SingleWin, internal.BaseBet)
		//			for _, awarddatavec := range platedata.AwardDataVec {
		//				ut.FloatPtrMul(awarddatavec.Win, bet)
		//				ut.FloatPtrDiv(awarddatavec.Win, internal.BaseBet)
		//			}
		//		}
		//		ut.FloatPtrMul(round.RoundWin, bet)
		//		ut.FloatPtrDiv(round.RoundWin, internal.BaseBet)
		//	}
		//	ut.FloatPtrMul(spinAck.TotalWin, bet)
		//	ut.FloatPtrDiv(spinAck.TotalWin, internal.BaseBet)
		//	for _, freeGameInfo := range spinAck.FreeGame {
		//		for _, queuedata := range freeGameInfo.FreeQueue {
		//			ut.FloatPtrMul(queuedata.RoundWin, bet)
		//			ut.FloatPtrDiv(queuedata.RoundWin, internal.BaseBet)
		//			for _, platedata := range queuedata.PlateQueue {
		//				ut.FloatPtrMul(platedata.SingleWin, bet)
		//				ut.FloatPtrDiv(platedata.SingleWin, internal.BaseBet)
		//				for _, awarddata := range platedata.AwardDataVec {
		//					ut.FloatPtrMul(awarddata.Win, bet)
		//					ut.FloatPtrDiv(awarddata.Win, internal.BaseBet)
		//				}
		//			}
		//		}
		//	}
		//
		//	realWin := spinAck.GetTotalWin()
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
		spinAck.NowMoney = proto.Float64(ut.Gold2Money(balanceGold))

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

		spinAck.ShowIndex = proto.String(fmt.Sprintf("%s-%s-%s", showIndex[0:5], showIndex[5:11], showIndex[11:19]))
		var spinAll = message.Bh_SpinAllData{
			Info: []*message.Bh_SpinAck{
				spinAck,
			}}
		encode, _ := proto.Marshal(&spinAll)

		var resData = message.Bh_ResData{
			Type:  proto.Int32(AckType["spin"]),
			Token: proto.String(token),
			Data: []*message.Bh_InfoData{
				{
					Encode: encode,
				},
			},
		}
		ret, _ = proto.Marshal(&resData)
		doc.HistoryRecord.RoundIndex = showIndex
		his := doc.HistoryRecord
		costMoney := ut.Gold2Money(costGold)
		his.Bet = ut.Ftoa(costMoney)
		his.Win = ut.Ftoa(totalWin)
		his.NetValue = ut.Ftoa(totalWin - costMoney)
		his.PreMoney = ut.Ftoa(balance + costMoney)
		his.PostMoney = ut.Ftoa(spinAck.GetNowMoney())
		his.CreateTime = now.UnixMilli()
		for _, sm := range doc.SingleRoundLogSummaries {
			for i, desc := range sm.Desc {
				// "Free game win: 0.4",
				// "Main game win: 1.6"
				rets := winReg.FindStringSubmatch(desc)
				if len(rets) == 3 {
					winstr := rets[2]
					ut.FloatStrPtrMul(&winstr, bet)
					ut.FloatStrPtrDiv(&winstr, internal.BaseBet)
					ret := rets[1] + winstr
					sm.Desc[i] = ret
				}
			}
		}
		for _, pinfos := range doc.LogPlateInfos {
			for _, pinfo := range pinfos {
				ut.FloatStrPtrMul(&pinfo.Win, bet)
				ut.FloatStrPtrDiv(&pinfo.Win, internal.BaseBet)
				for _, v := range pinfo.List {
					ut.FloatStrPtrMul(&v.W, bet)
					ut.FloatStrPtrDiv(&v.W, internal.BaseBet)
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

		if ps.Header.Get("jump_log") != "1" {
			db.Collection("BetHistory").InsertOne(context.TODO(), hdoc)
		}
		return nil
	})

	return
}
