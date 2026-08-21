package main

import (
	"context"
	"game/comm"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	RegMsgProc("/AdminInfo/GetGamesLoginList", "获取IP访问查询列表", "AdminInfo", GetGamesLoginList, gameLoginListParam{
		OperatorId: 100022,
		GameID:     "Olympus",
		StartTime:  0,
		EndTime:    0,
	})
}

type gameLoginListParam struct {
	OperatorId int64
	GameID     string
	StartTime  int64
	EndTime    int64
	pageSetting
}

type gameLoginListResults struct {
	pageSetting
	List []*comm.GameLoginDetail
}

func GetGamesLoginList(ctx *Context, ps gameLoginListParam, ret *gameLoginListResults) (err error) {
	if ps.PageNumber == 0 {
		ps.PageNumber = 1
	}
	if ps.PageSize == 0 {
		ps.PageSize = 20
	}
	filter := bson.M{}
	if ps.OperatorId != 0 {
		var operator = &comm.Operator{}
		err = CollAdminOperator.FindId(ps.OperatorId, operator)
		if err != nil {
			return
		}
		filter["Pid"] = operator.AppID
	}
	if len(ps.GameID) != 0 {
		filter["GameID"] = ps.GameID
	}
	if ps.StartTime != 0 {
		filter["LoginTime"] = bson.M{"$gte": ps.StartTime}
	}
	if ps.EndTime != 0 {
		filter["LoginTime"] = bson.M{"$lte": ps.EndTime}
	}
	if ps.StartTime != 0 && ps.EndTime != 0 {
		filter["LoginTime"] = bson.M{"$gte": ps.StartTime, "$lte": ps.EndTime}
	}
	findOptions := options.Find()
	skip := (ps.PageNumber - 1) * ps.PageSize
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"LoginTime": -1})
	findOptions.SetLimit(int64(ps.PageSize))

	err = CollGameLoginDetail.FindAllOpt(filter, &ret.List, findOptions)
	if err != nil {
		return err
	}
	ret.PageNumber = ps.PageNumber
	ret.PageSize = ps.PageSize
	c, _ := CollGameLoginDetail.CountDocuments(filter)
	ret.Count = int(c)
	return
}

func removeGameLoginDetail() {
	slog.Log(context.Background(), slog.LevelInfo, " >>>>>> removeGameLoginDetail 9.30 remove old records start")
	delTime := time.Now().AddDate(0, 0, -7).Unix()
	filter := bson.M{}
	filter["LoginTime"] = bson.M{"$lte": delTime}
	CollGameLoginDetail.Coll().DeleteMany(context.TODO(), filter)
	slog.Log(context.Background(), slog.LevelInfo, " <<<<<< removeGameLoginDetail 9.30 remove old records end")
}
