package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/AddPermission", "添加权限组", "AdminInfo", addPermission, addPermissionParams{
		OperatorId: 0,
		Name:       "",
		Remark:     "",
		MenuIds:    []int64{1, 2},
	})
}

type addPermissionParams struct {
	OperatorId int64
	Name       string
	Remark     string
	MenuIds    []int64
}

func addPermission(ctx *Context, ps addPermissionParams, ret *comm.Empty) (err error) {
	permissionId, err := NextID(CollAdminPermission, 100)
	if err != nil {
		return err
	}
	// 游戏菜单权限   所有全选可以添加
	ps.MenuIds = append(ps.MenuIds, 1)
	newPermission := &comm.Permission{
		Id:      permissionId,
		Name:    ps.Name,
		Remark:  ps.Remark,
		MenuIds: ps.MenuIds,
	}
	var user comm.User

	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	if user.PermissionId > 3 {

		// 查询主账号
		err = CollAdminUser.FindOne(bson.M{"AppID": user.AppID}, &user)
		if err != nil {
			return err
		}

	}

	var oldPermission *comm.Permission
	err = CollAdminPermission.FindOne(bson.M{"Name": ps.Name, "OperatorId": user.Id}, &oldPermission)
	if err == nil || oldPermission != nil {
		return errors.New(lang.GetLang(ctx.Lang, "名字重复"))
	}

	newPermission.OperatorId = user.Id

	err = CollAdminPermission.InsertOne(newPermission)
	return err

	//var operator comm.Operator
	//err = CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &operator)
	//if err != nil {
	//	return err
	//}
	//newPermission.OperatorId = operator.Id
	//for i := range ps.MenuIds {
	//	if !lo.Contains(operator.MenuIds, ps.MenuIds[i]) {
	//		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	//	}
	//}

	// 原来添加健全的方法
	/*
		if user.AppID != "admin" {
			var operator comm.Operator
			err = CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &operator)
			if err != nil {
				return err
			}
			newPermission.OperatorId = operator.Id
			for i := range ps.MenuIds {
				if !lo.Contains(operator.MenuIds, ps.MenuIds[i]) {
					return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
				}
			}
		} else {
			var operator comm.Operator
			err = CollAdminOperator.FindId(ps.OperatorId, &operator)
			if err != nil {
				return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
			}
			newPermission.OperatorId = ps.OperatorId
		}
	*/

}
