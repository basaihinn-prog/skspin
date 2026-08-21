package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	//游戏数据 / 数据汇总
	RegMsgProc("/AdminInfo/GetBetWinTotal", "获取下注赢分统计", "AdminInfo", getBetWinTotal, getBetWinTotalParams{
		Date:       "",
		OperatorId: 0,
	})
}

type getBetWinTotalParams struct {
	Date         string //日期
	GameListId   string //游戏id(多个)
	OperatorId   int64
	Need7        bool
	AppID        string `json:"AppID" bson:"AppID"`               //新加入字段  商户名称
	CurrencyCode string `json:"CurrencyCode" bson:"CurrencyCode"` //货币代码
	Manufacturer string `json:"Manufacturer"`
}

type gameBetWinInfo struct {
	//GameID string `bson:"_id"`
	Id  map[string]any `json:"-" bson:"_id"`
	Bet int64          `bson:"betamount"`
	Win int64          `bson:"winamount"`

	GameID       string `json:"GameID" bson:"gameid"`
	AppId        string `json:"AppId" bson:"appid"`
	CurrencyCode string `json:"CurrencyType" bson:"CurrencyKey"`  //货币名称
	CurrencyName string `json:"CurrencyName" bson:"CurrencyName"` //货币名称
}

type getBetWinTotalResults struct {
	NowDay    []gameBetWinInfo
	SevenDay  []gameBetWinInfo
	NowMonth  []gameBetWinInfo
	LastMonth []gameBetWinInfo
}

