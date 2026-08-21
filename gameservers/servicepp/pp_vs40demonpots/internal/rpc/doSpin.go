package rpc

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/lang"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicepp/pp_vs40demonpots/internal"
	"serve/servicepp/pp_vs40demonpots/internal/gamedata"
	"serve/servicepp/pp_vs40demonpots/internal/gendata"
	"serve/servicepp/ppcomm"
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	ppcomm.RegRpc("doSpin", doSpin)
	ppcomm.RegRpc("doBonus", doSpin)
}

/*
acci: 0
accm: cp
accv: 0
balance: 999,990.00
balance_bonus: 0.00
balance_cash: 999,990.00
bl: 0
c: 0.50
counter: 4
index: 2
l: 20
na: s
ntp: -10.00
reel_set: 1
rid: 38693909614091
s: 10,7,5,8,3,3,9,8,10,8,3,3,10,6,5,8,11,9,8,11,7,10,11,9,8,9,7,10,8,5
sa: 6,6,3,5,9,4
sb: 3,7,8,11,10,5
sh: 5
stime: 1725258160393
sver: 5
tw: 0.00
w: 0.00
*/
func doSpin(msg *nats.Msg) (ret []byte, err error) {
	// action=doSpin&symbol=vs20olympx&c=0.5&l=20&bl=0&index=2&counter=3&repeat=0&mgckey=AUTHTOKEN@6318b8f52aa27a735236b247dafd0e27037578637940d87ad54250bfeed0e431~stylename@hllgd_hollygod~SESSION@7dbfc452-ef7a-4b14-8069-a857e5dc165b~SN@e9d7b550

	ps := ppcomm.ParseVariables(string(msg.Data))
	var pid int64
	pid, err = jwtutil.ParseToken(ps.Str("mgckey"))
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *ppcomm.Player) error {
		//// 是否黑名单
		//if slotsmongo.IsBanned(lazy.ServiceName, plr.PID) {
		//	return define.PGInvalidSession
		//}

		// 参数处理
		c := ps.Float("c")            //下注金额
		action := ps.Str("action")    //action
		isDouble := ps.Int("bl") == 1 //开启双倍

		_, isBuy := ps["pur"]

		if isDouble {
			isBuy = false
		}

		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)

		App, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		isEnd, params := plr.IsEndO()
		if isEnd == true {
			plr.IsBuy = isBuy
		}
		plrCs, err := redisx.GetPlayerCs(plr.AppID, plr.PID, isEnd)
		if err != nil {
			slog.Error("doSpin", "GetPlayerCs", err)
			return err
		}
		betTotal := c * internal.Line
		if (plr.IsBuy && plr.C != 0) || (action == "doBonus" && c == 0) {
			betTotal = plr.C * internal.Line
		}
		if betTotal < plrCs[0]*internal.Line*curItem.Multi ||
			betTotal > plrCs[len(plrCs)-1]*internal.Line*curItem.Multi {
			//当触发这个错误的时候直接一刀切，把该用户的上一局历史删除
			plr.RewriteLastData()
			return lang.Error(plr.Language, "下注额非法")
		}

		//isEnd := plr.IsEnd    //上一轮是否是最后一轮
		//lastId := plr.LastSid //上一轮id
		//if lastId == "" {
		//	isEnd = true
		//}
		//LastIndex := plr.LastIndex //上一轮id
		//lastRound := true          //本轮是否是最后一轮
		//isEnd, params := plr.IsEndO()

		//multiply := (int)(plr.C*internal.Line) / internal.MinBet // 倍数，所有与钱相关的都在入库的时候乘以了100转换为整数
		//multiply := (int)(plr.C*internal.Line) / internal.MinBet // 倍数，所有与钱相关的都在入库的时候乘以了100转换为整数

		var doc *ppcomm.SimulateData
		bet := int64(0)
		balance := int64(0)
		costGold := int64(0)
		sf := ut.NewSnowflake()
		formatInt := strconv.FormatInt(sf.NextID(), 10)
		betId := fmt.Sprintf(`70%v000`, formatInt[len(formatInt)-9:])

		if isEnd {
			//LastIndex = -1
			// 需要读取新的数据
			plr.IsBuy = isBuy
			plr.IsCollect = false
			plr.C = c //下注金额
			plr.BetHistLastID = primitive.NewObjectID()
			plr.LastBetId = betId
			plr.BigReward = 0
			mul := float64(1)

			if plr.IsBuy && internal.BuyBetMulti != 0 {
				mul = float64(internal.GetFreeMultiply()[ps.Int("pur")])
			}

			if isDouble { //是否是双倍下注
				mul *= internal.Double
			}

			//多倍购买控制
			if !internal.DoubleBuy && plr.IsBuy && isDouble {
				return lang.Error(plr.Language, "下注额非法")
			}

			bet = ut.Money2Gold(c * internal.Line) //Line是固定倍数，因为页面传过来的值是  1/20  部分游戏可能是1/5    1/25
			if bet <= 0 {
				return lang.Error(plr.Language, "下注额非法")
			}

			costGold = int64(float64(bet) * mul) //本次实际下注金额
			plr.CostGold = costGold

			//调用游戏中心去修改玩家钱包信息，这里是扣除下注金额
			balance, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
				Pid:     plr.PID,
				Change:  -costGold,
				RoundID: plr.LastBetId,
				Reason:  slotsmongo.ReasonBet,
			})

			if err != nil {
				err = define.PGNotEnoughCashErr
				return err
			}

			var forcedKill, hitBigAward, buyKill bool
			selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(plr.PID)

			doc, hitBigAward, forcedKill, buyKill, err = ppcomm.NextPlayResp(&ppcomm.NextPlayRespParam{
				Player:                 plr,
				Bet:                    bet,
				AppStore:               App,
				SelfPoolGold:           selfPoolGold,
				IsBuy:                  plr.IsBuy,
				BuyMinAwardPercent:     App.BuyMinAwardPercent,
				Combine:                ppcomm.NewCombine(gendata.GetCombine()),
				NoAwardPercent:         gamedata.Settings.NoAwardPercent,
				HitBigAwardPercent:     gamedata.Settings.HitBigAwardPercent,
				NextFunc:               gendata.GCombineDataMng.Next,                  //正常轮次下一轮获取
				SimulateByBucketIdFunc: gendata.GCombineDataMng.SampleSimulate,        //根据既定bucketId获取
				BigRewardFunc:          gendata.GCombineDataMng.GetBigReward,          //大奖，从整体奖池里找即可
				BigRewardFunc2_5:       gendata.GCombineDataMng.GetBigReward2_5,       //大奖，从整体奖池里找即可
				MinBuyFunc:             gendata.GCombineDataMng.GetBuyMinData,         //杀！
				NextBuyFunc:            gendata.GCombineDataMng.NextBuy,               //购买游戏轮的次获取，包括首轮和后续轮次
				ControlNextData:        gendata.GCombineDataMng.ControlNextDataNormal, //控制游戏的下一次func
			}, internal.RandPlayResp)

			if err != nil {
				//回滚
				_, _ = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
					Pid:     plr.PID,
					Change:  costGold,
					RoundID: plr.LastBetId,
					Reason:  slotsmongo.ReasonRefund,
				})
				return err
			}

			// 应该是当获得下一把结果时，计算该结果是否满足特殊情况，进行特殊处理，例如杀/大奖
			//fmt.Println(doc.Id)
			// 自动多轮时，每一轮的tw是累计的，所以取最后一轮即可用于判断盈利
			//lastPan := doc.DropPan[len(doc.DropPan)-1]
			//allWin := ut.GetFloat(lastPan.Float("tw")) * plr.C / internal.Line
			// 金额数
			allWin := ut.GetFloat(doc.Times) * plr.C * internal.Line
			//fmt.Println("本次多轮allWin:", allWin)

			//winMaxLimit := slotspool.GetMaxWin(plr.AppID, int64(float64(plr.GetWinLose())/10000), curItem.Multi)
			//if winMaxLimit != nil && allWin > float64(winMaxLimit.MaxWin)*curItem.Multi {
			//	if isBuy {
			//		doc, err = gendata.GCombineDataMng.GetBuyMinData()
			//	} else {
			//		doc, err = gendata.GCombineDataMng.SampleSimulate(0)
			//	}
			//	if err != nil {
			//		//回滚下注信息
			//		_, _ = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			//			Pid:     plr.PID,
			//			Change:  costGold,
			//			RoundID: plr.BetHistLastID.Hex(),
			//			Reason:  slotsmongo.ReasonRefund,
			//		})
			//		return err
			//	}
			//
			//	lastPan = doc.DropPan[len(doc.DropPan)-1]
			//	//realWin := ut.GetFloat(lastPan.Float("tw")) * plr.C * internal.Line
			//	realWin := ut.GetFloat(lastPan.Float("tw")) * float64(multiply)
			//	// 大奖处理，报警
			//	slotsmongo.BigRewardDeal(&slotsmongo.BigRewardDealParams{
			//		Pid:         plr.PID,
			//		AppID:       plr.AppID,
			//		CurrencyKey: curItem.Key,
			//		ServiceName: lazy.ServiceName,
			//		OriginWin:   allWin,
			//		RealWin:     realWin,
			//		WinLose:     plr.GetWinLose(),
			//		Desc:        winMaxLimit.Desc,
			//	})
			//	allWin = realWin
			//	hitBigAward = false
			//	forcedKill = false
			//	buyKill = false
			//}

			if msg.Header != nil && msg.Header.Get("isstat") == "1" {
				msg.Header.Set("stat_bet", fmt.Sprintf("%d", costGold))
			}

			var toSelfAwardPool int64

			if isBuy {
				if buyKill {
					toSelfAwardPool += int64(float64(bet) * mul)
				}
			} else {
				if hitBigAward {
					toSelfAwardPool += -ut.Money2Gold(allWin)
					plr.BigReward = ut.Money2Gold(allWin)
				}
				//else {
				//	poolCost := gendata.GBuckets.GetPoolCost(doc.BucketId)
				//	if poolCost > 0 {
				//		toSelfAwardPool += -bet * int64(poolCost)
				//	}
				//}
				if isDouble {
					toSelfAwardPool += int64(float32(bet) * 0.25)
				}
			}

			if !isBuy || buyKill {
				plr.UpdatePool(&ppcomm.UpdatePoolParam{
					Bet:                  bet,
					App:                  App,
					SelfPoolGold:         selfPoolGold,
					ToSelfAwardPool:      toSelfAwardPool,
					ForcedKill:           forcedKill,
					IsBuy:                plr.IsBuy,
					RewardPercentLess100: 0,
				})
			}

			plr.OnSpinFinish(costGold, ut.Money2Gold(allWin), plr.IsBuy, false, internal.GetFreeMultiply()[ps.Int("pur")])

			//本轮是首轮，且包含多条子结果
			//lastRound = len(doc.DropPan) == 1

			//复制QueryString,从DropPan的首个数组中获取，统一供后续使用
			doc.QueryString = doc.DropPan[0]

			//fmt.Println(doc.QueryString)
		} else {
			//objId, _ := primitive.ObjectIDFromHex(lastId)
			////取的QueryString，写入doc
			//err = db.Collection("simulate").FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&doc)
			//if err != nil {
			//	return err
			//}
			//balance, err = slotsmongo.GetBalance(plr.PID)
			//if err != nil {
			//	return err
			//}
			//lastRound = (len(doc.DropPan) - 1) == (LastIndex + 1)
			//
			//doc.QueryString = doc.DropPan[LastIndex+1]
			objectId := params[0]
			num, _ := strconv.Atoi(params[1])
			findOptions := &options.FindOneOptions{}
			findOptions.SetProjection(db.D("droppan", bson.M{"$slice": []int{num, 1}}))
			objId, _ := primitive.ObjectIDFromHex(objectId)
			err = db.Collection("simulate").FindOne(context.TODO(), bson.M{"_id": objId}, findOptions).Decode(&doc)
			if err != nil {
				return err
			}
			balance, err = slotsmongo.GetBalance(plr.PID)
			if err != nil {
				return err
			}
		}

		//初始化本次plr数值
		//plr.LastSid = doc.Id.Hex()
		//plr.LastIndex = LastIndex + 1
		//plr.IsEnd = lastRound

		//sampleDoSpin := ppcomm.ParseVariables(s)
		data := doc.Deal2(plr.C, internal.Line, ut.Gold2Money(balance))

		data.SetStr("index", ps.Str("index"))
		data.SetInt("counter", ps.Int("counter")+1)
		//data.SetStr("c", fmt.Sprintf("%.2f", c))
		//
		//data.SetInt("1", internal.Line)

		plr.LastData = maps.Clone(data)

		isCompleted := false
		gid := data.Str("gid")
		gids := strings.Split(gid, "_")
		if gids[1] == gids[2] && plr.LastData.Str("na") != "c" { // 表示是最后一盘 "c" 还回执行collect操作
			isCompleted = true
		}

		if msg.Header != nil && msg.Header.Get("isstat") == "1" && data.Get("na") == "c" {
			msg.Header.Set("stat_win", fmt.Sprintf("%d", ut.Money2Gold(data.Currency("tw"))))
		}
		if msg.Header != nil && msg.Header.Get("jump_log") != "1" {
			ppcomm.InsertBetHistoryEvery(plr, ut.Gold2Money(costGold), isEnd, false, ps, data, curItem.Key, curItem.Symbol, betTotal, doc.Id, sampleDoInit.Encode(), betId)
		}
		//
		//obet := ut.Money2Gold(plr.C * internal.Line)
		//if plr.IsBuy {
		//	bet *= internal.BuyMul
		//	obet *= internal.BuyMul
		//}
		//if plr.IsDouble {
		//	bet = int64(float64(bet) * internal.Double)
		//	obet = int64(float64(bet) * internal.Double)
		//}
		//
		//win :=
		//aw := ut.Money2Gold(data.Float("tw"))
		if isCompleted {
			slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
				Pid:     plr.PID,
				Change:  0,
				RoundID: plr.LastBetId,
				Reason:  slotsmongo.ReasonWin,
				IsEnd:   isCompleted,
			})
		}
		betLog := &slotsmongo.AddBetLogParams{
			UserName:     "",
			CurrencyKey:  plr.CurrencyKey,
			ID:           betId,
			Pid:          plr.PID,
			Bet:          costGold,
			Win:          0,
			Balance:      balance,
			RoundID:      plr.LastBetId,
			Completed:    isCompleted,
			TotalWinLoss: 0,
			IsBuy:        plr.IsBuy,
			Grade:        int(plr.C * internal.Line),
			PGBetID:      plr.BetHistLastID,
			BigReward:    plr.BigReward,
			GameType:     define.GameType_Slot,
		}
		if msg.Header != nil && msg.Header.Get("jump_log") != "1" {
			slotsmongo.AddBetLog(betLog)
		}
		ret = data.Bytes()
		return nil
	})
	return
}
