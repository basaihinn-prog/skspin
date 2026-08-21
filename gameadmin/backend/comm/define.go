package comm

import (
	"game/duck/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PG   = "pg"
	PP   = "pp"
	JILI = "jili"
	TADA = "tada"
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
	Id                   int64              `json:"Pid" bson:"_id"`
	Uid                  string             `json:"Uid" bson:"Uid"`
	AppID                string             `json:"AppID" bson:"AppID" md:"运营商"`
	LoginAt              *mongodb.TimeStamp `json:"LoginAt" bson:"LoginAt"`
	CreateAt             *mongodb.TimeStamp `json:"CreateAt" bson:"CreateAt"`
	TypeInfo             map[int]*TypeInfo  `json:"TypeInfo" bson:"TypeInfo"`
	Bet                  int64              `json:"Bet" bson:"Bet"`
	Win                  int64              `json:"Win" bson:"Win"`
	Multi                float64            `json:"Multi" bson:"Multi" description:"玩家当前赢取倍数"`
	Status               int                `json:"Status" bson:"Status"`       //0 启用， 1 禁用
	CloseTime            *mongodb.TimeStamp `json:"CloseTime" bson:"CloseTime"` //禁用时间
	LastSpinTime         *mongodb.TimeStamp `json:"LastSpinTime" bson:"LastSpinTime"`
	RestrictionsStatus   int64              `json:"RestrictionsStatus" bson:"RestrictionsStatus" description:"玩家约束状态"`
	RestrictionsMaxWin   int64              `json:"RestrictionsMaxWin" bson:"RestrictionsMaxWin" description:"玩家约束最大赢取金额"`
	RestrictionsMaxMulti float64            `json:"RestrictionsMaxMulti" bson:"RestrictionsMaxMulti" description:"玩家约束最大赢取倍数"`
	RestrictionsTime     *mongodb.TimeStamp `json:"RestrictionsTime" bson:"RestrictionsTime" description:"玩家约束设置时间"`
}

// 账号表
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
	OperatorAdmin bool               `json:"OperatorAdmin" bson:"OperatorAdmin"`
	PermissionId  int64              `json:"PermissionId" bson:"PermissionId"`
}

const (
	WalletModeOld = iota
	WalletModeTransfer
	WalletModeSeamless
)

//	type Operator struct {
//		Id             int64              `json:"Id" bson:"_id"`
//		Name           string             `json:"Name" bson:"Name"`
//		AppID          string             `json:"AppID" bson:"AppID" md:"运营商"`
//		CreateTime     *mongodb.TimeStamp `json:"CreateTime" bson:"CreateTime"`
//		Status         int                `json:"Status" bson:"Status"`
//		MenuIds        []int64            `json:"MenuIds" bson:"MenuIds"`
//		PermissionId   int64              `json:"PermissionId" bson:"PermissionId"`
//		AppSecret      string             `json:"AppSecret" bson:"AppSecret"`
//		Address        string             `json:"Address" bson:"Address"`
//		WhiteIps       []string           `json:"WhiteIps" bson:"WhiteIps"`
//		WalletMode     int                `json:"WalletMode"  bson:"WalletMode"`
//		PublishHistory bool               `json:"PublishHistory" bson:"PublishHistory"`
//		CurrencyKey    string             `json:"CurrencyKey" bson:"CurrencyKey"`
//		ExcluedGameId  primitive.ObjectID `json:"ExcluedGameId" bson:"ExcluedGameId"`
//		ExcluedGameIds []string           `json:"ExcluedGameIds" bson:"ExcluedGameIds"`
//		CurrentRtp     int                `json:"CurrentRtp"  bson:"CurrentRtp"`
//		Robot          string             `json:"Robot" bson:"Robot"`   //telegram 机器人
//		ChatID         string             `json:"ChatID" bson:"ChatID"` //telegram 聊天群id
//	}
type Operator struct {
	Robot      string `json:"Robot" bson:"Robot"`   //telegram 机器人
	ChatID     string `json:"ChatID" bson:"ChatID"` //telegram 聊天群id
	CurrentRtp int    `json:"CurrentRtp"  bson:"CurrentRtp"`
	Id         int64  `json:"Id" bson:"_id"`

	Name    string  `json:"Name" bson:"Name"` //上级代理商名称
	MenuIds []int64 `json:"MenuIds" bson:"MenuIds"`

	//谷歌密钥和密码
	ExcluedGameId primitive.ObjectID `json:"ExcluedGameId" bson:"ExcluedGameId"`

	//谷歌密钥和密码
	ExcluedGameIds []string `json:"ExcluedGameIds" bson:"ExcluedGameIds"`

	//权限Id
	PermissionId int64 `json:"PermissionId" bson:"PermissionId"`

	//运营商名称   AdminOperator_copy1 / AppID
	AppID string `json:"AppID" bson:"AppID" md:"运营商"`

	//运营商编码   AdminOperator_copy1 / AppSecret	暂定和主账号一样
	AppSecret string `json:"AppCode" bson:"AppSecret"`

	//运营商账号/初始账号    AdminUser / Username
	UserName string `json:"UserName" bson:"UserName"`

	//平台费
	PlatformPay float64 `json:"PlatformPay" bson:"PlatformPay"` //todo:修改小数为数

	//合作模式   收益分成：1   购分方式：2   暂定只有收益分成模式 #2024.10.12
	CooperationType int `json:"CooperationType" bson:"CooperationType"`
	//本月费率
	PresentRate float64 `json:"PresentRate" bson:"PresentRate"`
	//下月费率
	NextRate float64 `json:"NextRate" bson:"NextRate"`

	//运营商类型
	OperatorType int64 `json:"OperatorType" bson:"OperatorType"`

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

	CurrencyVisibleOff map[string]int `json:"CurrencyVisibleOff" bson:"CurrencyVisibleOff"`
}

