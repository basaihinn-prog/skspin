package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/lazy"
	"game/duck/logger"
	"game/duck/mongodb"
	"sync"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB = mongodb.NewDB("GameAdmin")
var MenuList []*comm.PerMenu
var MenuListLock sync.Mutex

// 后台
var (
	CollAdminUser          = DB.Collection("AdminUser")          // 管理员表
	CollAdminConfigFileLog = DB.Collection("AdminConfigFileLog") // 后台记录一些提交的配置缓存数据
	//CollAdminAuth          = DB.Collection("AdminAuth")          // 权限组表(弃用)
	CollAdminPermission = DB.Collection("AdminPermission") // 权限组表
	CollAdminMenu       = DB.Collection("AdminMenu")       // 菜单表
	CollAdminOperator   = DB.Collection("AdminOperator")   // 运营商
	CollStore           = DB.Collection("Store")
	CollGameLoginDetail = DB.Collection("GameLoginDetail") //游戏登录日志
	//CollAdminOperatorV2 = DB.Collection("AdminOperator")   // 运营商  测试业务逻辑用
	CollGameConfig       = DB.Collection("GameConfig")
	CollPlayerRTPControl = DB.Collection("PlayerRTPControl")
	CollSystemConfig     = DB.Collection("SystemConfig")
)

const (
	key_idCounter      = 100001
	key_analysisCursor = "analysisCursor"
)

func NewOtherDB(dbName string) *mongodb.DB {
	return &mongodb.DB{Name: dbName, Client: DB.Client}
}

func NextID(targetColl *mongodb.Collection, min int64) (int64, error) {

	op := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	mp := map[string]int64{}
	err := CollStore.Coll().FindOneAndUpdate(context.TODO(),
		bson.M{
			"_id": key_idCounter,
		},
		bson.M{
			"$inc": bson.M{
				targetColl.Name: 1,
			},
		}, op).Decode(&mp)

	if err != nil {
		return 0, err
	}

	return mp[targetColl.Name] + min, err
}

// func ConnectMongoByConfig() {
// 	mongoAddr := lazy.GetAddr(lazy.ServiceName + ".mongo")
// 	if len(mongoAddr) == 0 {
// 		mongoAddr = "127.0.0.1:27017"
// 	}

// 	logger.Info("连接mongodb", mongoAddr)

// 	err := DB.Connect("mongodb://" + mongoAddr)

// 	if err != nil {
// 		logger.Info("连接mongodb失败", err)
// 	} else {
// 		logger.Info("连接mongodb成功")
// 	}
// }

func ConnectMongoByConfig() {
	mongoAddr := lo.Must(lazy.RouteFile.Get("mongo"))

	logger.Info("连接mongodb", mongoAddr)

	err := DB.Connect(mongoAddr)

	if err != nil {
		logger.Info("连接mongodb失败", err)
	} else {
		logger.Info("连接mongodb成功")
		db.SetupClient(DB.Client)
	}
}

//func initBase() {
//	// 创建一个基础的管理员
//	adminInfo := adminpb.DBAdminer{}
//	err := CollAdminUser.FindOne(bson.M{"Username": "admin"}, &adminInfo)
//	if err == mongo.ErrNoDocuments {
//		base := adminpb.DBAdminer{
//			ID:       int64(ut2.RandomInt(10000, 99999)),
//			Username: "admin",
//			Password: "123456",
//			GroupId:  1,
//		}
//		if err := CollAdminUser.InsertOne(&base); err != nil {
//			logger.Info(err)
//		}
//	}
//
//	group := adminpb.DBAuth{}
//	err = CollAdminAuth.FindId(1, &group)
//	if err == mongo.ErrNoDocuments {
//		DBAuth := adminpb.DBAuth{
//			ID:         1,
//			Name:       "超级管理",
//			Remark:     "超级管理组",
//			CreateTime: mongodb.NewTimeStamp(time.Now()),
//		}
//		if err := CollAdminAuth.InsertOne(&DBAuth); err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	count, err := CollAdminMenu.CountDocuments(bson.M{})
//	if err != nil {
//		logger.Err(err)
//		return
//	}
//	if count != 0 {
//		return
//	}
//
//	list := []*adminpb.DBMenu{}
//	json.Unmarshal([]byte(baseMenu()), &list)
//
//	CollAdminMenu.InsertMany(lo.ToAnySlice(list))
//}

