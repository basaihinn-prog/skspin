package main

import (
	"context"
	"fmt"
	"game/comm"
	"game/duck/mongodb"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/GetPokerOperatorReportData", "获取扑克运营商报表数据", "AdminInfo", getPokerOperatorReportData, getPokerOperatorReportDataParams{
		getOperatorReportDataParams: getOperatorReportDataParams{
			Date:      "",
			EndDate:   "",
			Operator:  1,
			PageIndex: 0,
			PageSize:  0,
			Type:      "",
		},
	})
}

type getPokerOperatorReportDataParams struct {
	getOperatorReportDataParams
}

type getPokerOperatorReportResults struct {
	List []*comm.PokerOperatorReportData
	All  int64
}

func getPokerOperatorReportData(ctx *Context, ps getOperatorReportDataParams, ret *getPokerOperatorReportResults) (err error) {
	ps.PageIndex = 0
	ps.PageSize = 2000
	if ps.Type == comm.Month {
		return getPokerOperatorReportMonthData(ctx, ps, ret)
	}
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
		curor, err := NewOtherDB("reports").Collection("BetPokerDailyReport").Coll().Aggregate(context.TODO(), []bson.M{
			{
				"$match": bson.M{
					"date": bson.M{
						"$gte": ps.Date,
						"$lte": ps.EndDate,
					},
				},
			},
			{
				"$group": bson.M{
					"_id":            "$date",
					"betamount":      bson.M{"$sum": "$betamount"},
					"spincount":      bson.M{"$sum": "$spincount"},
					"winamount":      bson.M{"$sum": "$winamount"},
					"enterplrcount":  bson.M{"$sum": "$enterplrcount"},
					"registplrcount": bson.M{"$sum": "$registplrcount"},
					"roundcount":     bson.M{"$sum": "$roundcount"},
					"flow":           bson.M{"$sum": "$flow"},
				},
			},
			{
				"$addFields": bson.M{
					"date": "$_id",
				},
			},
		})
		if err != nil {
			return err
		}
		err = curor.All(context.TODO(), &ret.List)
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
			query["appid"] = user.AppID
		}
		filter.Query = query
		count, err := NewOtherDB("reports").Collection("BetPokerDailyReport").FindPage(filter, &ret.List)
		if err != nil {
			return err
		}
		ret.All = count
	}

	sort.Slice(ret.List, func(i, j int) bool {
		return ret.List[i].Date > ret.List[j].Date
	})
	return
}

func getPokerOperatorReportMonthData(ctx *Context, ps getOperatorReportDataParams, ret *getPokerOperatorReportResults) (err error) {
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	if ps.Date == "" || ps.EndDate == "" {
		now := time.Now()
		if user.AppID == "admin" {
			ps.EndDate = fmt.Sprintf("%s31", now.Format("200601")) // 月末
			ps.Date = fmt.Sprintf("%s00", now.AddDate(0, -11, 0).Format("200601"))
		} else {
			ps.EndDate = fmt.Sprintf("%s31", now.Format("200601")) // 月末
			ps.Date = fmt.Sprintf("%s00", now.AddDate(0, -1, 0).Format("200601"))
		}
	}
	var match bson.M
	appid := ""
	if user.AppID == "admin" && ps.Operator == 0 {
		match = bson.M{
			"date": bson.M{
				"$gte": ps.Date,
				"$lte": ps.EndDate,
			},
		}
	} else {
		if user.AppID == "admin" {
			var operator *comm.Operator
			err = CollAdminOperator.FindId(ps.Operator, &operator)
			if err != nil {
				return err
			}
			appid = operator.AppID

		} else {
			appid = user.AppID
		}
		match = bson.M{
			"appid": appid,
			"date": bson.M{
				"$gte": ps.Date,
				"$lte": ps.EndDate,
			},
		}
	}
	curor, err := NewOtherDB("reports").Collection("BetPokerDailyReport").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$substr": []any{"$date", 0, 6},
				},
				"betamount":      bson.M{"$sum": "$betamount"},
				"spincount":      bson.M{"$sum": "$spincount"},
				"winamount":      bson.M{"$sum": "$winamount"},
				"enterplrcount":  bson.M{"$sum": "$enterplrcount"},
				"registplrcount": bson.M{"$sum": "$registplrcount"},
				"roundcount":     bson.M{"$sum": "$roundcount"},
				"flow":           bson.M{"$sum": "$flow"},
			},
		},
		{
			"$addFields": bson.M{
				"date": "$_id",
			},
		},
	})
	if err != nil {
		return err
	}
	err = curor.All(context.TODO(), &ret.List)
	if err != nil {
		return err
	}
	ret.All = int64(len(ret.List))
	sort.Slice(ret.List, func(i, j int) bool {
		return ret.List[i].Date > ret.List[j].Date
	})
	return nil
}
