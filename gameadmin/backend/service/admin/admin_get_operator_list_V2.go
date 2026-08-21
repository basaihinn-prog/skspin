package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/ut"
	"game/duck/lang"
	"game/duck/logger"
	"game/duck/mongodb"
	"game/service/admin/channel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetOperatorListV2", "获取营运商列表V2", "AdminInfo", getOperatorListV2, getOperatorListParamsV2{})
	RegMsgProc("/AdminInfo/AddOperatorV2", "添加运营商V2", "AdminInfo", addOperatorV2, addOperatorParamsV2{
		UserName:        "",
		AppID:           "",
		CooperationType: 1, //收益分成：1   购分方式：2
	})
	RegMsgProc("/AdminInfo/UpdataOperatorV2", "修改运营商V2", "AdminInfo", updataOperator, updateOperatorParams{})
	RegMsgProc("/AdminInfo/GetLineOperator", "获取线路商下的商户", "AdminInfo", getLineOperator, getOperatorListParamsV2{})
	RegMsgProc("/AdminInfo/getOperatorByAppId", "获取商户详情", "AdminInfo", getOperatorInfoByAppId, getOperatorListParamsV2{})
	mux.RegRpc("/AdminInfo/Interior/operatorVerifyIP", "内部接口-通过AppId与IP判断是否为白名单", "AdminInfo", operatorVerifyIp, operatorVerify{})
	RegMsgProc("/AdminInfo/GetMerchantReviewList", "获取线路商审核列表", "AdminInfo", getMerchantReviewList, reviewListParams{})
	RegMsgProc("/AdminInfo/UpdateMerchantReviewStatus", "修改线路商审核状态", "AdminInfo", updateMerchantReviewStatus, reviewStatusParams{})
	mux.RegRpc("/AdminInfo/Interior/operatorInfo", "内部接口-通过AppId获取当前用户的详情", "AdminInfo", operatorInfo, operatorVerify{})
	RegMsgProc("/AdminInfo/GetExtraRtpConfig", "获取额外rtp配置", "AdminInfo", getExtraRtpConfig, getExtraRtpConfigParams{})
	RegMsgProc("/AdminInfo/SetExtraRtpConfig", "修改额外rtp配置", "AdminInfo", updateExtraRtpConfig, updateExtraRtpConfigParams{})
	mux.RegRpc("/AdminInfo/Interior/setExtraRtpConfig", "内部接口-设置额外rtp配置", "AdminInfo", interUpdateExtraRtpConfig, updateExtraRtpConfigParams{})
	RegMsgProc("/AdminInfo/IncrementUpdataOperator", "增量修改运营商", "AdminInfo", incrementUpdataOperator, incrementUpdateOperatorParams{})
	RegMsgProc("/AdminInfo/SetGameMonitor", "更新游戏监控配置", "AdminInfo", SetGameMonitor, setGameMonitorParams{})
	mux.RegRpc("/AdminInfo/Interior/SetGameMonitor", "内部接口-更新游戏监控配置", "AdminInfo", interSetGameMonitor, setGameMonitorParams{})
}

var defaultMoniter = comm.MoniterConfig{
	IsMoniter:            0,
	MoniterNewbieNum:     50,
	MoniterRTPErrorValue: 5,
	MoniterNumCycle:      5,
	MoniterAddRTPRangeValue: []comm.MoniterRtpRangeValue{
		{
			RangeMinValue:  1,
			RangeMaxValue:  6,
			NewbieValue:    0,
			NotNewbieValue: 0,
		},
		{
			RangeMinValue:  6,
			RangeMaxValue:  11,
			NewbieValue:    4,
			NotNewbieValue: 4,
		},
		{
			RangeMinValue:  11,
			RangeMaxValue:  16,
			NewbieValue:    8,
			NotNewbieValue: 8,
		},
		{
			RangeMinValue:  16,
			RangeMaxValue:  21,
			NewbieValue:    12,
			NotNewbieValue: 12,
		},
		{
			RangeMinValue:  21,
			RangeMaxValue:  9999999,
			NewbieValue:    15,
			NotNewbieValue: 15,
		},
	},
	MoniterReduceRTPRangeValue: []comm.MoniterRtpRangeValue{
		{
			RangeMinValue:  1,
			RangeMaxValue:  6,
			NewbieValue:    0,
			NotNewbieValue: 1,
		},
		{
			RangeMinValue:  6,
			RangeMaxValue:  11,
			NewbieValue:    5,
			NotNewbieValue: 5,
		},
		{
			RangeMinValue:  11,
			RangeMaxValue:  16,
			NewbieValue:    10,
			NotNewbieValue: 10,
		},
		{
			RangeMinValue:  16,
			RangeMaxValue:  21,
			NewbieValue:    15,
			NotNewbieValue: 15,
		},
		{
			RangeMinValue:  21,
			RangeMaxValue:  9999999,
			NewbieValue:    20,
			NotNewbieValue: 20,
		},
	},
}

type operatorVerify struct {
	IP    string
	AppID string
}

type getOperatorListParamsV2 struct {
	Status           int64  `json:"Status" bson:"Status"`             //状态
	AppID            string `json:"Name" bson:"Name"`                 //商户名称
	RealAppID        string `json:"AppID" bson:"AppSecret"`           //商户名称
	OperatorType     int    `json:"OperatorType" bson:"OperatorType"` //商户类型
	CreatedStartTime int64  `json:"CreatedStartTime"`
	CreatedEndTime   int64  `json:"CreatedEndTime"`
	PageIndex        int64  `json:"PageIndex"`
	PageSize         int64  `json:"PageSize"`
	AppSecret        string `json:"AppSecret"`
}

type getOperatorListResultV2 struct {
	List     []*comm.Operator_V2
	AllCount int64
}