// 玩家RTP控制
type PlayerRTPControlModel struct {
	ID                 int64     `bson:"_id"`
	AppID              string    `bson:"AppID"`  //运营商  关联运营商
	GameID             string    `bson:"GameID"` //游戏编号 关联游戏
	CreateAt           time.Time `bson:"CreateAt"`
	Pid                int64     `bson:"Pid"`                                                   //用户账号
	GameName           string    `bson:"GameName"`                                              //游戏名称
	PlayerRTP          int64     `bson:"PlayerRTP"`                                             //玩家历史RTP RTP控制之前的
	ControlRTP         float64   `bson:"ControlRTP"`                                            //账号RTP
	RewardPercent      int64     `bson:"RewardPercent"`                                         //账号RTP
	NoAwardPercent     int64     `bson:"NoAwardPercent"`                                        //账号RTP
	AutoRemoveRTP      float64   `bson:"AutoRemoveRTP"`                                         //自动解除RTP
	AutoRewardPercent  int64     `bson:"AutoRewardPercent"`                                     //自动解除RTP
	AutoNoAwardPercent int64     `bson:"AutoNoAwardPercent"`                                    //自动解除RTP
	Status             int       `bson:"Status"`                                                //状态
	FromPlant          string    `bson:"FromPlant" description:"数据来源的平台: admin:后台  demo:'试玩站'"` //状态
	BuyRTP             int       `bson:"BuyRTP"`                                                //购买游戏RTP
	BuyMinAwardPercent int       `bson:"BuyMinAwardPercent"`
	PersonWinMaxMult   int       `bson:"PersonWinMaxMult"`  // 玩家最大赢分倍数
	PersonWinMaxScore  int64     `bson:"PersonWinMaxScore"` // 玩家最大赢分
}

