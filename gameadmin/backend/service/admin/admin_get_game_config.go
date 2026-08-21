package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"strconv"
	"strings"
)

func init() {
	RegMsgProc("/AdminInfo/GameConfig", "游戏配置", "AdminInfo", getGameCofig, getGameCofigParams{})
	RegMsgProc("/AdminInfo/UpdateGameConfig", "修改游戏配置", "AdminInfo", updataGameConfig, updataGameCofigParams{
		AppID: "",
	})
	RegMsgProc("/AdminInfo/BatchSetGameBet", "批量修改Bet", "AdminInfo", batchSetGameBet, betUpdateConfig{})
	RegMsgProc("/AdminInfo/BatchSetGameBetMult", "批量修改Bet倍数", "AdminInfo", batchSetGameBetMult, betUpdateConfigMult{})
	mux.RegRpc("/AdminInfo/Interior/getGame", "内部调用MQ-商户游戏数据", "AdminInfo", getPrivateGameList, mqOperatorGameConfigRequest{})
	RegMsgProc("/AdminInfo/UpdateGameUniversalConfig", "更新通用配置", "AdminInfo", updataGameUniversalConfig, gameUniversalConfigParams{})
	RegMsgProc("/AdminInfo/BatchSetGameOn", "批量修改游戏状态", "AdminInfo", batchSetGameOn, gameOnUpdateConfig{})
}

type getGameCofigParams struct {
	//运营商名称
	UserName string `json:"UserName"`
	//游戏分类
	Gametype int `json:"Gametype" bson:"Gametype"`
	//游戏名称
	GameName string `json:"GameName" bson:"GameName"`
	//游戏状态
	GameOn int64 `json:"GameOn" bson:"GameOn"`
	//游戏厂商
	GameManufacturer []string `json:"GameManufacturer" bson:"ManufacturerName"`

	PageIndex int64 `json:"PageIndex"`
	PageSize  int64 `json:"PageSize"`
}
type GameCofigParams struct {
	//运营商名称
	AppID string `json:"AppID"`
	//游戏分类
	Gametype int `json:"Gametype" bson:"Gametype"`
	//游戏名称
	GameName string `json:"GameName" bson:"GameName"`
	//游戏状态
	GameOn int64 `json:"GameOn" bson:"GameOn"`

	PageIndex int64 `json:"PageIndex"`
	PageSize  int64 `json:"PageSize"`
}

// 运营商游戏设置
type operatorGameConfig struct {
	//游戏编号	game/Games/_Id
	GameId string `json:"GameId" bson:"GameId"`

	GameIdd string `json:"GameIdd" bson:"_id"`
	//游戏分类	game/Games/Type
	Gametype int `json:"Gametype" bson:"Type"`

	//游戏名称	game/Games/Name
	GameName       string                 `json:"GameName" bson:"Name"`
	ChangeBetOff   int64                  `json:"ChangeBetOff" bson:"ChangeBetOff"`
	GameNameConfig map[string]*GameConfig `json:"GameNameConfig" bson:"GameNameConfig"`
	GameStatus     int64                  `json:"Status" bson:"Status"`

	GameManufacturer string `json:"GameManufacturer" bson:"ManufacturerName"`
	//配置文件
	ConfigPath string `json:"ConfigPath" bson:"ConfigPath"`

	//RTP设置
	RTP float64 `json:"RTP" bson:"RTP"`

	//止盈止损开关
	StopLoss int64 `json:"StopLoss" bson:"StopLoss"`

	//赢取最高押注倍数
	MaxMultiple int64 `json:"MaxMultiple" bson:"MaxMultiple"`

	//游戏投注
	BetBase string `json:"BetBase" bson:"BetBase"`

	//游戏类型
	GamePattern int `json:"GamePattern" bson:"GamePattern"`

	//预设面额
	Preset float32 `json:"Preset" bson:"Preset"`

	//游戏状态
	GameOn int64 `json:"GameOn" bson:"GameOn"`

	//免费游戏开关
	FreeGameOff int64 `json:"FreeGameOff" bson:"FreeGameOff"`

	//试玩链接
	GameDemo string `json:"GameDemo" bson:"GameDemo"`

	//RTP设置功能开关
	RTPOff int64 `json:"RTPOff" bson:"RTPOff"`
	//止损止赢功能开关
	StopLossOff int64 `json:"StopLossOff" bson:"StopLossOff"`
	//赢取最高押注倍数设置功能开关
	MaxMultipleOff int64 `json:"MaxMultipleOff" bson:"MaxMultipleOff"`

	//名字和时间显示
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff" bson:"ShowNameAndTimeOff"`
	//退出按钮显示
	ShowExitBtnOff int64 `json:"ShowExitBtnOff" bson:"ShowExitBtnOff"`
	//退出按钮链接
	ExitLink     string  `json:"ExitLink" bson:"ExitLink"`
	MaxWinPoints float64 `json:"MaxWinPoints" bson:"MaxWinPoints"`
	// 默认cs
	DefaultCs float64 `json:"DefaultCs" bson:"DefaultCs"`
	// 默认押注倍数
	DefaultBetLevel int64 `json:"DefaultBetLevel" bson:"DefaultBetLevel"`

	// 显示背包
	ShowBag int64 `json:"ShowBag" bson:"ShowBag"`

	OnlineUpNum   float64 `json:"OnlineUpNum" bson:"OnlineUpNum"`
	OnlineDownNum float64 `json:"OnlineDownNum" bson:"OnlineDownNum"`
	BuyRTP        int     `json:"BuyRTP" bson:"BuyRTP"`
	BuyType       int     `json:"BuyType" bson:"BuyType"`
	// 投注倍数
	//BetMult float64 `json:"BetMult" bson:"BetMult"`
	// 利润率
	ProfitMargin float64 `json:"ProfitMargin" bson:"ProfitMargin"`
	// 坠毁概率
	CrashRate float64 `json:"CrashRate" bson:"CrashRate"`
	// 刻度
	Scale float64 `json:"Scale" bson:"Scale"`
}

