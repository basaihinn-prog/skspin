package main

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/ResetPassword", "获取权限组列表", "AdminInfo", resetPassword, ResetPasswordParams{})
}

type ResetPasswordParams struct {
	ID    int64  `json:"ID"`
	AppID string `json:"AppID"`
}

func resetPassword(ctx *Context, ps ResetPasswordParams, ret *getPermissionListResult) (err error) {
	//user, err := GetUser(ctx)
	//if err != nil {
	//	return
	//}
	//
	//oper := comm.Operator_V2{}
	//if ps.AppID != "admin" {
	//	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &oper)
	//	if err != nil {
	//		return
	//	}
	//}
	//if oper != nil {
	//	if user.UserName != oper.UserName {
	//		return errors.New("请使用主账号登录")
	//	}
	//}
	filter := bson.M{"_id": ps.ID}
	update := bson.M{"$set": bson.M{"Password": "Aa123456"}}

	err = CollAdminUser.UpdateOne(filter, update)
	if err != nil {
		return errors.New("修改失败")
	}
	return

}
