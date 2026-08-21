package rpc

import (
	"context"
	"encoding/json"
	"log/slog"
	"serve/comm/redisx"
	"serve/service/pgcomm"

	"serve/comm"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lang"
	"serve/comm/lazy"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_120/internal/gendata"
	"serve/service/pg_120/internal/models"

	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.RegRpc("/game-api/120/v2/Spin", "spin", "game-api", db.WrapRpcPlayer(spin), nil)
	mux.RegRpc("/game-api/120/v2/spin", "spin", "game-api", db.WrapRpcPlayer(spin), nil)
}

// https://api.pg-demo.com/game-api/piggy-gold/v2/Spin?traceId=IEUAPU12
func spin(plr *models.Player, ps define.PGParams, ret *json.RawMessage) (err error) {
	cs := ps.GetFloat("cs")
	ml := ps.GetFloat("ml")
	id := ps.Get("id")
	isBuy := ps.Get("fb") == "2"

	curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
	App, err := redisx.LoadAppIdCache(plr.AppID)
	if err != nil {
		return err
	}
	isEnd2, _ := plr.IsEndO()
	c, err := redisx.GetPlayerCs(plr.AppID, plr.PID, isEnd2)
	if !ut.FloatInArr(ut.FloatArrMul(c, curItem.Multi), cs) ||
		!ut.FloatInArr(Ml, ml) {
		//当触发这个错误的时候直接一刀切，把该用户的上一局历史删除
		plr.RewriteLastData()
		slog.Error("spin in cs", cs, c)
		return lang.Error(plr.Language, "下注额非法1")
	}

	params := strings.Split(id, "_")
	isEnd := false
	if id == "0" {
		isEnd = true
	} else {
		if len(params) != 3 || id != plr.LastSid {
			return lang.Error(plr.Language, "参数错误")
		}
		num, _ := strconv.Atoi(params[1])
		all, _ := strconv.Atoi(params[2])
		if num >= all {
			isEnd = true
		}
	}
	var doc *gendata.SimulateData
	num := 0
	bet := int64(0)
	bigReward := int64(0)
	balance := int64(0)
	if isEnd {
		// 新的一轮，需要扣除
		plr.IsBuy = isBuy
		plr.BdRecords = plr.BdRecords[:0]
		plr.BetHistLastID = primitive.NewObjectID()
		mul := int64(1)
		if plr.IsBuy {
			mul *= gendata.BuyMul
		}
		bet = ut.Money2Gold(cs * ml * gendata.Line)
		if bet <= 0 {
			return lang.Error(plr.Language, "下注额非法")
		}
		//_, err = slotsmongo.ModifyGold(plr.PID, -bet*mul, "下注")
		balance, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  -bet * mul,
			RoundID: plr.BetHistLastID.Hex(),
			Reason:  slotsmongo.ReasonBet,
		})
		if err != nil {
			err = define.PGNotEnoughCashErr

			si := getInitSi(App.DefaultCs, App.DefaultBetLevel)
			gold, _ := slotsmongo.GetBalance(plr.PID)
			balance := ut.Gold2Money(gold)
			si["bl"] = balance
			si["blb"] = balance
			si["blab"] = balance
			si["cs"] = cs
			si["ml"] = ml

			sibuf, _ := ut.GetJsonRaw(map[string]any{"si": si})
			*ret = sibuf
			return err
		}
		var forcedKill, hitBigAward, buyKill bool
		selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(plr.PID)
		//doc, hitBigAward, forcedKill, buyKill, err = nextPlayResp(plr, bet, selfPoolGold, isBuy, App)
		doc, hitBigAward, forcedKill, err = pgcomm.PGNextPlayRespFunc(
			pgcomm.PGNextPlayRespStruct{
				Player:             (*pgcomm.PGPlayer)(plr),
				Bet:                bet,
				IsBuy:              isBuy,
				SelfPoolGold:       selfPoolGold,
				HitBigAwardPercent: gendata.HitBigAwardPercent,
				App:                App,
			},
			pgcomm.GetBigRewardStruct_1{
				Next:            gendata.GCombineDataMng.Next,
				NextBuy:         gendata.GCombineDataMng.NextBuy,
				ControlNextData: gendata.GCombineDataMng.ControlNextData,
				SampleSimulate:  gendata.GCombineDataMng.SampleSimulate,
				GetBigReward:    gendata.GCombineDataMng.GetBigReward,
				GetBigReward2_5: gendata.GCombineDataMng.GetBigReward2_5,
				GetBuyMinData:   gendata.GCombineDataMng.GetBuyMinData,
			},
		)
		if err != nil {
			//slotsmongo.ModifyGold(plr.PID, bet*mul, "下注-退回, err:"+err.Error())
			slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
				Pid:     plr.PID,
				Change:  bet * mul,
				RoundID: plr.BetHistLastID.Hex(),
				Reason:  slotsmongo.ReasonRefund,
			})
			return err
		}
		lastPan := doc.DropPan[len(doc.DropPan)-1]
		allWin := ut.GetFloat(lastPan["aw"]) * cs * ml * gendata.Line
		//bigReward = ut.Money2Gold(allWin)
		if ps.Header != nil {
			ps.Header["stat_bet"] = bet
		}
		var toSelfAwardPool int64
		if isBuy {
			if buyKill {
				toSelfAwardPool += bet * mul
			}
		} else {
			if hitBigAward {
				toSelfAwardPool += -ut.Money2Gold(allWin)
				bigReward = ut.Money2Gold(allWin)
			} /* else {
				poolCost := gendata.GBuckets.GetPoolCost(doc.BucketId)
				if poolCost > 0 {
					toSelfAwardPool += -bet * int64(poolCost)
				}
			} */
		}
		if !isBuy || buyKill {
			plr.UpdatePool(bet, selfPoolGold, toSelfAwardPool, forcedKill, isBuy, App)
		}

		plr.Ml = ml
		plr.Cs = cs

		plr.OnSpinFinish(bet, ut.Money2Gold(allWin), plr.IsBuy, false, gendata.BuyMul)
	} else {
		objectId := params[0]
		num, _ = strconv.Atoi(params[1])
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

	data := doc.Deal(num, ut.Gold2Money(balance), plr.Cs, plr.Ml, gendata.Line, gendata.BuyMul, plr.IsBuy)

	isCompleted := false
	plr.LastSid = data["sid"].(string)
	sid := strings.Split(plr.LastSid, "_")
	if sid[1] == sid[2] { // 表示是最后一盘
		isCompleted = true
	}

	win := ut.Money2Gold(ut.GetFloat(data["tw"]))
	//balance, _ = slotsmongo.ModifyGold(plr.PID, win, "赢分")
	balance, _ = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
		Pid:     plr.PID,
		Change:  win,
		RoundID: plr.BetHistLastID.Hex(),
		Reason:  slotsmongo.ReasonWin,
		IsEnd:   isCompleted,
	})

	if ps.Header != nil {
		ps.Header["stat_win"] = win
		ps.Header["LastSid"] = plr.LastSid
	}
	num, err = strconv.Atoi(sid[1])
	if err != nil {
		return
	}
	sf := ut.NewSnowflake()
	betId := strconv.Itoa(int(sf.NextID()))
	detailsId := primitive.ObjectID{}
	if ps.Header["jump_log"] == nil {
		detailsId = slotsmongo.InsertBetHistoryEvery(plr.BetHistLastID, plr.PID, 120, num, data, betId, plr.GetCurrencyOrTHB(), isEnd)
	}
	plr.BetHistLastID = detailsId
	obet := ut.Money2Gold(cs * ml * gendata.Line)
	if plr.IsBuy {
		bet *= gendata.BuyMul
		obet *= gendata.BuyMul
	}
	aw := ut.Money2Gold(data["aw"].(float64))
	betLog := &slotsmongo.AddBetLogParams{
		UserName:     "",
		CurrencyKey:  plr.CurrencyKey,
		ID:           betId,
		Pid:          plr.PID,
		Bet:          bet,
		Win:          win,
		Balance:      balance,
		RoundID:      plr.BetHistLastID.Hex(),
		Completed:    isCompleted,
		TotalWinLoss: aw - obet,
		IsBuy:        plr.IsBuy,
		Grade:        int(plr.Cs*1000 + plr.Ml),
		PGBetID:      plr.BetHistLastID,
		BigReward:    bigReward,
		GameType:     comm.GameType_Slot,
	}
	if ps.Header["jump_log"] == nil {
		slotsmongo.AddBetLog(betLog)
	}
	sibuf, _ := ut.GetJsonRaw(map[string]any{"si": data})
	*ret = sibuf
	plr.LS = string(sibuf)

	return nil
}
