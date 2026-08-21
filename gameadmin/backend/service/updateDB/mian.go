package main

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/duck/lazy"
	"game/duck/logger"
	"game/duck/mongodb"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
	"time"
)

var DB = mongodb.NewDB("GameAdmin")
var CollAdminOperator = DB.Collection("AdminOperator")
var CollGameConfig = DB.Collection("GameConfig")

func NewOtherDB(dbName string) *mongodb.DB {
	return &mongodb.DB{Name: dbName, Client: DB.Client}
}

type getAppGameInfoData struct {
	AppId        string  `bson:"appid"`
	OperatorType int     `bson:"operator_type" json:"Type"`
	Bet          int64   `bson:"betamount"`
	PlantRate    float64 `bson:"plantrate"`
	Profit       float64 `bson:"-"`
	Win          int64   `bson:"winamount"`
	Date         string  `bson:"date"`
	GGR          float64 `bson:"GGR"`
	InCome       float64 `bson:"InCome"`
}

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

func getHistoryMonthlyReport() ([]*getAppGameInfoData, error) {

	var gameInfoData []*getAppGameInfoData
	filter := bson.M{}
	err := NewOtherDB("reports").Collection("OperatorMonthlyGGR").FindAll(filter, &gameInfoData)
	if err != nil {
		return nil, err
	}
	return gameInfoData, err
}

// 更新收益
func updateProfit() {
	data, err := getHistoryMonthlyReport()
	if err != nil {
		return
	}
	oper := []*comm.Operator_V2{}
	CollAdminOperator.FindAll(bson.M{}, &oper)
	mapRale := make(map[string]float64)
	mapParent := make(map[string]string)
	for _, op := range oper {
		mapRale[op.AppID] = op.PlatformPay
		mapParent[op.AppID] = op.Name
	}
	for _, v := range data {
		rate := mapRale[v.AppId]

		InCome := float64(0.0)
		if mapParent[v.AppId] != "admin" {
			parentRate := mapRale[mapParent[v.AppId]]
			InCome = float64(v.Bet-v.Win) * ((rate - parentRate) * 0.01)
		}
		filter := bson.M{
			"appid": v.AppId,
			"date":  v.Date,
		}

		update := bson.M{
			"$set": bson.M{
				"InCome": InCome,
			},
		}
		coll := NewOtherDB("reports").Collection("OperatorMonthlyGGR").Coll()
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			slog.Info("err", err)
		}
	}
}

// 添加每日汇总余额
func addDailyBalance() {
	coll := NewOtherDB("reports").Collection("OperatorDailyGGR").Coll()
	_, err := coll.UpdateMany(context.TODO(), bson.M{}, bson.M{"$set": bson.M{"OperatorBalance": 5000}})
	if err != nil {
		slog.Info("err", err)
	}
	date := time.Now().Format("20060102")
	coll1 := NewOtherDB("reports").Collection("OperatorDailyReport").Coll()
	_, err = coll1.UpdateMany(context.TODO(), bson.M{"date": date}, bson.M{"$set": bson.M{"OperatorBalance": 5000}})
	if err != nil {
		slog.Info("err", err)
	}
}

// 添加商户表需要添加的默认数据
func addOperatorDefaultData() {
	coll := NewOtherDB("GameAdmin").Collection("AdminOperator").Coll()
	_, err := coll.UpdateMany(context.TODO(), bson.M{}, bson.M{"$set": bson.M{"ReviewStatus": 1,
		"Reviewer":                 "System",
		"Balance":                  5000,
		"ReviewType":               1,
		"BalanceThreshold":         1500,
		"BalanceThresholdInterval": 1,
		"ArrearsThresholdInterval": 3600000,
		"Lang":                     "zh"}})
	if err != nil {
		slog.Info("err", err)
	}
}

