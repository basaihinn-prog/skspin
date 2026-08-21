package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"serve/comm/redisx"
	"strconv"
	"strings"

	"serve/comm/slotspool"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lang"
	"serve/comm/lazy"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_39/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	mux.RegRpc("/game-api/39/v2/Spin", "spin", "game-api", db.WrapRpcPlayer(spin), nil)
	mux.RegRpc("/game-api/39/v2/spin", "spin", "game-api", db.WrapRpcPlayer(spin), nil)
}

// https://api.pg-demo.com/game-api/piggy-gold/v2/Spin?traceId=IEUAPU12
func spin(plr *models.Player, ps define.PGParams, ret *json.RawMessage) (err error) {
	pid := plr.PID
	// cs=10&ml=7&pf=1&id=1768821035597177344&wk=0_C&btt=1&atk=E4A5F4F4-3756-4268-BBDA-1363F9E8BEA6

	// jsonstr := `{"si":{"wp":{"1":[0,1,2]},"lw":{"1":6.0},"frl":[4,4,6,5,4,3,1,4,2],"pc":4,"wm":null,"tnbwm":null,"gwt":2,"fb":null,"ctw":6.0,"pmt":null,"cwc":1,"fstc":null,"pcwc":1,"rwsp":{"1":10.0},"hashr":"0:4;4;4#R#4#001020#MV#0.6#MT#1#MG#6.0#","ml":2,"cs":0.3,"rl":[4,4,4],"sid":"1766020646670630400","psid":"1766020646670630400","st":1,"nst":1,"pf":1,"aw":6.00,"wid":0,"wt":"C","wk":"0_C","wbn":null,"wfg":null,"blb":99998.20,"blab":99997.60,"bl":100003.60,"tb":0.60,"tbb":0.60,"tw":6.00,"np":5.40,"ocr":null,"mr":null,"ge":[3,11]}}`
	// *ret = json.RawMessage(jsonstr)

	cs := ps.GetFloat("cs")
	ml := ps.GetFloat("ml")
	// id := ps.Get("id")
	// id := strconv.Itoa(int(time.Now().UnixNano()))
	sf := ut.NewSnowflake()
	id := strconv.Itoa(int(sf.NextID()))

	curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
	App, err := redisx.LoadAppIdCache(plr.AppID)
	if err != nil {
		return err
	}
	if !ut.FloatInArr(ut.FloatArrMul(App.Cs, curItem.Multi), cs) ||
		!ut.FloatInArr(Ml, ml) {
		//当触发这个错误的时候直接一刀切，把该用户的上一局历史删除
		plr.RewriteLastData()
		slog.Error("spin in cs", cs, App.Cs)
		return lang.Error(plr.Language, "下注额非法1")
	}

	bet := ut.Money2Gold(cs * ml)
	// balance, err := slotsmongo.ModifyGold(plr.PID, -bet, "下注")
	balance, err := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
		Pid:     plr.PID,
		Change:  -bet,
		Reason:  slotsmongo.ReasonBet,
		RoundID: id,
	})
	if err != nil {
		err = define.PGNotEnoughCashErr

		si := getInitSi(App.DefaultCs, App.DefaultBetLevel)
		gold, _ := slotsmongo.GetBalance(plr.PID)
		balance := ut.Gold2Money(gold)
		si["bl"] = balance
		si["blb"] = balance
		si["blab"] = balance
		si["sid"] = plr.LastID
		si["psid"] = plr.LastID
		si["cs"] = cs
		si["ml"] = ml

		sibuf, _ := ut.GetJsonRaw(M{"si": si})
		*ret = sibuf

		return
	}

	selfPoolGold, _ := slotsmongo.GetSelfSlotsPool(pid)
	/*
		var pan core.Pan
		pan.RandomFill(&core.SimulateParams{
			Weight: core.SimulateParamsWeight{
				Wild:         1,
				Pig:          1,
				Gold:         1,
				Cabbage:      1,
				Firecracker3: 1,
				Firecracker2: 1,
				Firecracker1: 1,
				Nothing:      1,
			},
			Multi: [4]int{50, 30, 20, 10},
		})
		pgmsg := pan.ToPGMsg(cs, ml)
	*/

	plr.SpinCount++

	// pgmsg, err := sampleFromDB(cs, ml)
	// playResp, , err := nextPlayResp(player, bet, selfPoolGold)
	playResp, hitBigAward, forcedKill, err := nextPlayResp(plr, bet, selfPoolGold, App)
	if err != nil {
		// slotsmongo.ModifyGold(pid, bet, "下注-退回, err:"+err.Error())
		slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     pid,
			Change:  bet,
			Comment: "下注-退回, err:" + err.Error(),
			Reason:  slotsmongo.ReasonRefund,
			RoundID: id,
		})
		return
	}
	if ps.Header != nil {
		ps.Header["stat_bet"] = bet
	}
	pgmsg := *playResp

	{
		bet := cs * ml
		pgmsg["cs"] = cs
		pgmsg["ml"] = ml
		pgmsg["tb"] = bet
		pgmsg["tbb"] = bet
		mkmul(pgmsg, "np", bet)
		mkmul(pgmsg, "ctw", bet)
		mkmul(pgmsg, "aw", bet)
		mkmul(pgmsg, "tw", bet)
		mk2mul(pgmsg, "lw", "1", bet)

		hashrpart := fmt.Sprintf("#MV#%v#MT#1#MG#%v#", bet, pgmsg["aw"].(float64))

		hashr := pgmsg["hashr"].(string)

		pos := strings.Index(hashr, "#MV#")
		hashr = hashr[:pos] + hashrpart
		pgmsg["hashr"] = hashr
	}
	pgmsg["blb"] = ut.Gold2Money(balance + bet)
	pgmsg["blab"] = ut.Gold2Money(balance)
	pgmsg["sid"] = id
	pgmsg["psid"] = id
	win := pgmsg["aw"].(float64)
	wingold := ut.Money2Gold(win)

	if win >= 0 {
		// balance, err = slotsmongo.ModifyGold(plr.PID, wingold, "赢分")
		balance, err = slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
			Pid:     plr.PID,
			Change:  wingold,
			Reason:  slotsmongo.ReasonWin,
			RoundID: id,
			IsEnd:   true,
		})
		if err != nil {
			return
		}
		// pgmsg["bl"] = ut.Gold2Money(balance)
	}
	pgmsg["bl"] = ut.Gold2Money(balance)

	poolchange := calcPoolPlus(plr, bet, plr.SpinCount, forcedKill, App)

	var bigReward int64
	if hitBigAward {
		bigReward = ut.Money2Gold(win)
		poolchange -= bigReward
	} /* else {
		bound := gendata.GBuckets.GetBound(ut.GetInt(pgmsg["bucketid"]))
		if bound.PoolCost > 0 {
			// slotsmongo.IncSelfSlotsPool(pid, -(Grades[ps.Grade] * int64(bound.PoolCost)))
			poolchange -= (bet * int64(bound.PoolCost))
		}
	} */

	slotsmongo.IncSelfSlotsPool(pid, poolchange)

	msgbuf, err := ut.GetJsonRaw(pgmsg)
	if err != nil {
		return
	}
	detailsId := primitive.ObjectID{}
	if ps.Header["jump_log"] == nil {
		detailsId = insertBetHistory(plr, pgmsg, msgbuf)
	}

	plr.BetHistLastID = detailsId
	sibuf, _ := ut.GetJsonRaw(M{"si": msgbuf})
	*ret = sibuf
	plr.LS = string(sibuf)
	plr.LastID = id
	if ps.Header != nil {
		ps.Header["stat_win"] = wingold
		ps.Header["LastSid"] = plr.LastID
	}
	addRecord(plr, bet, wingold, balance /*, pgmsg["_id"].(primitive.ObjectID).Hex()*/, false, false, int(cs*1000+ml), detailsId, bigReward, id)

	return
}