type addOperatorParamsV2 struct {
	MenuIds []int64

	ExcluedGameId  primitive.ObjectID
	ExcluedGameIds []string

	//运营商类型
	OperatorType int64 `json:"OperatorType" bson:"OperatorType"`

	//运营商名称
	AppID string `json:"AppID" bson:"AppID" md:"运营商"`

	//运营商主账号/初始账号
	UserName string `json:"UserName" bson:"Username"`

	//平台费
	PlatformPay float64 `json:"PlatformPay" bson:"PlatformPay"`

	//合作模式	//暂定只有一个值收益分成 #2024.10.12  收益分成：1   购分方式：2
	CooperationType int `json:"CooperationType" bson:"CooperationType"`

	//预付款金额
	Advance float32 `json:"Advance" bson:"Advance"`

	//商户币种
	CurrencyKey string `json:"CurrencyKey" bson:"CurrencyKey"`

	//联系方式
	Contact string `json:"Contact" bson:"Contact"`

	//下面的是线路商没有,运营商有的字段
	//钱包类型
	WalletMode int `json:"WalletMode" bson:"WalletMode"`

	//会员前缀
	Surname string `json:"Surname" bson:"Surname"`

	//客户端默认语言
	Lang string `json:"Lang" bson:"Lang"`

	//服务器地址
	ServiceIp string `json:"ServiceIp" bson:"ServiceIp"`

	//服务器IP白名单
	WhiteIps []string `json:"WhiteIps" bson:"WhiteIps"`

	//服务器回调地址
	Address string `json:"Address" bson:"Address"`

	//用户白名单
	UserWhite string `json:"UserWhite" bson:"UserWhite"`

	//用来区分线路商和总控 程序运行需要
	LineMerchant int32 `json:"_LineMerchant" bson:"_LineMerchant"`

	Remark           string `description:"备注"`
	BelongingCountry string `description:"归属国家"`

	Robot  string `description:"Telegram机器人信息"`
	ChatID string `description:"Telegram机器人信息"`
	// 线路商户ID
	LineMerchantID string `json:"LineMerchantID"`
	// 商户余额
	Balance float64 `json:"Balance"`
	// 审核类型
	ReviewType int64 `json:"ReviewType"`
	// 余额不足阈值
	BalanceThreshold float64 `json:"BalanceThreshold"`
	// 余额不足间隔
	BalanceThresholdInterval int64 `json:"BalanceThresholdInterval"`
	// 欠费提示间隔
	ArrearsThresholdInterval int64 `json:"ArrearsThresholdInterval"`
	// 密码
	Password string `json:"Password"`

	// 新增参数
	ShowExitBtnOff     int64 `json:"ShowExitBtnOff" bson:"ShowExitBtnOff"`
	CurrencyVisibleOff int   `description:"是否显示币种"`
	//RTP设置
	RTPOff int64 `json:"RTPOff" bson:"RTPOff"`
	//赢取最高押注倍数  	只是一个开关  具体数值是针对某个游戏的，在游戏配置数据库中
	MaxMultipleOff      int64  `json:"MaxMultipleOff" bson:"MaxMultipleOff"`
	ExitLink            string `json:"ExitLink" bson:"ExitLink"`
	PlayerRTPSettingOff int64  `json:"PlayerRTPSettingOff"`
	MaxWinPointsOff     int64  `description:"最大赢钱数"`
	CurrencyCtrlStatus  int    `json:"CurrencyCtrlStatus"`
	// 高rtp显示开关
	HighRTPOff int64 `json:"HighRTPOff"`
	// 流水费
	TurnoverPay float64 `json:"TurnoverPay"`
	// 默认厂商开启
	DefaultManufacturerOn []string `json:"DefaultManufacturerOn"`
	// PG配置
	PGConfig comm.PGConfig `json:"PGConfig"`
	// JILI配置
	JILIConfig comm.JILIConfig `json:"JILIConfig"`
	// PP配置
	PPConfig comm.PPConfig `json:"PPConfig"`
	// TADA配置
	TADAConfig comm.TADAConfig `json:"TADAConfig"`
	// 购买Rtp开关
	BuyRTPOff int64 `json:"BuyRTPOff"`
}

type incrementUpdateOperatorParams struct {
	AppID          string         `json:"AppID"`
	IncrementKey   []string       `json:"IncrementKey"`
	IncrementValue map[string]any `json:"IncrementValue"`
}

type setGameMonitorParams struct {
	AppID                      string                      `json:"AppID"`
	MoniterType                int                         `json:"MoniterType"`
	IsMoniter                  int64                       `json:"IsMoniter"`
	MoniterNewbieNum           int64                       `json:"MoniterNewbieNum"`
	MoniterRTPErrorValue       int64                       `json:"MoniterRTPErrorValue"`
	MoniterNumCycle            int64                       `json:"MoniterNumCycle"`
	MoniterAddRTPRangeValue    []comm.MoniterRtpRangeValue `json:"MoniterAddRTPRangeValue"`
	MoniterReduceRTPRangeValue []comm.MoniterRtpRangeValue `json:"MoniterReduceRTPRangeValue"`
}
type updateOperatorParams struct {
	UserName string `json:"UserName" bson:"UserName"`
	AppID    string `json:"AppID" bson:"AppID"`
	//运营商状态   AdminOperator_copy1 / Status
	Status int64 `json:"Status" bson:"Status"`
	//合作模式   收益分成：1   购分方式：2   暂定只有收益分成模式 #2024.10.12
	CooperationType int `json:"CooperationType" bson:"CooperationType"`
	//下月费率
	NextRate float64 `json:"NextRate" bson:"NextRate"`

	//平台费
	PlatformPay float64 `json:"PlatformPay" bson:"PlatformPay"`
	//钱包类型        AdminOperator_copy1 / WalletMode
	WalletMode int `json:"WalletMode" bson:"WalletMode"`
	//商户币种        AdminOperator_copy1 / CurrencyKey
	CurrencyKey string `json:"CurrencyKey" bson:"CurrencyKey"`
	//客户端默认语言
	Lang string `json:"Lang" bson:"Lang"`
	//预付款金额
	Advance float32 `json:"Advance" bson:"Advance"`
	//联系方式           为多条数据
	Contact string `json:"Contact" bson:"Contact"`
	//服务器回调地址
	Address string `json:"Address" bson:"Address"`
	//服务器地址
	ServiceIp string `json:"ServiceIp" bson:"ServiceIp"`
	//记录开启
	LoginOff int64 `json:"LoginOff" bson:"LoginOff"`
	//断复开关
	RestoreOff int64 `json:"RestoreOff" bson:"RestoreOff"`
	//防休眠开启
	DormancyOff int64 `json:"DormancyOff" bson:"DormancyOff"`
	//免游游戏开关板
	FreeOff int64 `json:"FreeOff" bson:"FreeOff"`
	//消息推送
	MassageOff int64 `json:"MassageOff" bson:"MassageOff"`
	//消息推送地址
	MassageIp string `json:"MassageIp" bson:"MassageIp"`
	//游戏新型预设
	NewGameDefaulOff int64 `json:"NewGameDefaul" bson:"NewGameDefaul"`
	//手动全屏开启
	ManualFullScreenOff int64 `json:"ManualFullScreen" bson:"ManualFullScreenOff"`
	//RTP设置
	RTPOff int64 `json:"RTPOff" bson:"RTPOff"`
	//止损止盈开关
	StopLoss int64 `json:"StopLoss" bson:"StopLoss"`
	//赢取最高押注倍数  	只是一个开关  具体数值是针对某个游戏的，在游戏配置数据库中
	MaxMultipleOff      int64 `json:"MaxMultipleOff" bson:"MaxMultipleOff"`
	ShowNameAndTimeOff  int64 `json:"ShowNameAndTimeOff"`
	UserWhite           string
	PlayerRTPSettingOff int64  `json:"PlayerRTPSettingOff"`
	ShowExitBtnOff      int64  `json:"ShowExitBtnOff" bson:"ShowExitBtnOff"`
	ExitLink            string `json:"ExitLink" bson:"ExitLink"`

	Remark             string `description:"备注"`
	BelongingCountry   string `description:"归属国家"`
	MaxWinPointsOff    int64  `description:"最大赢钱数"`
	Robot              string `description:"Telegram信息"`
	ChatID             string `description:"Telegram信息"`
	CurrencyVisibleOff int    `description:"是否显示币种"`
	// 商户余额
	Balance float64 `json:"Balance"`
	// 审核类型
	ReviewType int64 `json:"ReviewType"`
	// 余额不足阈值
	BalanceThreshold float64 `json:"BalanceThreshold"`
	// 余额不足间隔
	BalanceThresholdInterval int64 `json:"BalanceThresholdInterval"`
	// 欠费提示间隔
	ArrearsThresholdInterval int64 `json:"ArrearsThresholdInterval"`
	// 高rtp显示开关
	HighRTPOff int64 `json:"HighRTPOff"`
	// 流水费
	TurnoverPay float64 `json:"TurnoverPay"`
	// 默认厂商开启
	DefaultManufacturerOn []string `json:"DefaultManufacturerOn"`
	//总控机器人
	AdminRobot string `json:"AdminRobot"`
	//总控机器人的chatid
	AdminChatID string `json:"AdminChatID"`
	//余额告警阈值
	BalanceAlert float64 `json:"BalanceAlert"`
	//余额告警时间 暂时固定为24每天0点进行检测并触发余额告警
	BalanceAlertTimeInterval int64 `json:"BalanceAlertTimeInterval"`
	// 购买Rtp开关
	BuyRTPOff int64 `json:"BuyRTPOff"`
}
type getLineOpertorreq struct {
	AppID     string `json:"AppID"`
	PageIndex int64
	PageSize  int64
}