// 其中两个字段获取 game/ game 中的数据
type respGameCofig struct {
	List     []*operatorGameConfig
	CountAll int64
}

type updataGameCofigParams struct {
	//运营商名称
	AppID string `json:"AppID" bson:"Operator"`

	//游戏编号
	GameId string `json:"GameId" bson:"GameId"`

	//预设面额
	Preset float32 `json:"Preset" bson:"Preset"`

	//止盈止损开关
	StopLoss int64 `json:"StopLoss" bson:"StopLoss"`

	//游戏类型
	GamePattern int `json:"GamePattern" bson:"GamePattern"`

	//免费游戏开关
	FreeGameOff int64 `json:"FreeGameOff" bson:"FreeGameOff"`

	//RTP设置
	RTP float64 `json:"RTP" bson:"RTP"`

	//赢取最高押注倍数
	MaxMultiple string `json:"MaxMultiple" bson:"MaxMultiple"`

	//游戏投注
	BetBase string `json:"BetBase" bson:"BetBase"`

	//游戏状态
	GameOn int64 `json:"GameOn" bson:"GameOn"`

	RewardPercent  int64 `json:"RewardPercent" bson:"RewardPercent"`
	NoAwardPercent int64 `json:"NoAwardPercent" bson:"NoAwardPercent"`

	//名字和时间显示
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff" bson:"ShowNameAndTimeOff"`

	//退出按钮显示
	ShowExitBtnOff int64 `json:"ShowExitBtnOff" bson:"ShowExitBtnOff"`
	//退出按钮链接
	ExitLink     string  `json:"ExitLink" bson:"ExitLink"`
	MaxWinPoints float64 `json:"MaxWinPoints" bson:"MaxWinPoints"`
	// 默认cs
	DefaultCs float64 `json:"DefaultCs" bson:"DefaultCs"`
	// 默认押注倍数
	DefaultBetLevel int64 `json:"DefaultBetLevel" bson:"DefaultBetLevel"`
	// 显示背包
	ShowBag       int64   `json:"ShowBag" bson:"ShowBag"`
	OnlineUpNum   float64 `json:"OnlineUpNum" bson:"OnlineUpNum"`
	OnlineDownNum float64 `json:"OnlineDownNum" bson:"OnlineDownNum"`

	// 购买游戏rtp
	BuyRTP             int `json:"BuyRTP" bson:"BuyRTP"`
	BuyMinAwardPercent int `json:"BuyMinAwardPercent" bson:"BuyMinAwardPercent"`
	// 投注倍数
	//BetMult float64 `json:"BetMult" bson:"BetMult"`
	// 利润率
	ProfitMargin float64 `json:"ProfitMargin" bson:"ProfitMargin"`
	// 坠毁概率
	CrashRate float64 `json:"CrashRate" bson:"CrashRate"`
	// 刻度
	Scale float64 `json:"Scale" bson:"Scale"`
}

type betUpdateConfig struct {
	AppID         string  `description:"商户"`
	Manufacturer  string  `description:"游戏厂商"`
	GameIds       string  `description:"游戏"`
	OnlineUpNum   float64 `description:"Bet上线"`
	OnlineDownNum float64 `description:"Bet下线"`
	BetBase       string  `description:"投注额"`
}
type betUpdateConfigMult struct {
	AppID        string  `description:"商户"`
	Manufacturer string  `description:"游戏厂商"`
	GameIds      string  `description:"游戏"`
	BetMult      float64 `description:"投注倍数"`
}
type gameOnUpdateConfig struct {
	AppID   string   `json:"AppID"`
	GameIds []string `json:"GameIds" description:"游戏"`
	GameOn  int64    `json:"GameOn" description:"游戏状态"`
}

type gameUniversalConfigParams struct {
	//运营商名称
	AppID      string          `json:"AppID"`
	PGConfig   comm.PGConfig   `json:"PGConfig"`
	JILIConfig comm.JILIConfig `json:"JILIConfig"`
	PPConfig   comm.PPConfig   `json:"PPConfig"`
	TADAConfig comm.TADAConfig `json:"TADAConfig"`
}

type getGameUniversalConfigParams struct {
	//运营商名称
	AppID string `json:"AppID"`
}