func getBetWinTotal(ctx *Context, ps getBetWinTotalParams, ret *getBetWinTotalResults) (err error) {

	if ps.Date == "" { // 默认当天
		ps.Date = time.Now().Format("20060102")
	}
	operck, err := GetCurrencyList(ctx)
	if err != nil {
		return
	}
	getData := func(match bson.M) ([]gameBetWinInfo, error) {
		res := make([]gameBetWinInfo, 0, 10)

		match = bson.M{"$match": match}

		groupfilter := bson.M{
			"_id": bson.M{
				"game":  "$game",
				"appid": "$appid",
			},
			"winamount": bson.M{"$sum": "$winamount"},
			"betamount": bson.M{"$sum": "$betamount"},
		}

		group := bson.M{"$group": groupfilter}

		sort := bson.M{
			"$sort": bson.M{
				"_id": -1,
			}}

		addfile := bson.M{
			"$addFields": bson.M{
				"gameid": "$_id.game",
				"appid":  "$_id.appid",
			}}

		var aggregate []bson.M
		aggregate = append(aggregate, match, group, sort, addfile)

		cursor, err := db.Collection2("reports", "BetDailyReport").Aggregate(context.TODO(), aggregate)
		if err != nil {
			return res, err
		}
		err = cursor.All(context.TODO(), &res)

		//币种
		for i, re := range res {
			if _, ok := operck[re.AppId]; !ok {
				operck[re.AppId] = &comm.CurrencyType{
					CurrencyCode:   "/",
					CurrencyName:   "/",
					CurrencySymbol: "/",
				}
			}

			if ctx.Lang != "zh" {
				res[i].CurrencyName = operck[re.AppId].CurrencyCode
				res[i].CurrencyCode = operck[re.AppId].CurrencyCode
			} else {
				res[i].CurrencyName = operck[re.AppId].CurrencyName
				res[i].CurrencyCode = operck[re.AppId].CurrencyCode
			}

		}

		return res, nil
	}

	match := bson.M{}

	//商户和币种改为非必填  2024 11 28
	//if ps.CurrencyCode != "" && ps.OperatorId != 0 {
	//	return errors.New(lang.GetLang(ctx.Lang, "请选择商户"))
	//} else if ps.CurrencyCode == "" && ps.OperatorId == 0 {
	//	return errors.New(lang.GetLang(ctx.Lang, "请在商戶和币种选择其中一个"))
	//}
	appidlist := []string{}
	user, err := GetUser(ctx)
	if err != nil {
		return
	}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}
	if ps.OperatorId > 0 && user.GroupId != 3 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.OperatorId, &operator)
		if err != nil {
			return err
		}
		temp := false
		for _, appid := range appidlist {
			if operator.AppID == appid {
				temp = true
			}
		}
		if temp == false {
			return errors.New(lang.Get(ctx.Lang, "未代理该商户"))
		}
		appidlist = appidlist[:0]
		appidlist = append(appidlist, operator.AppID)
	}

	if ps.CurrencyCode != "" {
		operlist := []*comm.Operator{}
		filter := bson.M{
			"CurrencyKey": ps.CurrencyCode,
			"AppID": bson.M{
				"$in": appidlist,
			}}
		err = CollAdminOperator.FindAll(filter, &operlist)
		if err != nil {
			return
		}
		appidlist = appidlist[:0]
		for _, operator := range operlist {
			appidlist = append(appidlist, operator.AppID)
		}
	}

	match["appid"] = bson.M{"$in": appidlist}

	if ps.GameListId != "ALL" && ps.GameListId != "" {
		match["game"] = ps.GameListId
	} else {
		if ps.Manufacturer == "" {
			ps.Manufacturer = "ALL"
		}
		list, err := GetManufacturerGame(ps.Manufacturer)
		if err != nil {
			return err
		}
		gIdlist := []string{}
		for _, item := range list {
			gIdlist = append(gIdlist, item.ID)
		}
		match["game"] = bson.M{"$in": gIdlist}

	}

	//NowDay
	match["date"] = ps.Date

	ret.NowDay, err = getData(match)
	if err != nil {
		return err
	}

	t, err := time.Parse("20060102", ps.Date)
	if err != nil {
		return err
	}

	// SevenDay
	if ps.Need7 {
		nowB7 := t.AddDate(0, 0, -6).Format("20060102")
		match["date"] = bson.M{
			"$gte": nowB7,
			"$lte": ps.Date,
		}
		ret.SevenDay, err = getData(match)
		if err != nil {
			return err
		}
	}

	// nowMonth
	monthFirst := fmt.Sprintf("%s00", t.Format("200601"))
	monthEnd := fmt.Sprintf("%s31", t.Format("200601"))
	match["date"] = bson.M{
		"$gte": monthFirst,
		"$lte": monthEnd,
	}
	ret.NowMonth, err = getData(match)
	if err != nil {
		return err
	}

	// lastMonth
	lm := t.AddDate(0, -1, 0).Format("200601")
	lastMonthFirst := fmt.Sprintf("%s00", lm)
	lastMonthEnd := fmt.Sprintf("%s31", lm)
	match["date"] = bson.M{
		"$gte": lastMonthFirst,
		"$lte": lastMonthEnd,
	}
	ret.LastMonth, err = getData(match)
	if err != nil {
		return err
	}

	return err
}

const (
	LotteryData_NowDay = iota
	LotteryData_SevenDay
	LotteryData_NowMonth
	LotteryData_LastMonth
)

// 彩票新加
// 1期买彩票
func getLottery1Data(date string, appId string, queryType int) (res []gameBetWinInfo, err error) {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return
	}

	match := bson.M{}
	if appId != "" {
		match["appID"] = appId
	}
	queryDate := t.Format("2006-01-02")

	if queryType == LotteryData_NowDay {
		match["buyDay"] = queryDate
	} else if queryType == LotteryData_SevenDay {
		nowB7 := t.AddDate(0, 0, -6).Format("2006-01-02")
		match["buyDay"] = bson.M{
			"$gte": nowB7,
		}
	} else if queryType == LotteryData_NowMonth {
		monthFirst := fmt.Sprintf("%s-00", t.Format("2006-01"))
		match["buyDay"] = bson.M{
			"$gte": monthFirst,
		}
	} else if queryType == LotteryData_LastMonth {
		lm := t.AddDate(0, -1, 0).Format("2006-01")
		lastMonthFirst := fmt.Sprintf("%s-00", lm)
		lastMonthEnd := fmt.Sprintf("%s-31", lm)
		match["buyDay"] = bson.M{
			"$gte": lastMonthFirst,
			"$lte": lastMonthEnd,
		}
	}

	info, err := getLottery1BetAndWin(match)
	if err != nil {
		return
	}
	if info != nil {
		res = append(res, *info)
	}
	return
}

