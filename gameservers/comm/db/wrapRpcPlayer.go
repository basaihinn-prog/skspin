package db

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"serve/comm/lazy"
	"serve/comm/ut"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//	type RpcReqContext struct {
//		Pid      int64
//		Language string
//	}
type RpcReqContext = WSReqContext

// func play(pid int64, ps playPs, ret *playRet) (err error)
// func play(player , ps playPs, ret *playRet) (err error)

type IGetPid interface {
	GetPid() int64
}

func WrapRpcPlayer[PLR any, Params IGetPid, RET any](callee func(*PLR, Params, *RET) error) func(Params, *RET) error {
	return func(ps Params /*ps define.PGParams*/, ret *RET) error {
		// player := &models.Player{
		// 	PID: pid,
		// }
		//pgps := define.PGParams{}
		//fmt.Println(pgps)
		pid := /*ps.Pid*/ ps.GetPid()

		/*
			var player = new(PLR)

			coll := Collection("players")
			err := coll.FindOne(context.TODO(), ID(pid)).Decode(player)
			if err == mongo.ErrNoDocuments {
				err = nil
				doc, err := getDocPlayer(pid)
				if err != nil {
					err = fmt.Errorf("cannot found player[%v] data, err [%s]", pid, err.Error())
					return err
				}

				// doc.Language = ctx.Language
				ut.SetField(player, "DocPlayer", *doc)
			}

			if err != nil {
				return err
			}

			ut.InitNilFields(player)
		*/
		var player, err = GetPlr[PLR](pid)
		if err != nil {
			return err
		}

		// player.Language = ctx.Language
		// ut.SetField(player, "Language", ctx.Language)

		//装载 playerControl 相关参数
		GetPlayerRTPControl(player)
		//callback
		err = callee(player, ps, ret)
		if err != nil {
			return err
		}
		coll := Collection("players")
		_, err = coll.ReplaceOne(context.TODO(), ID(pid), player, options.Replace().SetUpsert(true))
		MqSendGC(player)
		return err
	}
}

func WrapRpcRtpPlayer[PLR any, Params IGetPid, RET any](callee func(*PLR, Params, *RET) error) func(Params, *RET) error {
	return func(ps Params /*ps define.PGParams*/, ret *RET) error {
		// player := &models.Player{
		// 	PID: pid,
		// }
		//pgps := define.PGParams{}
		//fmt.Println(pgps)
		//pid := /*ps.Pid*/ ps.GetPid()

		/*
			var player = new(PLR)

			coll := Collection("players")
			err := coll.FindOne(context.TODO(), ID(pid)).Decode(player)
			if err == mongo.ErrNoDocuments {
				err = nil
				doc, err := getDocPlayer(pid)
				if err != nil {
					err = fmt.Errorf("cannot found player[%v] data, err [%s]", pid, err.Error())
					return err
				}

				// doc.Language = ctx.Language
				ut.SetField(player, "DocPlayer", *doc)
			}

			if err != nil {
				return err
			}

			ut.InitNilFields(player)
		*/
		//var player, err = GetPlr[PLR](pid)
		//if err != nil {
		//	return err
		//}

		// player.Language = ctx.Language
		// ut.SetField(player, "Language", ctx.Language)
		err := callee(nil, ps, ret)
		if err != nil {
			return err
		}

		//coll := Collection("players")
		//_, err = coll.ReplaceOne(context.TODO(), ID(pid), player, options.Replace().SetUpsert(true))
		return err
	}
}