type reviewListParams struct {
	AppID       string `json:"AppID"`
	ParentAppID string `json:"ParentAppID"`
	StartTime   int64  `json:"StartTime"`
	EndTime     int64  `json:"EndTime"`
	Page        int64  `json:"Page"`
	PageSize    int64  `json:"PageSize"`
}

type reviewStatusParams struct {
	OperatorID int64 `json:"OperatorID"`
}

type getReviewListResults struct {
	List []*comm.Operator_V2
	All  int64
}

type getExtraRtpConfigParams struct {
	Page         int64  `json:"Page"`
	PageSize     int64  `json:"PageSize"`
	AppID        string `json:"AppID"`
	IsProtection int    `json:"IsProtection"`
}

type updateExtraRtpConfigParams struct {
	AppID string `json:"AppID"`
	// 保护开关
	IsProtection int `json:"IsProtection""`
	// 保护次数
	ProtectionRotateCount int `json:"ProtectionRotateCount"`
	// 保护内获取比率
	ProtectionRewardPercentLess int `json:"ProtectionRewardPercentLess"`
}

type getExtraRtpConfigResults struct {
	List []*comm.Operator_V2
	All  int64
}

func addOperatorV2(ctx *Context, ps addOperatorParamsV2, ret *comm.Empty) (err error) {
	//运营商名称校验
	if ps.UserName == "admin" || ps.AppID == "admin" {
		return errors.New(lang.GetLang(ctx.Lang, "不能添加admin商户"))
	}
	if ps.UserName == "" || ps.AppID == "" {
		return errors.New(lang.GetLang(ctx.Lang, "请输入账号和商户名称"))
	}
	user, _ := IsAdminUser(ctx)
	if user.GroupId == 3 {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	ps.AppID = removeSpecialCharacters(ps.AppID)
	if len(ps.AppID) == 0 {
		return errors.New(lang.GetLang(ctx.Lang, "商户名称不合法"))
	}
	//appSecret := removeSpecialCharacters(ps.UserName)
	//if len(appSecret) == 0 {
	//	return errors.New(lang.GetLang(ctx.Lang, "输入不合法"))
	//}
	//线路商不可创建线路商
	reviewStatus := int64(1)
	reviewer := user.UserName
	if user.GroupId == 2 {
		if ps.OperatorType == 1 || ps.LineMerchantID != "" || ps.ReviewType == 2 {
			return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
		} else if ps.ReviewType == 0 {
			reviewStatus = 0
			reviewer = ""
		}
	}
	partntAppId := ""
	if ps.LineMerchantID != "" {
		partntAppId = ps.LineMerchantID
	} else {
		partntAppId = user.AppID
	}

	//判断重复
	var result bson.M
	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &result)
	if err == nil {
		return errors.New(lang.GetLang(ctx.Lang, "输入不合法或商户名称重复"))
	}
	err = CollAdminUser.FindOne(bson.M{"Username": ps.UserName}, &result)
	if err == nil {
		return errors.New(lang.GetLang(ctx.Lang, "账号已存在，修改后重试"))
	}

	id, err := NextID(CollAdminOperator, 100000)
	if err != nil {
		return err
	}

	//获取默认权限组
	var minudis comm.Permission
	CollAdminPermission.FindOne(bson.M{"_id": ps.OperatorType + 1}, &minudis)

	var operator *comm.Operator_V2

	//构造数据
	operator = &comm.Operator_V2{
		Id:             id,
		Name:           partntAppId,
		MenuIds:        minudis.MenuIds,
		PermissionId:   ps.OperatorType,
		ExcluedGameIds: nil,
		//下面这些字段是新添加的
		AppID:           ps.AppID,
		AppSecret:       uuid.NewString(), //进入游戏使用 2024 11 20
		UserName:        ps.UserName,
		PlatformPay:     ps.PlatformPay,
		CooperationType: ps.CooperationType,
		PresentRate:     ps.PlatformPay, //默认费率  为 5
		NextRate:        ps.PlatformPay,
		OperatorType:    ps.OperatorType,
		Advance:         ps.Advance,
		Status:          1,
		CreateTime:      mongodb.NewTimeStamp(time.Now()),
		TokenExpireAt:   nil,
		CurrencyKey:     ps.CurrencyKey,
		Contact:         ps.Contact,
		//运营商独有字段
		WalletMode: ps.WalletMode,
		Surname:    ps.Surname,
		Lang:       ps.Lang,

		ServiceIp:           ps.ServiceIp,
		WhiteIps:            ps.WhiteIps,
		Address:             ps.Address,
		UserWhite:           ps.UserWhite,
		LoginOff:            0,
		FreeOff:             0,
		DormancyOff:         0,
		RestoreOff:          0,
		ManualFullScreenOff: 0,
		NewGameDefaulOff:    0,
		MassageOff:          0,
		MassageIp:           "",
		RTPOff:              ps.RTPOff,
		StopLoss:            0,
		MaxMultipleOff:      ps.MaxMultipleOff,
		CurrencyManufactureVisibleOff: map[string]int{
			"jili": ps.CurrencyVisibleOff,
			"pp":   ps.CurrencyVisibleOff,
			"pg":   ps.CurrencyVisibleOff,
			"tada": ps.CurrencyVisibleOff,
		},
		PlayerRTPSettingOff: ps.PlayerRTPSettingOff,
		LineMerchant:        ps.OperatorType, //暂时无用 #2024.10.12
		ShowExitBtnOff:      ps.ShowExitBtnOff,
		ExitLink:            ps.ExitLink,

		// 二期需求需要添加的字段
		Remark:           ps.Remark,
		BelongingCountry: ps.BelongingCountry,
		MaxWinPointsOff:  ps.MaxWinPointsOff,

		Robot:                       ps.Robot,
		ChatID:                      ps.ChatID,
		ReviewStatus:                reviewStatus,
		Reviewer:                    reviewer,
		Balance:                     ps.Balance,
		ReviewType:                  ps.ReviewType,
		BalanceThreshold:            ps.BalanceThreshold,
		BalanceThresholdInterval:    ps.BalanceThresholdInterval,
		ArrearsThresholdInterval:    ps.ArrearsThresholdInterval,
		HighRTPOff:                  ps.HighRTPOff,
		IsProtection:                0,
		ProtectionRotateCount:       0,
		ProtectionRewardPercentLess: 0,
		CurrencyCtrlStatus:          ps.CurrencyCtrlStatus,
		TurnoverPay:                 ps.TurnoverPay,
		DefaultManufacturerOn:       ps.DefaultManufacturerOn,
		PGConfig:                    ps.PGConfig,
		PPConfig:                    ps.PPConfig,
		JILIConfig:                  ps.JILIConfig,
		TADAConfig:                  ps.TADAConfig,
		BalanceAlertTimeInterval:    24,
		BalanceAlert:                0,
		BuyRTPOff:                   ps.BuyRTPOff,
		MoniterConfig:               defaultMoniter,
		PersonalMoniterConfig:       defaultMoniter,
	}

	//添加游戏配置表  Game / ps.AppId
	var list_OperatorGameConfig []*operatorGameConfig
	err = NewOtherDB("game").Collection("Games").FindAll(bson.M{}, &list_OperatorGameConfig)
	if err != nil {
		return err
	}
	insertList := []bson.M{}
	// 留个口子，以后币种显示需要根据游戏处理只需要调用方处理
	for _, config := range list_OperatorGameConfig {
		setGameConfig := operatorSetGameConfig(ps.AppID, ps.UserName, *config)
		switch strings.ToLower(config.GameManufacturer) {
		case comm.PG:
			setGameConfig["CarouselOff"] = ps.PGConfig.CarouselOff
			setGameConfig["ShowNameAndTimeOff"] = ps.PGConfig.ShowNameAndTimeOff
			setGameConfig["StopLoss"] = ps.PGConfig.StopLoss
			setGameConfig["ExitBtnOff"] = ps.PGConfig.ExitBtnOff
			setGameConfig["ExitLink"] = ps.PGConfig.ExitLink
			setGameConfig["OfficialVerify"] = ps.PGConfig.OfficialVerify
		case comm.PP:
			setGameConfig["ShowNameAndTimeOff"] = ps.PPConfig.ShowNameAndTimeOff
		case comm.JILI:
			setGameConfig["BackPackOff"] = ps.JILIConfig.BackPackOff
			setGameConfig["OpenScreenOff"] = ps.JILIConfig.OpenScreenOff
			setGameConfig["SidebarOff"] = ps.JILIConfig.SidebarOff
		case comm.TADA:
			setGameConfig["BackPackOff"] = ps.TADAConfig.BackPackOff
			setGameConfig["OpenScreenOff"] = ps.TADAConfig.OpenScreenOff
			setGameConfig["SidebarOff"] = ps.TADAConfig.SidebarOff
		}
		//setGameConfig["CurrencyVisibleOff"] = ps.CurrencyVisibleOff
		insertList = append(insertList, setGameConfig)
	}

	//注册运营商管理员
	var operatorAdmin comm.User
	userId, err := NextID(CollAdminUser, 1000)
	if err != nil {
		return err
	}
	date := make(map[string]any)
	date["$date"] = mongodb.NewTimeStamp(time.Now()).AsTime().UTC().String()
	// rsa 解密
	//rsaBytes, err := base64.StdEncoding.DecodeString(ps.Password)
	//if err != nil {
	//	return err
	//}
	//passWord, err := comm.RsaDecrypt(ps.Password)
	//if err != nil {
	//	return err
	//}
	//md5Password := comm.Md5Encrypt(passWord)

	operatorAdmin = comm.User{
		Id:            userId,
		AppID:         ps.AppID,
		UserName:      ps.UserName,
		PassWord:      ps.Password,
		GoogleCode:    "",
		Qrcode:        "",
		Avatar:        "",
		IsOpenGoogle:  false,
		Status:        comm.Operator_Status_Normal,
		CreateAt:      mongodb.NewTimeStamp(time.Now()),
		LoginAt:       nil,
		GroupId:       ps.OperatorType + 1,
		Token:         "",
		TokenExpireAt: nil,
		OperatorAdmin: true,
		PermissionId:  ps.OperatorType + 1,
	}

	//添加运营商或商户
	if err = CollAdminOperator.InsertOne(operator); err != nil {
		return err
	}
	//添加运营商游戏配置
	if operator.OperatorType == 2 {
		if err = CollGameConfig.InsertMany(lo.ToAnySlice(insertList)); err != nil {
			return err
		}
	}
	//添加主账号
	if err = CollAdminUser.InsertOne(&operatorAdmin); err != nil {
		return err
	}

	return nil
}

