package comm

import (
	"game/duck/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TypeInfo struct {
	SpinCnt    int64 `json:"SpinCnt" bson:"SpinCnt"`
	Bet        int64 `json:"Bet" bson:"Bet"`
	Win        int64 `json:"Win" bson:"Win"`
	BuyGame    int64 `json:"BuyGame" bson:"BuyGame"`
	BuyGameBet int64 `json:"BuyGameBet" bson:"BuyGameBet"`
	BuyGameWin int64 `json:"BuyGameWin" bson:"BuyGameWin"`
	AvgBet     int64 `json:"AvgBet" bson:"-"`
}

type Player struct {
	Id       int64              `json:"Pid" bson:"_id"`
	Uid      string             `json:"Uid" bson:"Uid"`
	AppID    string             `json:"AppID" bson:"AppID" md:"运营商"`
	LoginAt  *mongodb.TimeStamp `json:"LoginAt" bson:"LoginAt"`
	CreateAt *mongodb.TimeStamp `json:"CreateAt" bson:"CreateAt"`
	TypeInfo map[int]*TypeInfo  `json:"TypeInfo" bson:"TypeInfo"`
	Bet      int64              `json:"Bet" bson:"Bet"`
	Win      int64              `json:"Win" bson:"Win"`
}

type User struct {
	Id            int64              `json:"Id" bson:"_id"`
	AppID         string             `json:"AppID" bson:"AppID" md:"运营商"`
	UserName      string             `json:"UserName" bson:"Username"`
	PassWord      string             `json:"-" bson:"Password"`
	GoogleCode    string             `json:"GoogleCode" bson:"GoogleCode"`
	Qrcode        string             `json:"Qrcode" bson:"Qrcode"`
	Avatar        string             `json:"Avatar" bson:"Avatar"`
	IsOpenGoogle  bool               `json:"IsOpenGoogle" bson:"IsOpenGoogle"`
	Status        int                `json:"Status" bson:"Status"`
	CreateAt      *mongodb.TimeStamp `json:"CreateAt" bson:"CreateAt"`
	LoginAt       *mongodb.TimeStamp `json:"Login_At" bson:"LoginAt"`
	GroupId       int64              `json:"GroupId" bson:"GroupId"`
	Token         string             `json:"Token" bson:"Token"`
	TokenExpireAt *mongodb.TimeStamp `json:"TokenExpireAt" bson:"TokenExpireAt"`
	OperatorAdmin bool               `json:"-" bson:"OperatorAdmin"`
	PermissionId  int64              `json:"PermissionId" bson:"PermissionId"`
}

const (
	WalletModeOld = iota
	WalletModeTransfer
	WalletModeSeamless
)

type Operator struct {
	Id                  int64              `json:"Id" bson:"_id"`
	Name                string             `json:"Name" bson:"Name"`
	AppID               string             `json:"AppID" bson:"AppID" md:"运营商"`
	CreateTime          *mongodb.TimeStamp `json:"CreateTime" bson:"CreateTime"`
	Status              int                `json:"Status" bson:"Status"`
	MenuIds             []int64            `json:"MenuIds" bson:"MenuIds"`
	PermissionId        int64              `json:"PermissionId" bson:"PermissionId"`
	AppSecret           string             `json:"AppSecret" bson:"AppSecret"`
	Address             string             `json:"Address" bson:"Address"`
	WhiteIps            []string           `json:"WhiteIps" bson:"WhiteIps"`
	WalletMode          int                `json:"WalletMode"  bson:"WalletMode"`
	PublishHistory      bool               `json:"PublishHistory" bson:"PublishHistory"`
	CurrencyKey         string             `json:"CurrencyKey" bson:"CurrencyKey"`
	ExcluedGameId       primitive.ObjectID `json:"ExcluedGameId" bson:"ExcluedGameId"`
	ExcluedGameIds      []string           `json:"ExcluedGameIds" bson:"ExcluedGameIds"`
	OperatorType        int                `json:"OperatorType" bson:"OperatorType"`
	ReviewStatus        int64              `json:"ReviewStatus" bson:"ReviewStatus"`
	CurrencyCtrlStatus  int                `json:"CurrencyCtrlStatus" bson:"CurrencyCtrlStatus"`
	RTPOff              int64              `json:"RTPOff" bson:"RTPOff"`
	PlayerRTPSettingOff int64              `json:"PlayerRTPSettingOff" bson:"PlayerRTPSettingOff"`
	HighRTPOff          int64              `json:"HighRTPOff" bson:"HighRTPOff"`
}

type Permission struct {
	Id         int64              `json:"Id" bson:"_id"`
	OperatorId int64              `json:"OperatorId" bson:"OperatorId"`
	Name       string             `json:"Name" bson:"Name"`
	Remark     string             `json:"Remark" bson:"Remark"`
	MenuIds    []int64            `json:"MenuIds" bson:"MenuIds" md:"子菜单"`
	CreateTime *mongodb.TimeStamp `json:"CreateTime" bson:"CreateTime"`
	UpdateTime *mongodb.TimeStamp `json:"UpdateTime" bson:"UpdateTime"`
}

