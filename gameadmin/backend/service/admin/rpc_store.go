package main

import (
	"errors"
	"game/comm"
	"game/comm/mux"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
)

func init() {
	mux.RegRpc("/gamecenter/AppStore/getAppInfo", "获取商户配置", "player", getAppInfo, getAs{"100001", "89"})
}

// type GameRpcPlayer struct{}

type getAs struct {
	AppID  string `json:"app_id"`
	GameId string `json:"game_id"`
}

// AppStore 所需商户表信息
type AppStoreRes struct {
	RewardPercent  int     `json:"reward_percent"`
	NoAwardPercent int     `json:"no_award_percent"`
	AppId          string  `json:"appId"`
	GamePatten     float64 `json:"gamePatten"`
	MaxMultiple    int64   `json:"maxMultiple"`
	//最大赢分 bat * time < limit_max_point
	MaxWinPoints float64   `json:"MaxWinPoints" bson:"MaxWinPoints"`
	GameId       string    `json:"gameId"`
	Cs           []float64 `json:"cs"`
	StopLoss     int64     `json:"StopLoss"` //止盈止损开关
	//名字和时间显示
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff" `

	//退出按钮显示
	ShowExitBtnOff int64 `json:"ShowExitBtnOff" `
	// 默认cs
	DefaultCs float64 `json:"DefaultCs"`
	// 默认押注倍数
	DefaultBetLevel             int64 `json:"DefaultBetLevel"`
	IsProtection                int   `json:"IsProtection"`                //是否开启 0 关 1开
	ProtectionRotateCount       int   `json:"ProtectionRotateCount"`       // 保护次数
	ProtectionRewardPercentLess int   `json:"ProtectionRewardPercentLess"` // 保护内获取比率
	BuyMinAwardPercent          int   `json:"BuyMinAwardPercent"`          // 购买时最小获取比率
	// 利润率
	ProfitMargin float64 `json:"ProfitMargin" bson:"ProfitMargin"`
	// 坠毁概率
	CrashRate float64 `json:"CrashRate" bson:"CrashRate"`
	// 刻度
	Scale float64 `json:"Scale" bson:"Scale"`
	//游戏监控
	MoniterConfig comm.MoniterConfig `json:"MoniterConfig"`
	// 个人监控
	PersonalMoniterConfig comm.MoniterConfig `json:"PersonalMoniterConfig"`
	// Rtp
	Rtp           float64 `json:"Rtp" bson:"Rtp"`
	OnlineUpNum   float64 `json:"OnlineUpNum"`
	OnlineDownNum float64 `json:"OnlineDownNum"`
}

//	type getGainRet struct {
//		Gain int64
//		Uid  string
//	}
//

// 运营商游戏设置
type GameCofig struct {
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
	//最大赢分 bat * time < limit_max_point
	MaxWinPoints float64 `json:"MaxWinPoints" bson:"MaxWinPoints"`

	//游戏投注
	BetBase string `json:"BetBase" bson:"BetBase"`

	//游戏类型
	GamePattern float64 `json:"GamePattern" bson:"GamePattern"`

	//预设面额
	Preset float32 `json:"Preset" bson:"Preset"`

	//游戏状态
	GameOn int64 `json:"GameOn" bson:"GameOn"`

	//免费游戏开关
	FreeGameOff int64 `json:"FreeGameOff" bson:"FreeGameOff"`

	//试玩链接
	GameDemo string `json:"GameDemo" bson:"GameDemo"`

	//下注的3% 个人奖池累计配置千分比 所有slots 通用,  普通, 购买不算
	RewardPercent int `json:"RewardPercent" bson:"RewardPercent"`

	//个人奖池为负数，不中奖概率千分比  20%不中奖,  优先级最高
	NoAwardPercent int `json:"NoAwardPercent" bson:"NoAwardPercent"`

	//名字和时间显示
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff" bson:"ShowNameAndTimeOff"`

	//退出按钮显示
	ShowExitBtnOff int64 `json:"ShowExitBtnOff" bson:"ShowExitBtnOff"`

	// 默认cs
	DefaultCs float64 `json:"DefaultCs" bson:"DefaultCs"`
	// 默认押注倍数
	DefaultBetLevel    int64 `json:"DefaultBetLevel" bson:"DefaultBetLevel"`
	BuyMinAwardPercent int   `json:"BuyMinAwardPercent" bson:"BuyMinAwardPercent"`
	// 投注倍数
	//BetMult float64 `json:"BetMult" bson:"BetMult"`
	// 利润率
	ProfitMargin float64 `json:"ProfitMargin" bson:"ProfitMargin"`
	// 坠毁概率
	CrashRate float64 `json:"CrashRate" bson:"CrashRate"`
	// 刻度
	Scale         float64 `json:"Scale" bson:"Scale"`
	OnlineUpNum   float64 `json:"OnlineUpNum" bson:"OnlineUpNum"`
	OnlineDownNum float64 `json:"OnlineDownNum" bson:"OnlineDownNum"`
}

