package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/GetAdminInfo", "获取用户详情", "AdminInfo", getAdminInfo, getAdminInfoParams{
		PageIndex: 1,
		PageSize:  10,
	})
}

type getAdminInfoParams struct {
	PageIndex int64
	PageSize  int64
}

type getAdminInfoResult struct {
	ID                  int64
	Username            string
	PermissionId        int64
	MenuList            []*comm.PerMenu
	OperatorId          int64
	BusinessesId        int64
	Businesses          *comm.Operator_V2
	OperatorName        string
	AppID               string
	IsDefaultPermission bool
	GroupId             int64
	TokenExpireAt       *mongodb.TimeStamp
}

func getAdminInfo(ctx *Context, ps getPlayerListParams, ret *getAdminInfoResult) (err error) {
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	ret.GroupId = user.GroupId

	if user.PermissionId == 0 {
		return errors.New(lang.GetLang(ctx.Lang, "账户未设置权限"))
	}

	var businesses comm.Operator_V2

	CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &businesses)

	ret.BusinessesId = businesses.Id
	ret.Businesses = &businesses

	var permission comm.Permission
	err = CollAdminPermission.FindId(user.PermissionId, &permission)
	if err != nil {
		return err
	}
	query := bson.M{}
	if permission.Id != 1 {
		permission.MenuIds = append(permission.MenuIds, 1)
		query["_id"] = bson.M{"$in": permission.MenuIds}
	}

	menuList := make([]*comm.PerMenu, 0)
	filter := mongodb.FindPageOpt{
		Page:     1,
		PageSize: 1000,
		Sort:     bson.M{"Sort": 1}, //可以为空
		Query:    query,
	}
	_, err = CollAdminMenu.FindPage(filter, &menuList)
	if err != nil {
		return err
	}

	ret.ID = ctx.PID
	ret.Username = ctx.Username
	ret.OperatorId = user.Id
	ret.PermissionId = permission.Id
	ret.AppID = user.AppID
	ret.TokenExpireAt = user.TokenExpireAt

	deleteIds := []int64{31}
	for i := len(menuList) - 1; i >= 0; i-- {
		if lo.Contains(deleteIds, menuList[i].ID) {
			menuList = append(menuList[:i], menuList[i+1:]...)
		}
	}
	ret.MenuList = menuList
	return nil
}
