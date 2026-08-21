package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/mongodb"
	"game/duck/notification"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/addNotification", "添加代理商公告", "AdminInfo", addNotification, &notificationParams{})
	RegMsgProc("/AdminInfo/deleteNotification", "删除代理商公告", "AdminInfo", deleteNotification, &notificationParams{})
	RegMsgProc("/AdminInfo/editNotification", "修改代理商公告", "AdminInfo", editNotification, &notificationParams{})
	RegMsgProc("/AdminInfo/getNotification", "查询代理商公告", "AdminInfo", getNotification, &notificationParams{})
	RegMsgProc("/AdminInfo/getEffectiveNotify", "查询有效的代理商公告", "AdminInfo", getEffectiveNotify, &notificationParams{})

}

type notificationParams struct {
	Id              string `json:"_id"`
	NotifyTitle     string `json:"notifyTitle"`
	NotifyContext   string `json:"notifyContext"`
	NotifyType      string `json:"notifyType"`
	OnlineLimitTime string `json:"onlineLimitTime"`
	LowerLimitTime  string `json:"lowerLimitTime"`
	LanguageCode    string `json:"languageCode"`
	Sort            int32  `json:"Sort"`
	LanguageContext string `json:"languageContext"`
	CreatedTime     string `json:"created_time"`
	Page            int64
	PageSize        int64
}

type uploadParams struct {
	FileName string `json:"file"`
}

type notificationResults struct {
	Count         int64
	Notifications []map[string]interface{}
}

func addNotification(ctx *Context, ps notificationParams, ret *notificationResults) (err error) {

	currentUser, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	json, err := notification.UnmarshalJSON([]byte(ps.LanguageContext))
	if err != nil {
		return err
	}

	doc := notification.Notification{
		NotifyTitle:     ps.NotifyTitle,
		NotifyType:      ps.NotifyType,
		OnlineLimitTime: ps.OnlineLimitTime,
		LowerLimitTime:  ps.LowerLimitTime,
		Sort:            1,
		CreatedTime:     time.Now().Format("2006-01-02 15:04:05"),
		Status:          "1",
		Operator:        currentUser.Id,
	}

	// 返回数据库字段
	baseData, err := db.Collection2("GameAdmin", "notify").InsertOne(context.TODO(), doc)
	// 通过存入数据库的数据进行数据修改
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "添加失败"))
	}
	_, err = db.Collection2("GameAdmin", "notify").UpdateOne(context.TODO(), bson.M{"_id": baseData.InsertedID}, bson.D{{"$set", json}})

	return err
}

func deleteNotification(ctx *Context, ps notificationParams, ret *notificationResults) (err error) {

	if _, ok := IsAdminUser(ctx); !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	hex, err := mongodb.ObjectIDFromHex(ps.Id)
	if err != nil {
		return err
	}
	_, err = db.Collection2("GameAdmin", "notify").DeleteOne(context.TODO(), bson.M{"_id": hex})

	return err
}

func editNotification(ctx *Context, ps notificationParams, ret *notificationResults) (err error) {

	currentUser, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	json, err := notification.UnmarshalJSON([]byte(ps.LanguageContext))
	if err != nil {
		return err
	}

	doc := notification.Notification{

		NotifyTitle:     ps.NotifyTitle,
		NotifyType:      ps.NotifyType,
		OnlineLimitTime: ps.OnlineLimitTime,
		LowerLimitTime:  ps.LowerLimitTime,
		Sort:            1,
		Status:          "1",
		CreatedTime:     ps.CreatedTime,
		Operator:        currentUser.Id,
	}
	hex, err := mongodb.ObjectIDFromHex(ps.Id)
	if err != nil {
		return err
	}
	// 返回数据库字段
	_, err = db.Collection2("GameAdmin", "notify").UpdateOne(context.TODO(), bson.M{"_id": hex}, bson.D{{"$set", doc}})
	// 通过存入数据库的数据进行数据修改
	if err != nil {
		fmt.Println(err)
		return errors.New(lang.GetLang(ctx.Lang, "修改失败"))
	}
	_, err = db.Collection2("GameAdmin", "notify").UpdateOne(context.TODO(), bson.M{"_id": hex}, bson.D{{"$set", json}})

	return err
}

func getNotification(ctx *Context, ps notificationParams, ret *notificationResults) (err error) {

	opt := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"sort": -1},
		Query:    bson.M{},
	}

	ret.Count, err = NewOtherDB("GameAdmin").Collection("notify").FindPage(opt, &ret.Notifications)

	return err
}
func getEffectiveNotify(ctx *Context, ps notificationParams, ret *notificationResults) (err error) {

	opt := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"sort": -1},
		Query:    bson.M{"notifyStatus": "1"},
	}

	ret.Count, err = NewOtherDB("GameAdmin").Collection("notify").FindPage(opt, &ret.Notifications)

	return err
}
