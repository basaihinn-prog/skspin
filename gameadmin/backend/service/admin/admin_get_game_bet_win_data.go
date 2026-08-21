package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetGameBetWinData", "获取游戏下注赢分数据", "AdminInfo", getGameBetWinData, getGameBetWinDataParams{})
}

type getGameBetWinDataParams struct {
	OperatorID  int64  `json:"OperatorID" bson:"operator_id"`
	CurrencyKey string `json:"CurrencyKey" bson:"currencyKey"`
}

type getGameBetWinDataResults struct {
	NowDay   []gameBetWinData
	SevenDay []gameBetWinData
	NowMonth []gameBetWinData
}

type gameBetWinData struct {
	GameID string `bson:"_id"`
	Bet    int64  `bson:"betamount"`
	Win    int64  `bson:"winamount"`
}

func getGameBetWinData(ctx *Context, ps getGameBetWinDataParams, ret *getGameBetWinDataResults) (err error) {

	if ps.CurrencyKey != "" && ps.OperatorID != 0 {
		return errors.New(lang.Get(ctx.Lang, "商戶和币种只能选择其中一个"))
	} else if ps.CurrencyKey == "" && ps.OperatorID == 0 {
		return errors.New(lang.Get(ctx.Lang, "请在商戶和币种选择其中一个"))
	}

	var results []struct {
		Appid string `bson:"AppID"`
	}
	//ps.CurrencyKey = "VND"
	appidlist := []string{}
	//取appid
	if ps.CurrencyKey != "" {
		err = CollAdminOperator.FindAll(bson.M{"CurrencyKey": ps.CurrencyKey}, &results)
		if err != nil {
			return
		}
		for _, res := range results {
			appidlist = append(appidlist, res.Appid)
		}
	} else {
		var oper *comm.Operator_V2
		err = CollAdminOperator.FindId(ps.OperatorID, &oper)
		if err != nil {
			return err
		}
		appidlist = append(appidlist, oper.AppID)
	}

	coll := db.Collection2("reports", "BetDailyReport")

	getData := func(startT, endT string) []gameBetWinData {
		res := make([]gameBetWinData, 0, 10)
		match := bson.M{
			"date": bson.M{
				"$gte": startT,
				"$lte": endT,
			},
			"appid": bson.M{
				"$in": appidlist,
			},
		}
		cursor, err := coll.Aggregate(context.TODO(), []bson.M{
			{
				"$match": match,
			}, {
				"$group": bson.M{
					"_id":       "$game",
					"winamount": bson.M{"$sum": "$winamount"},
					"betamount": bson.M{"$sum": "$betamount"},
				},
			},
		})
		err = cursor.All(context.TODO(), &res)
		if err != nil {
			return res
		}
		return res
	}

	now := time.Now()
	//当天
	nowd := now.Format("20060102")
	ret.NowDay = getData(nowd, nowd)
	//7天之前
	nowB7 := now.AddDate(0, 0, -6).Format("20060102")
	ret.SevenDay = getData(nowB7, nowd)

	//本月
	monthFirst := fmt.Sprintf("%s00", now.Format("200601"))
	monthEnd := fmt.Sprintf("%s31", now.Format("200601"))
	ret.NowMonth = getData(monthFirst, monthEnd)

	return nil
}
