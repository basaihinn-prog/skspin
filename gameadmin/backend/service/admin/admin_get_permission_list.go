package main

import (
	"game/comm"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/GetPermissionList", "获取权限组列表", "AdminInfo", getPermissionList, getPermissionListParams{
		PageIndex:  1,
		PageSize:   10,
		OperatorId: 0,
	})
	RegMsgProc("/AdminInfo/GetPermissionListAll", "获取权限组列表", "AdminInfo", getPermissionListAll, getPermissionListParams{
		PageIndex:  1,
		PageSize:   10,
		OperatorId: 0,
	})
}

type getPermissionListParams struct {
	PageIndex  int64
	PageSize   int64
	OperatorId int64
}

type getPermissionListResult struct {
	List     []*comm.Permission
	AllCount int64
}

func getPermissionList(ctx *Context, ps getPermissionListParams, ret *getPermissionListResult) (err error) {
	query := bson.M{}
	var user comm.User
	var parentUser comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	if user.PermissionId > 3 {

		err = CollAdminUser.FindOne(bson.M{"AppID": user.AppID}, &parentUser)

		query["OperatorId"] = parentUser.Id
	} else {

		query["OperatorId"] = user.Id
	}

	err, list, count := getPermissionData(ps.PageIndex, ps.PageSize, query)
	if err != nil {
		return err
	}

	ret.AllCount = count
	ret.List = list

	return err
}

func getPermissionListAll(ctx *Context, ps getPermissionListParams, ret *getPermissionListResult) (err error) {

	query := bson.M{}

	// 非Admin只查询自己所关联的权限

	err, list, count := getPermissionData(ps.PageIndex, ps.PageSize, query)
	if err != nil {
		return err
	}

	ret.AllCount = count
	ret.List = list

	return err
}

func getPermissionData(PageIndex int64, PageSize int64, query bson.M) (err error, list []*comm.Permission, AllCount int64) {

	query["_id"] = bson.M{"$ne": 1}
	filter := mongodb.FindPageOpt{
		Page:     PageIndex,
		PageSize: PageSize,
		Sort:     bson.M{"_id": -1},
		Query:    query,
	}

	count, err := CollAdminPermission.FindPage(filter, &list)

	if err != nil {
		return err, nil, 0
	}

	return nil, list, count
}