func getGameCofig(ctx *Context, ps getGameCofigParams, ret *respGameCofig) (err error) {
	if ps.PageSize <= 0 || ps.PageIndex <= 0 {
		ps.PageIndex = 1
		ps.PageSize = 20
	}

	filter := bson.M{}
	var or []interface{}
	if ps.GameName != "" {
		filterStr := fmt.Sprintf("GameNameConfig.%s.GameName", ctx.Lang)

		or = append(or, bson.M{filterStr: bson.M{"$regex": ps.GameName}})
		or = append(or, bson.M{"_id": bson.M{"$regex": ps.GameName}})
		or = append(or, bson.M{"GameId": bson.M{"$regex": ps.GameName}})
		if len(ps.GameManufacturer) == 0 {
			filter["$or"] = or
		}
	}
	if len(ps.GameManufacturer) != 0 {
		and := []interface{}{
			bson.M{"ManufacturerName": bson.M{"$in": ps.GameManufacturer}},
		}
		if len(or) != 0 {
			and = append(and, bson.M{"$or": or})
		}

		filter["$and"] = and
	}
	err, games := GetGameListFormatName(ctx.Lang, filter)
	gameIDlist := make([]string, 10)
	mapBuytype := map[string]int{}
	for _, game2 := range games {
		gameIDlist = append(gameIDlist, game2.ID)
		mapBuytype[game2.ID] = game2.BuyType
	}

	query := bson.M{}
	query["AppID"] = ps.UserName

	// 当游戏为未关闭时
	query["GameOn"] = bson.M{
		"$ne": 2,
	}

	if ps.GameOn != -1 {
		query["GameOn"] = ps.GameOn
	}
	query["GameId"] = bson.M{"$in": gameIDlist}

	gameNameMap := map[string]*comm.Game2{}

	find := mongodb.FindPageOpt{
		Page:       ps.PageIndex,
		PageSize:   ps.PageSize,
		Sort:       bson.M{"GameId": 1},
		Query:      query,
		Projection: bson.M{"_id": 0},
	}

	var list_OperatorGameConfig []*operatorGameConfig
	count, err := CollGameConfig.FindPage(find, &list_OperatorGameConfig)

	for i, config := range list_OperatorGameConfig {
		ga, ok := MapGameBet[config.GameId]
		if ok == false {
			list_OperatorGameConfig[i].GameName = ""
		} else {
			list_OperatorGameConfig[i].GameName = ga.GameName
		}
		list_OperatorGameConfig[i].BuyType = mapBuytype[config.GameId]
	}
	ret.List = list_OperatorGameConfig
	ret.CountAll = count

	for _, game := range games {
		gameNameMap[game.ID] = game
	}
	for i, play := range ret.List {
		_, ok := gameNameMap[play.GameId]
		if ok == false {
			gameNameMap[play.GameId] = nil
		}

		if gameNameMap[play.GameId] == nil {

			ret.List[i].GameName = "/"
			ret.List[i].GameManufacturer = ""
			ret.List[i].ChangeBetOff = 0
		} else {

			ret.List[i].GameName = gameNameMap[play.GameId].Name
			ret.List[i].GameManufacturer = gameNameMap[play.GameId].ManufacturerName
			ret.List[i].ChangeBetOff = gameNameMap[play.GameId].ChangeBetOff
		}

	}

	return
}

func GetOperatorGameConfig(ps getGameCofigParams, ret *respGameCofig) (err error) {
	var operator comm.Operator_V2
	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.UserName}, &operator)
	if err != nil {
		return
	}

	query := bson.M{}
	query["AppID"] = ps.UserName

	if ps.GameOn != -1 {
		query["GameOn"] = ps.GameOn
	}

	//游戏名称 模糊查询
	if ps.GameName != "" {
		//query["GameName"] = bson.M{"$regex": "/" + ps.GameName + "/", "$options": 'i'}
		query["GameName"] = primitive.Regex{
			Pattern: ps.GameName,
			Options: "i",
		}
	}

	find := mongodb.FindPageOpt{
		Page:       ps.PageIndex,
		PageSize:   ps.PageSize,
		Sort:       bson.M{"GameId": 1},
		Query:      query,
		Projection: bson.M{"_id": 0},
	}

	var list_OperatorGameConfig []*operatorGameConfig
	count, err := CollGameConfig.FindPage(find, &list_OperatorGameConfig)
	//for _, config := range list_OperatorGameConfig {
	//	_ = NewOtherDB("game").Collection("Games").FindOne(bson.M{"_id": config.GameId}, config)
	//}

	for i, config := range list_OperatorGameConfig {
		ga, ok := MapGameBet[config.GameId]
		if ok == false {
			list_OperatorGameConfig[i].GameName = ""
		} else {
			list_OperatorGameConfig[i].GameName = ga.GameName
		}
	}
	ret.List = list_OperatorGameConfig
	ret.CountAll = count
	return err
}

func isElementInArray(target float64, nums []float64) bool {
	for _, num := range nums {
		if num == target {
			return true
		}
	}
	return false
}

