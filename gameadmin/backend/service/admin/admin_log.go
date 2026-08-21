package main

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/duck/mongodb"
	"game/duck/ut2/httputil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/AddLog", "管理添加日志", "AdminInfo", AddLog, addLogParams{
		OperateMod:     0,
		OperateType:    0,
		OperateContent: "",
	})
	RegMsgProc("/AdminInfo/GetLog", "管理获取日志", "AdminInfo", GetLog, getLogParams{
		OperatorType: 0,
		OperateMod:   0,
		StartTime:    0,
		EndTime:      0,
		OperateType:  0,
		OperatorName: "",
		Page:         0,
		PageSize:     0,
	})
}

type addLogParams struct {
	OperateMod     int64
	OperateType    int64
	OperateContent string
}
type getLogParams struct {
	AppID          string
	OperatorType   int64
	OperateMod     int64
	StartTime      int64
	EndTime        int64
	OperateType    int64
	OperatorName   string
	Page           int64
	PageSize       int64
	OperateContent string // hml添加方便完成预存额日志中通过商户查询的功能添加
	//Type      string
}

type getLogResults struct {
	List []*comm.AdminLog
	All  int64
}

func AddLog(ctx *Context, ps addLogParams, ret *comm.Empty) (err error) {
	user, err := GetUser(ctx)
	if err != nil {
		return
	}
	coll := db.Collection2("GameAdmin", "AdminLog")
	ip := httputil.GetIPFromRequest(ctx.Request)
	one := comm.AdminLog{
		ID:             primitive.NewObjectID(),
		OperatorID:     user.Id,
		OperatorName:   user.UserName,
		OperatorType:   user.GroupId,
		OperateMod:     ps.OperateMod,
		OperateType:    ps.OperateType,
		IpAdress:       ip,
		OperateContent: ps.OperateContent,
		OperateTime:    mongodb.NowTimeStamp(),
		AppID:          user.AppID,
	}
	_, err = coll.InsertOne(context.TODO(), one)
	if err != nil {
		return
	}

	return
}

func GetLog(ctx *Context, ps getLogParams, ret *getLogResults) (err error) {
	query := bson.M{}
	appidlist, err := GetOperatopAppID(ctx)
	if err != nil {
		return
	}
	if ps.StartTime != 0 {
		startTime := time.Unix(ps.StartTime, 0)
		endTime := time.Unix(ps.EndTime, 0)
		query["operatetime"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	filter := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"_id": -1},
	}

	user, err := GetUser(ctx)
	if err != nil {
		return
	}
	if user.GroupId != 1 {
		if user.GroupId == 2 {
			appidlist = append(appidlist, user.AppID)
		}
		query["appid"] = bson.M{"$in": appidlist}
	}

	if ps.OperatorType > 0 {
		query["operatortype"] = ps.OperatorType
	}
	if ps.OperateMod > 0 {
		query["operatemod"] = ps.OperateMod
	}
	if ps.OperateType > 0 {
		query["operatetype"] = ps.OperateType
	}
	if ps.OperatorName != "" {
		query["operatorname"] = ps.OperatorName
	}

	// =====傻逼前端写的傻逼代码与后端无关=====
	// =====前端传数据时将修改的商户写入字符串对象中了只能这么写判断条件=====
	// =====2025-01-08 落款：hml====
	if ps.OperateContent != "" {
		query["operatecontent"] = bson.M{"$regex": ps.OperateContent}
	}

	filter.Query = query
	ret.All, err = NewOtherDB("GameAdmin").Collection("AdminLog").FindPage(filter, &ret.List)
	if err != nil {
		return err
	}

	return
}
