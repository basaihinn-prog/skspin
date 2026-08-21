package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/AddOperatorAdmin", "添加运营商管理员", "AdminInfo", addOperatorAdmin, addOperatorAdminParams{
		AppID:    "",
		UserName: "",
		Password: "",
	})
}

type addOperatorAdminParams struct {
	UserName string
	Password string
	AppID    string
}

func addOperatorAdmin(ctx *Context, ps addOperatorAdminParams, ret *comm.Empty) (err error) {
	user, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	if ps.UserName == "" || ps.Password == "" {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	var operator *comm.Operator
	err = CollAdminOperator.FindOne(bson.M{"Name": ps.AppID}, &operator)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	user = comm.User{}
	err = CollAdminUser.FindOne(bson.M{"Username": ps.UserName, "AppID": ps.AppID}, &user)
	if err == nil {
		return errors.New(lang.GetLang(ctx.Lang, "名字重复"))
	}

	if err != mongo.ErrNoDocuments {
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
		PermissionId:  operator.PermissionId,
	}

	err = CollAdminUser.InsertOne(user)
	return err
}
