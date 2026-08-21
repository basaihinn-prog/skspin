package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"slices"
	"strings"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetPlayerReportData", "获取玩家报表数据", "AdminInfo", getPlayerReportDataV2, getPlayerReportDataParams{
		Date:      "",
		EndDate:   "",
		Operator:  1,
		PageIndex: 0,
		PageSize:  0,
		//Type:      "",
	})
	RegMsgProc("/AdminInfo/GetRankingList", "获取玩家排行榜", "AdminInfo", getRankingList, getRankingListParams{
		Page:     0,
		PageSize: 0,
		//Type:      "",
	})
}

type getPlayerReportDataParams struct {
	Date      string
	EndDate   string
	Operator  int64
	PageIndex int64
	PageSize  int64
	//Type      string
}
type PlayerDailyDoc struct {
	ID         string `json:"-" bson:"_id"`
	BetAmount  int64  `bson:"betamount"` // 支出
	WinAmount  int64  `bson:"winamount"` // 赚取
	SpinCount  int    `bson:"spincount"` // 转动次数
	RoundCount int    `bson:"roundcount"`

	Date        string  `bson:"date"`
	AppID       string  `bson:"appid"`
	Pid         int64   `bson:"pid"`
	PresentRate float64 `bson:"PresentRate"`
}

type getPlayerReportDataResults struct {
	List []*PlayerDailyDoc
	All  int64
}

func getPlayerReportDataV2(ctx *Context, ps getPlayerReportDataParams, ret *getPlayerReportDataResults) (err error) {
	//ps.PageIndex = 0
	//ps.PageSize = 2000
	if ps.Date == "" || ps.EndDate == "" {
		now := time.Now()
		ps.EndDate = now.Format("20060102")
		ps.Date = now.AddDate(0, 0, -15).Format("20060102")
	}

	//var user comm.User
	//user, err = GetUser(ctx)
	//if err != nil {
	//	return err
	//}
	appidlist, err := GetOperatopAppID(ctx)
	if err != nil {
		return
	}

	filter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     bson.M{"_id": -1},
	}
	query := bson.M{
		"date": bson.M{
			"$gte": ps.Date,
			"$lte": ps.EndDate,
		},
	}
	var operator *comm.Operator

	if ps.Operator > 0 {
		err = CollAdminOperator.FindId(ps.Operator, &operator)

		if err != nil {
			return err
		}

		if slices.Contains(appidlist, operator.AppID) == false {
			return errors.New(lang.GetLang(ctx.Lang, "未代理该商户"))
		}
		query["appid"] = operator.AppID
	} else {
		query["appid"] = bson.M{"$in": appidlist}
	}

	filter.Query = query
	ret.All, err = NewOtherDB("reports").Collection("PlayerDailyReport").FindPage(filter, &ret.List)
	if err != nil {
		return err
	}
	oper := []*comm.Operator_V2{}
	err = CollAdminOperator.FindAll(bson.M{}, &oper)
	ratelist := make(map[string]float64)
	for _, vl := range oper {
		ratelist[vl.AppID] = vl.PresentRate
	}

	for i, vl := range ret.List {
		ret.List[i].PresentRate = ratelist[vl.AppID]
	}

	//if user.AppID != "admin" || ps.Operator != 0 {
	//	filter := mongodb.FindPageOpt{
	//		Page:     ps.PageIndex,
	//		PageSize: ps.PageSize,
	//		Sort:     bson.M{"_id": -1},
	//	}
	//	query := bson.M{
	//		"date": bson.M{"$gte": ps.Date, "$lte": ps.EndDate},
	//	}
	//	if user.AppID == "admin" {
	//		var operator *comm.Operator
	//		err = CollAdminOperator.FindId(ps.Operator, &operator)
	//		if err != nil {
	//			return err
	//		}
	//		query["appid"] = operator.AppID
	//	} else {
	//		var user *comm.User
	//		err = CollAdminUser.FindId(ctx.PID, &user)
	//		if err != nil {
	//			return err
	//		}
	//		query["appid"] = user.AppID
	//	}
	//	filter.Query = query
	//	_, err := NewOtherDB("reports").Collection("PlayerDailyReport").FindPage(filter, &ret.List, false)
	//	if err != nil {
	//		return err
	//	}
	//	ret.All = int64(len(ret.List))
	//}
	return
}