func sampleFromDB() (spindata bson.M, err error) {
	coll := db.Collection("pgSpinData")
	cursor, err := coll.Aggregate(context.TODO(), bson.A{
		bson.D{{"$sample", bson.D{{"size", 1}}}},
	})
	// id, _ := primitive.ObjectIDFromHex("6604d1cce04f4526f342ac8a")
	// cursor, _ := coll.Find(context.TODO(), db.ID(id))
	if err != nil {
		return
	}

	// cursor.Next()
	var docs []bson.M
	cursor.All(context.TODO(), &docs)
	spindata = docs[0]

	return
}

func mkmul(m bson.M, k string, mul float64) {
	if v := m[k]; v != nil {
		fv := 0.0
		switch v.(type) {
		case int32:
			fv = float64(v.(int32))
		case int64:
			fv = float64(v.(int64))
		case float64:
			fv = v.(float64)
		default:
			panic("unsporrted type")
		}
		m[k] = fv * mul
	}
}

func mk2mul(m bson.M, k1, k2 string, mul float64) {
	if m1 := m[k1]; m1 != nil {
		mkmul(m1.(bson.M), k2, mul)
	}
}

func mk3mul(m bson.M, k1, k2, k3 string, mul float64) {
	if m1 := m[k1]; m1 != nil {
		mk2mul(m1.(bson.M), k2, k3, mul)
	}
}

func calcPoolPlus(p *models.Player, bet int64, spinCount int, forcedKill bool, app *redisx.AppStore) (toSelfAwardPool int64) {
	if forcedKill {
		toSelfAwardPool = bet
	} else {
		rotateCount := p.SpinCount
		//走平台的
		var betMul int
		betMul = redisx.LoadAwardPercent(p.AppID) + slotspool.GetSlotsPool(p.AppID).Value
		if p.RewardPercent != 0 {
			//走用户的
			betMul = p.RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
		}

		if app.IsProtection == 1 && rotateCount < app.ProtectionRotateCount {
			betMul += app.ProtectionRewardPercentLess
		}
		toSelfAwardPool = bet * int64(betMul) / 1000
	}
	return
}
