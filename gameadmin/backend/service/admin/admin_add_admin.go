package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/AddAdmin", "添加超级管理员", "AdminInfo", addAdmin, addAdminParams{
		UserName: "",
		Password: "",
	})
	RegMsgProc("/AdminInfo/EditAdmin", "修改超级管理员", "AdminInfo", editAdmin, editAdminParams{
		PermissionId: 0,
	})

	RegMsgProc("/AdminInfo/editConfig", "修改管理员通用配置", "AdminInfo", editConfig, ConfigParam{
		HalfMaintenance:   0,
		EntireMaintenance: 0,
	})

	RegMsgProc("/AdminInfo/getConfig", "获取管理员通用配置", "AdminInfo", getConfig, ConfigParam{
		HalfMaintenance:   0,
		EntireMaintenance: 0,
	})

	RegMsgProc("/AdminInfo/getMaintenanceLog", "获取维护记录信息", "AdminInfo", getMaintenanceLog, MaintenanceLogParams{
		PageIndex: 1,
		PageSize:  20,
	})
}

type addAdminParams struct {
	UserName     string
	Password     string
	PermissionId int64
}

type editAdminParams struct {
	Id           int64
	PermissionId int64
}

type ConfigParam struct {
	HalfMaintenance   int64
	EntireMaintenance int64
	StartTime         int64
	EndTime           int64
}

type MaintenanceLogParams struct {
	PageIndex int64 `json:"PageIndex"`
	PageSize  int64 `json:"PageSize"`
}

type MaintenanceLogRep struct {
	List  []*MaintenanceLog
	Count int64
}

type MaintenanceLog struct {
	ID                primitive.ObjectID `json:"ID" bson:"_id,omitempty"`
	MaintenanceMod    int                `json:"MaintenanceMod" bson:"MaintenanceMod"`
	MaintenanceStatus int                `json:"MaintenanceStatus" bson:"MaintenanceStatus"`
	StartTime         int64              `json:"StartTime" bson:"StarTime"`
	EndTime           int64              `json:"EndTime" bson:"EndTime"`
	RealEndTime       int64              `json:"RealEndTime" bson:"RealEndTime"`
	OperatorId        string             `json:"OperatorId" bson:"OperatorId"`
	OperatorTime      int64              `json:"OperatorTime" bson:"OperatorTime"`
}

func addAdmin(ctx *Context, ps addAdminParams, ret *comm.Empty) (err error) {

	//todo:逻辑需要改  账号有主次分
	if ps.UserName == "" || ps.Password == "" {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	user := comm.User{}
	err = CollAdminUser.FindOne(bson.M{"Username": ps.UserName}, &user)
	if err == nil {
		return errors.New(lang.GetLang(ctx.Lang, "名字重复"))
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	secret := NewGoogleAuth().GetSecret()
	qrcode := NewGoogleAuth().GetQrcodeUrl(ps.UserName, secret)
	id, err := NextID(CollAdminUser, 1000)
	user, _ = IsAdminUser(ctx)
	if err != nil {
		return err
	}
	user = comm.User{
		Id:            id,
		AppID:         user.AppID,
		UserName:      ps.UserName,
		PassWord:      ps.Password,
		GoogleCode:    secret,
		Qrcode:        qrcode,
		Avatar:        "",
		IsOpenGoogle:  false,
		GroupId:       user.GroupId,
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

func editAdmin(ctx *Context, ps editAdminParams, ret *comm.Empty) (err error) {

	query := bson.M{}
	if ps.PermissionId != 0 {
		query["PermissionId"] = ps.PermissionId
	}

	err = CollAdminUser.Update(bson.M{"_id": ps.Id}, bson.M{"$set": query})

	return nil
}

func editConfig(ctx *Context, ps ConfigParam, ret *comm.Empty) (err error) {

	user, _ := IsAdminUser(ctx)

	if user.GroupId > 1 {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	if ps.HalfMaintenance == 1 && ps.EntireMaintenance == 1 {
		return errors.New(lang.GetLang(ctx.Lang, "半维护状态与全维护状态不可同时开启"))
	}

	var logStartTime int64
	var logMaintenanceMod int

	var dbData comm.SystemConfig

	dbData.ID = user.AppID

	CollConfig := db.Collection2("GameAdmin", "SystemConfig")
	err = CollConfig.FindOne(context.TODO(), bson.M{"AppId": user.AppID}).Decode(&dbData)

	logStartTime = dbData.StartTime
	logMaintenanceMod = dbData.Maintenance

	status := 0

	// 维护状态为 半维护状态
	if ps.HalfMaintenance == 1 {
		status = 1
	}
	// 维护状态为 全维护状态
	if ps.EntireMaintenance == 1 {
		status = 2
	}

	dbData.Maintenance = status

	// 如果维护关闭的情况下，则清空时间
	if status == 0 {
		dbData.StartTime = 0
		dbData.EndTime = 0
	}

	dbData.StartTime = ps.StartTime
	dbData.EndTime = ps.EndTime

	if err != nil {

		_, err = CollConfig.InsertOne(context.TODO(), dbData)
	} else {
		_, err = CollConfig.UpdateOne(context.TODO(), bson.M{"AppId": user.AppID}, bson.M{"$set": dbData})

	}

	if err != nil {
		fmt.Println(err.Error())
		return errors.New(lang.GetLang(ctx.Lang, "修改错误："+err.Error()))
	}

	err = addLog(logStartTime, logMaintenanceMod, dbData, ctx)
	if err != nil {
		return err
	}

	return nil

}

func getConfig(ctx *Context, ps ConfigParam, ret *comm.SystemConfig) (err error) {

	user, _ := IsAdminUser(ctx)

	if user.GroupId > 1 {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	CollConfig := db.Collection2("GameAdmin", "SystemConfig")
	err = CollConfig.FindOne(context.TODO(), bson.M{"AppId": user.AppID}).Decode(ret)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return errors.New(lang.GetLang(ctx.Lang, "获取配置失败"))
		}
		return err
	}

	return err
}

func addLog(startTime int64, logMaintenanceMod int, SystemConfig comm.SystemConfig, ctx *Context) (err error) {

	CollConfig := db.Collection2("GameAdmin", "MaintenanceLog")

	var log MaintenanceLog

	log.MaintenanceMod = logMaintenanceMod

	log.StartTime = SystemConfig.StartTime
	log.EndTime = SystemConfig.EndTime
	log.OperatorTime = time.Now().Unix()
	log.OperatorId = ctx.Username

	if SystemConfig.Maintenance != 0 {
		log.MaintenanceStatus = 0
	} else {
		log.ID = primitive.NewObjectID()
		log.RealEndTime = time.Now().Unix()
		log.StartTime = startTime
		log.MaintenanceStatus = 1

	}
	_, err = CollConfig.InsertOne(context.TODO(), log)

	return err

}

func getMaintenanceLog(ctx *Context, ps MaintenanceLogParams, ret *MaintenanceLogRep) (err error) {

	CollConfig := NewOtherDB("GameAdmin").Collection("MaintenanceLog")

	filter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     bson.M{"OperatorTime": -1},
		Query:    bson.M{"MaintenanceStatus": 1},
	}
	ret.Count, err = CollConfig.FindPage(filter, &ret.List)

	return err
}