func initBase1() {
	// 创建一个超级管理员
	adminInfo := comm.User{}
	err := CollAdminUser.FindOne(bson.M{"Username": "admin"}, &adminInfo)
	if err == mongo.ErrNoDocuments {
		secret := NewGoogleAuth().GetSecret()
		qrcode := NewGoogleAuth().GetQrcodeUrl("admin", secret)

		base := comm.User{
			Id:           1,
			UserName:     "admin",
			PassWord:     comm.RandPassword(),
			GroupId:      1,
			PermissionId: 1,
			AppID:        "admin",
			GoogleCode:   secret,
			Qrcode:       qrcode,
			IsOpenGoogle: false,
		}
		if err := CollAdminUser.InsertOne(&base); err != nil {
			logger.Info(err)
		}
	}

	//group := adminpb.DBAuth{}
	//err = CollAdminAuth.FindId(1, &group)
	//if err == mongo.ErrNoDocuments {
	//	DBAuth := adminpb.DBAuth{
	//		ID:         1,
	//		Name:       "超级管理",
	//		Remark:     "超级管理组",
	//		CreateTime: mongodb.NewTimeStamp(time.Now()),
	//	}
	//	if err := CollAdminAuth.InsertOne(&DBAuth); err != nil {
	//		fmt.Println(err)
	//	}
	//}

	list := []*comm.PerMenu{}
	json.Unmarshal([]byte(baseMenu()), &list)

	//添加权限组
	var permission *comm.Permission
	err = CollAdminPermission.FindId(1, &permission)
	if err == mongo.ErrNoDocuments {
		menuIds := make([]int64, 0, len(list))
		for i := range list {
			menuIds = append(menuIds, list[i].ID)
		}
		//添加管理员权限组
		permission := &comm.Permission{
			Id:         1,
			OperatorId: 0,
			Name:       "กลุ่มสิทธิ์ admin",
			Remark:     "",
			MenuIds:    menuIds,
		}
		err = CollAdminPermission.InsertOne(permission)
		if err != nil {
			fmt.Println(err)
		}
	}

	//添加菜单列表
	count, err := CollAdminMenu.CountDocuments(bson.M{})
	if err != nil {
		logger.Err(err)
		return
	}
	if count != 0 {
		return
	}

	CollAdminMenu.InsertMany(lo.ToAnySlice(list))

}

func UpdateMenList() {
	MenuListLock.Lock()
	defer MenuListLock.Unlock()

	CollAdminMenu.FindAll(bson.M{}, &MenuList)
}