func WrapJiLIRpcRtpPlayer[PLR any, Params nats.Msg, RET any](callee func(*PLR, Params, *RET) error) func(Params, *RET) error {
	return func(ps Params /*ps define.PGParams*/, ret *RET) error {
		// player := &models.Player{
		// 	PID: pid,
		// }
		//pgps := define.PGParams{}
		//fmt.Println(pgps)
		//pid := /*ps.Pid*/ ps.GetPid()

		/*
			var player = new(PLR)

			coll := Collection("players")
			err := coll.FindOne(context.TODO(), ID(pid)).Decode(player)
			if err == mongo.ErrNoDocuments {
				err = nil
				doc, err := getDocPlayer(pid)
				if err != nil {
					err = fmt.Errorf("cannot found player[%v] data, err [%s]", pid, err.Error())
					return err
				}

				// doc.Language = ctx.Language
				ut.SetField(player, "DocPlayer", *doc)
			}

			if err != nil {
				return err
			}

			ut.InitNilFields(player)
		*/
		//var player, err = GetPlr[PLR](pid)
		//if err != nil {
		//	return err
		//}

		// player.Language = ctx.Language
		// ut.SetField(player, "Language", ctx.Language)
		err := callee(nil, ps, ret)
		if err != nil {
			return err
		}

		//coll := Collection("players")
		//_, err = coll.ReplaceOne(context.TODO(), ID(pid), player, options.Replace().SetUpsert(true))
		return err
	}
}

func GetPlr[PLR any](pid int64) (player *PLR, err error) {
	err = IsBanned(lazy.ServiceName, pid)
	if err != nil {
		return nil, err
	}

	player = new(PLR)

	coll := Collection("players")

	err = coll.FindOne(context.TODO(), ID(pid)).Decode(player)

	if err == mongo.ErrNoDocuments {
		var doc *DocPlayer
		doc, err = getDocPlayer(pid)
		if err != nil {
			err = fmt.Errorf("cannot found player[%v] data, err [%s]", pid, err.Error())
			return
		}

		// doc.Language = ctx.Language
		ut.SetField(player, "DocPlayer", *doc)
	}

	ut.InitNilFields(player)
	return
}

func CallWithPlayer[PLR any](pid int64, callee func(*PLR) error) error {
	var player, err = GetPlr[PLR](pid)
	if err != nil {
		slog.Error("CallWithPlayer::GetPlr", "err", err)
		return err
	}

	//装载 playerControl 相关参数
	GetPlayerRTPControl(player)

	err = callee(player)
	if err != nil {
		slog.Error("CallWithPlayer::callee", "err", err)
		return err
	}

	coll := Collection("players")
	_, err = coll.ReplaceOne(context.TODO(), ID(pid), player, options.Replace().SetUpsert(true))
	MqSendGC(player)
	return err
}

const (
	TargetRTP          = "TargetRTP"
	RelieveRTP         = "RelieveRTP"
	BetAmount          = "BetAmount"
	WinAmount          = "WinAmount"
	PID                = "PID"
	RewardPercent      = "RewardPercent"
	NoAwardPercent     = "NoAwardPercent"
	GameAdmin          = "GameAdmin"
	PlayerRTPControl   = "PlayerRTPControl"
	BuyMinAwardPercent = "BuyMinAwardPercent"
	PersonWinMaxMult   = "PersonWinMaxMult"
	PersonWinMaxScore  = "PersonWinMaxScore"
	IsControlRTP       = 1 //开启RTP 控制
	KeepService        = -1
)