func operatorSetGameConfig(AppID, UserName string, config operatorGameConfig) bson.M {

	set := &GameBet{}
	set, ok := MapGameBet[config.GameIdd]
	if ok == false {
		set = &GameBet{}
		set.Bet = "0.1,0.2,0.3"
		set.GameName = ""
		set.DefaultBet = 0.1
		set.DefaultBetLevel = 1
	}

	BetList := strings.Split(set.Bet, ",")

	OnlineUpNum, err := strconv.ParseFloat(BetList[0], 64)
	if err != nil {
		OnlineUpNum = 0
		logger.Err(err.Error())
	}
	OnlineDownNum, err := strconv.ParseFloat(BetList[len(BetList)-1], 64)
	if err != nil {
		OnlineDownNum = 0
		logger.Err(err.Error())
	}

	operatorConfig := bson.M{}
	operatorConfig["AppID"] = AppID
	operatorConfig["AppName"] = UserName         //todo 预留插槽
	operatorConfig["GameId"] = config.GameIdd    //游戏编号
	operatorConfig["GameName"] = set.GameName    //游戏名称
	operatorConfig["ConfigPath"] = "config/path" //配置文件
	operatorConfig["StopLoss"] = 1               //止盈止损开关   0 开启  1 关闭
	operatorConfig["MaxMultipleOff"] = 0         //赢取最高押注倍数
	operatorConfig["MaxMultiple"] = 100
	operatorConfig["MaxWinPoints"] = 1000000                 // 赢取最高分
	operatorConfig["BetBase"] = set.Bet                      //游戏投注
	operatorConfig["GamePattern"] = BUCKET_STABLE            //游戏类型,默认值改为3
	operatorConfig["Preset"] = 10                            //预设面额
	operatorConfig["GameOn"] = 0                             //游戏状态  0 开启  1 关闭
	operatorConfig["ShowNameAndTimeOff"] = 1                 //显示游戏名称和时间
	operatorConfig["ShowExitBtnOff"] = 0                     //显示退出按钮 todo：修改到operator表中
	operatorConfig["ExitLink"] = ""                          //退出按钮链接 todo：修改到operator表中
	operatorConfig["RTP"] = 93                               //RTP设置
	operatorConfig["BuyRTP"] = 90                            //购买游戏RTP(显示的值)
	operatorConfig["BuyMinAwardPercent"] = MapBuyGameRTP[90] //购买游戏RTP(实际游戏用的值)
	operatorConfig["OnlineUpNum"] = OnlineUpNum              //Bet值上限
	operatorConfig["OnlineDownNum"] = OnlineDownNum          //Bet值下限
	operatorConfig["ProfitMargin"] = 3.00                    // 利润率
	operatorConfig["CrashRate"] = 1.00                       //坠毁概率
	operatorConfig["Scale"] = 3.00                           //刻度
	// pg游戏添加默认押注
	if strings.HasPrefix(config.GameIdd, "pg_") {
		operatorConfig["DefaultCs"] = set.DefaultBet            //默认押注
		operatorConfig["DefaultBetLevel"] = set.DefaultBetLevel //默认押注等级
	} else if strings.HasPrefix(config.GameIdd, "hacksaw_") {
		operatorConfig["DefaultCs"] = set.DefaultBet //默认押注
	}
	rw, nw := GetGameNwandAw(config.GameIdd, 96)
	operatorConfig["RewardPercent"] = rw
	operatorConfig["NoAwardPercent"] = nw

	return operatorConfig
}

