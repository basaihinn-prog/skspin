package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/duck/lang"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	//每日平台汇总
	//每月平台汇总
	//平台汇总
	RegMsgProc("/AdminInfo/GetOperatorReportData", "获取运营商报表数据", "AdminInfo", getOperatorReportData, getOperatorReportDataParams{
		Date:      "",
		EndDate:   "",
		Operator:  1,
		PageIndex: 0,
		PageSize:  0,
		Type:      "",
	})
}

type getOperatorReportDataParams struct {
	Date      string
	EndDate   string
	Operator  int64
	PageIndex int64
	PageSize  int64
	Type      string
}
type opDailyDoc struct {
	ID             map[string]any `json:"-" bson:"_id"`
	BetAmount      int64          `bson:"betamount"` // 支出
	WinAmount      int64          `bson:"winamount"` // 赚取
	SpinCount      int            `bson:"spincount"` // 转动次数
	RoundCount     int            `bson:"roundcount"`
	LoginPlrCount  int            `bson:"loginplrcount"`
	RegistPlrCount int            `bson:"registplrcount"`

	Date            string  `bson:"date"`
	AppID           string  `bson:"appid"`
	GGR             float64 `bson:"GGR" json:"GGR"`
	OperatorBalance float64 `bson:"OperatorBalance"`
	CurrencyKey     string  `json:"CurrencyKey"`
	CurrencyName    string  `json:"CurrencyName"`
	CurrencySymbol  string  `json:"CurrencySymbol"`
}

type getOperatorReportResults struct {
	List []*opDailyDoc
	All  int64
}

func getOperatorReportAllData(ctx *Context, ps getOperatorReportDataParams, ret *getOperatorReportResults) (err error) {
	//db：reports / OperatorDailyReport
	if ps.PageIndex == 0 || ps.PageSize == 0 {
		ps.PageIndex = 1
		ps.PageSize = 20
	}

	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	query := bson.M{}

	//日期过滤
	RDate := ps.Date + ps.EndDate
	switch RDate {
	case "":
	case ps.Date:
		query["date"] = bson.M{"$gte": ps.Date}
	case ps.EndDate:
		query["date"] = bson.M{"$lte": ps.EndDate}
	default:
		if len(RDate) > 0 {
			query["date"] = bson.M{
				"$gte": ps.Date,
				"$lte": ps.EndDate,
			}
		}
	}

	appidlist := []string{}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}
	if ps.Operator > 0 && user.GroupId != 3 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.Operator, &operator)
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
	query["appid"] = bson.M{"$in": appidlist}

	match := bson.M{"$match": query}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				//"date":  "$date",
				"appid": "$appid",
			},
			//"date":           bson.M{"$first": "$date"},
			"betamount":      bson.M{"$sum": "$betamount"},
			"spincount":      bson.M{"$sum": "$spincount"},
			"winamount":      bson.M{"$sum": "$winamount"},
			"loginplrcount":  bson.M{"$sum": "$loginplrcount"},
			"registplrcount": bson.M{"$sum": "$registplrcount"},
			"roundcount":     bson.M{"$sum": "$roundcount"},
		}}

	addfile := bson.M{
		"$addFields": bson.M{
			//"date":  "$_id.date",
			"appid": "$_id.appid",
		}}

	sortfile := bson.M{
		"$sort": bson.M{
			"date":  -1,
			"appid": 1,
		}}

	//pageSkip := bson.M{"$skip": (ps.PageIndex - 1) * ps.PageSize}
	//
	//pageLimit := bson.M{"$limit": ps.PageSize}

	var aggregate []bson.M
	//aggregate = append(aggregate, match, group, addfile, sortfile, pageSkip, pageLimit)
	aggregate = append(aggregate, match, group, addfile, sortfile)

	curor, err := NewOtherDB("reports").Collection("OperatorDailyReport").Coll().Aggregate(context.TODO(), aggregate)
	err = curor.All(context.TODO(), &ret.List)
	if err != nil {
		return err
	}
	ret.All = int64(len(ret.List))
	start := (ps.PageIndex - 1) * ps.PageSize
	end := start + ps.PageSize
	if end > ret.All {
		end = ret.All
	}
	ret.List = ret.List[start:end]

	return
}

