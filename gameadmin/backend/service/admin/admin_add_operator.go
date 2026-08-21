package main

import (
	"errors"
	"fmt"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	RegMsgProc("/AdminInfo/AddOperator", "添加运营商", "AdminInfo", addOperator, addOperatorParams{
		Name:  "",
		AppID: "",
	})
}

type addOperatorParams struct {
	Name           string
	AppID          string
	MenuIds        []int64
	Remark         string
	AppSecret      string
	Address        string
	WhiteIps       []string
	CurrencyKey    string
	WalletMode     int
	ExcluedGameId  primitive.ObjectID
	ExcluedGameIds []string
	Robot          string //telegram 机器人
	ChatID         string //telegram 聊天群id
}

func addOperator(ctx *Context, ps addOperatorParams, ret *comm.Empty) (err error) {
	if ps.Name == "admin" || ps.AppID == "admin" {
		return errors.New(lang.GetLang(ctx.Lang, "不能添加admin商户"))
	}
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	var operator *comm.Operator

	if ps.Name == "" || ps.AppID == "" {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	err = CollAdminOperator.FindOne(bson.M{"Name": ps.Name}, &operator)
	if err == nil {
		return errors.New(lang.GetLang(ctx.Lang, "名字重复"))
	}

	if /*err != nil &&*/ err != mongo.ErrNoDocuments {
		return err
	}

	id, err := NextID(CollAdminOperator, 100000)
	if err != nil {
		return err
	}
	permissionId, err := NextID(CollAdminPermission, 100)
	if err != nil {
		return nil
	}

	ps.MenuIds = NormalMenu(ps.MenuIds)

	//注册运营商
	operator = &comm.Operator{
		Id:             id,
		Name:           ps.Name,
		AppID:          ps.AppID,
		CreateTime:     mongodb.NewTimeStamp(time.Now()),
		Status:         comm.Operator_Status_Normal,
		MenuIds:        ps.MenuIds,
		PermissionId:   permissionId,
		AppSecret:      uuid.NewString(),
		Address:        ps.Address,
		WhiteIps:       ps.WhiteIps,
		WalletMode:     ps.WalletMode,
		CurrencyKey:    ps.CurrencyKey,
		ExcluedGameId:  ps.ExcluedGameId,
		ExcluedGameIds: ps.ExcluedGameIds,
		CurrentRtp:     95,
		Robot:          ps.Robot,
		ChatID:         ps.ChatID,
	}
	err = CollAdminOperator.InsertOne(operator)
	if err != nil {
		return err
	}
	//添加权限
	timeStamp := mongodb.NewTimeStamp(time.Now())
	permission := &comm.Permission{
		Id:         permissionId,
		OperatorId: id,
		Name:       fmt.Sprintf("%s%s", ps.Name, lang.Get(ctx.Lang, "默认权限组")),
		Remark:     ps.Remark,
		MenuIds:    ps.MenuIds,
		CreateTime: timeStamp,
		UpdateTime: timeStamp,
	}
	err = CollAdminPermission.InsertOne(permission)
	if err != nil {
		return err
	}

	//注册运营商管理员
	userName := fmt.Sprintf("admin_%s", ps.Name)
	secret := NewGoogleAuth().GetSecret()
	qrcode := NewGoogleAuth().GetQrcodeUrl(userName, secret)
	var operatorAdmin comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": userName, "AppID": ps.AppID}, operatorAdmin)
	if err == nil {
		//找到了就不用添加了
		return
	}
	userId, err := NextID(CollAdminUser, 1000)
	if err != nil {
		return err
	}
	operatorAdmin = comm.User{
		Id:           userId,
		AppID:        ps.AppID,
		UserName:     userName,
		PassWord:     comm.RandPassword(),
		GoogleCode:   secret,
		Qrcode:       qrcode,
		Avatar:       "",
		IsOpenGoogle: false,
		Status:       comm.Operator_Status_Normal,
		CreateAt:     mongodb.NewTimeStamp(time.Now()),
		LoginAt:      mongodb.NewTimeStamp(time.Now()),
		//GroupId:       0,
		Token:         "",
		TokenExpireAt: nil,
		OperatorAdmin: true,
		PermissionId:  permissionId,
	}
	if err = CollAdminUser.InsertOne(&operatorAdmin); err != nil {
		return err
	}

	return nil
}

func addFakeOPerator() (err error) {

	var operator comm.Operator
	err = CollAdminOperator.FindOne(bson.M{"Name": "faketrans"}, &operator)
	if err == nil || err != mongo.ErrNoDocuments {
		return
	}

	id, err := NextID(CollAdminOperator, 100000)
	if err != nil {
		return err
	}

	// menuIds := NormalMenu(ps.MenuIds)

	//注册运营商
	operator = comm.Operator{
		Id:         id,
		Name:       "faketrans",
		AppID:      "faketrans",
		CreateTime: mongodb.NewTimeStamp(time.Now()),
		Status:     comm.Operator_Status_Normal,
		AppSecret:  "b6337af9-a91a-4085-b1f2-466923470735",
		WalletMode: comm.WalletModeTransfer,
	}
	err = CollAdminOperator.InsertOne(operator)
	if err != nil {
		return err
	}

	return
}