// 获取商户上月和本月收益
func GetOperatorProfit(ctx *Context, OperatorList []*comm.Operator_V2) ([]*comm.Operator_V2, error) {
	ids := []string{}
	for _, item := range OperatorList {
		ids = append(ids, item.AppID)
	}
	lastDate := time.Now().UTC().AddDate(0, -1, 0).Format("200601")
	lastGameInfoData, err := getHistoryMonthlyReport(lastDate, ids)
	if err != nil {
		return OperatorList, err
	}
	lastMap := map[string]float64{}
	for _, item := range lastGameInfoData {
		lastMap[item.AppId] = item.InCome
	}
	gameInfoData, err := getNowMonthlyReport(ids)
	nowMap := map[string]float64{}
	for _, item := range gameInfoData {
		nowMap[item.AppId] = item.InCome
	}
	if err != nil {
		return OperatorList, err
	}
	// 获取货币
	operck, err := GetCurrencyList(ctx)
	for _, item := range OperatorList {
		item.ThisMonthProfit = nowMap[item.AppID]
		item.LastMonthProfit = lastMap[item.AppID]
		if _, ok := operck[item.AppID]; !ok {
			operck[item.AppID] = &comm.CurrencyType{
				CurrencyCode:   "/",
				CurrencyName:   "/",
				CurrencySymbol: "/",
			}
		}

		if ctx.Lang != "zh" {
			item.CurrencyName = operck[item.AppID].CurrencyCode
			item.CurrencyKey = operck[item.AppID].CurrencyCode
		} else {
			item.CurrencyName = operck[item.AppID].CurrencyName
			item.CurrencyKey = operck[item.AppID].CurrencyCode
		}
	}
	return OperatorList, nil
}

func getOperatorListV2(ctx *Context, ps getOperatorListParamsV2, ret *getOperatorListResultV2) (err error) {

	user, _ := IsAdminUser(ctx)
	var list []*comm.Operator_V2
	query := bson.M{}
	if ps.CreatedStartTime != 0 {
		startTime := mongodb.NewTimeStamp(time.UnixMilli(ps.CreatedStartTime))
		endTime := mongodb.NewTimeStamp(time.UnixMilli(ps.CreatedEndTime))
		query = bson.M{
			//"Name":user.AppID,
			"CreateTime": bson.M{
				"$gte": startTime, // 大于或等于开始时间
				"$lte": endTime,   // 小于或等于结束时间
			},
		}
	}
	if ps.AppID != "" { //模糊查询
		//query["AppID"] = bson.M{"$regex": "/" + ps.AppID + "/", "$options": 'i'}
		query["AppID"] = primitive.Regex{
			Pattern: ps.AppID,
			Options: "i",
		}
	}

	if ps.RealAppID != "" { //模糊查询
		//query["AppID"] = bson.M{"$regex": "/" + ps.AppID + "/", "$options": 'i'}
		query["AppSecret"] = primitive.Regex{
			Pattern: ps.RealAppID,
			Options: "i",
		}
	}

	if ps.Status != -1 {
		query["Status"] = ps.Status
	}

	if ps.OperatorType != 0 {
		query["OperatorType"] = ps.OperatorType
	}

	if ps.AppSecret != "" {
		query["AppSecret"] = ps.AppSecret
	}
	//大于200 填充下拉列表
	if ps.PageSize > 200 {
		switch user.GroupId {
		case 0:
			//query["Name"] = "admin"
		case 1:
			//query["Name"] = "admin"
		case 2:
			query["Name"] = user.AppID
			//query["$or"] = []bson.M{
			//	{"Name": user.AppID},
			//	{"AppID": user.AppID},
			//}
		case 3:
			query["AppID"] = user.AppID
		default:
		}
	}

	pageFileter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     bson.M{"CreateTime": -1},
		Query:    query,
	}

	ret.AllCount, err = CollAdminOperator.FindPage(pageFileter, &list)
	if err != nil {
		return err
	}
	ret.List, err = GetOperatorProfit(ctx, list)
	if err != nil {
		return err
	}

	return err
}

func GetOperatorInfo(query bson.M, list *[]*comm.Operator_V2) (err error) {

	err = CollAdminOperator.FindAll(query, list)
	if err != nil {
		return
	}

	for _, operatorv2 := range *list {
		if operatorv2.Status == 0 {
			file := bson.M{"AppID": operatorv2.AppID}
			setInfo := bson.M{"Status": 1}
			err = CollAdminUser.Update(file, bson.M{"$set": setInfo})
			if err != nil {
				return
			}
		}
	}
	return err
}