type PerMenu struct {
	ID      int64  `json:"ID" bson:"_id"`
	Title   string `json:"Title" bson:"Title"`
	THTitle string `json:"THTitle" bson:"THTitle"`
	Pid     int64  `json:"Pid" bson:"Pid"`
	Url     string `json:"Url" bson:"Url"`
	Icon    string `json:"Icon" bson:"Icon"`
	Sort    int64  `json:"Sort" bson:"Sort"`
}

type SlotsPool struct {
	ID   string `json:"ID" bson:"_id"`
	Pid  int64  `json:"Pid" bson:"Pid"`
	Type int64  `json:"Type" bson:"Type"`
	Gold int64  `json:"Gold" bson:"Gold"`
}

type Game struct {
	ID          string `json:"ID" bson:"_id"`
	Name        string `json:"Name" bson:"Name"`
	Type        int    `json:"Type" bson:"Type"`
	Status      int    `json:"Status" bson:"Status"`
	IconUrl     string
	WebIconUrl  string
	ExtraData   map[string]any `bson:"-"`
	BuyBetMulti int            `json:"BuyBetMulti" bson:"BuyBetMulti"`
	Bet         int64          `json:"-" bson:"-"`
}

type GameNew struct {
	ID               string `json:"ID" bson:"_id"`
	GameID           string `json:"GameID" bson:"GameID"`
	Name             string `json:"Name" bson:"Name"`
	Type             int    `json:"Type" bson:"Type"`
	Status           int    `json:"Status" bson:"Status"`
	ManufacturerName string `json:"ManufacturerName" bson:"ManufacturerName"`
	IconUrl          string
	WebIconUrl       string
	ExtraData        map[string]any `bson:"-"`
	BuyBetMulti      int            `json:"BuyBetMulti" bson:"BuyBetMulti"`
	Bet              int64          `json:"-" bson:"Bet"`
	NoAwardPercent   int            `json:"NoAwardPercent" bson:"NoAwardPercent"`
}

type Empty struct {
}

type GameLoginDetail struct {
	ID        primitive.ObjectID `json:"ID" bson:"_id"`
	Pid       string             `json:"Pid" bson:"Pid"`       // 渠道
	UserID    string             `json:"UserID" bson:"UserID"` //玩家账号
	GameID    string             `json:"GameID" bson:"GameID"` //游戏名
	Ip        string             `json:"Ip" bson:"Ip"`         //登录ip
	Loc       string             `json:"Loc" bson:"Loc"`       //ip区域码
	LoginTime int64              `json:"LoginTime" bson:"LoginTime"`
}

type SystemConfig struct {
	ID          string `json:"ID" bson:"AppId"`
	Maintenance int    `json:"Maintenance" bson:"Maintenance"`
	StartTime   int64  `json:"StartTime" bson:"StartTime"`
	EndTime     int64  `json:"EndTime" bson:"EndTime"`
	ExtraMaxRtp int    `json:"ExtraMaxRtp" bson:"ExtraMaxRtp"`
}

type NextMult struct {
	MinMult float64 `json:"MinMult"`
	MaxMult float64 `json:"MaxMult"`
}

