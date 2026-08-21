package main

import (
	"game/comm"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/UserConfig", "用户账号配置", "AdminInfo", getUser, getUserParams{})
	RegMsgProc("/AdminInfo/UpdataUserConfig", "修改用户账号配置", "AdminInfo", updataUser, getUserParams{})
}

type getUserParams struct {
}
type respUser struct {
	List     []*comm.User
	AllCount int
}
type updataUserParams struct {
	Id           int    `json:"Id" bson:"Id" description:"ID"`
	UserName     string `json:"UserName" bson:"UserName" description:"用户名"`
	Status       int    `json:"Status" bson:"Status" description:"状态"`
	PermissionId int    `json:"PermissionId" bson:"PermissionId" description:"角色ID"`
	PassWord     string `json:"PassWord" bson:"PassWord" description:"密码"`
}

func updataUser(ctx *Context, ps updataUserParams, ret *respUser) (err error) {

	var setingInfo bson.M
	if ps.UserName == ctx.Username {
		return err
	}
	setingInfo["UserName"] = ps.UserName
	setingInfo["Status"] = ps.Status
	setingInfo["PermissionId"] = ps.PermissionId
	setingInfo["PassWord"] = ps.PassWord

	filter := bson.M{"id": ps.Id}
	update := bson.M{"$set": setingInfo}
	err = CollAdminUser.Update(filter, update)
	if err != nil {
		return
	}
	return
}

func getUser(ctx *Context, ps getGameCofigParams, ret *respUser) (err error) {

	var list []*comm.User

	err = CollAdminUser.FindAll(bson.M{}, &list)
	if err != nil {
		return err
	}

	ret.List = list
	ret.AllCount = len(list)
	return
}
