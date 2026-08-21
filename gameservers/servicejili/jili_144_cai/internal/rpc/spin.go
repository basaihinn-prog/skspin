package rpc

import (
	"context"
	"fmt"
	"math/rand"
	"serve/comm/redisx"
	"strconv"

	"serve/comm/define"
	"serve/comm/lazy"
	"serve/servicejili/jili_144_cai/internal"
	"serve/servicejili/jili_144_cai/internal/message"
	"serve/servicejili/jili_144_cai/internal/models"

	"regexp"
	"time"

	"serve/comm/lang"

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

var winReg = regexp.MustCompile(`^(.+ win: )([0-9.+-]+)$`)

func spin(ps *nats.Msg) (ret []byte, err error) {
	var gameReq message.Cai_GameReqData
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

	var spinReq message.Cai_SpinReq
	err = proto.Unmarshal(gameReq.Encode, &spinReq)
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		extra := spinReq.GetExtra()
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
		baseBet := float64(1)
		ty := internal.GameTypeNormal
		if extra == 1 {
			baseBet = 2
			gold = ut.Money2Gold(bet / 10)
			costGold = gold * 5
			ty = internal.GameTypeExtra
		}
		if extra == 2 {
			baseBet = 10
			gold = ut.Money2Gold(bet / 10)
			costGold = gold * 50
			ty = internal.GameTypeExtraPlus
		}
		lo.Must0(bet > 0)

		now := time.Now()
		logid := primitive.NewObjectIDFromTimestamp(now)
		showIndex := fmt.Sprintf("%d%d", now.UnixNano(), rand.Int()%10)
		//balanceGold, err := slotsmongo.ModifyGold(plr.PID, -gold, "下注")
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

		var forcedKill, hitBigAward bool
		var doc *models.RawSpin
		selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(plr.PID)
		doc, hitBigAward, forcedKill, err = nextPlayResp(plr, costGold, selfPoolGold, ty, App)
		if err != nil {
			//slotsmongo.ModifyGold(plr.PID, -gold, "下注-退回, err:"+err.Error())
			slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
				Pid:     plr.PID,
				Change:  costGold,
				RoundID: showIndex,
				Reason:  slotsmongo.ReasonRefund,
			})
			return err
		}
		fmt.Println("转动数据：", "pid:", plr.PID)
		spinAck := doc.Data
		// 处理数据
		for _, round := range spinAck.RoundQueue {
			for i := range round.AwardDataVec {
				ut.FloatPtrMul(round.AwardDataVec[i].Win, bet)
				ut.FloatPtrDiv(round.AwardDataVec[i].Win, baseBet)
			}
			ut.FloatPtrMul(round.RoundWin, bet)
			ut.FloatPtrDiv(round.RoundWin, baseBet)
		}
		ut.FloatPtrMul(spinAck.TotalWin, bet)
		ut.FloatPtrDiv(spinAck.TotalWin, baseBet)
		totalWin := spinAck.GetTotalWin()

		//winMaxLimit := slotspool.GetMaxWin(plr.AppID, int64(float64(plr.GetWinLose())/10000), curItem.Multi)
		//if winMaxLimit != nil && totalWin > float64(winMaxLimit.MaxWin)*curItem.Multi {
		//	doc, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
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
		//	for _, round := range spinAck.RoundQueue {
		//		for i := range round.AwardDataVec {
		//			ut.FloatPtrMul(round.AwardDataVec[i].Win, bet)
		//			ut.FloatPtrDiv(round.AwardDataVec[i].Win, baseBet)
		//		}
		//		ut.FloatPtrMul(round.RoundWin, bet)
		//		ut.FloatPtrDiv(round.RoundWin, baseBet)
		//	}
		//	ut.FloatPtrMul(spinAck.TotalWin, bet)
		//	ut.FloatPtrDiv(spinAck.TotalWin, baseBet)
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
		if hitBigAward {
			toSelfAwardPool += -ut.Money2Gold(totalWin)
			bigReward = ut.Money2Gold(totalWin)
		}
		// else {
		//			poolCost := gendata.GBuckets.GetPoolCost(doc.BucketId)
		//			if poolCost > 0 {
		//				toSelfAwardPool += -gold * int64(poolCost)
		//			}
		//		}

		plr.UpdatePool(costGold, selfPoolGold, toSelfAwardPool, forcedKill, ty, App)
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
			Extra:        extra > 0,
		}
		if ps.Header.Get("isstat") == "1" {
			ps.Header.Set("stat_bet", strconv.Itoa(int(betLog.Bet)))
			ps.Header.Set("stat_win", strconv.Itoa(int(betLog.Win)))
		}
		if ps.Header.Get("jump_log") != "1" {
			slotsmongo.AddBetLog(betLog)
		}
		spinAck.ShowIndex = proto.String(fmt.Sprintf("%s-%s-%s", showIndex[0:5], showIndex[5:11], showIndex[11:19]))
		var spinAll = message.Cai_SpinAllData{
			Info: []*message.Cai_SpinAck{
				spinAck,
			}}
		encode, _ := proto.Marshal(&spinAll)

		var resData = message.Cai_ResData{
			Type:  proto.Int32(AckType["spin"]),
			Token: proto.String(token),
			Data: []*message.Cai_InfoData{
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
				rets := winReg.FindStringSubmatch(desc)
				if len(rets) == 3 {
					winstr := rets[2]
					ut.FloatStrPtrMul(&winstr, bet)
					ut.FloatStrPtrDiv(&winstr, baseBet)
					ret := rets[1] + winstr
					sm.Desc[i] = ret
				}
			}
		}

		for _, pinfos := range doc.LogPlateInfos {
			for _, pinfo := range pinfos {
				//ut.FloatStrPtrMul(&pinfo.Win, bet)
				//ut.FloatStrPtrDiv(&pinfo.Win, baseBet)
				for _, v := range pinfo.List {
					ut.FloatStrPtrMul(&v.W, bet)
					ut.FloatStrPtrDiv(&v.W, baseBet)
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