// 运营商表
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
	// 线路商审核状态
	ReviewStatus int64 `json:"ReviewStatus" bson:"ReviewStatus"`
	// 审核人
	Reviewer string `json:"Reviewer" bson:"Reviewer"`

	CurrencyVisibleOff            int            `bson:"-" description:"是否显示币种"`                                                       //  商户类型
	CurrencyManufactureVisibleOff map[string]int `bson:"CurrencyVisibleOff" json:"CurrencyManufactureVisibleOff" description:"是否显示币种"` //  商户类型

	HasChildren bool `json:"HasChildren"`
	// 商户余额
	Balance float64 `json:"Balance" bson:"Balance"`
	// 审核类型 0：需要审核，1：免除审核，2：不允许开户
	ReviewType int64 `json:"ReviewType" bson:"ReviewType"`
	// 余额不足阈值
	BalanceThreshold float64 `json:"BalanceThreshold" bson:"BalanceThreshold"`
	// 余额不足间隔
	BalanceThresholdInterval int64 `json:"BalanceThresholdInterval" bson:"BalanceThresholdInterval"`
	// 欠费提示间隔
	ArrearsThresholdInterval int64 `json:"ArrearsThresholdInterval" bson:"ArrearsThresholdInterval"`
	// 这个月收益
	ThisMonthProfit float64 `json:"ThisMonthProfit" bson:"-"`
	// 上月收益
	LastMonthProfit float64 `json:"LastMonthProfit" bson:"-"`
	// 币种名字
	CurrencyName string `json:"CurrencyName" bson:"-"`
	// 高rtp显示开关
	HighRTPOff int64 `json:"HighRTPOff" bson:"HighRTPOff"`
	// 保护开关
	IsProtection int `json:"IsProtection" bson:"IsProtection"`
	// 保护次数
	ProtectionRotateCount int `json:"ProtectionRotateCount" bson:"ProtectionRotateCount"`
	// 保护内获取比率
	ProtectionRewardPercentLess int `json:"ProtectionRewardPercentLess" bson:"ProtectionRewardPercentLess"`
	// 货币单位
	CurrencyCtrlStatus int `json:"CurrencyCtrlStatus" bson:"CurrencyCtrlStatus"`
	// 流水费
	TurnoverPay float64 `json:"TurnoverPay" bson:"TurnoverPay"`
	// 默认厂商开启
	DefaultManufacturerOn []string `json:"DefaultManufacturerOn" bson:"DefaultManufacturerOn"`
	// PG游戏配置
	PGConfig PGConfig `json:"PGConfig" bson:"PGConfig"`
	// JILI游戏配置
	JILIConfig JILIConfig `json:"JILIConfig" bson:"JILIConfig"`
	// PP游戏配置
	PPConfig PPConfig `json:"PPConfig" bson:"PPConfig"`
	// TADA游戏配置
	TADAConfig TADAConfig `json:"TADAConfig" bson:"TADAConfig"`
	//总控机器人
	AdminRobot string `json:"AdminRobot" bson:"AdminRobot"`
	//总控机器人的chatid
	AdminChatID string `json:"AdminChatID" bson:"AdminChatID"`
	//余额告警阈值
	BalanceAlert float64 `json:"BalanceAlert" bson:"BalanceAlert"`
	//余额告警时间 暂时固定为24每天0点进行检测并触发余额告警
	BalanceAlertTimeInterval int64 `json:"BalanceAlertTimeInterval" bson:"BalanceAlertTimeInterval"`
	// 购买Rtp开关
	BuyRTPOff int64 `json:"BuyRTPOff" bson:"BuyRTPOff"`
	// 游戏监控
	MoniterConfig MoniterConfig `json:"MoniterConfig" bson:"MoniterConfig"`
	// 个人监控
	PersonalMoniterConfig MoniterConfig `json:"PersonalMoniterConfig" bson:"PersonalMoniterConfig"`
}

type MoniterConfig struct {
	//游戏监控
	IsMoniter int64 `json:"IsMoniter" bson:"IsMoniter"`
	// 监控新手数量
	MoniterNewbieNum int64 `json:"MoniterNewbieNum" bson:"MoniterNewbieNum"`
	// 监控RTP错误值
	MoniterRTPErrorValue int64 `json:"MoniterRTPErrorValue" bson:"MoniterRTPErrorValue"`
	// 监控周期
	MoniterNumCycle int64 `json:"MoniterNumCycle" bson:"MoniterNumCycle"`
	// 监控增加RTP范围值
	MoniterAddRTPRangeValue []MoniterRtpRangeValue `json:"MoniterAddRTPRangeValue" bson:"MoniterAddRTPRangeValue"`
	// 监控减少RTP范围值
	MoniterReduceRTPRangeValue []MoniterRtpRangeValue `json:"MoniterReduceRTPRangeValue" bson:"MoniterReduceRTPRangeValue"`
}