func RelievePlayer(pid int64) {
	//回调修改
	coll := Collection2(GameAdmin, PlayerRTPControl)
	Matched, err := coll.UpdateOne(context.TODO(), bson.M{"GameID": lazy.ServiceName, "Pid": pid}, bson.M{"$set": bson.M{"Status": 0}})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	//没找到匹配的就接着控制 保底
	if Matched.ModifiedCount > 0 {
		//修改自己
		coll2 := Collection("players")
		_, err = coll2.UpdateOne(context.TODO(), ID(pid), bson.M{"$set": bson.M{"reward_percent": 0, "no_award_percent": 0}})
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

// 给plr 添加贯穿上下文值
func GetPlayerRTPControl(plr any) {
	v := reflect.ValueOf(plr).Elem()
	pid := v.FieldByName(PID).Int()

	coll2 := Collection2("game", "Players")
	var doc DocPlayer
	err := coll2.FindOne(context.TODO(), ID(pid), options.FindOne().SetProjection(plrprojection)).Decode(&doc)
	if err != nil {
		return
	}

	v.FieldByName("RestrictionsStatus").SetInt(doc.RestrictionsStatus)
	v.FieldByName("Win").SetInt(doc.Win)
	v.FieldByName("Multi").SetFloat(doc.Multi)
	v.FieldByName("RestrictionsMaxWin").SetInt(doc.RestrictionsMaxWin)
	v.FieldByName("RestrictionsMaxMulti").SetFloat(doc.RestrictionsMaxMulti)

	var prc PlayerRTPControlModel
	coll := Collection2(GameAdmin, PlayerRTPControl)
	err = coll.FindOne(context.TODO(), bson.M{"GameID": lazy.ServiceName, "Pid": pid}).Decode(&prc)
	if err == mongo.ErrNoDocuments {
		return
	}

	if prc.Status == IsControlRTP {
		targetRTPField := v.FieldByName(TargetRTP)
		relieveRTPField := v.FieldByName(RelieveRTP)
		rewardPercent := v.FieldByName(RewardPercent)
		noAwardPercent := v.FieldByName(NoAwardPercent)
		//IsRtpControl := true
		//// 检查字段是否存在
		//if !targetRTPField.IsValid() || !relieveRTPField.IsValid() || !rewardPercent.IsValid() || !noAwardPercent.IsValid() {
		//	IsRtpControl = false
		//}
		//// 设置  TargetRTP 的值
		//if IsRtpControl {
		//	targetRTPField.SetInt(int64(prc.ControlRTP))
		//	relieveRTPField.SetInt(int64(prc.AutoRemoveRTP))
		//	rewardPercent.SetInt(prc.RewardPercent)
		//	noAwardPercent.SetInt(prc.NoAwardPercent)
		//}

		if targetRTPField.IsValid() && relieveRTPField.IsValid() && rewardPercent.IsValid() && noAwardPercent.IsValid() {
			targetRTPField.SetInt(int64(prc.ControlRTP))
			relieveRTPField.SetInt(int64(prc.AutoRemoveRTP))
			rewardPercent.SetInt(prc.RewardPercent)
			noAwardPercent.SetInt(prc.NoAwardPercent)
		}

		buyMinAwardPercent := v.FieldByName(BuyMinAwardPercent)
		//isBuyControl := true
		//// 检查字段是否存在
		//if !buyMinAwardPercent.IsValid() {
		//	isBuyControl = false
		//}
		//// 设置
		//if isBuyControl {
		//	buyMinAwardPercent.SetInt(prc.BuyMinAwardPercent)
		//}
		if buyMinAwardPercent.IsValid() {
			buyMinAwardPercent.SetInt(prc.BuyMinAwardPercent)
		}

		personWinMaxMult := v.FieldByName(PersonWinMaxMult)
		personWinMaxScore := v.FieldByName(PersonWinMaxScore)

		if personWinMaxMult.IsValid() && personWinMaxScore.IsValid() {
			personWinMaxMult.SetInt(int64(prc.PersonWinMaxMult))
			personWinMaxScore.SetInt(int64(prc.PersonWinMaxScore))
		}

	}

}

// MqSendGC 判断是否需要释放RTP控制
func MqSendGC(ptr any) {
	v := reflect.ValueOf(ptr).Elem()

	targetRTP := v.FieldByName(TargetRTP)
	//不存在直接退出
	if !targetRTP.IsValid() || targetRTP.IsZero() {
		return
	}
	relieveRTP := v.FieldByName(RelieveRTP)
	//新特性，如果是 0 ，那么继续控制并且永不释放
	if !relieveRTP.IsValid() || relieveRTP.IsZero() {
		return
	}
	betAmount := v.FieldByName(BetAmount)
	winAmount := v.FieldByName(WinAmount)

	CurrentRTP := float64(winAmount.Int()) / float64(betAmount.Int()) * 100
	if CurrentRTP <= 0 {
		return
	}
	relieve := false
	//解除限制 两个范围控制  t__c__r__c__t
	if float64(targetRTP.Int()) <= CurrentRTP && CurrentRTP <= float64(relieveRTP.Int()) {
		relieve = true
	} else if float64(targetRTP.Int()) >= CurrentRTP && CurrentRTP >= float64(relieveRTP.Int()) {
		relieve = true
	}
	if relieve {
		pid := v.FieldByName(PID).Int()
		RelievePlayer(pid)
	}
}