func updataGameConfig(ctx *Context, ps updataGameCofigParams, ret *comm.Empty) (err error) {

	// 验证RTP可以设置的最大值
	err, b := comm.VerifyMaxRTP(ps.RTP)
	if b == false {
		return err
	}

	var result comm.Operator_V2

	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New(lang.GetLang(ctx.Lang, "商户不存在"))
	}
	if err != nil {
		return err
	}
	sp := strings.Split(ps.BetBase, ",")
	cs := make([]float64, len(sp))
	for i := 0; i < len(sp); i++ {
		cs[i], _ = strconv.ParseFloat(sp[i], 64)
	}
	settingInfo := bson.M{}
	if strings.HasPrefix(ps.GameId, "pg_") {
		if !isElementInArray(ps.DefaultCs, cs) {
			return errors.New(lang.GetLang(ctx.Lang, "默认押注金额和当前投注额不符"))
		}
		settingInfo["DefaultCs"] = ps.DefaultCs
		settingInfo["DefaultBetLevel"] = ps.DefaultBetLevel
	}
	if strings.HasPrefix(ps.GameId, "spribe_") || strings.HasPrefix(ps.GameId, "hacksaw_") {
		settingInfo["DefaultCs"] = ps.DefaultCs
	}
	settingInfo["Preset"] = ps.Preset
	settingInfo["StopLoss"] = ps.StopLoss
	settingInfo["GamePattern"] = ps.GamePattern
	settingInfo["FreeGameOff"] = ps.FreeGameOff
	settingInfo["RTP"] = ps.RTP
	settingInfo["MaxMultiple"], _ = strconv.ParseFloat(ps.MaxMultiple, 64)
	settingInfo["MaxWinPoints"] = ps.MaxWinPoints
	settingInfo["BetBase"] = ps.BetBase
	settingInfo["GameOn"] = ps.GameOn
	settingInfo["ShowNameAndTimeOff"] = ps.ShowNameAndTimeOff
	settingInfo["ShowExitBtnOff"] = ps.ShowExitBtnOff
	settingInfo["OnlineUpNum"] = ps.OnlineUpNum
	settingInfo["OnlineDownNum"] = ps.OnlineDownNum
	settingInfo["ShowBag"] = ps.ShowBag
	settingInfo["ExitLink"] = ps.ExitLink

	// 特殊处理，当购买rtp开关关闭时不修改购买rtp
	if result.BuyRTPOff != 0 {
		settingInfo["BuyRTP"] = ps.BuyRTP
		settingInfo["BuyMinAwardPercent"] = MapBuyGameRTP[ps.BuyRTP]
	}
	rw, nw := GetGameNwandAw(ps.GameId, ps.RTP)
	settingInfo["RewardPercent"] = rw
	settingInfo["NoAwardPercent"] = nw
	settingInfo["ProfitMargin"] = ps.ProfitMargin
	settingInfo["CrashRate"] = ps.CrashRate
	settingInfo["Scale"] = ps.Scale
	//gameConfig := comm.GameConfig{}
	//CollGameConfig.FindOne(bson.M{"GameId": ps.GameId, "AppID": ps.AppID}, &gameConfig)
	//if gameConfig.BetMult == 0 {
	//	settingInfo["BetMult"] = ps.BetMult
	//} else {
	//	settingInfo["BetMult"] = gameConfig.BetMult * ps.BetMult
	//}

	set := bson.M{
		"$set": settingInfo,
	}

	err = CollGameConfig.Update(bson.M{"GameId": ps.GameId, "AppID": ps.AppID}, set)

	PushMsgStoreSetInfo(ps, rw, nw, cs)

	return
}