type MoniterRtpRangeValue struct {
	RangeMinValue  int64   `json:"RangeMinValue" bson:"RangeMinValue"`
	RangeMaxValue  int64   `json:"RangeMaxValue" bson:"RangeMaxValue"`
	NewbieValue    float64 `json:"NewbieValue" bson:"NewbieValue"`
	NotNewbieValue float64 `json:"NotNewbieValue" bson:"NotNewbieValue"`
}

// 运营商游戏设置
type GameConfig struct {
	//游戏编号
	GameId string `json:"GameId" bson:"GameId"`

	//运营商
	AppID string `json:"AppID" bson:"AppID" md:"运营商"`

	//配置文件
	ConfigPath string `json:"ConfigPath" bson:"ConfigPath"`

	//RTP设置
	RTP float64 `json:"RTP" bson:"RTP"`

	//止盈止损开关
	StopLoss int64 `json:"StopLoss" bson:"StopLoss"`

	//赢取最高押注倍数
	MaxMultipleOff int64 `json:"MaxMultipleOff" bson:"MaxMultipleOff"`
	MaxMultiple    int64 `json:"MaxMultiple" bson:"MaxMultiple"`

	//游戏投注
	BetBase string `json:"BetBase" bson:"BetBase"`

	//游戏类型
	GamePattern int `json:"GamePattern" bson:"GamePattern"`

	//预设面额
	Preset float32 `json:"Preset" bson:"Preset"`

	//游戏状态  0 开启  1 关闭
	GameOn int64 `json:"GameOn" bson:"GameOn"`

	//免费游戏开关
	FreeGameOff int64 `json:"FreeGameOff" bson:"FreeGameOff"`

	//试玩链接
	GameDemo string `json:"GameDemo" bson:"GameDemo"`

	//下注的3% 个人奖池累计配置千分比 所有slots 通用,  普通, 购买不算
	RewardPercent int `json:"RewardPercent" bson:"RewardPercent"`

	//个人奖池为负数，不中奖概率千分比  20%不中奖,  优先级最高
	NoAwardPercent     int     `json:"NoAwardPercent" bson:"NoAwardPercent"`
	MaxWinPoints       float64 `json:"MaxWinPoints" bson:"MaxWinPoints"`
	BuyRTP             int     `json:"BuyRTP" bson:"BuyRTP"`
	BuyMinAwardPercent int     `json:"BuyMinAwardPercent" bson:"BuyMinAwardPercent"`
	// 投注倍数
	//BetMult float64 `json:"BetMult" bson:"BetMult"`
	// 默认cs
	DefaultCs float64 `json:"DefaultCs" bson:"DefaultCs"`
	// 默认押注倍数
	DefaultBetLevel int64   `json:"DefaultBetLevel" bson:"DefaultBetLevel"`
	OnlineUpNum     float64 `json:"OnlineUpNum" bson:"OnlineUpNum"`
	OnlineDownNum   float64 `json:"OnlineDownNum" bson:"OnlineDownNum"`
	// 利润率
	ProfitMargin float64 `json:"ProfitMargin" bson:"ProfitMargin"`
	// 坠毁概率
	CrashRate float64 `json:"CrashRate" bson:"CrashRate"`
	// 刻度
	Scale float64 `json:"Scale" bson:"Scale"`
}