func getPlayerReportData(ctx *Context, ps getPlayerReportDataParams, ret *getPlayerReportDataResults) (err error) {
	ps.PageIndex = 0
	ps.PageSize = 2000
	if ps.Date == "" || ps.EndDate == "" {
		now := time.Now()
		ps.EndDate = now.Format("20060102")
		ps.Date = now.AddDate(0, 0, -15).Format("20060102")
	}

	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}
	if user.AppID == "admin" && ps.Operator == 0 { //获取所有运营商的
		filter := mongodb.FindPageOpt{
			Page:     ps.PageIndex,
			PageSize: ps.PageSize,
			Sort:     bson.M{"_id": -1},
		}
		query := bson.M{
			"date": bson.M{
				"$gte": ps.Date,
				"$lte": ps.EndDate,
			},
		}
		filter.Query = query
		_, err := NewOtherDB("reports").Collection("PlayerDailyReport").FindPage(filter, &ret.List, false)
		if err != nil {
			return err
		}
		ret.All = int64(len(ret.List))
	} else {
		filter := mongodb.FindPageOpt{
			Page:     ps.PageIndex,
			PageSize: ps.PageSize,
			Sort:     bson.M{"_id": -1},
		}
		query := bson.M{
			"date": bson.M{"$gte": ps.Date, "$lte": ps.EndDate},
		}
		if user.AppID == "admin" {
			var operator *comm.Operator
			err = CollAdminOperator.FindId(ps.Operator, &operator)
			if err != nil {
				return err
			}
			query["appid"] = operator.AppID
		} else {
			var user *comm.User
			err = CollAdminUser.FindId(ctx.PID, &user)
			if err != nil {
				return err
			}
			query["appid"] = user.AppID
		}
		filter.Query = query
		_, err := NewOtherDB("reports").Collection("PlayerDailyReport").FindPage(filter, &ret.List, false)
		if err != nil {
			return err
		}
		ret.All = int64(len(ret.List))
	}
	return
}

type rankingList struct {
	Pid            uint64             `ch:"Pid"`
	TotalWinLoss   int64              `ch:"totalWinLose"`
	TotalWin       int64              `ch:"totalWin"`
	HistoryWinLoss int64              `ch:"historyWinLoss"`
	Balance        int64              `ch:"Balance"`
	SpinCount      int64              `ch:"spinCount"`
	TotalBet       int64              `ch:"totalBet"`
	AppID          string             `ch:"AppID"`
	WinRate        float64            `ch:"winRate"`
	CurrencyCode   string             `json:"CurrencyType" bson:"CurrencyKey"`  //货币名称
	CurrencyName   string             `json:"CurrencyName" bson:"CurrencyName"` //货币名称
	LoginAt        *mongodb.TimeStamp `json:"LoginAt" bson:"LoginAt"`
	CreateAt       *mongodb.TimeStamp `json:"CreateAt" bson:"CreateAt"`
}
type getRankingListResults struct {
	List  []*rankingList
	Count int64
}
type getRankingListParams struct {
	OperatorId      int64 //商户ID
	StartTime       int64
	EndTime         int64
	GameID          string //游戏ID
	PageSize        int64
	Page            int64
	Manufacturer    string
	RankingListType int64
	GameType        int
}