func updataOperator(ctx *Context, ps updateOperatorParams, ret *comm.Empty) (err error) {
	settingInfo := bson.M{}
	settingInfo["Status"] = ps.Status
	settingInfo["CooperationType"] = ps.CooperationType
	settingInfo["WalletMode"] = ps.WalletMode
	settingInfo["CurrencyKey"] = ps.CurrencyKey
	settingInfo["Lang"] = ps.Lang
	settingInfo["Advance"] = ps.Advance
	settingInfo["Contact"] = ps.Contact
	settingInfo["Address"] = ps.Address
	settingInfo["ServiceIp"] = ps.ServiceIp
	settingInfo["LoginOff"] = ps.LoginOff
	settingInfo["RestoreOff"] = ps.RestoreOff
	settingInfo["PlatformPay"] = ps.PlatformPay
	settingInfo["TurnoverPay"] = ps.TurnoverPay
	settingInfo["DormancyOff"] = ps.DormancyOff
	settingInfo["FreeOff"] = ps.FreeOff
	settingInfo["MassageOff"] = ps.MassageOff
	settingInfo["MassageIp"] = ps.MassageIp
	settingInfo["NewGameDefaul"] = ps.NewGameDefaulOff
	settingInfo["ManualFullScreenOff"] = ps.ManualFullScreenOff
	settingInfo["RTPOff"] = ps.RTPOff
	settingInfo["StopLoss"] = ps.StopLoss
	settingInfo["MaxMultipleOff"] = ps.MaxMultipleOff
	settingInfo["NextRate"] = ps.NextRate //更新下月费率
	settingInfo["ShowNameAndTimeOff"] = ps.ShowNameAndTimeOff
	settingInfo["PlayerRTPSettingOff"] = ps.PlayerRTPSettingOff
	settingInfo["ShowExitBtnOff"] = ps.ShowExitBtnOff
	settingInfo["ExitLink"] = ps.ExitLink   //退出按钮链接
	settingInfo["UserWhite"] = ps.UserWhite //退出按钮链接

	// hml添加的最高赢钱数
	settingInfo["Remark"] = ps.Remark                                     //退出按钮链接
	settingInfo["BelongingCountry"] = ps.BelongingCountry                 //退出按钮链接
	settingInfo["MaxWinPointsOff"] = ps.MaxWinPointsOff                   //最高赢钱数
	settingInfo["Robot"] = ps.Robot                                       //最高赢钱数
	settingInfo["ChatID"] = ps.ChatID                                     //最高赢钱数
	settingInfo["AdminRobot"] = ps.AdminRobot                             //总控机器人
	settingInfo["AdminChatID"] = ps.AdminChatID                           //总控机器人的chatid
	settingInfo["BalanceAlert"] = ps.BalanceAlert                         //余额告警
	settingInfo["BalanceAlertTimeInterval"] = ps.BalanceAlertTimeInterval //余额告警时间

	// 游戏是否显示币种
	// 不同厂商分开控制
	settingInfo["CurrencyVisibleOff"] = map[string]int{
		"jili": ps.CurrencyVisibleOff,
		"pp":   ps.CurrencyVisibleOff,
		"pg":   ps.CurrencyVisibleOff,
		"tada": ps.CurrencyVisibleOff,
	} //最高赢钱数
	// 余额审核相关参数
	settingInfo["Balance"] = ps.Balance
	settingInfo["ReviewType"] = ps.ReviewType
	settingInfo["BalanceThreshold"] = ps.BalanceThreshold
	settingInfo["BalanceThresholdInterval"] = ps.BalanceThresholdInterval
	settingInfo["ArrearsThresholdInterval"] = ps.ArrearsThresholdInterval
	settingInfo["HighRTPOff"] = ps.HighRTPOff
	settingInfo["TurnoverPay"] = ps.TurnoverPay
	settingInfo["DefaultManufacturerOn"] = ps.DefaultManufacturerOn
	settingInfo["BuyRTPOff"] = ps.BuyRTPOff

	coll := DB.Collection("PlayerRTPControl")
	if ps.PlayerRTPSettingOff == 0 {
		err = coll.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"Status": 0}})
		if err != nil {
			return
		}
		go func() {

			contr := []*comm.PlayerRTPControlModel{}
			err = coll.FindAll(bson.M{"AppID": ps.AppID}, &contr)
			if err != nil {
				slog.Error("findAll PlayerRTPControl Error appid ==", ps.AppID)
			}
			for _, model := range contr {
				var oper = channel.RtpSettings{
					GameId:         model.GameID,
					RewardPercent:  0,
					NoAwardPercent: 0,
					PlayerId:       model.Pid,
				}
				channel.TaskQueue <- oper
			}
		}()
	}
	set := bson.M{
		"$set": settingInfo,
	}
	err = CollAdminOperator.Update(bson.M{"AppID": ps.AppID}, set)
	if err != nil {
		return
	}
	// 留个口子，以后币种显示需要根据游戏处理只需要调用方处理
	//err = CollGameConfig.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"ShowNameAndTimeOff": ps.ShowNameAndTimeOff, "CurrencyVisibleOff": ps.CurrencyVisibleOff}})
	err = CollGameConfig.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"ShowNameAndTimeOff": ps.ShowNameAndTimeOff}})
	if err != nil {
		return
	}
	var s int
	if ps.Status == 1 {
		s = 0
	} else {
		s = 1
	}
	err = CollAdminUser.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"Status": s}})
	if err != nil {
		return errors.New("参数错误")
	}
	dateStr := time.Now().UTC().Format("20060102")
	id := ut.JoinStr(':', dateStr, ps.AppID)
	update := db.D(
		"$set", db.D(
			"OperatorBalance", ps.Balance,
		),
	)
	collRep := db.Collection2("reports", "OperatorDailyReport")
	collRep.UpdateByID(context.TODO(), id, update)

	return
}
func incrementUpdataOperator(ctx *Context, ps incrementUpdateOperatorParams, ret *comm.Empty) (err error) {
	settingInfo := bson.M{}
	for _, k := range ps.IncrementKey {
		settingInfo[k] = ps.IncrementValue[k]
	}

	coll := DB.Collection("PlayerRTPControl")
	value, exists := ps.IncrementValue["PlayerRTPSettingOff"]
	if exists && value == 0 {
		err = coll.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"Status": 0}})
		if err != nil {
			return
		}
		go func() {

			contr := []*comm.PlayerRTPControlModel{}
			err = coll.FindAll(bson.M{"AppID": ps.AppID}, &contr)
			if err != nil {
				slog.Error("findAll PlayerRTPControl Error appid ==", ps.AppID)
			}
			for _, model := range contr {
				var oper = channel.RtpSettings{
					GameId:         model.GameID,
					RewardPercent:  0,
					NoAwardPercent: 0,
					PlayerId:       model.Pid,
				}
				channel.TaskQueue <- oper
			}
		}()
	}
	set := bson.M{
		"$set": settingInfo,
	}
	err = CollAdminOperator.Update(bson.M{"AppID": ps.AppID}, set)
	if err != nil {
		return
	}
	value, exists = ps.IncrementValue["ShowNameAndTimeOff"]
	if exists {
		err = CollGameConfig.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"ShowNameAndTimeOff": value}})
		if err != nil {
			return
		}
	}

	value, exists = ps.IncrementValue["Status"]
	if exists {
		var s int
		if value == 1 {
			s = 0
		} else {
			s = 1
		}
		err = CollAdminUser.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"Status": s}})
		if err != nil {
			return errors.New("参数错误")
		}
	}
	value, exists = ps.IncrementValue["Balance"]
	if exists {
		dateStr := time.Now().UTC().Format("20060102")
		id := ut.JoinStr(':', dateStr, ps.AppID)
		update := db.D(
			"$set", db.D(
				"OperatorBalance", value,
			),
		)
		collRep := db.Collection2("reports", "OperatorDailyReport")
		collRep.UpdateByID(context.TODO(), id, update)
	}

	return
}

func removeSpecialCharacters(input string) string {
	// 定义通用的特殊字符字典，匹配所有非字母数字字符
	re := regexp.MustCompile(`[^\w]`)
	// 将所有特殊字符替换为空字符串
	result := re.ReplaceAllString(input, "")
	return result
}

func GetOperatopAppID(ctx *Context) (listappid []string, err error) {
	listappid = make([]string, 0)
	user, _ := IsAdminUser(ctx)
	query := bson.M{}

	switch user.GroupId {
	case 2:
		query["Name"] = user.AppID
	case 3:
		query["AppID"] = user.AppID
	default:
	}
	if ctx.CurrencyKey != "" {
		query["CurrencyKey"] = ctx.CurrencyKey
	}

	var list []*comm.Operator_V2

	err = CollAdminOperator.FindAll(query, &list)

	if len(list) == 0 {
		return listappid, err
	}
	for _, operatorV2 := range list {
		listappid = append(listappid, operatorV2.AppID)
	}

	return
}
func GetOperatopList(ctx *Context) (list []*comm.Operator, err error) {
	user, _ := IsAdminUser(ctx)
	query := bson.M{}

	switch user.GroupId {
	case 2:
		query["Name"] = user.AppID
	case 3:
		query["AppID"] = user.AppID
	default:
	}
	if ctx.CurrencyKey != "" {
		query["CurrencyKey"] = ctx.CurrencyKey
	}

	err = CollAdminOperator.FindAll(query, &list)

	return
}