func getOperatorReportData(ctx *Context, ps getOperatorReportDataParams, ret *getOperatorReportResults) (err error) {

	switch ps.Type {
	case comm.ALL:
		//var all getOperatorReportAllResults
		return getOperatorReportAllData(ctx, ps, ret)
	case comm.Month:
		return getOperatorReportMonthData(ctx, ps, ret)
	}

	now := time.Now().UTC()
	if ps.Date == "" || ps.EndDate == "" {
		ps.EndDate = now.Format("20060102")
		ps.Date = now.AddDate(0, 0, -15).Format("20060102")
	}
	user, err := GetUser(ctx)
	if err != nil {
		return
	}

	appidlist := []string{}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}
	if ps.Operator > 0 && user.GroupId != 3 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.Operator, &operator)
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

	match := bson.M{
		"$match": bson.M{
			"date": bson.M{
				"$gte": ps.Date,
				"$lte": ps.EndDate,
			},
			"appid": bson.M{
				"$in": appidlist,
			},
		}}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"date":  "$date",
				"appid": "$appid",
			},
			"betamount":       bson.M{"$sum": "$betamount"},
			"spincount":       bson.M{"$sum": "$spincount"},
			"winamount":       bson.M{"$sum": "$winamount"},
			"loginplrcount":   bson.M{"$sum": "$loginplrcount"},
			"registplrcount":  bson.M{"$sum": "$registplrcount"},
			"roundcount":      bson.M{"$sum": "$roundcount"},
			"GGR":             bson.M{"$sum": "$GGR"},
			"OperatorBalance": bson.M{"$first": "$OperatorBalance"},
		}}

	addfile := bson.M{
		"$addFields": bson.M{
			"date":  "$_id.date",
			"appid": "$_id.appid",
		}}

	sortfile := bson.M{
		"$sort": bson.M{
			"date":  -1,
			"appid": 1,
		}}

	var aggregate []bson.M
	aggregate = append(aggregate, match, group, addfile, sortfile)
	curor, err := NewOtherDB("reports").Collection("OperatorDailyGGR").Coll().Aggregate(context.TODO(), aggregate)
	if err != nil {
		return err
	}
	list := []*opDailyDoc{}
	err = curor.All(context.TODO(), &list)

	if err != nil {
		return err
	}

	nowdate := time.Now().UTC().Format("20060102")

	if ps.Date == nowdate || ps.EndDate == nowdate {
		aggregate[0] = bson.M{
			"$match": bson.M{
				"date": nowdate,
				"appid": bson.M{
					"$in": appidlist,
				},
			}}
		nowDateGGR, err := getNowDayReport(aggregate, appidlist)
		if err != nil {
			return err
		}
		list = append(list, nowDateGGR...)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Date > list[j].Date // 按时间升序排序
	})

	ret.All = int64(len(list))
	// 计算起始和结束索引
	start := (ps.PageIndex - 1) * ps.PageSize
	end := start + ps.PageSize

	// 兜底机制：确保 start 和 end 在合理范围内
	if start > ret.All {
		start = ret.All // 如果起始索引超过总记录数，则设置为总记录数
	}
	if end > ret.All {
		end = ret.All // 如果结束索引超过总记录数，则设置为总记录数
	}
	// 如果总记录数小于 PageSize，则调整 end
	if ret.All < ps.PageSize {
		end = ret.All
	}

	ret.List = list[start:end]
	operck, err := GetCurrencyList(ctx)
	for _, item := range ret.List {
		if _, ok := operck[item.AppID]; !ok {
			operck[item.AppID] = &comm.CurrencyType{
				CurrencyCode:   "/",
				CurrencyName:   "/",
				CurrencySymbol: "/",
			}
		}

		if ctx.Lang != "zh" {
			item.CurrencyName = operck[item.AppID].CurrencyCode
			item.CurrencyKey = operck[item.AppID].CurrencyCode
			item.CurrencySymbol = operck[item.AppID].CurrencySymbol
		} else {
			item.CurrencyName = operck[item.AppID].CurrencyName
			item.CurrencyKey = operck[item.AppID].CurrencyCode
			item.CurrencySymbol = operck[item.AppID].CurrencySymbol
		}
	}
	return
}

func getNowDayReport(aggregate []bson.M, appidlist []string) (op []*opDailyDoc, err error) {

	curor, err := NewOtherDB("reports").Collection("OperatorDailyReport").Coll().Aggregate(context.TODO(), aggregate)
	if err != nil {
		return nil, err
	}
	err = curor.All(context.TODO(), &op)
	operlist := []*comm.Operator_V2{}
	filter := bson.M{
		"AppID": bson.M{
			"$in": appidlist,
		}}
	err = CollAdminOperator.FindAll(filter, &operlist)

	ratelist := make(map[string]float64)
	mapCooperstion := make(map[string]int)
	mapTurnoverRale := make(map[string]float64)
	for _, vl := range operlist {
		ratelist[vl.AppID] = vl.PlatformPay
		mapCooperstion[vl.AppID] = vl.CooperationType
		mapTurnoverRale[vl.AppID] = vl.TurnoverPay
	}

	for i, vl := range op {
		if mapCooperstion[vl.AppID] == 1 {
			op[i].GGR = float64(vl.BetAmount-vl.WinAmount) * (ratelist[vl.AppID] * 0.01)
		} else {
			op[i].GGR = float64(vl.BetAmount) * (mapTurnoverRale[vl.AppID] * 0.01)
		}
	}
	return
}