func generateRankingQuery(tableName string, ps getRankingListParams, appidlist []string) (queryStr, historyArgStr string, countCache int64, ArgList []any, HistoryArgList []any, err error) {
	queryStr = "select Pid, sum(Win) as totalWin, sum(Win)-sum(Bet) as totalWinLose, sum(Bet) as totalBet, countIf(Bet > 0) as spinCount, CASE  WHEN sum(Bet) = 0 THEN 0 ELSE sum(Win) / sum(Bet) END AS winRate from " + tableName + " where"
	countStr := "select COUNT(DISTINCT Pid)  from " + tableName + " where"
	argStr := ""
	historyArgStr = ""
	historyArgList := []any{}
	argList := []any{}
	argStr += " LogType = ? "
	historyArgStr += " LogType = ? "
	argList = append(argList, 0)
	historyArgList = append(historyArgList, 0)
	//argStr += " and Bet > ?"
	//argList = append(argList, 0)
	if tableName == "aviatorgamelogs" {
		argStr += " and FinishType = ?"
		historyArgStr += " and FinishType = ?"
		argList = append(argList, 1)
		historyArgList = append(historyArgList, 1)
	}
	if ps.Manufacturer != "" {
		argStr += " and ManufacturerName = ?"
		historyArgStr += " and ManufacturerName = ?"
		argList = append(argList, ps.Manufacturer)
		historyArgList = append(historyArgList, ps.Manufacturer)
	}

	if ps.GameID != "ALL" && ps.GameID != "" {
		argStr += " and GameID = ?"
		historyArgStr += " and GameID = ?"
		argList = append(argList, ps.GameID)
		historyArgList = append(historyArgList, ps.GameID)
	}

	sorts := bson.D{}
	sorts = append(sorts, bson.E{Key: "_id", Value: -1})

	RDate := ps.EndTime + ps.StartTime
	switch RDate {
	case 0:
	case ps.EndTime:
		argStr += " and InsertTime <= ?"
		argList = append(argList, ps.EndTime)
	case ps.StartTime:
		argStr += " and InsertTime >= ?"
		argList = append(argList, ps.StartTime)
	default:
		argStr += " and InsertTime >= ? and InsertTime <= ?"
		argList = append(argList, ps.StartTime, ps.EndTime)
	}
	conditions := []string{}
	for _, appid := range appidlist {
		argList = append(argList, appid)
		conditions = append(conditions, "?")
	}

	argStr += " and AppID in (" + strings.Join(conditions, ",") + ")"

	conn, err := db.ClickHouseCollection("")
	if err != nil {
		return
	}
	countCache = int64(0)
	err = conn.QueryRow(countStr+argStr, argList...).Scan(&countCache)
	if err != nil {
		return
	}
	pagestart := (ps.Page - 1) * ps.PageSize
	pageend := pagestart + ps.PageSize
	if pagestart > countCache {
		return
	}

	if pageend > countCache {
		pageend = countCache
	}
	if err != nil {
		return
	}
	sort := "ORDER BY totalWinLose desc"
	switch ps.RankingListType {
	case 1:
		sort = "ORDER BY totalWinLose asc"
	case 2:
		sort = "ORDER BY totalBet desc"
	case 3:
		sort = "ORDER BY winRate desc"
	}

	argList = append(argList, ps.PageSize, (ps.Page-1)*ps.PageSize)
	queryStr = queryStr + argStr + " group by Pid " + sort + " LIMIT ? OFFSET ?"
	return queryStr, historyArgStr, countCache, argList, historyArgList, err
}