func batchSetGameBet(ctx *Context, ps betUpdateConfig, ret *comm.Empty) (err error) {

	user, _ := IsAdminUser(ctx)

	if user.GroupId == 2 {
		return errors.New("权限不足")
	}
	gameIds := strings.Split(ps.GameIds, ",")

	gameMap := map[string]*comm.Game2{}
	updateModel := []mongo.WriteModel{}

	gameFilter := bson.M{}

	if strings.ToUpper(ps.Manufacturer) != "PP" {
		gameFilter["ChangeBetOff"] = 1
	}

	if ps.GameIds == "ALL" {
		gameFilter["ManufacturerName"] = ps.Manufacturer
		err, gameList := GetGameListFormatName(ctx.Lang, gameFilter)

		if err != nil {
			return err
		}

		gameIds = nil

		for _, game := range gameList {
			gameMap[game.ID] = game
			gameIds = append(gameIds, game.ID)
		}

	} else {
		gameFilter["_id"] = bson.M{"$in": gameIds}
		err, gameList := GetGameListFormatName(ctx.Lang, gameFilter)

		if err != nil {
			return err
		}

		for _, game := range gameList {
			gameMap[game.ID] = game
		}

	}

	betBase := convertBetFloat(ps.BetBase)

	if gameIds != nil && len(gameIds) > 0 {
		if strings.ToUpper(ps.Manufacturer) == "JILI" || strings.ToUpper(ps.Manufacturer) == "TADA" {

			var gameLineNum float64
			var gameBet []string
			model := &mongo.UpdateManyModel{}
			for _, id := range gameIds {

				//=======================================JILI游戏bet信息每个值需要乘线数配置=======================================

				gameLineNum = gameMap[id].LineNum // 当前游戏的线数
				gameBet = []string{}              // 重置bet数组

				for _, i2 := range betBase {
					r := strconv.FormatFloat(i2*gameLineNum, 'f', 2, 64)
					gameBet = append(gameBet, r)
				}

				model = &mongo.UpdateManyModel{
					Filter: bson.M{
						"AppID":  ps.AppID,
						"GameId": id,
					},
					Update: bson.M{"$set": bson.M{
						"OnlineUpNum":   ps.OnlineUpNum,
						"OnlineDownNum": ps.OnlineDownNum,
						"BetBase":       strings.Join(gameBet, ","),
						//"BetMult":       1,
					}},
				}

				updateModel = append(updateModel, model)
				model = nil
			}

		} else {
			filterBson := bson.M{
				"AppID":  ps.AppID,
				"GameId": bson.M{"$in": gameIds},
			}
			setBson := bson.M{}

			if strings.ToLower(ps.Manufacturer) == "pg" {
				if !isElementInArray(betBase[0], betBase) {
					return errors.New(lang.GetLang(ctx.Lang, "默认押注金额和当前投注额不符"))
				}
				setBson["DefaultCs"] = betBase[0]
				setBson["DefaultBetLevel"] = 10
			}

			//setBson["BetMult"] = 1
			setBson["OnlineDownNum"] = ps.OnlineDownNum
			setBson["OnlineUpNum"] = ps.OnlineUpNum
			setBson["BetBase"] = ps.BetBase

			model := &mongo.UpdateManyModel{
				Filter: filterBson,
				Update: bson.M{
					"$set": setBson,
				},
			}

			updateModel = append(updateModel, model)
		}
		_, err := CollGameConfig.Coll().BulkWrite(context.TODO(), updateModel, options.BulkWrite().SetOrdered(false))
		if err != nil {
			return err
		}

		go func() {
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
	}

	return err
}

func batchSetGameBetMult(ctx *Context, ps betUpdateConfigMult, ret *comm.Empty) (err error) {

	user, _ := IsAdminUser(ctx)

	MIN_BET := 0.005

	if user.GroupId == 2 {
		return errors.New("权限不足")
	}
	gameIds := strings.Split(ps.GameIds, ",")

	gameMap := map[string]*comm.Game2{}
	updateModel := []mongo.WriteModel{}

	gameFilter := bson.M{}

	if ps.GameIds == "ALL" {
		if ps.Manufacturer != "" && ps.Manufacturer != "ALL" {
			gameFilter["ManufacturerName"] = ps.Manufacturer
		}
		gameFilter["Type"] = 0
		err, gameList := GetGameListFormatName(ctx.Lang, gameFilter)

		if err != nil {
			return err
		}

		gameIds = nil

		for _, game := range gameList {
			gameMap[game.ID] = game
			gameIds = append(gameIds, game.ID)
		}

	} else {
		gameFilter["_id"] = bson.M{"$in": gameIds}
		err, gameList := GetGameListFormatName(ctx.Lang, gameFilter)

		if err != nil {
			return err
		}

		for _, game := range gameList {
			gameMap[game.ID] = game
		}

	}
	if ps.BetMult == 0 {
		ps.BetMult = 1
	}

	if gameIds != nil && len(gameIds) > 0 {
		for _, id := range gameIds {
			gameConfig := comm.GameConfig{}
			CollGameConfig.FindOne(bson.M{"AppID": ps.AppID, "GameId": id}, &gameConfig)
			model := &mongo.UpdateManyModel{}
			betBase := convertBetFloat(gameConfig.BetBase)
			gameBet := []string{} // 重置bet数组
			for _, i2 := range betBase {
				i2 = i2 * ps.BetMult

				// 限值bet值不能过小
				if i2 < MIN_BET {
					i2 = MIN_BET
				}
				r := strconv.FormatFloat(i2, 'f', 5, 64)
				gameBet = append(gameBet, r)
			}

			ComputedDefaultCs := math.Round(gameConfig.DefaultCs*ps.BetMult*100000) / 100000
			ComputedOnlineUpNum := math.Round(gameConfig.OnlineUpNum*ps.BetMult*100000) / 100000
			ComputedOnlineDownNum := gameConfig.OnlineDownNum

			// 如果计算后默认值不能小于 最小值
			if ComputedDefaultCs < MIN_BET {
				ComputedDefaultCs = MIN_BET
			}

			// 如果计算后投注上限不能小于 最小值
			if ComputedOnlineUpNum < MIN_BET {
				ComputedOnlineUpNum = MIN_BET
			}

			if ComputedOnlineUpNum != MIN_BET {
				ComputedOnlineDownNum = math.Round(ComputedOnlineDownNum*ps.BetMult*100000) / 100000
			}

			// 如果计算后投注下限不能小于 最小值
			if ComputedOnlineDownNum < MIN_BET {
				ComputedOnlineDownNum = MIN_BET
			}
			gameConfig.DefaultCs = ComputedDefaultCs
			gameConfig.OnlineUpNum = ComputedOnlineUpNum
			gameConfig.OnlineDownNum = ComputedOnlineDownNum

			model = &mongo.UpdateManyModel{
				Filter: bson.M{
					"AppID":  ps.AppID,
					"GameId": id,
				},
				Update: bson.M{"$set": bson.M{
					//"BetMult": gameConfig.BetMult * ps.BetMult,
					"BetBase":       strings.Join(gameBet, ","),
					"DefaultCs":     gameConfig.DefaultCs,
					"OnlineUpNum":   gameConfig.OnlineUpNum,
					"OnlineDownNum": gameConfig.OnlineDownNum,
				}},
			}
			updateModel = append(updateModel, model)
		}
		_, err := CollGameConfig.Coll().BulkWrite(context.TODO(), updateModel, options.BulkWrite().SetOrdered(false))
		if err != nil {
			return err
		}

		go func() {
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
	}

	return err
}

type AppStore struct {
	RewardPercent  int       `json:"reward_percent"`
	NoAwardPercent int       `json:"no_award_percent"`
	AppId          string    `json:"appId"`
	GamePatten     int       `json:"gamePatten"`
	MaxMultiple    float64   `json:"maxMultiple"`
	MaxWinPoints   float64   `json:"maxWinPoints"`
	GameId         string    `json:"gameId"`
	Cs             []float64 `json:"cs"`
	StopLoss       int64     `json:"StopLoss"` //止盈止损开关
	//名字和时间显示
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff" `
	//退出按钮显示
	ShowExitBtnOff  int64   `json:"ShowExitBtnOff" `
	DefaultCs       float64 `json:"DefaultCs"`
	DefaultBetLevel int64   `json:"DefaultBetLevel"`
	//新增
	IsProtection                int `json:"IsProtection"`                //是否开启 0 关 1开
	ProtectionRotateCount       int `json:"ProtectionRotateCount"`       // 保护次数
	ProtectionRewardPercentLess int `json:"ProtectionRewardPercentLess"` // 保护内获取比率
	BuyMinAwardPercent          int `json:"BuyMinAwardPercent"`          // 购买时最小获取比率
	// 利润率
	ProfitMargin float64 `json:"ProfitMargin"`
	// 坠毁概率
	CrashRate float64 `json:"CrashRate"`
	// 刻度
	Scale float64 `json:"Scale"`
	//游戏监控
	MoniterConfig comm.MoniterConfig `json:"MoniterConfig"`
	// 个人监控
	PersonalMoniterConfig comm.MoniterConfig `json:"PersonalMoniterConfig"`
	// Rtp
	RTP           float64 `json:"RTP"`
	OnlineUpNum   float64 `json:"OnlineUpNum"`
	OnlineDownNum float64 `json:"OnlineDownNum"`
}

func PushMsgStoreSetInfo(ps updataGameCofigParams, rw, nw int64, cs []float64) {
	var ret AppStore
	ret.AppId = ps.AppID
	ret.GameId = ps.GameId
	//ret.RewardPercent = int(ps.RewardPercent)
	//ret.NoAwardPercent = 200
	ret.RewardPercent = int(rw)
	ret.NoAwardPercent = int(nw)
	//if ps.BetMult != 0 {
	//	for i, _ := range cs {
	//		cs[i] = cs[i] * ps.BetMult
	//	}
	if strings.HasPrefix(ps.GameId, "jdb_") || strings.HasPrefix(ps.GameId, "hacksaw_") {
		for i, _ := range cs {
			cs[i] = cs[i] * 100
		}
	}
	ret.Cs = cs
	ret.GamePatten = ps.GamePattern
	ret.MaxMultiple, _ = strconv.ParseFloat(ps.MaxMultiple, 64)
	ret.MaxWinPoints = float64(ps.MaxWinPoints)
	ret.StopLoss = ps.StopLoss
	ret.ShowNameAndTimeOff = ps.ShowNameAndTimeOff
	ret.ShowExitBtnOff = ps.ShowExitBtnOff
	ret.ProfitMargin = ps.ProfitMargin
	ret.CrashRate = ps.CrashRate
	ret.Scale = ps.Scale
	ret.RTP = ps.RTP
	if strings.HasPrefix(ps.GameId, "pg_") {
		ret.DefaultCs = ps.DefaultCs
		//if ps.BetMult != 0 {
		//	ret.DefaultCs = ps.DefaultCs * ps.BetMult
		//}
		ret.DefaultBetLevel = ps.DefaultBetLevel
	}
	if strings.HasPrefix(ps.GameId, "pp_") {
		ret.DefaultCs = ps.DefaultCs
		//if ps.BetMult != 0 {
		//	ret.DefaultCs = ps.DefaultCs * ps.BetMult
		//}
	}
	if strings.HasPrefix(ps.GameId, "spribe_") {
		ret.DefaultCs = ps.DefaultCs
		ret.OnlineUpNum = ps.OnlineUpNum
		ret.OnlineDownNum = ps.OnlineDownNum
	}
	if strings.HasPrefix(ps.GameId, "hacksaw_") {
		ret.DefaultCs = ps.DefaultCs
	}
	operator := &comm.Operator_V2{}
	CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, operator)
	// BuyMinAwardPercent值不准确当为-1不接收方接收到消息应该忽略BuyMinAwardPercent的值使用旧值或者数据库中的值
	if operator.BuyRTPOff != 0 {
		ret.BuyMinAwardPercent = ps.BuyMinAwardPercent
	} else {
		ret.BuyMinAwardPercent = -1
	}
	ret.IsProtection = operator.IsProtection
	ret.ProtectionRotateCount = operator.ProtectionRotateCount
	ret.ProtectionRewardPercentLess = operator.ProtectionRewardPercentLess * 10
	ret.MoniterConfig = operator.MoniterConfig
	ret.PersonalMoniterConfig = operator.PersonalMoniterConfig
	sub := "/store/setInfo_" + ps.GameId
	_ = mq.PublishMsg(sub, ret)

	return
}

//
//type GameNow struct {
//	GameID        string `json:"gameId"`
//	RewardPercent int64  `json:"reward_percent"`
//}

//var MapGameNow = map[string]int64{}

func GetGameNwandAw(gameId string, RTP float64) (rw, nw int64) {
	//var gameinfo *comm.Game
	//err := NewOtherDB("game").Collection("Games").FindOne(bson.M{"_id": gameId}, &gameinfo)
	//fileBytes, err := os.ReadFile("./config/output.json")
	//
	//var data []*GameNow
	//
	//err = json.Unmarshal(fileBytes, &data)
	//if err != nil {
	//	log.Fatalf("Error unmarshalling JSON: %v", err)
	//}
	//
	//for _, vl := range data {
	//	MapGameNow[vl.GameID] = vl.RewardPercent
	//}

	var k int64
	switch RTP {
	case 10:
		k = -8400
		nw = 1001
	case 20:
		k = -3700
		nw = 1001
	case 30:
		k = -2130
		nw = 1001
	case 40:
		k = -1350
		nw = 1001
	case 50:
		k = -870
		nw = 1001
	case 60:
		k = -570
		nw = 700
	case 70:
		k = -340
		nw = 600
	case 80:
		k = -170
		nw = 400
	case 85:
		k = -105
		nw = 300
	case 90:
		k = -44
		nw = 200
	case 91:
		k = -32
		nw = 200
	case 92:
		k = -21
		nw = 200
	case 93:
		k = -10
		nw = 200
	case 94:
		k = 0
		nw = 200
	case 95:
		k = 10
		nw = 200
	case 96:
		k = 20
		nw = 200
	case 97:
		k = 30
		nw = 200
	case 98:
		k = 40
		nw = 200
	case 99:
		k = 50
		nw = 200
	case 100:
		k = 60
		nw = 200
	case 102:
		k = 80
		nw = 200
	case 105:
		k = 110
		nw = 200
	case 110:
		k = 160
		nw = 200
	case 120:
		k = 260
		nw = 200
	case 130:
		k = 360
		nw = 200
	case 140:
		k = 460
		nw = 200
	case 150:
		k = 560
		nw = 200
	case 160:
		k = 660
		nw = 200
	case 170:
		k = 760
		nw = 200
	case 180:
		k = 860
		nw = 200
	case 190:
		k = 960
		nw = 200
	case 200:
		k = 1060
		nw = 200
	case 220:
		k = 1260
		nw = 200
	case 240:
		k = 1460
		nw = 200
	case 260:
		k = 1660
		nw = 200
	case 280:
		k = 1860
		nw = 200
	case 300:
		k = 2060
		nw = 200
	}
	_, ok := MapGameNow[gameId]
	if ok == false {
		MapGameNow[gameId] = 0
	}
	rw = k + MapGameNow[gameId]
	return
}

var MapBuyGameRTP = map[int]int{0: 0, 10: 950, 20: 830, 30: 720, 40: 610, 50: 500, 60: 390, 70: 300, 80: 280, 85: 100, 90: 50, 95: 1}

// -------------------------------------------内部MQ使用------------------------------------------------

type mqOperatorGameConfigRequest struct {
	Language string
	AppID    string
	Status   int64
	getGameCofigParams
}
type mqOperatorGameConfigResponse struct {
	CountAll int64
	List     []*comm.Game2
}
type gameConfig struct {
	GameName string
	Icon     string
}

func getPrivateGameList(ps mqOperatorGameConfigRequest, ret *mqOperatorGameConfigResponse) (err error) {

	ps.UserName = ps.AppID
	ps.GameOn = 0
	ps.PageIndex = 1
	ps.PageSize = 30000

	var operator comm.Operator_V2
	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.UserName}, &operator)
	if err != nil {
		return
	}
	GameConfig := NewOtherDB("game").Collection("Games")

	gameMap := map[string]*comm.Game2{}
	var allGame []*comm.Game2

	GameConfigQuery := bson.M{
		"Status": 0,
	}

	if len(operator.DefaultManufacturerOn) > 0 {

		GameConfigQuery["ManufacturerName"] = bson.M{"$in": operator.DefaultManufacturerOn}

	}

	err = GameConfig.FindAll(GameConfigQuery, &allGame)

	if err != nil {
		return err
	}
	for _, v := range allGame {
		gameMap[v.ID] = v
	}

	var para respGameCofig
	err = GetOperatorGameConfig(ps.getGameCofigParams, &para)

	if err != nil {
		return err
	}

	//图标

	for _, v := range para.List {

		if gameMap[v.GameId] == nil {

			continue
		}

		gameNameConfig := gameMap[v.GameId].GameNameConfig
		if gameNameConfig != nil {

			if gameNameConfig[strings.ToLower(ps.Language)].GameName != "" {
				gameMap[v.GameId].Name = gameNameConfig[strings.ToLower(ps.Language)].GameName
			}
			if gameNameConfig[strings.ToLower(ps.Language)].Icon != "" {
				gameMap[v.GameId].IconUrl = gameNameConfig[strings.ToLower(ps.Language)].Icon
			}

		}

		ret.List = append(ret.List, gameMap[v.GameId])
	}

	return err
}