func updataGameUniversalConfig() (err error) {

	var result comm.Operator_V2
	pgConfig := comm.PGConfig{
		CarouselOff:        1,
		ShowNameAndTimeOff: 1,
		ExitLink:           "",
		ExitBtnOff:         1,
		StopLoss:           1,
	}
	ppConfig := comm.PPConfig{
		ShowNameAndTimeOff: 1,
	}
	jiliConfig := comm.JILIConfig{
		BackPackOff:   1,
		OpenScreenOff: 1,
		SidebarOff:    1,
	}
	tadaConfig := comm.TADAConfig{
		BackPackOff:   1,
		OpenScreenOff: 1,
		SidebarOff:    1,
	}

	settingPPInfo := bson.M{}
	settingPGInfo := bson.M{}
	settingJILIInfo := bson.M{}
	settingTADAInfo := bson.M{}
	settingPPInfo["ShowNameAndTimeOff"] = ppConfig.ShowNameAndTimeOff
	//settingPPInfo["CurrencyVisibleOff"] = ps.PPConfig.CurrencyVisibleOff
	settingPGInfo["CarouselOff"] = pgConfig.CarouselOff
	settingPGInfo["ShowNameAndTimeOff"] = pgConfig.ShowNameAndTimeOff
	//settingPGInfo["CurrencyVisibleOff"] = ps.PGConfig.CurrencyVisibleOff
	//settingPGInfo["StopLoss"] = pgConfig.StopLoss
	settingPGInfo["ExitBtnOff"] = pgConfig.ExitBtnOff
	settingPGInfo["ExitLink"] = pgConfig.ExitLink
	settingJILIInfo["BackPackOff"] = jiliConfig.BackPackOff
	settingJILIInfo["OpenScreenOff"] = jiliConfig.OpenScreenOff
	settingJILIInfo["SidebarOff"] = jiliConfig.SidebarOff
	//settingJILIInfo["CurrencyVisibleOff"] = ps.JILIConfig.CurrencyVisibleOff
	settingTADAInfo["BackPackOff"] = tadaConfig.BackPackOff
	settingTADAInfo["OpenScreenOff"] = tadaConfig.OpenScreenOff
	settingTADAInfo["SidebarOff"] = tadaConfig.SidebarOff
	//settingTADAInfo["CurrencyVisibleOff"] = ps.TADAConfig.CurrencyVisibleOff

	setPP := bson.M{
		"$set": settingPPInfo,
	}
	setPG := bson.M{
		"$set": settingPGInfo,
	}
	setJILI := bson.M{
		"$set": settingJILIInfo,
	}
	setTADA := bson.M{
		"$set": settingTADAInfo,
	}

	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^pp"}}, setPP)
	if err != nil {
		return err
	}
	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^pg"}}, setPG)
	if err != nil {
		return err
	}
	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^jili"}}, setJILI)
	if err != nil {
		return err
	}
	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^tada"}}, setTADA)
	if err != nil {
		return err
	}
	//result.CurrencyManufactureVisibleOff["pg"] = int(ps.PGConfig.CurrencyVisibleOff)
	//result.CurrencyManufactureVisibleOff["pp"] = int(ps.PPConfig.CurrencyVisibleOff)
	//result.CurrencyManufactureVisibleOff["jili"] = int(ps.JILIConfig.CurrencyVisibleOff)
	//result.CurrencyManufactureVisibleOff["tada"] = int(ps.TADAConfig.CurrencyVisibleOff)
	result.PGConfig = pgConfig
	result.PPConfig = ppConfig
	result.JILIConfig = jiliConfig
	result.TADAConfig = tadaConfig
	settingInfo := bson.M{}
	//settingInfo["CurrencyVisibleOff"] = result.CurrencyManufactureVisibleOff
	settingInfo["PGConfig"] = result.PGConfig
	settingInfo["PPConfig"] = result.PPConfig
	settingInfo["JILIConfig"] = result.JILIConfig
	settingInfo["TADAConfig"] = result.TADAConfig
	update := bson.M{
		"$set": settingInfo,
	}
	CollAdminOperator.Update(bson.M{}, update)

	return
}

// 更新PG官方验证配置,和默认余额告警配置
func updataPGVerifyConfig() (err error) {
	settingPGInfo := bson.M{}
	settingPGInfo["OfficialVerify"] = 1

	setPG := bson.M{
		"$set": settingPGInfo,
	}
	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^pg"}}, setPG)
	if err != nil {
		return err
	}
	settingInfo := bson.M{}

	settingInfo["PGConfig.OfficialVerify"] = 1
	settingInfo["BalanceAlertTimeInterval"] = 24
	settingInfo["BalanceAlert"] = 0
	update := bson.M{
		"$set": settingInfo,
	}
	CollAdminOperator.Update(bson.M{}, update)

	return
}