func getRankingList(ctx *Context, ps getRankingListParams, ret *getRankingListResults) (err error) {

	user, err := GetUser(ctx)
	if err != nil {
		return err
	}

	//appId := user.AppID
	appidlist := []string{}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}
	if len(appidlist) == 0 {
		return lang.Error(ctx.Lang, "这个币种下没有商户")
	}

	if ps.OperatorId != 0 {
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
			return errors.New(lang.GetLang(ctx.Lang, "未代理该商户"))
		}
		appidlist = appidlist[:0]
		appidlist = append(appidlist, operator.AppID)

	}

	if len(appidlist) == 0 {
		return lang.Error(ctx.Lang, "参数错误")
	}
	tableName := ""
	if ps.GameType == 0 {
		tableName = "gamelogs"
	} else {
		tableName = "aviatorgamelogs"
	}
	conn, err := db.ClickHouseCollection("")
	if err != nil {
		return
	}
	argList := []any{}
	historyArgList := []any{}
	queryStr, historyArgStr := "", ""
	queryStr, historyArgStr, ret.Count, argList, historyArgList, err = generateRankingQuery(tableName, ps, appidlist)

	rows, err := conn.Query(queryStr, argList...)
	if err != nil {
		return err
	}
	defer rows.Close()
	list := []*rankingList{}
	conditions := []string{}
	for rows.Next() {
		betlog := &rankingList{}
		//if err = rows.Scan(&betlog.ID, &betlog.Pid, &betlog.UserID, &betlog.GameID, &betlog.GameType, &betlog.RoundID, &betlog.Bet, &betlog.Win, &betlog.Comment, &betlog.InsertTime, &betlog.Balance, &betlog.Completed, &betlog.TotalWinLoss, &betlog.WinLose, &betlog.SpinDetailsJson, &betlog.LogType, &betlog.PGBetID, &betlog.AppID, &betlog.CurrencyCode, &betlog.CurrencyName, &betlog.Grade, &betlog.HitBigReward, &betlog.LogType, &betlog.TransferAmount); err != nil {
		//	return err
		//}
		if err = rows.Scan(
			&betlog.Pid,
			&betlog.TotalWin,
			&betlog.TotalWinLoss,
			&betlog.TotalBet,
			&betlog.SpinCount,
			&betlog.WinRate,
		); err != nil {
			return err
		}
		list = append(list, betlog)
		conditions = append(conditions, "?")
		historyArgList = append(historyArgList, betlog.Pid)

	}
	if len(conditions) == 0 {
		ret.List = list
		return
	}
	historyArgStr += " and Pid in (" + strings.Join(conditions, ",") + ")"
	//historyArgStr += " and Pid in (" + strings.Join(conditions, ",") + ")"
	//fmt.Println(err)

	if err != nil {
		return
	}
	historyWinlossQuery := "select Pid, sum(Win)-sum(Bet) as historyWinLoss from " + tableName + " where"
	rows, err = conn.Query(historyWinlossQuery+historyArgStr+" and InsertTime >= now() - INTERVAL 30 DAY group by Pid ", historyArgList...)
	if err != nil {
		return err
	}
	defer rows.Close()
	mapHistoryWinLoss := map[uint64]int64{}
	for rows.Next() {

		betlog := &rankingList{}
		//if err = rows.Scan(&betlog.ID, &betlog.Pid, &betlog.UserID, &betlog.GameID, &betlog.GameType, &betlog.RoundID, &betlog.Bet, &betlog.Win, &betlog.Comment, &betlog.InsertTime, &betlog.Balance, &betlog.Completed, &betlog.TotalWinLoss, &betlog.WinLose, &betlog.SpinDetailsJson, &betlog.LogType, &betlog.PGBetID, &betlog.AppID, &betlog.CurrencyCode, &betlog.CurrencyName, &betlog.Grade, &betlog.HitBigReward, &betlog.LogType, &betlog.TransferAmount); err != nil {
		//	return err
		//}
		if err = rows.Scan(
			&betlog.Pid,
			&betlog.HistoryWinLoss,
		); err != nil {
			return err
		}
		mapHistoryWinLoss[betlog.Pid] = betlog.HistoryWinLoss
	}

	//币种
	operck, err := GetCurrencyList(ctx)
	for i, log := range list {
		// -----------------------查询玩家信息--------------------
		var playerInfo player
		playerColl := db.Collection2("game", "Players")
		cour := playerColl.FindOne(context.TODO(), bson.M{"_id": log.Pid})
		err = cour.Decode(&playerInfo)
		log.AppID = playerInfo.AppID
		log.Balance = playerInfo.Balance
		log.LoginAt = playerInfo.LoginAt
		log.CreateAt = playerInfo.CreateAt
		log.WinRate = log.WinRate * 100
		log.HistoryWinLoss = mapHistoryWinLoss[log.Pid]
		if ctx.Lang != "zh" {
			list[i].CurrencyName = operck[log.AppID].CurrencyCode
			list[i].CurrencyCode = operck[log.AppID].CurrencyCode
		} else {
			list[i].CurrencyName = operck[log.AppID].CurrencyName
			list[i].CurrencyCode = operck[log.AppID].CurrencyCode
		}
	}
	ret.List = list

	return err
}
