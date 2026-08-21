package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/DeletePermission", "删除权限组", "AdminInfo", deletePermission, deletePermissionParams{
		PermissionId: 1,
	})
}

type deletePermissionParams struct {
	PermissionId int64
}

func deletePermission(ctx *Context, ps deletePermissionParams, ret *comm.Empty) (err error) {
	var permission comm.Permission
	err = CollAdminPermission.FindId(ps.PermissionId, &permission)
	if err != nil {
		return err
	}

	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	if user.AppID != "admin" {
		if !user.OperatorAdmin {
			return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
		}
		var operator comm.Operator
		err = CollAdminOperator.FindId(permission.OperatorId, &operator)
		if err != nil {
			return err
		}
		if operator.AppID != user.AppID {
			return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
		}
	}

	CollAdminPermission.DeleteId(ps.PermissionId)
	return
}