// todo 給admin后台干这活
func getAppInfo(req getAs, ret *AppStoreRes) (err error) {
	var doc []*GameCofig

	err = CollGameConfig.FindAll(bson.M{"AppID": req.AppID, "GameId": req.GameId}, &doc)
	if err != nil {
		return err
	}
	if len(doc) == 0 {
		return errors.New("can not find app config")
	}
	operator := &comm.Operator_V2{}
	CollAdminOperator.FindOne(bson.M{"AppID": req.AppID}, operator)
	if len(doc) > 0 {
		ret.AppId = doc[0].AppID
		ret.GameId = doc[0].GameId
		ret.RewardPercent = doc[0].RewardPercent
		ret.NoAwardPercent = doc[0].NoAwardPercent
		sp := strings.Split(doc[0].BetBase, ",")
		cs := make([]float64, len(sp))
		for i := 0; i < len(sp); i++ {
			cs[i], _ = strconv.ParseFloat(sp[i], 64)
			//if doc[0].BetMult != 0 {
			//	cs[i] = cs[i] * doc[0].BetMult
			//}
			if strings.HasPrefix(doc[0].GameId, "jdb_") || strings.HasPrefix(doc[0].GameId, "hacksaw_") {
				cs[i] = cs[i] * 100
			}
		}
		ret.Cs = cs
		ret.GamePatten = doc[0].GamePattern
		ret.MaxMultiple = doc[0].MaxMultiple
		ret.MaxWinPoints = doc[0].MaxWinPoints
		ret.StopLoss = doc[0].StopLoss
		ret.ShowNameAndTimeOff = doc[0].ShowNameAndTimeOff
		ret.ShowExitBtnOff = doc[0].ShowExitBtnOff
		ret.ProfitMargin = doc[0].ProfitMargin
		ret.CrashRate = doc[0].CrashRate
		ret.Scale = doc[0].Scale
		ret.Rtp = doc[0].RTP
		if strings.HasPrefix(doc[0].GameId, "pg_") {
			ret.DefaultCs = doc[0].DefaultCs
			//if doc[0].BetMult != 0 {
			//	ret.DefaultCs = doc[0].DefaultCs * doc[0].BetMult
			//}
			ret.DefaultBetLevel = doc[0].DefaultBetLevel
		}
		if strings.HasPrefix(doc[0].GameId, "spribe_") {
			ret.DefaultCs = doc[0].DefaultCs
			ret.OnlineUpNum = doc[0].OnlineUpNum
			ret.OnlineDownNum = doc[0].OnlineDownNum
		}
		if strings.HasPrefix(doc[0].GameId, "hacksaw_") {
			ret.DefaultCs = doc[0].DefaultCs
		}
		ret.BuyMinAwardPercent = doc[0].BuyMinAwardPercent
		ret.IsProtection = operator.IsProtection
		ret.ProtectionRewardPercentLess = operator.ProtectionRewardPercentLess * 10
		ret.ProtectionRotateCount = operator.ProtectionRotateCount
		ret.MoniterConfig = operator.MoniterConfig
		ret.PersonalMoniterConfig = operator.PersonalMoniterConfig
	}

	return
}