// 1期买彩票
func getLottery1BetAndWin(match bson.M) (ret *gameBetWinInfo, err error) {
	cursor, err := db.Collection2("CaiPiao", "Report_Day").Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":                bson.M{},
				"soldGold":           bson.M{"$sum": "$soldGold"},           //一期
				"soldGuessGold":      bson.M{"$sum": "$soldGuessGold"},      //二期
				"periodProfits":      bson.M{"$sum": "$periodProfits"},      //一期
				"periodGuessProfits": bson.M{"$sum": "$periodGuessProfits"}, //二期
			},
		},
	})

	if err != nil {
		return
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	if len(results) <= 0 {
		return nil, err
	}
	ret = &gameBetWinInfo{
		GameID: "lottery",
	}
	for _, result := range results {
		gold1, err := strconv.Atoi(fmt.Sprintf("%v", result["soldGold"]))
		if err != nil {
			break
		}

		gold3, err := strconv.Atoi(fmt.Sprintf("%v", result["periodProfits"]))
		if err != nil {
			break
		}

		ret.Bet += int64(gold1)

		ret.Win += int64(gold3)
	}
	return
}

// 2期猜号码
func getLottery2Data(date string, appId string, queryType int) (res []gameBetWinInfo, err error) {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return
	}

	match := bson.M{}
	if appId != "" {
		match["appID"] = appId
	}
	queryDate := t.Format("2006-01-02")

	if queryType == LotteryData_NowDay {
		match["buyDay"] = queryDate
	} else if queryType == LotteryData_SevenDay {
		nowB7 := t.AddDate(0, 0, -6).Format("2006-01-02")
		match["buyDay"] = bson.M{
			"$gte": nowB7,
		}
	} else if queryType == LotteryData_NowMonth {
		monthFirst := fmt.Sprintf("%s-00", t.Format("2006-01"))
		match["buyDay"] = bson.M{
			"$gte": monthFirst,
		}
	} else if queryType == LotteryData_LastMonth {
		lm := t.AddDate(0, -1, 0).Format("2006-01")
		lastMonthFirst := fmt.Sprintf("%s-00", lm)
		lastMonthEnd := fmt.Sprintf("%s-31", lm)
		match["buyDay"] = bson.M{
			"$gte": lastMonthFirst,
			"$lte": lastMonthEnd,
		}
	}

	info, err := getLottery2BetAndWin(match)
	if err != nil {
		return
	}
	if info != nil {
		res = append(res, *info)
	}
	return
}

// 2期猜号码
func getLottery2BetAndWin(match bson.M) (ret *gameBetWinInfo, err error) {
	cursor, err := db.Collection2("CaiPiao", "Report_Day_Guess").Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":                bson.M{},
				"soldGuessGold":      bson.M{"$sum": "$soldGuessGold"},      //二期
				"periodGuessProfits": bson.M{"$sum": "$periodGuessProfits"}, //二期
			},
		},
	})

	if err != nil {
		return
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	if len(results) <= 0 {
		return nil, err
	}
	ret = &gameBetWinInfo{
		GameID: "lottery_guess",
	}
	for _, result := range results {
		gold2, err := strconv.Atoi(fmt.Sprintf("%v", result["soldGuessGold"]))
		if err != nil {
			break
		}

		gold4, err := strconv.Atoi(fmt.Sprintf("%v", result["periodGuessProfits"]))
		if err != nil {
			break
		}

		ret.Bet += int64(gold2)

		ret.Win += int64(gold4)
	}
	return
}

// 获取某个游戏厂商下的游戏
func GetManufacturerGame(Manufacturer string) (list []*comm.Game2, err error) {

	coll := NewOtherDB("game").Collection("Games")

	filter := bson.M{}

	if Manufacturer != "ALL" {
		filter["ManufacturerName"] = Manufacturer
	}

	err = coll.FindAll(filter, &list)
	if err != nil {
		return nil, err
	}
	return
}
