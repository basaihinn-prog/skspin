package main

import (
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"slices"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetPlayerRetentionReportData", "获取玩家留存报表数据", "AdminInfo", getPlayerRetentionReport, getPlayerRetentionReportParams{
		ReportId:  "",
		StartTime: 0,
		EndTime:   0,
		Operator:  1,
		Page:      0,
		PageSize:  0,
		//Type:      "",
	})
	RegMsgProc("/AdminInfo/GetPlayerRetentionGameReportData", "获取玩家游戏留存报表数据", "AdminInfo", getPlayerRetentionGameReport, getPlayerRetentionGameReportParams{
		ReportId: "",
	})
}

type getPlayerRetentionReportParams struct {
	ReportId  string
	StartTime int64
	EndTime   int64
	Operator  int64
	Page      int64
	PageSize  int64
	//Type      string
}
type getPlayerRetentionGameReportParams struct {
	ReportId string
}

type getPlayerRetentionReportResults struct {
	List []*comm.PlayerRetentionReport
	All  int64
}

type getPlayerRetentionGameReportResults struct {
	List []*comm.PlayerRetentionGameReport
	All  int64
}

func getPlayerRetentionReport(ctx *Context, ps getPlayerRetentionReportParams, ret *getPlayerRetentionReportResults) (err error) {
	query := bson.M{}
	appidlist, err := GetOperatopAppID(ctx)
	if err != nil {
		return
	}
	if ps.StartTime != 0 {
		startTime := time.Unix(ps.StartTime, 0)
		endTime := time.Unix(ps.EndTime, 0)
		query["date"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	filter := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"_id": -1},
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
	if ps.ReportId != "" {
		query["_id"] = ps.ReportId
	}
	query["gameId"] = bson.M{"$exists": false}

	filter.Query = query
	ret.All, err = NewOtherDB("reports").Collection("PlayerRetentionReport").FindPage(filter, &ret.List)
	if err != nil {
		return err
	}

	return
}

func getPlayerRetentionGameReport(ctx *Context, ps getPlayerRetentionGameReportParams, ret *getPlayerRetentionGameReportResults) (err error) {

	query := bson.M{}
	query["_id"] = primitive.Regex{
		Pattern: "^" + ps.ReportId + ":",
		Options: "i",
	}

	err = NewOtherDB("reports").Collection("PlayerRetentionReport").FindAll(query, &ret.List)
	if err != nil {
		return err
	}
	ret.All = int64(len(ret.List))

	return
}