// 玩家RTP控制
type PlayerRtpControl struct {
	AppID    string //运营商
	GameID   string //游戏编号
	CreateAt time.Time
	Uid      string //用户账号
	GameName string //游戏名称
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

type CurrencyType struct {
	//币种名称
	CurrencyName string `json:"CurrencyName" bson:"CurrencyName"`
	//币种代码
	CurrencyCode string `json:"CurrencyCode" bson:"CurrencyCode"`
	//币种符号
	CurrencySymbol string `json:"CurrencySymbol" bson:"CurrencySymbol"`
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
	ID               string `json:"ID" bson:"_id"`
	GameID           string `json:"GameID" bson:"GameID"`
	Name             string `json:"Name" bson:"Name"`
	Type             int    `json:"Type" bson:"Type"`
	Status           int    `json:"Status" bson:"Status"`
	ManufacturerName string `json:"ManufacturerName" bson:"ManufacturerName"`
	IconUrl          string
	WebIconUrl       string
	ExtraData        map[string]any             `bson:"-"`
	BuyBetMulti      int                        `json:"BuyBetMulti" bson:"BuyBetMulti"`
	Bet              int64                      `json:"-" bson:"Bet"`
	GameNameConfig   map[string]*GameConfigItem `json:"GameNameConfig" bson:"GameNameConfig"`
	NoAwardPercent   int64                      `json:"NoAwardPercent" bson:"NoAwardPercent"`
}

type GameConfigItem struct {
	GameName string
	Icon     string
}
type Game2 struct {
	ID               string                     `json:"ID" bson:"_id"`
	GameID           string                     `json:"GameId" bson:"GameId"`
	Name             string                     `json:"Name" bson:"Name"`
	LineNum          float64                    `json:"LineNum" bson:"LineNum"`
	Type             int                        `json:"Type" bson:"Type"`
	Status           int                        `json:"Status" bson:"Status"`
	ManufacturerName string                     `json:"ManufacturerName" bson:"ManufacturerName"`
	ExtraData        map[string]any             `bson:"-"`
	Bet              int64                      `json:"Bet" bson:"Bet"`
	BetList          string                     `json:"BetList" bson:"BetList"`
	ChangeBetOff     int64                      `json:"ChangeBetOff" bson:"ChangeBetOff"`
	GameNameConfig   map[string]*GameConfigItem `json:"GameNameConfig" bson:"GameNameConfig"`
	NoAwardPercent   int64                      `json:"NoAwardPercent" bson:"NoAwardPercent"`
	BuyBetMulti      int                        `json:"BuyBetMulti" bson:"BuyBetMulti"`
	OriginName       string                     `json:"OriginName" bson:"-"`
	DefaultBet       float64                    `json:"DefaultBet" bson:"DefaultBet"`
	DefaultBetLevel  int64                      `json:"DefaultBetLevel" bson:"DefaultBetLevel"`
	RewardPercent    int64                      `json:"RewardPercent" bson:"RewardPercent"`
	IconUrl          string
	WebIconUrl       string
	BuyType          int `json:"BuyType" bson:"BuyType"`
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

type PokerOperatorReportData struct { //扑克平台汇总
	ID             string `json:"-" bson:"_id"`
	Date           string `bson:"date"`
	AppID          string `bson:"appid"`          // 渠道
	BetAmount      int64  `bson:"betamount"`      // 支出
	WinAmount      int64  `bson:"winamount"`      // 赚取
	SpinCount      int    `bson:"spincount"`      // 总游戏次数
	EnterPlrCount  int    `bson:"enterplrcount"`  //老玩家进入次数
	RegistPlrCount int    `bson:"registplrcount"` //新玩家进入次数
	Flow           int    `bson:"flow"`           //流水(玩家输赢的绝对值之和)
}

// todo：奖池
type SlotsPoolHistory struct { // 奖池(未转移金额、Slot奖池和百人奖池)
	ID           string             `bson:"_id"`
	OpName       string             `bson:"OpName"`       //操作人
	OpPid        int64              `bson:"OpPid"`        //操作人id
	AnimUserPid  int64              `bson:"AnimUserPid"`  //玩家id(唯一标识)
	AnimUserName string             `bson:"AnimUserName"` //玩家名称(玩家ID)
	Type         int64              `bson:"Type"`         // 1:Slot奖池,2:百人奖池,3:未转移金额,4:玩家状态
	EnsureStatus int                `bson:"EnsureStatus"` //0:未审核,1已审核
	OldGold      int64              `bson:"OldGold"`      //修改前
	NewGold      int64              `bson:"NewGold"`      //修改后
	Change       int64              `bson:"Change"`       //变化值
	Time         *mongodb.TimeStamp `bson:"time"`         //变化时间
	AppID        string             `bson:"AppID"`        //运营商
	Currency     string             `bson:"Currency"`     //币种
	Remark       string             `bson:"Remark"`       //审核备注
	EnsureOpName string             `bson:"EnsureOpName"` //审核人
	EnsureOpPid  int64              `bson:"EnsureOpPid"`  //审核id
	EnsureTime   *mongodb.TimeStamp `bson:"EnsureTime"`   //审核时间
}

type SlotWinLoseLimit struct { // 玩家slots净输赢记录
	ID           string     `bson:"_id"`
	Pid          int64      `bson:"ID"`           //玩家id(唯一标识)
	AppID        string     `bson:"AppID"`        //运营商
	Currency     string     `bson:"Currency"`     //币种
	GameID       string     `bson:"GameID"`       //游戏ID
	WinLose      int64      `bson:"WinLose"`      //净输赢
	OriginWin    int64      `bson:"OriginWin"`    //原本的输赢
	RealWin      int64      `bson:"RealWin"`      //实际的输赢
	Desc         string     `bson:"Desc"`         //详细说明
	EnsureStatus int        `bson:"EnsureStatus"` //0:未审核,1已审核
	Time         time.Time  `bson:"Time"`         //记录时间
	Remark       *string    `bson:"Remark"`       //审核备注
	EnsureOpName *string    `bson:"EnsureOpName"` //审核人
	EnsureOpPid  *int64     `bson:"EnsureOpPid"`  //审核id
	EnsureTime   *time.Time `bson:"EnsureTime"`   //审核时间
}

type RTPGear struct {
	ID  int64 `bson:"_id"`
	RTP int64 `bson:"RTP"`
}

type BetDailyGGR struct {
	ID            string `bson:"_id"`
	AppID         string `bson:"appid"`
	Date          string `bson:"date"`
	Enterplrcount int    `bson:"enterplrcount"`
	GameID        string `bson:"game"`
	BetAmount     int64  `bson:"betamount"` // 支出
	Bigreward     int64  `bson:"bigreward"`
	BuyBetAmount  int64  `bson:"buybetamount"`
	ButWinamount  int64  `bson:"butwinamount"`
	SpinCount     int    `bson:"spincount"` // 转动次数
	Spinplrcount  int    `bson:"spinplrcount"`
	WinAmount     int64  `bson:"winamount"` // 赚取

	GGR         float64 `bson:"GGR"`
	PresentRate float64 `bson:"PresentRate"`
}

type OperatorDailyGGR struct {
	ID             string `bson:"_id"`
	AppID          string `bson:"appid"`
	Date           string `bson:"date"`
	Enterplrcount  int    `bson:"enterplrcount"`
	LoginPlrCount  int    `bson:"loginplrcount"`
	RegistPlrCount int    `bson:"registplrcount"`
	BetAmount      int64  `bson:"betamount"` // 支出
	Bigreward      int64  `bson:"bigreward"`
	BuyBetAmount   int64  `bson:"buybetamount"`
	ButWinamount   int64  `bson:"butwinamount"`
	SpinCount      int    `bson:"spincount"` // 转动次数
	Spinplrcount   int    `bson:"spinplrcount"`
	WinAmount      int64  `bson:"winamount"` // 赚取

	GGR             float64 `bson:"GGR"`
	PresentRate     float64 `bson:"PresentRate"`
	OperatorBalance float64 `bson:"OperatorBalance"`
}

type OperatorMonthlyGGR struct {
	ID        string  `bson:"ID"`
	AppID     string  `bson:"appid"`
	Date      string  `bson:"date"`
	BetAmount string  `bson:"betamount"`
	WinAmount string  `bson:"winamount"`
	PlantRate float64 `bson:"plantrate"`
	GGR       int64   `bson:"GGR"`
}

type PlayerRetentionReport struct {
	ID                   string  `bson:"_id"`
	RetentionPlayerCount int64   `bson:"RetentionPlayerCount"`
	RetentionPlayer1d    float64 `bson:"RetentionPlayer1d"`  // 次日留存率
	RetentionPlayer3d    float64 `bson:"RetentionPlayer3d"`  // 3日留存率
	RetentionPlayer7d    float64 `bson:"RetentionPlayer7d"`  // 7日留存率
	RetentionPlayer14d   float64 `bson:"RetentionPlayer14d"` // 14日留存率

	RetentionPlayer30d float64            `bson:"RetentionPlayer30d"` //30日留存率
	Date               *mongodb.TimeStamp `bson:"date"`
	AppID              string             `bson:"appid"`
}

type PlayerRetentionGameReport struct {
	ID                   string  `bson:"_id"`
	RetentionPlayerCount int64   `bson:"RetentionPlayerCount"`
	RetentionPlayer1d    float64 `bson:"RetentionPlayer1d"`  // 次日留存率
	RetentionPlayer3d    float64 `bson:"RetentionPlayer3d"`  // 3日留存率
	RetentionPlayer7d    float64 `bson:"RetentionPlayer7d"`  // 7日留存率
	RetentionPlayer14d   float64 `bson:"RetentionPlayer14d"` // 14日留存率

	RetentionPlayer30d float64            `bson:"RetentionPlayer30d"` //30日留存率
	Date               *mongodb.TimeStamp `bson:"date"`
	GameID             string             `bson:"gameId"`
}

type RetentionReport struct {
	ID                   string  `json:"-" bson:"_id"`
	RetentionPlayerCount int64   `bson:"RetentionPlayerCount"`
	RetentionPlayer1d    float64 `bson:"RetentionPlayer1d"`  // 次日留存率
	RetentionPlayer3d    float64 `bson:"RetentionPlayer3d"`  // 3日留存率
	RetentionPlayer7d    float64 `bson:"RetentionPlayer7d"`  // 7日留存率
	RetentionPlayer14d   float64 `bson:"RetentionPlayer14d"` // 14日留存率

	RetentionPlayer30d float64            `bson:"RetentionPlayer30d"` //30日留存率
	AppID              string             `bson:"appid"`
	Date               *mongodb.TimeStamp `bson:"date"`
	GameID             string             `bson:"gameId"`
}

type AdminLog struct {
	ID             primitive.ObjectID `bson:"_id"`
	OperatorID     int64              `bson:"operatorid"`     //操作人id
	OperatorName   string             `bson:"operatorname"`   //操作人名称
	OperatorType   int64              `bson:"operatortype"`   //操作人类型
	OperateType    int64              `bson:"operatetype"`    //操作类型
	OperateContent string             `bson:"operatecontent"` //操作内容
	OperateTime    *mongodb.TimeStamp `bson:"operatetime"`    //操作时间
	IpAdress       string             `bson:"ipadress"`       //操作ip
	OperateMod     int64              `bson:"operatemod"`     //操作模块
	AppID          string             `bson:"appid"`          //运营商
}

// 修改商户余额
type OperatorBalance struct {
	AppID  string
	Change float64
	Bet    float64
}

// pp游戏通用配置
type PPConfig struct {
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff"`
	//CurrencyVisibleOff int64 `json:"CurrencyVisibleOff"`
}

// pg游戏通用配置
type PGConfig struct {
	CarouselOff        int64 `json:"CarouselOff"`
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff"`
	//CurrencyVisibleOff int64  `json:"CurrencyVisibleOff"`
	StopLoss       int64  `json:"StopLoss"`
	ExitBtnOff     int64  `json:"ExitBtnOff"`
	ExitLink       string `json:"ExitLink"`
	OfficialVerify int64  `json:"OfficialVerify"`
}

// jili游戏通用配置
type JILIConfig struct {
	BackPackOff int64 `json:"BackPackOff"`
	// 开屏弹窗开关
	OpenScreenOff int64 `json:"OpenScreenOff"`
	// 侧边栏开关
	SidebarOff int64 `json:"SidebarOff"`
	//CurrencyVisibleOff int64 `json:"CurrencyVisibleOff"`
}

// TADA游戏通用配置
type TADAConfig struct {
	BackPackOff int64 `json:"BackPackOff"`
	// 开屏弹窗开关
	OpenScreenOff int64 `json:"OpenScreenOff"`
	// 侧边栏开关
	SidebarOff int64 `json:"SidebarOff"`
	//CurrencyVisibleOff int64 `json:"CurrencyVisibleOff"`
}

// 余额告警参数
type BalanceAlertPs struct {
	AppID           string
	CooperationType int
	WalletType      int
	MerchantBalance float64
}
