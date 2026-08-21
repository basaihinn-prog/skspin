package main

import (
	"game/comm"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/UpdateUserGoogleCode", "是否开始google验证", "AdminInfo", updateUserGoogleCode, updateUserGoogleCodeParams{
		Id:   0,
		Open: true,
	})
}

type updateUserGoogleCodeParams struct {
	Id   int64
	Open bool
}

func updateUserGoogleCode(ctx *Context, ps updateUserGoogleCodeParams, ret *comm.Empty) (err error) {

	var currentUser *comm.User

	err = CollAdminUser.FindOne(bson.M{"_id": ps.Id}, &currentUser)
	if err != nil {
		return err
	}
	setValue := bson.M{
		"IsOpenGoogle": ps.Open,
	}

	secret := NewGoogleAuth().GetSecret()
	qrcode := NewGoogleAuth().GetQrcodeUrl(currentUser.UserName, secret)

	if currentUser.GoogleCode == "" && currentUser.Qrcode == "" {

		setValue["GoogleCode"] = secret
		setValue["Qrcode"] = qrcode
	}

	err = CollAdminUser.UpdateId(ps.Id, bson.M{
		"$set": setValue,
	})
	return err
}