// 获取商户对应的币种信息
func GetCurrencyList(ctx *Context) (list map[string]*comm.CurrencyType, err error) {
	oplist, err := GetOperatopList(ctx)
	if err != nil {
		return nil, err
	}

	ck := []string{}
	for _, operator := range oplist {
		ck = append(ck, operator.CurrencyKey)
	}
	filter := bson.M{
		"CurrencyCode": bson.M{
			"$in": ck,
		}}
	cklist := []*comm.CurrencyType{}
	err = NewOtherDB("GameAdmin").Collection("CurrencyType").FindAll(filter, &cklist)
	if err != nil {
		return nil, err
	}

	ckMap := map[string]*comm.CurrencyType{}
	for _, c := range cklist {
		ckMap[c.CurrencyCode] = c
	}
	list = make(map[string]*comm.CurrencyType)
	for _, operator := range oplist {
		list[operator.AppID] = ckMap[operator.CurrencyKey]
	}

	return
}

func getLineOperator(ctx *Context, ps getLineOpertorreq, ret *getOperatorListResultV2) (err error) {
	var list []*comm.Operator_V2
	query := bson.M{}

	query["Name"] = ps.AppID

	pageFileter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     bson.M{"CreateTime": -1},
		Query:    query,
	}
	ret.AllCount, err = CollAdminOperator.FindPage(pageFileter, &list)
	if err != nil {
		return
	}

	ret.List, err = GetOperatorProfit(ctx, list)
	if err != nil {
		return err
	}

	return

	//for _, operatorv2 := range list {
	//	if operatorv2.Status == 0 {
	//		file := bson.M{"AppID": operatorv2.AppID}
	//		setInfo := bson.M{"Status": 1}
	//		err = CollAdminUser.Update(file, bson.M{"$set": setInfo})
	//		if err != nil {
	//			return
	//		}
	//	}
	//}
	//if err != nil {
	//	return err
	//}
}

func getOperatorInfoByAppId(ctx *Context, ps getLineOpertorreq, ret *comm.Operator_V2) (err error) {

	id, err := GetOperatopAppID(ctx)
	if err != nil {
		return err
	}

	temp := false
	for _, appid := range id {
		if ps.AppID == appid {
			temp = true
		}
	}
	if temp == false {
		return errors.New(lang.GetLang(ctx.Lang, "请选择商户"))
	}

	query := bson.M{}

	query = bson.M{"AppID": ps.AppID}

	err = CollAdminOperator.FindOne(query, &ret)
	if err != nil {
		return err
	}

	return nil
}

func operatorVerifyIp(ps operatorVerify, ret *ExistsWhitResponse) (err error) {

	query := bson.M{
		"AppID": ps.AppID,
	}

	var operator comm.Operator

	err = CollAdminOperator.FindOne(query, &operator)

	if err != nil {
		ret.Allow = false
		return nil
	}

	if operator.UserWhite == "" {
		ret.Allow = true
		return
	}

	whitIpList := strings.Split(operator.UserWhite, ",")

	flag := false
	for _, ip := range whitIpList {

		if ip == ps.IP {
			flag = true
			break
		}

	}

	ret.Allow = flag

	return nil
}