func baseMenu() string {
	return `[
		{
			"ID": 2,
			"Title": "产品列表",
			"Pid": 21,
			"Url": "/appList",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 3,
			"Title": "玩家管理",
			"Pid": 0,
			"Url": "/player",
			"Icon": "User",
			"Sort": 0
		},
		{
			"ID": 4,
			"Title": "玩家信息",
			"Pid": 3,
			"Url": "/playerInfo",
			"Icon": "Female",
			"Sort": 0
		},
		{
			"ID": 6,
			"Title": "游戏运行时",
			"Pid": 0,
			"Url": "/runtime",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 8,
			"Title": "数据统计",
			"Pid": 0,
			"Url": "/data",
			"Icon": "DocumentCopy",
			"Sort": 0
		},
		{
			"ID": 9,
			"Title": "玩家数据",
			"Pid": 8,
			"Url": "/playerAnalysis",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 10,
			"Title": "每日游戏详情",
			"Pid": 8,
			"Url": "/gameAnalysis",
			"Icon": "Histogram",
			"Sort": 0
		},
		{
			"ID": 11,
			"Title": "产品数据",
			"Pid": 8,
			"Url": "/appAnalysis",
			"Icon": "Histogram",
			"Sort": 0
		},
		{
			"ID": 12,
			"Title": "下注历史",
			"Pid": 8,
			"Url": "/betLog",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 18,
			"Title": "slots生成",
			"Pid": 0,
			"Url": "/slotsgen",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 20,
			"Title": "配置文件",
			"Pid": 0,
			"Url": "/config",
			"Icon": "Paperclip",
			"Sort": 2
		},
		{
			"ID": 21,
			"Title": "产品管理",
			"Pid": 0,
			"Url": "/app",
			"Icon": "Basketball",
			"Sort": 5
		},
		{
			"ID": 22,
			"Title": "权限组列表",
			"Pid": 1,
			"Url": "/groupList",
			"Icon": "Setting",
			"Sort": 10
		},
		{
			"ID": 23,
			"Title": "管理员列表",
			"Pid": 1,
			"Url": "/adminList",
			"Icon": "Setting",
			"Sort": 20
		},
		{
			"ID": 24,
			"Title": "后台管理",
			"Pid": 0,
			"Url": "/admin",
			"Icon": "Setting",
			"Sort": 50
		},
		{
			"ID": 25,
			"Title": "游戏",
			"Pid": 21,
			"Url": "/game_list",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 26,
			"Title": "运营商维护",
			"Pid": 24,
			"Url": "/operatorMaintenance",
			"Icon": "Avatar",
			"Sort": 0
		},
		{
			"ID": 27,
			"Title": "每日平台汇总",
			"Pid": 8,
			"Url": "/daily_platform_summary",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 28,
			"Title": "每月平台汇总",
			"Pid": 8,
			"Url": "/monthly_platform_summary",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 29,
			"Title": "平台玩家汇总",
			"Pid": 8,
			"Url": "/platform_player_summary",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 30,
			"Title": "功能菜单维护",
			"Pid": 24,
			"Url": "/menuList",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 31,
			"Title": "游戏运行数据统计",
			"Pid": 0,
			"Url": "/runtime",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 32,
			"Title": "游戏列表",
			"Pid": 21,
			"Url": "/game_list",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 33,
			"Title": "平台汇总",
			"Pid": 8,
			"Url": "/platform_summary",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 34,
			"Title": "玩家信息（运营商）",
			"Pid": 3,
			"Url": "/playerInfoOperator",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 35,
			"Title": "功能权限维护",
			"Pid": 24,
			"Url": "/groupList",
			"Icon": "Setting",
			"Sort": 0
		},
		{
			"ID": 36,
			"Title": "用户维护",
			"Pid": 24,
			"Url": "/adminList",
			"Icon": "Setting",
			"Sort": 0
		}

	]`
}

// =========== Help =====
func IsAdminUser(ctx *Context) (comm.User, bool) {
	var user comm.User
	err := CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil || user.GroupId > 1 {
		return user, false
	}
	return user, true
}

func GetUser(ctx *Context) (comm.User, error) {
	var user comm.User
	err := CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

var rtpList1 = []int{10, 20, 30, 40, 50, 60, 70, 80, 85, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99}
var rtpList2 = []int{93, 94, 95, 96, 97, 98, 99, 100, 102, 105, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 220, 240, 260, 280, 300}
var rtpList3 = []int{10, 20, 30, 40, 50, 60, 70, 80, 85, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 102, 105, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 220, 240, 260, 280, 300}
var rtpListMap = map[int][]int{
	1: rtpList1,
	2: rtpList2,
	3: rtpList3,
}

func Contains(slice []int, number int) bool {
	for _, v := range slice {
		if v == number {
			return true
		}
	}
	return false
}

// 获取当前登录商户是否有修改个人RTP的权限
func VerifyAuth(ctx *Context, ps CreatePlayerRTP) (bool, error) {
	user, err := GetUser(ctx)
	if err != nil {
		return false, err
	}
	oper := comm.Operator_V2{}
	CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &oper)

	//线 1   运 2
	if oper.PlayerRTPSettingOff == 0 && oper.OperatorType == 2 {
		return false, errors.New(lang.GetLang(ctx.Lang, "个人RTP权限已关闭"))
	}
	if !Contains(rtpListMap[int(oper.HighRTPOff)], int(ps.ContrllRTP)) && oper.OperatorType == 2 || (ps.AutoRemoveRTP != 0 && !Contains(rtpListMap[3], int(ps.AutoRemoveRTP))) && oper.OperatorType == 2 {
		return false, errors.New(lang.GetLang(ctx.Lang, "RTP值不在允许范围内"))
	}
	return true, nil
}
