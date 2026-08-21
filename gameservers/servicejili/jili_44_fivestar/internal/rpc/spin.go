package rpc

import (
	"context"
	"fmt"
	"regexp"
	"serve/comm/redisx"
	"serve/servicejili/jili_44_fivestar/internal/message"
	"strconv"
	"time"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lang"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_44_fivestar/internal"
	"serve/servicejili/jili_44_fivestar/internal/models"
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
var winReg2 = regexp.MustCompile(`^(\w+ win: )([0-9.+-]+)$`)
var myVar = int64(1000) //该游戏特殊处理 用来对一些数据进行处理

func spin(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	var spinReq message.Server_SpinReq
	err = proto.Unmarshal(data, &spinReq)
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {

		bet := float64(spinReq.GetBet()) / float64(myVar)
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)

		App, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}

		//bet := float64(betss)
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
		ty := internal.GameTypeNormal
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

		var forcedKill, hitBigAward bool
		var doc *models.RawSpin
		selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(plr.PID)
		doc, hitBigAward, forcedKill, err = nextPlayResp(plr, gold, selfPoolGold, ty, App)
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
		fivestarData := doc.FivestarData
		for _, roundInfo := range fivestarData.GetPlate() {
			for _, award2 := range roundInfo.Award {
				*award2.Win = int64(float64(*award2.Win) * bet)
			}
			if roundInfo.Win != nil {
				*roundInfo.Win = int64(float64(*roundInfo.Win) * bet)
			}
			if roundInfo.DiamondBonus != nil {
				*roundInfo.DiamondBonus = int64(float64(*roundInfo.DiamondBonus) * bet)
			}
		}
		// 通用数据
		jiliUtData := doc.UtData
		if jiliUtData.FullJpWin != nil {
			*jiliUtData.FullJpWin = int64(float64(*jiliUtData.FullJpWin) * bet)

		}
		//lo.Must0(ut.FloatEQ(spinAck.GetTotalWin(), totalWin))
		totalWin := float64(jiliUtData.GetTotalWin() / myVar)

		//winMaxLimit := slotspool.GetMaxWin(plr.AppID, int64(float64(plr.GetWinLose())/10000), curItem.Multi)
		//if winMaxLimit != nil && float64(totalWin) > float64(winMaxLimit.MaxWin)*curItem.Multi {
		//	doc, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
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
		//	fivestarData = doc.FivestarData
		//	for _, roundInfo := range fivestarData.GetPlate() {
		//		for _, award2 := range roundInfo.Award {
		//			*award2.Win = int64(float64(*award2.Win) * bet)
		//		}
		//		if roundInfo.Win != nil {
		//			*roundInfo.Win = int64(float64(*roundInfo.Win) * bet)
		//		}
		//		if roundInfo.DiamondBonus != nil {
		//			*roundInfo.DiamondBonus = int64(float64(*roundInfo.DiamondBonus) * bet)
		//		}
		//	}
		//	// 通用数据
		//	jiliUtData = doc.UtData
		//	if jiliUtData.FullJpWin != nil {
		//		*jiliUtData.FullJpWin = int64(float64(*jiliUtData.FullJpWin) * bet)
		//	}
		//	realWin := jiliUtData.GetTotalWin() / myVar
		//	// 大奖处理
		//	slotsmongo.BigRewardDeal(&slotsmongo.BigRewardDealParams{
		//		Pid:         plr.PID,
		//		AppID:       plr.AppID,
		//		CurrencyKey: curItem.Key,
		//		ServiceName: lazy.ServiceName,
		//		OriginWin:   float64(totalWin),
		//		RealWin:     float64(realWin),
		//		WinLose:     plr.GetWinLose(),
		//		Desc:        winMaxLimit.Desc,
		//	})
		//	totalWin = realWin
		//	hitBigAward = false
		//	forcedKill = false
		//}
		totalWin *= bet
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
		//fivestarData.NowMoney = proto.Float64(ut.Gold2Money(balanceGold))
		jiliUtData.PostMoney = proto.Int64(int64(ut.Gold2Money(balanceGold * myVar)))
		jiliUtData.PreMoney = proto.Int64(int64((balance + costMoney) * float64(myVar)))
		jiliUtData.AllPlate, _ = proto.Marshal(fivestarData)

		var toSelfAwardPool, bigReward int64

		if hitBigAward {
			toSelfAwardPool = -ut.Money2Gold(totalWin)
			bigReward = ut.Money2Gold(totalWin)
		}
		//else {
		//			poolCost := gendata.GBuckets.GetPoolCost(doc.BucketId)
		//			if poolCost > 0 {
		//				toSelfAwardPool = -gold * int64(poolCost)
		//			}
		//		}

		plr.UpdatePool(costGold, selfPoolGold, toSelfAwardPool, forcedKill, App)
		plr.OnSpinFinish(costGold, ut.Money2Gold(totalWin*float64(myVar)), false, false, 1)
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

		jiliUtData.RoundIndex = proto.Int64(roundIndexV2)
		ret = jiliUtData

		doc.HistoryRecord.RoundIndex = showIndex
		his := doc.HistoryRecord
		his.Bet = ut.Ftoa(costMoney)
		his.Win = ut.Ftoa(float64(totalWin))
		his.NetValue = ut.Ftoa(float64(totalWin) - costMoney)
		his.PreMoney = ut.Ftoa(balance + costMoney)
		his.PostMoney = ut.Ftoa(float64(jiliUtData.GetPostMoney()))
		/*	his.Bet = ut.Ftoa(costMoney * float64(myVar))
			his.Win = ut.Ftoa(float64(totalWin * myVar))
			his.NetValue = ut.Ftoa((float64(totalWin) - costMoney) * float64(myVar))
			his.PreMoney = ut.Ftoa((balance + costMoney) * float64(myVar))*/

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
				} else {
					rets := winReg2.FindStringSubmatch(desc)
					if len(rets) == 3 {
						winstr := rets[2]
						ut.FloatStrPtrMul(&winstr, bet)
						ret := rets[1] + winstr
						sm.Desc[i] = ret
					}
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