type Operator_V2 struct {
	Id int64 `json:"Id" bson:"_id"`

	Name    string  `json:"Name" bson:"Name"` //上级代理商名称
	MenuIds []int64 `json:"MenuIds" bson:"MenuIds"`

	//谷歌密钥和密码
	//ExcluedGameId primitive.ObjectID `json:"ExcluedGameId" bson:"ExcluedGameId"`

	//谷歌密钥和密码
	ExcluedGameIds []string `json:"ExcluedGameIds" bson:"ExcluedGameIds" description:"谷歌密钥和密码"`

	//权限Id
	PermissionId int64 `json:"PermissionId" bson:"PermissionId" description:"权限Id"`

	//运营商名称   AdminOperator_copy1 / AppID
	AppID string `json:"AppID" bson:"AppID" md:"运营商"`

	//运营商编码   AdminOperator_copy1 / AppSecret	暂定和主账号一样
	AppSecret string `json:"AppCode" bson:"AppSecret" description:"运营商编码"`

	//运营商账号/初始账号    AdminUser / Username
	UserName string `json:"UserName" bson:"UserName" description:"运营商账号"`

	//平台费
	PlatformPay float64 `json:"PlatformPay" bson:"PlatformPay" description:"平台费"`

	//合作模式   收益分成：1   购分方式：2   暂定只有收益分成模式 #2024.10.12
	CooperationType int `json:"CooperationType" bson:"CooperationType" description:"合作模式  收益分成：1   购分方式：2   暂定只有收益分成模式 #2024.10.12"`
	//本月费率
	PresentRate float64 `json:"PresentRate" bson:"PresentRate" description:"本月费率"`
	//下月费率
	NextRate float64 `json:"NextRate" bson:"NextRate"  description:"下月费率"`

	//运营商类型
	OperatorType int64 `json:"OperatorType" bson:"OperatorType"  description:"运营商类型"`

	//预付款金额
	Advance float32 `json:"Advance" bson:"Advance"`

	//运营商状态   AdminOperator_copy1 / Status
	Status int64 `json:"Status" bson:"Status"`

	//创建时间     AdminOperator_copy1 /CreateTime
	CreateTime *mongodb.TimeStamp `json:"CreateTime" bson:"CreateTime"`

	//最后登录时间 Adminuser / TokenExpireAt	todo：最后登录时间在登录接口中写入
	TokenExpireAt *mongodb.TimeStamp `json:"TokenExpireAt" bson:"TokenExpireAt"`

	//商户币种        AdminOperator_copy1 / CurrencyKey
	CurrencyKey string `json:"CurrencyKey" bson:"CurrencyKey"`

	//联系方式           为多条数据
	Contact string `json:"Contact" bson:"Contact"`

	//下面的是线路商没有,运营商有的字段
	//钱包类型        AdminOperator_copy1 / WalletMode
	WalletMode int `json:"WalletMode" bson:"WalletMode"`

	//会员前缀
	Surname string `json:"Surname" bson:"Surname"`

	//客户端默认语言
	Lang string `json:"Lang" bson:"Lang"`

	//服务器地址
	ServiceIp string `json:"ServiceIp" bson:"ServiceIp"`

	//服务器IP白名单  AdminOperator_copy1 / WhiteIps
	WhiteIps []string `json:"WhiteIps" bson:"WhiteIps"`

	//服务器回调地址  AdminOperator_copy1 / Address
	Address string `json:"Address" bson:"Address"`

	//用户白名单
	UserWhite string `json:"UserWhite" bson:"UserWhite"`

	//记录开启
	LoginOff int64 `json:"LoginOff" bson:"LoginOff"`

	//免游游戏开关板
	FreeOff int64 `json:"FreeOff" bson:"FreeOff"`

	//防休眠开启
	DormancyOff int64 `json:"DormancyOff" bson:"DormancyOff"`

	//断复功能开关
	RestoreOff int64 `json:"RestoreOff" bson:"RestoreOff"`

	//手动全屏开启
	ManualFullScreenOff int64 `json:"ManualFullScreen" bson:"ManualFullScreenOff"`

	//游戏新型预设
	NewGameDefaulOff int64 `json:"NewGameDefaul" bson:"NewGameDefaul"`

	//消息推送
	MassageOff int64 `json:"MassageOff" bson:"MassageOff"`
	//消息推送地址
	MassageIp string `json:"MassageIp" bson:"MassageIp"`

	//RTP设置功能开关
	RTPOff int64 `json:"RTPOff" bson:"RTPOff"`

	//止损止盈功能开关
	StopLoss int64 `json:"StopLoss" bson:"StopLoss"`

	//赢取最高押注倍数功能开关  	只是一个开关  具体数值是针对某个游戏的，在游戏配置数据库中
	MaxMultipleOff int64 `json:"MaxMultipleOff" bson:"MaxMultipleOff"`

	//用来区分线路商和总控 程序运行需要
	LineMerchant int64 `json:"_LineMerchant" bson:"_LineMerchant"`

	PlayerRTPSettingOff int64 `json:"PlayerRTPSettingOff" bson:"PlayerRTPSettingOff"`

	//显示游戏名称和时间
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff" bson:"ShowNameAndTimeOff"`

	//退出按钮显示
	ShowExitBtnOff int64  `json:"ShowExitBtnOff" bson:"ShowExitBtnOff"`
	ExitLink       string `json:"ExitLink" bson:"ExitLink"`

	// 二期需求添加的字段
	Remark           string `bson:"Remark" description:"备注"`
	BelongingCountry string `bson:"BelongingCountry" description:"国家"`
	MaxWinPointsOff  int64  `bson:"MaxWinPointsOff" description:"最高赢钱数"`

	Robot  string `bson:"Robot" description:"Telegram机器人信息"`
	ChatID string `bson:"ChatID" description:"Telegram机器人信息"`
	// 购买Rtp开关
	BuyRTPOff int64 `json:"BuyRTPOff" bson:"BuyRTPOff"`
	// 商户余额
	Balance float64 `json:"Balance" bson:"Balance"`
}

type MoniterRtpRangeValue struct {
	RangeMinValue  int64 `json:"RangeMinValue" bson:"RangeMinValue"`
	RangeMaxValue  int64 `json:"RangeMaxValue" bson:"RangeMaxValue"`
	NewbieValue    int64 `json:"NewbieValue" bson:"NewbieValue"`
	NotNewbieValue int64 `json:"NotNewbieValue" bson:"NotNewbieValue"`
}
