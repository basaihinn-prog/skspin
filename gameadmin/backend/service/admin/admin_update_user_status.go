package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/UpdateUserStatus", "修改用户状态", "AdminInfo", updateUserStatus, updateUserStatusParams{
		Id:     0,
		Status: 0,
	})
}

type updateUserStatusParams struct {
	Id     int64
	Status int
}

func updateUserStatus(ctx *Context, ps updateUserStatusParams, ret *comm.Empty) (err error) {
	var user *comm.User
	err = CollAdminUser.FindId(ctx.PID, &user)
	if err != nil {
		return err
	}
	var animUser *comm.User
	err = CollAdminUser.FindId(ps.Id, &animUser)
	if err != nil {
		return err
	}
	if user.AppID != "admin" { //不是超级管理员
		if !user.OperatorAdmin { //不是运营商管理员
			return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
		}
		if animUser.OperatorAdmin { //非超级管理员不能修改管理员状态
			return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
		}
		if user.AppID != animUser.AppID {
			return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
		}
	}
	//if animUser.AppID != "admin" {
	//
	//	var user *comm.User
	//	err = CollAdminUser.FindId(ctx.PID, &user)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if !user.OperatorAdmin {
	//
	//	}
	//
	//	if animUser.AppID != user.AppID {
	//		return errors.New(lang.Get(ctx.Lang, "权限不足"))
	//	}
	//}
	err = CollAdminUser.UpdateId(ps.Id, bson.M{
		"$set": bson.M{
			"Status": ps.Status,
		},
	})
	if ps.Status == 1 {
		tokenMap.Delete(animUser.Token)
	}

	return err
}