func getMerchantReviewList(ctx *Context, ps reviewListParams, ret *getReviewListResults) (err error) {
	query := bson.M{}
	user, err := GetUser(ctx)
	if err != nil {
		return
	}
	if user.GroupId > 1 {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	if err != nil {
		return
	}
	if ps.StartTime != 0 {
		startTime := time.Unix(ps.StartTime, 0)
		endTime := time.Unix(ps.EndTime, 0)
		query["CreateTime"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	filter := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     db.D("ReviewStatus", 1, "CreateTime", -1),
	}

	if ps.AppID != "" {
		query["AppID"] = ps.AppID
	}
	if ps.ParentAppID != "" {
		if ps.ParentAppID == "admin" {
			return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
		} else {
			query["Name"] = ps.ParentAppID
		}
	} else {
		query["Name"] = bson.M{"$ne": "admin"}
	}
	query["OperatorType"] = 2
	query["ReviewType"] = 0

	filter.Query = query
	ret.All, err = NewOtherDB("GameAdmin").Collection("AdminOperator").FindPage(filter, &ret.List)
	if err != nil {
		return err
	}

	return
}

func updateMerchantReviewStatus(ctx *Context, ps reviewStatusParams, ret *comm.Empty) (err error) {
	user, err := GetUser(ctx)
	if err != nil {
		return
	}
	if user.GroupId > 1 {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	if ps.OperatorID <= 0 {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	var operator *comm.Operator
	err = CollAdminOperator.FindId(ps.OperatorID, &operator)

	if err != nil {
		return err
	}

	if operator.OperatorType == 1 {
		slog.Error("updateMerchantReviewStatus", "参数错误")
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	coll := db.Collection2("GameAdmin", "AdminOperator")
	updateRetention := db.D(
		"$set", db.D(
			"ReviewStatus", 1,
			"Reviewer", user.UserName,
		))

	_, err = coll.UpdateByID(context.TODO(), ps.OperatorID, updateRetention)
	return
}

func operatorInfo(ps operatorVerify, ret *comm.Operator_V2) (err error) {

	query := bson.M{
		"AppID": ps.AppID,
	}

	err = CollAdminOperator.FindOne(query, &ret)

	return nil
}

// 修改商户余额，监听商户余额修改消息处理
func UpdateOperatorBalance(ps *comm.OperatorBalance) (err error) {
	if ps.AppID == "" {
		return errors.New("参数错误")
	}
	// 获取当前商户信息
	operator := &comm.Operator_V2{}
	if err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, operator); err != nil {
		return err
	}
	balance := operator.Balance + ps.Change*(operator.PlatformPay/100)
	if operator.CooperationType == 2 {
		balance = operator.Balance - ps.Bet*(operator.TurnoverPay/100)
	}

	err = CollAdminOperator.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": bson.M{"Balance": balance}})
	if err != nil {
		return err
	}
	dateStr := time.Now().UTC().Format("20060102")
	id := ut.JoinStr(':', dateStr, ps.AppID)
	update := db.D(
		"$set", db.D(
			"OperatorBalance", balance,
		),
	)
	coll := db.Collection2("reports", "OperatorDailyReport")
	coll.UpdateByID(context.TODO(), id, update)
	return
}

func updateExtraRtpConfig(ctx *Context, ps updateExtraRtpConfigParams, ret *comm.Empty) (err error) {
	id, err := GetOperatopAppID(ctx)
	if err != nil {
		return err
	}

	temp := false
	for _, appid := range id {
		if ps.AppID == appid {
			temp = true
		}
	}
	if temp == false {
		return errors.New(lang.GetLang(ctx.Lang, "请选择商户"))
	}
	sysconfig := &comm.SystemConfig{}
	err = CollSystemConfig.FindOne(bson.M{}, &sysconfig)
	if err != nil {
		return err
	}
	if ps.ProtectionRewardPercentLess > sysconfig.ExtraMaxRtp {
		return errors.New("超出RTP最大值")
	}
	interUpdateExtraRtpConfig(ps, ret)

	return
}
func interUpdateExtraRtpConfig(ps updateExtraRtpConfigParams, ret *comm.Empty) (err error) {
	update := db.D(
		"$set", db.D(
			"IsProtection", ps.IsProtection,
			"ProtectionRewardPercentLess", ps.ProtectionRewardPercentLess,
			"ProtectionRotateCount", ps.ProtectionRotateCount,
		),
	)
	err = CollAdminOperator.Update(bson.M{"AppID": ps.AppID}, update)
	if err != nil {
		return err
	}
	gameAllList := []*comm.Game2{}

	err, gameAllList = GetGameListFormatName("", bson.M{})

	if err != nil {
		return err
	}
	go func() {
		gameIds := []string{}
		for _, game := range gameAllList {
			gameIds = append(gameIds, game.ID)
		}
		gameList := []*comm.GameConfig{}
		CollGameConfig.FindAll(bson.M{"AppID": ps.AppID, "GameId": bson.M{"$in": gameIds}}, &gameList)
		for _, config := range gameList {
			req := updataGameCofigParams{
				AppID:              ps.AppID,
				GameId:             config.GameId,
				Preset:             config.Preset,
				StopLoss:           config.StopLoss,
				GamePattern:        config.GamePattern,
				FreeGameOff:        config.FreeGameOff,
				RTP:                config.RTP,
				BuyRTP:             config.BuyRTP,
				BuyMinAwardPercent: config.BuyMinAwardPercent,
				MaxMultiple:        strconv.Itoa(int(config.MaxMultiple)),
				MaxWinPoints:       config.MaxWinPoints,
				BetBase:            config.BetBase,
				GameOn:             config.GameOn,
				//BetMult:            config.BetMult,
				DefaultCs:       config.DefaultCs,
				DefaultBetLevel: config.DefaultBetLevel,
				ProfitMargin:    config.ProfitMargin,
				CrashRate:       config.CrashRate,
				Scale:           config.Scale,
				OnlineUpNum:     config.OnlineUpNum,
				OnlineDownNum:   config.OnlineDownNum,
			}
			sp := strings.Split(config.BetBase, ",")
			cs := make([]float64, len(sp))
			for i := 0; i < len(sp); i++ {
				cs[i], _ = strconv.ParseFloat(sp[i], 64)
			}
			rw := int64(config.RewardPercent)
			nw := int64(config.NoAwardPercent)
			PushMsgStoreSetInfo(req, rw, nw, cs)
		}
	}()
	return
}

func getExtraRtpConfig(ctx *Context, ps getExtraRtpConfigParams, ret *getExtraRtpConfigResults) (err error) {
	id, err := GetOperatopAppID(ctx)
	if err != nil {
		return err
	}

	temp := false
	for _, appid := range id {
		if ps.AppID == appid {
			temp = true
		}
	}
	if ps.AppID == "" {
		temp = true
	}
	if temp == false {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	var list []*comm.Operator_V2
	query := bson.M{}
	if ps.AppID != "" { //模糊查询
		//query["AppID"] = bson.M{"$regex": "/" + ps.AppID + "/", "$options": 'i'}
		query["AppID"] = ps.AppID
	} else {
		query["AppID"] = bson.M{"$in": id}
	}
	if ps.IsProtection != -1 {
		query["IsProtection"] = ps.IsProtection
	}
	pageFileter := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"CreateTime": -1},
		Query:    query,
	}

	ret.All, err = CollAdminOperator.FindPage(pageFileter, &list)
	if err != nil {
		return err
	}
	ret.List = list
	if err != nil {
		return err
	}

	return err
}

func SetGameMonitor(ctx *Context, ps setGameMonitorParams, ret *comm.Empty) (err error) {
	id, err := GetOperatopAppID(ctx)
	if err != nil {
		return err
	}

	temp := false
	for _, appid := range id {
		if ps.AppID == appid {
			temp = true
		}
	}
	if temp == false {
		return errors.New(lang.GetLang(ctx.Lang, "请选择商户"))
	}
	interSetGameMonitor(ps, ret)
	return
}

func interSetGameMonitor(ps setGameMonitorParams, ret *comm.Empty) (err error) {
	moniterConfig := &comm.MoniterConfig{}
	moniterConfig.IsMoniter = ps.IsMoniter
	moniterConfig.MoniterNewbieNum = ps.MoniterNewbieNum
	moniterConfig.MoniterNumCycle = ps.MoniterNumCycle
	moniterConfig.MoniterRTPErrorValue = ps.MoniterRTPErrorValue
	moniterConfig.MoniterReduceRTPRangeValue = ps.MoniterReduceRTPRangeValue
	moniterConfig.MoniterAddRTPRangeValue = ps.MoniterAddRTPRangeValue
	update := db.D(
		"$set", db.D(
			"MoniterConfig", moniterConfig,
		),
	)
	if ps.MoniterType == 1 {
		update = db.D(
			"$set", db.D(
				"PersonalMoniterConfig", moniterConfig,
			),
		)
	}
	err = CollAdminOperator.Update(bson.M{"AppID": ps.AppID}, update)
	if err != nil {
		return err
	}
	gameAllList := []*comm.Game2{}

	err, gameAllList = GetGameListFormatName("", bson.M{})

	if err != nil {
		return err
	}
	go func() {
		gameIds := []string{}
		for _, game := range gameAllList {
			gameIds = append(gameIds, game.ID)
		}
		gameList := []*comm.GameConfig{}
		CollGameConfig.FindAll(bson.M{"AppID": ps.AppID, "GameId": bson.M{"$in": gameIds}}, &gameList)
		for _, config := range gameList {
			req := updataGameCofigParams{
				AppID:              ps.AppID,
				GameId:             config.GameId,
				Preset:             config.Preset,
				StopLoss:           config.StopLoss,
				GamePattern:        config.GamePattern,
				FreeGameOff:        config.FreeGameOff,
				RTP:                config.RTP,
				BuyRTP:             config.BuyRTP,
				BuyMinAwardPercent: config.BuyMinAwardPercent,
				MaxMultiple:        strconv.Itoa(int(config.MaxMultiple)),
				MaxWinPoints:       config.MaxWinPoints,
				BetBase:            config.BetBase,
				GameOn:             config.GameOn,
				//BetMult:            config.BetMult,
				DefaultCs:       config.DefaultCs,
				DefaultBetLevel: config.DefaultBetLevel,
				ProfitMargin:    config.ProfitMargin,
				CrashRate:       config.CrashRate,
				Scale:           config.Scale,
				OnlineUpNum:     config.OnlineUpNum,
				OnlineDownNum:   config.OnlineDownNum,
			}
			sp := strings.Split(config.BetBase, ",")
			cs := make([]float64, len(sp))
			for i := 0; i < len(sp); i++ {
				cs[i], _ = strconv.ParseFloat(sp[i], 64)
			}
			rw := int64(config.RewardPercent)
			nw := int64(config.NoAwardPercent)
			PushMsgStoreSetInfo(req, rw, nw, cs)
		}
	}()
	return
}