// 更新游戏类型,旧数据统一修改为 3:仿正版
func updataConfigGameType() (err error) {
	settingInfo := bson.M{}
	settingInfo["GamePattern"] = 3

	set := bson.M{
		"$set": settingInfo,
	}
	err = CollGameConfig.Update(bson.M{}, set)
	if err != nil {
		return err
	}

	return
}

// 更新商户表更新购买rtp开关为开启
func updataBuyRtpOff() (err error) {
	settingInfo := bson.M{}
	settingInfo["BuyRTPOff"] = 1

	set := bson.M{
		"$set": settingInfo,
	}
	err = CollAdminOperator.Update(bson.M{}, set)
	if err != nil {
		return err
	}

	return
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

func addMonitorDefaultData() (err error) {
	settingInfo := bson.M{}
	settingInfo["MoniterConfig"] = defaultMoniter
	settingInfo["PersonalMoniterConfig"] = defaultMoniter

	set := bson.M{
		"$set": settingInfo,
	}
	err = CollAdminOperator.Update(bson.M{}, set)
	if err != nil {
		return err
	}
	return
}

func addSpribe1DefaultData() (err error) {
	settingInfo := bson.M{}
	settingInfo["ProfitMargin"] = 3.00 // 利润率
	settingInfo["CrashRate"] = 1.00    //坠毁概率
	settingInfo["Scale"] = 3.00        // 刻度

	set := bson.M{
		"$set": settingInfo,
	}
	err = CollGameConfig.Update(bson.M{}, set)
	if err != nil {
		return err
	}
	return
}

func updateUpDown() (err error) {
	settingInfo := bson.M{}
	settingInfo["OnlineUpNum"] = 1
	settingInfo["OnlineDownNum"] = 500

	set := bson.M{
		"$set": settingInfo,
	}
	err = CollGameConfig.Update(bson.M{"GameId": "spribe_01"}, set)
	if err != nil {
		return err
	}
	return
}

func main() {
	lazy.Init("admin")
	ConnectMongoByConfig()
	// 1.0
	// 添加收益
	//slog.Info("更新收益开始")
	//updateProfit()
	//slog.Info("更新收益结束")
	//// 添加每日汇总余额
	//slog.Info("添加每日汇总余额开始")
	//addDailyBalance()
	//slog.Info("添加每日汇总余额结束")
	//// 添加商户表需要添加的默认数据
	//slog.Info("添加商户表需要添加的默认数据开始")
	//addOperatorDefaultData()
	//slog.Info("添加商户表需要添加的默认数据结束")
	// 2.0 添加商户表需要添加的通用配置
	//slog.Info("添加商户表需要添加的通用配置开始")
	//updataGameUniversalConfig()
	//slog.Info("添加商户表需要添加的通用配置结束")
	// 3.0 添加PG商户官方验证开关配置
	//slog.Info("添加PG商户官方验证开关和默认余额告警配置开始")
	//updataPGVerifyConfig()
	//slog.Info("添加PG商户官方验证开关和默认余额告警配置结束")
	// 4.0 更新游戏类型,旧数据统一修改为 3:仿正版，更新商户表更新购买rtp开关为开启
	//slog.Info("更新游戏类型,旧数据统一修改为 3:仿正版开始")
	//updataConfigGameType()
	//slog.Info("更新游戏类型,旧数据统一修改为 3:仿正版结束")
	//slog.Info("更新商户表更新购买rtp开关为开启开始")
	//updataBuyRtpOff()
	//slog.Info("更新商户表更新购买rtp开关为开启结束")
	// 5.0 旧商户增加监控默认数据
	slog.Info("旧商户增加监控默认数据开始")
	addMonitorDefaultData()
	slog.Info("旧商户增加监控默认数据结束")
	slog.Info("旧商户添加红飞机默认数据开始")
	addSpribe1DefaultData()
	slog.Info("旧商户添加红飞机默认数据结束")
}
