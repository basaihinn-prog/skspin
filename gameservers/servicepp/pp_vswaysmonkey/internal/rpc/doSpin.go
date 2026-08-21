package rpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"serve/comm/redisx"
	"strconv"
	"strings"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/lang"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicepp/pp_vswaysmonkey/internal"
	"serve/servicepp/pp_vswaysmonkey/internal/gamedata"
	"serve/servicepp/pp_vswaysmonkey/internal/gendata"
	"serve/servicepp/ppcomm"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	ppcomm.RegRpc("doSpin", doSpin)
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
		c := ps.Float("c")
		isDouble := ps.Int("bl") == 1
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
		if betTotal < plrCs[0]*internal.Line*curItem.Multi ||
			betTotal > plrCs[len(plrCs)-1]*internal.Line*curItem.Multi {
			//当触发这个错误的时候直接一刀切，把该用户的上一局历史删除
			plr.RewriteLastData()
			return lang.Error(plr.Language, "下注额非法")
		}

		//if betTotal < internal.BetMin*curItem.Multi ||
		//	betTotal > internal.BetMax*curItem.Multi {
		//	return lang.Error(plr.Language, "下注额非法")
		//}

		//isEnd, params := plr.IsEndO()

		var doc *ppcomm.SimulateData
		bet := int64(0)
		balance := int64(0)
		costGold := int64(0)

		frbdata := plr.FRBData
		if frbdata != nil {
			if !ut.FloatEQ(betTotal, frbdata.Config.TotalBet) {
				return errors.New("Invalid betTotal")
			}

			if isBuy || isDouble {
				return errors.New("Invalid spin params")
			}
		}
		sf := ut.NewSnowflake()
		formatInt := strconv.FormatInt(sf.NextID(), 10)
		betId := fmt.Sprintf(`70%v000`, formatInt[len(formatInt)-9:])

		if isEnd {
			// 需要读取新的数据
			plr.IsBuy = isBuy
			plr.IsCollect = false
			plr.C = c
			plr.BetHistLastID = primitive.NewObjectID()
			plr.LastBetId = betId
			plr.BigReward = 0
			mul := float64(1)
			if plr.IsBuy {
				mul *= internal.BuyMul
			}
			if isDouble {
				mul *= internal.Double
			}
			bet = ut.Money2Gold(c * internal.Line)
			if bet <= 0 {
				return lang.Error(plr.Language, "下注额非法")
			}
			costGold = int64(float64(bet) * mul)
			plr.CostGold = costGold

			if frbdata != nil {
				frbdata.Frn--
				lo.Must0(frbdata.Frn >= 0)

				balance, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
					Pid:     plr.PID,
					Change:  0,
					RoundID: plr.LastBetId,
					Reason:  slotsmongo.ReasonFRBBet,
					Comment: fmt.Sprintf("costGold:%d,frn:%d, fra:%.2f,id:%s,bonuscode:%s", costGold, frbdata.Frn, frbdata.Fra, frbdata.Config.ID.Hex(), frbdata.Config.BonusCode),
				})
			} else {
				balance, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
					Pid:     plr.PID,
					Change:  -costGold,
					RoundID: plr.LastBetId,
					Reason:  slotsmongo.ReasonBet,
				})
			}
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
				NextFunc:               gendata.GCombineDataMng.Next,
				SimulateByBucketIdFunc: gendata.GCombineDataMng.SampleSimulate,
				BigRewardFunc:          gendata.GCombineDataMng.GetBigReward,
				MinBuyFunc:             gendata.GCombineDataMng.GetBuyMinData,
				NextBuyFunc:            gendata.GCombineDataMng.NextBuy,
				ControlNextData:        gendata.GCombineDataMng.ControlNextDataNormal, //控制游戏的下一次func
				BigRewardFunc2_5:       gendata.GCombineDataMng.GetBigReward2_5,
			}, internal.RandPlayResp)
			if err != nil {
				//slotsmongo.ModifyGold(plr.PID, bet*mul, "下注-退回, err:"+err.Error())

				if frbdata == nil {
					slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
						Pid:     plr.PID,
						Change:  costGold,
						RoundID: plr.LastBetId,
						Reason:  slotsmongo.ReasonRefund,
					})
				} else {
					frbdata.Frn++
				}
				return err
			}
			lastPan := doc.DropPan[len(doc.DropPan)-1]
			allWin := ut.GetFloat(lastPan.Currency("tw")) * plr.C * internal.Line
			//新逻辑，与现有限制最大应分冲突且会影响rtp计算，预备弃用
			//winMaxLimit := slotspool.GetMaxWin(plr.AppID, int64(float64(plr.GetWinLose())/10000), curItem.Multi)
			//if winMaxLimit != nil && allWin > float64(winMaxLimit.MaxWin)*curItem.Multi {
			//	if isBuy {
			//		doc, err = gendata.GCombineDataMng.GetBuyMinData()
			//	} else {
			//		doc, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
			//	}
			//	if err != nil {
			//		if frbdata == nil {
			//			slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			//				Pid:     plr.PID,
			//				Change:  costGold,
			//				RoundID: plr.BetHistLastID.Hex(),
			//				Reason:  slotsmongo.ReasonRefund,
			//			})
			//		} else {
			//			frbdata.Frn++
			//		}
			//		return err
			//	}
			//
			//	lastPan = doc.DropPan[len(doc.DropPan)-1]
			//	realWin := ut.GetFloat(lastPan.Currency("tw")) * plr.C * internal.Line
			//	// 大奖处理
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

			plr.OnSpinFinish(costGold, ut.Money2Gold(allWin), plr.IsBuy, false, internal.BuyMul)

		} else {
			objectId := params[0]
			num, _ := strconv.Atoi(params[1])
			findOptions := &options.FindOneOptions{}
			findOptions.SetProjection(db.D("droppan", bson.M{"$slice": []int{num, 1}}))
			objId, _ := primitive.ObjectIDFromHex(objectId)
			//simulate
			err = db.Collection("simulate").FindOne(context.TODO(), bson.M{"_id": objId}, findOptions).Decode(&doc)
			if err != nil {
				return err
			}
			balance, err = slotsmongo.GetBalance(plr.PID)
			if err != nil {
				return err
			}
		}

		//sampleDoSpin := ppcomm.ParseVariables(s)
		data := doc.Deal2(plr.C, internal.Line, ut.Gold2Money(balance))

		data.SetStr("index", ps.Str("index"))
		data.SetInt("counter", ps.Int("counter")+1)

		plr.LastData = maps.Clone(data)
		gid := data.Str("gid")
		gids := strings.Split(gid, "_")
		//sf := ut.NewSnowflake()
		//betId := strconv.Itoa(int(sf.NextID()))
		isCompleted := false
		if gids[1] == gids[2] && plr.LastData.Str("na") != "c" { // 表示是最后一盘 "c" 还回执行collect操作
			isCompleted = true

			lo.Must0(plr.LastData.Str("na") == "s")
		}
		if msg.Header != nil && msg.Header.Get("isstat") == "1" && data.Get("na") == "c" {
			msg.Header.Set("stat_win", fmt.Sprintf("%d", ut.Money2Gold(data.Currency("tw"))))
			msg.Header.Set("lastId", gid)
		}
		ppcomm.FRBOnSpin(plr, data, isCompleted)

		delete(data, "gid")
		if msg.Header != nil && msg.Header.Get("jump_log") != "1" {
			ppcomm.InsertBetHistoryEvery(plr, ut.Gold2Money(costGold), isEnd, frbdata != nil, ps, data, curItem.Key, curItem.Symbol, betTotal, doc.Id, sampleDoInit.Encode(), betId)
		}

		bonusCode := ""
		if frbdata != nil {
			bonusCode = frbdata.Config.BonusCode
		}
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
			Frb:          frbdata != nil,
			BonusCode:    bonusCode,
		}
		if msg.Header != nil && msg.Header.Get("jump_log") != "1" {
			slotsmongo.AddBetLog(betLog)
		}
		ret = data.Bytes()
		return nil
	})
	return
}