func getOperatorReportMonthData(ctx *Context, ps getOperatorReportDataParams, ret *getOperatorReportResults) (err error) {

	now := time.Now()
	if ps.Date == "" || ps.EndDate == "" {
		ps.EndDate = fmt.Sprintf("%s31", now.Format("200601")) // 月末
		ps.Date = fmt.Sprintf("%s00", now.Format("200601"))
	}

	user, err := GetUser(ctx)
	if err != nil {
		return
	}

	appidlist := []string{}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}
	if ps.Operator > 0 && user.GroupId != 3 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.Operator, &operator)
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

	match := bson.M{
		"$match": bson.M{
			"date": bson.M{
				"$gte": ps.Date,
				"$lte": ps.EndDate,
			},
			"appid": bson.M{
				"$in": appidlist,
			},
		}}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"date":  bson.M{"$substr": []interface{}{"$date", 0, 6}},
				"appid": "$appid",
			},
			"betamount":       bson.M{"$sum": "$betamount"},
			"spincount":       bson.M{"$sum": "$spincount"},
			"winamount":       bson.M{"$sum": "$winamount"},
			"loginplrcount":   bson.M{"$sum": "$loginplrcount"},
			"registplrcount":  bson.M{"$sum": "$registplrcount"},
			"roundcount":      bson.M{"$sum": "$roundcount"},
			"GGR":             bson.M{"$sum": "$GGR"},
			"OperatorBalance": bson.M{"$first": "$OperatorBalance"},
		}}

	addfile := bson.M{
		"$addFields": bson.M{
			"date":  "$_id.date",
			"appid": "$_id.appid",
		}}

	sortfile := bson.M{
		"$sort": bson.M{
			"date":  -1,
			"appid": 1,
		}}

	var aggregate []bson.M
	aggregate = append(aggregate, match, group, addfile, sortfile)

	curor, err := NewOtherDB("reports").Collection("OperatorDailyGGR").Coll().Aggregate(context.TODO(), aggregate)
	err = curor.All(context.TODO(), &ret.List)
	if err != nil {
		return
	}

	psMonth := ps.Date[:6]
	nowMonth := now.Format("200601")
	if psMonth == nowMonth {

		aggregate[0] = bson.M{"$match": bson.M{
			"date": time.Now().UTC().Format("20060102"),
			"appid": bson.M{
				"$in": appidlist,
			},
		}}

		nowlist, err := getNowDayReport(aggregate, appidlist)
		if err != nil {
			return err
		}

		mapnowOp := map[string]*opDailyDoc{}
		for _, doc := range nowlist {
			mapnowOp[doc.AppID] = doc
		}
		if ret.List == nil {
			ret.List = nowlist
		} else {
			for i, doc := range ret.List {
				if _, ok := mapnowOp[doc.AppID]; !ok {
					mapnowOp[doc.AppID] = &opDailyDoc{}
					mapnowOp[doc.AppID].OperatorBalance = doc.OperatorBalance
				}
				ret.List[i].BetAmount += mapnowOp[doc.AppID].BetAmount
				ret.List[i].WinAmount += mapnowOp[doc.AppID].WinAmount
				ret.List[i].SpinCount += mapnowOp[doc.AppID].SpinCount
				ret.List[i].RoundCount += mapnowOp[doc.AppID].RoundCount
				ret.List[i].LoginPlrCount += mapnowOp[doc.AppID].LoginPlrCount
				ret.List[i].RegistPlrCount += mapnowOp[doc.AppID].RegistPlrCount
				ret.List[i].GGR += mapnowOp[doc.AppID].GGR
				ret.List[i].OperatorBalance = mapnowOp[doc.AppID].OperatorBalance
				delete(mapnowOp, doc.AppID)
			}
			for _, doc := range mapnowOp {
				ret.List = append(ret.List, doc)
			}
		}
	}

	ret.All = int64(len(ret.List))
	// 计算起始和结束索引
	start := (ps.PageIndex - 1) * ps.PageSize
	end := start + ps.PageSize

	// 兜底机制：确保 start 和 end 在合理范围内
	if start > ret.All {
		start = ret.All // 如果起始索引超过总记录数，则设置为总记录数
	}
	if end > ret.All {
		end = ret.All // 如果结束索引超过总记录数，则设置为总记录数
	}

	// 如果总记录数小于 PageSize，则调整 end
	if ret.All < ps.PageSize {
		end = ret.All
	}
	operck, err := GetCurrencyList(ctx)
	for _, item := range ret.List {
		if _, ok := operck[item.AppID]; !ok {
			operck[item.AppID] = &comm.CurrencyType{
				CurrencyCode:   "/",
				CurrencyName:   "/",
				CurrencySymbol: "/",
			}
		}

		if ctx.Lang != "zh" {
			item.CurrencyName = operck[item.AppID].CurrencyCode
			item.CurrencyKey = operck[item.AppID].CurrencyCode
			item.CurrencySymbol = operck[item.AppID].CurrencySymbol
		} else {
			item.CurrencyName = operck[item.AppID].CurrencyName
			item.CurrencyKey = operck[item.AppID].CurrencyCode
			item.CurrencySymbol = operck[item.AppID].CurrencySymbol
		}
	}

	return nil
}

