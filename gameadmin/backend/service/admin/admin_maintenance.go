package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/add_whitList", "添加白名单", "maintenance", addMaintenanceWhiteList, WhitRequest{
		UserIp: "",
		Remark: "",
	})
	RegMsgProc("/AdminInfo/get_whitList", "获取白名单", "maintenance", getMaintenanceWhiteList, WhitRequest{
		UserIp:   "",
		Remark:   "",
		PageSize: 20,
		Page:     1,
	})
	RegMsgProc("/AdminInfo/del_system_maintenance", "删除白名单", "maintenance", deleteMaintenance, DeleteWhitRequest{
		Ids: "",
	})

	// 注册MQ仅供内部调用
	mux.RegRpc("/AdminInfo/Interior/whiteList_exists", "判断当前维护状态是否能通过验证白名单", "maintenance", getPrivateWhiteList, ExistsWhitRequest{
		Ip: "",
	})
}

type ExistsWhitResponse struct {
	Allow bool `json:"allow"`
}
type ExistsWhitRequest struct {
	Ip    string
	AppId string
}
type WhitListResponse struct {
	List  []WhitData
	Count int64
}
type DeleteWhitRequest struct {
	Ids string
}

type WhitRequest struct {
	UserIp   string
	Remark   string
	PageSize int64
	Page     int64
}

type WhitData struct {
	Id           primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	UserIp       string             `json:"UserIp" bson:"UserIp"`
	Remark       string             `json:"Remark" bson:"Remark"`
	OperatorTime time.Time          `json:"OperatorTime" bson:"OperatorTime"`
	OperatorName string             `json:"OperatorName" bson:"OperatorName"`
}

func addMaintenanceWhiteList(ctx *Context, ps WhitRequest, ret *comm.Empty) (err error) {
	user, _ := IsAdminUser(ctx)
	if user.AppID != "admin" {
		return fmt.Errorf("权限不足")
	}
	var whit WhitData
	if ps.UserIp == "" {
		matchString, err := regexp.MatchString("/^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})\\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})$/", ps.UserIp)

		if err != nil && !matchString {

			return fmt.Errorf(lang.GetLang(ctx.Lang, "请输入白名单IP"))
		}

	}

	whit.UserIp = ps.UserIp
	whit.Remark = ps.Remark
	whit.OperatorTime = time.Now()
	whit.OperatorName = ctx.Username

	coll := NewOtherDB("GameAdmin").Collection("WhitUser")

	exists, err := coll.Exists(bson.M{"UserIp": ps.UserIp})
	if err != nil {
		return err
	}

	if exists {
		return errors.New(lang.GetLang(ctx.Lang, "该IP已存在"))
	}

	err = coll.InsertOne(whit)

	return err
}

func getMaintenanceWhiteList(ctx *Context, ps WhitRequest, ret *WhitListResponse) (err error) {
	user, _ := IsAdminUser(ctx)
	if user.AppID != "admin" {
		return fmt.Errorf("权限不足")
	}
	err = getWhitUsers(ps, ret)
	return err
}
func deleteMaintenance(ctx *Context, ps DeleteWhitRequest, ret *comm.Empty) (err error) {
	user, _ := IsAdminUser(ctx)
	if user.AppID != "admin" {
		return fmt.Errorf("权限不足")
	}
	var split []primitive.ObjectID
	for i := range strings.Split(ps.Ids, ",") {

		hex, err := primitive.ObjectIDFromHex(strings.Split(ps.Ids, ",")[i])
		if err != nil {
			break
		}
		split = append(split, hex)
	}

	coll := db.Collection2("GameAdmin", "WhitUser")
	_, err = coll.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": split}})

	return err
}

func getPrivateWhiteList(ps ExistsWhitRequest, ret *ExistsWhitResponse) (err error) {

	systemConfigColl := NewOtherDB("GameAdmin").Collection("SystemConfig")
	whitUserColl := NewOtherDB("GameAdmin").Collection("WhitUser")

	if ps.AppId == "" {
		ps.AppId = "admin"
	}

	var systemConfig comm.SystemConfig
	err = systemConfigColl.FindOne(bson.M{"AppId": ps.AppId}, &systemConfig)

	// 报错或者未维护
	if err != nil || systemConfig.Maintenance == 0 {
		ret.Allow = true
		return
	}

	// 全部维护
	if systemConfig.Maintenance == 2 {
		return
	}

	// 半维护状态需要验证IP是否填写
	if ps.Ip == "" {
		return
	}

	exists, err := whitUserColl.Exists(bson.M{
		"UserIp": ps.Ip,
	})
	// 当报错或不存在时，不允许通过
	if err != nil || !exists {
		return
	}
	ret.Allow = true
	return err
}

func getWhitUsers(ps WhitRequest, ret *WhitListResponse) (err error) {

	filter := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"OperatorTime": -1},
		Query:    bson.M{},
	}

	coll := NewOtherDB("GameAdmin").Collection("WhitUser")

	ret.Count, err = coll.FindPage(filter, &ret.List)

	return err
}