func updataGameUniversalConfig(ctx *Context, ps gameUniversalConfigParams, ret *comm.Empty) (err error) {

	var result comm.Operator_V2

	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New(lang.GetLang(ctx.Lang, "商户不存在"))
	}
	if err != nil {
		return err
	}
	settingPPInfo := bson.M{}
	settingPGInfo := bson.M{}
	settingJILIInfo := bson.M{}
	settingTADAInfo := bson.M{}
	settingPPInfo["ShowNameAndTimeOff"] = ps.PPConfig.ShowNameAndTimeOff
	//settingPPInfo["CurrencyVisibleOff"] = ps.PPConfig.CurrencyVisibleOff
	settingPGInfo["CarouselOff"] = ps.PGConfig.CarouselOff
	settingPGInfo["ShowNameAndTimeOff"] = ps.PGConfig.ShowNameAndTimeOff
	//settingPGInfo["CurrencyVisibleOff"] = ps.PGConfig.CurrencyVisibleOff
	settingPGInfo["StopLoss"] = ps.PGConfig.StopLoss
	settingPGInfo["ExitBtnOff"] = ps.PGConfig.ExitBtnOff
	settingPGInfo["ExitLink"] = ps.PGConfig.ExitLink
	settingPGInfo["OfficialVerify"] = ps.PGConfig.OfficialVerify
	settingJILIInfo["BackPackOff"] = ps.JILIConfig.BackPackOff
	settingJILIInfo["OpenScreenOff"] = ps.JILIConfig.OpenScreenOff
	settingJILIInfo["SidebarOff"] = ps.JILIConfig.SidebarOff
	//settingJILIInfo["CurrencyVisibleOff"] = ps.JILIConfig.CurrencyVisibleOff
	settingTADAInfo["BackPackOff"] = ps.TADAConfig.BackPackOff
	settingTADAInfo["OpenScreenOff"] = ps.TADAConfig.OpenScreenOff
	settingTADAInfo["SidebarOff"] = ps.TADAConfig.SidebarOff
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

	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^pp"}, "AppID": ps.AppID}, setPP)
	if err != nil {
		return err
	}
	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^pg"}, "AppID": ps.AppID}, setPG)
	if err != nil {
		return err
	}
	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^jili"}, "AppID": ps.AppID}, setJILI)
	if err != nil {
		return err
	}
	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$regex": "^tada"}, "AppID": ps.AppID}, setTADA)
	if err != nil {
		return err
	}
	//result.CurrencyManufactureVisibleOff["pg"] = int(ps.PGConfig.CurrencyVisibleOff)
	//result.CurrencyManufactureVisibleOff["pp"] = int(ps.PPConfig.CurrencyVisibleOff)
	//result.CurrencyManufactureVisibleOff["jili"] = int(ps.JILIConfig.CurrencyVisibleOff)
	//result.CurrencyManufactureVisibleOff["tada"] = int(ps.TADAConfig.CurrencyVisibleOff)
	result.PGConfig = ps.PGConfig
	result.PPConfig = ps.PPConfig
	result.JILIConfig = ps.JILIConfig
	result.TADAConfig = ps.TADAConfig
	settingInfo := bson.M{}
	//settingInfo["CurrencyVisibleOff"] = result.CurrencyManufactureVisibleOff
	settingInfo["PGConfig"] = result.PGConfig
	settingInfo["PPConfig"] = result.PPConfig
	settingInfo["JILIConfig"] = result.JILIConfig
	settingInfo["TADAConfig"] = result.TADAConfig
	update := bson.M{
		"$set": settingInfo,
	}
	CollAdminOperator.UpdateOne(bson.M{"AppID": ps.AppID}, update)

	return
}

func batchSetGameOn(ctx *Context, ps gameOnUpdateConfig, ret *comm.Empty) (err error) {
	var result comm.Operator_V2
	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New(lang.GetLang(ctx.Lang, "商户不存在"))
	}
	if err != nil {
		return err
	}
	settingInfo := bson.M{}

	settingInfo["GameOn"] = ps.GameOn

	set := bson.M{
		"$set": settingInfo,
	}

	err = CollGameConfig.Update(bson.M{"GameId": bson.M{"$in": ps.GameIds}, "AppID": ps.AppID}, set)

	return
}
