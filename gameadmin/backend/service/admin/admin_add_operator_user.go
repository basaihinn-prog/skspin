package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/AddOperatorUser", "添加运营商用户", "AdminInfo", addOperatorUser, addOperatorUserParams{
		UserName:     "",
		Password:     "",
		AppID:        "",
		PermissionId: 0,
	})
}

type addOperatorUserParams struct {
	UserName     string
	Password     string
	AppID        string
	PermissionId int64
}

func addOperatorUser(ctx *Context, ps addOperatorUserParams, ret *comm.Empty) (err error) {
	var admin comm.User
	err = CollAdminUser.FindOne(bson.M{"_id": ctx.PID}, &admin)
	if err != nil {
		return err
	}
	_, ok := IsAdminUser(ctx)
	if (!admin.OperatorAdmin || ps.AppID != admin.AppID) && !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ps.UserName, "AppID": ps.AppID}, &user)
	if err == nil {
		return errors.New(lang.GetLang(ctx.Lang, "名字重复"))
	}

	var operator *comm.Operator
	err = CollAdminOperator.FindOne(bson.M{"Name": ps.AppID}, &operator)
	if err != nil {
		return err
	}

	var permission *comm.Permission
	err = CollAdminPermission.FindOne(bson.M{
		"_id":        ps.PermissionId,
		"OperatorId": operator.Id,
	}, &permission)
	if err != nil {
		return err
	}

	secret := NewGoogleAuth().GetSecret()
	qrcode := NewGoogleAuth().GetQrcodeUrl(ps.UserName, secret)
	id, err := NextID(CollAdminUser, 1000)
	if err != nil {
		return err
	}
	user = comm.User{
		Id:            id,
		AppID:         ps.AppID,
		UserName:      ps.UserName,
		PassWord:      ps.Password,
		GoogleCode:    secret,
		Qrcode:        qrcode,
		Avatar:        "",
		IsOpenGoogle:  false,
		Status:        comm.Operator_Status_Normal,
		CreateAt:      mongodb.NewTimeStamp(time.Now()),
		LoginAt:       mongodb.NewTimeStamp(time.Now()),
		Token:         "",
		TokenExpireAt: nil,
		OperatorAdmin: false,
		PermissionId:  ps.PermissionId,
	}
	err = CollAdminUser.InsertOne(user)
	return err
}