type LotteryReportInfo struct {
	Id             int64 `json:"-" bson:"_id"`
	WinAmount      int64 `bson:"winamount"` // 产出
	BetAmount      int64 `bson:"betamount"` // 流水
	SpinCount      int   `bson:"spincount"` // 转动次数
	LoginPlrCount  int   `bson:"loginplrcount"`
	RegistPlrCount int   `bson:"registplrcount"`
}

// 1期买彩票
func getLottery1OperatorReportData(date string, appId string) (ret LotteryReportInfo, err error) {
	t, _ := time.Parse("20060102", date)
	match := bson.M{
		"buyDay": t.Format("2006-01-02"),
	}
	if appId != "" {
		match["appID"] = appId
	}
	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":            1,
				"winamount":      bson.M{"$sum": "$periodProfits"},
				"betamount":      bson.M{"$sum": "$soldGold"},
				"loginplrcount":  bson.M{"$sum": "$loginplrcount"},
				"registplrcount": bson.M{"$sum": "$registplrcount"},
				"spincount":      bson.M{"$sum": "$spincount"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []LotteryReportInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

func getLottery2OperatorReportData(date string, appId string) (ret LotteryReportInfo, err error) {
	t, _ := time.Parse("20060102", date)
	match := bson.M{
		"buyDay": t.Format("2006-01-02"),
	}
	if appId != "" {
		match["appID"] = appId
	}
	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day_Guess").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":            1,
				"winamount":      bson.M{"$sum": "$periodGuessProfits"},
				"betamount":      bson.M{"$sum": "$soldGuessGold"},
				"loginplrcount":  bson.M{"$sum": "$loginplrcount"},
				"registplrcount": bson.M{"$sum": "$registplrcount"},
				"spincount":      bson.M{"$sum": "$spincount"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []LotteryReportInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

func getLottery1OperatorReportMonthData(yearMonth string, appId string) (ret LotteryReportInfo, err error) {
	start, end := getMonthFirstAndLastDay(yearMonth)
	match := bson.M{
		"buyDay": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	//fmt.Println("getLottery1OperatorReportMonthData invoke start:", start, " end:", end)
	if appId != "" {
		match["appID"] = appId
	}
	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":            1,
				"winamount":      bson.M{"$sum": "$periodProfits"},
				"betamount":      bson.M{"$sum": "$soldGold"},
				"loginplrcount":  bson.M{"$sum": "$loginplrcount"},
				"registplrcount": bson.M{"$sum": "$registplrcount"},
				"spincount":      bson.M{"$sum": "$spincount"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []LotteryReportInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

func getLottery2OperatorReportMonthData(yearMonth string, appId string) (ret LotteryReportInfo, err error) {
	start, end := getMonthFirstAndLastDay(yearMonth)
	match := bson.M{
		"buyDay": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	//fmt.Println("getLottery2OperatorReportMonthData invoke start:", start, " end:", end)
	if appId != "" {
		match["appID"] = appId
	}
	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day_Guess").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":            1,
				"winamount":      bson.M{"$sum": "$periodGuessProfits"},
				"betamount":      bson.M{"$sum": "$soldGuessGold"},
				"loginplrcount":  bson.M{"$sum": "$loginplrcount"},
				"registplrcount": bson.M{"$sum": "$registplrcount"},
				"spincount":      bson.M{"$sum": "$spincount"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []LotteryReportInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

func getMonthFirstAndLastDay(date string) (day1st, daylast string) {
	year, _ := time.Parse("200601", date)
	month := year.Month()

	// 获取该月的第一天
	firstDay := time.Date(year.Year(), month, 1, 0, 0, 0, 0, time.UTC)

	// 获取该月的下一个月的第一天，然后减去一天得到该月的最后一天
	nextMonth := firstDay.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	fmt.Println("该月的第一天是：", firstDay)
	fmt.Println("该月的最后一天是：", lastDay)

	return firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02")
}
