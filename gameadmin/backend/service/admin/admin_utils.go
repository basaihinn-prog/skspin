package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/duck/logger"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

func init() {
	RegMsgProc("/AdminInfo/SyncOperatorBet", "同步默认BET值", "AdminInfo", syncOperatorBet, comm.Empty{})
	//RegMsgProc("/AdminInfo/OperatorSyncGames", "同步默认BET值", "AdminInfo", operatorSyncGames, comm.Empty{})
}

type syncGameCofigParams struct {
	MaxMultiple int64 `bson:"MaxMultiple" json:"MaxMultiple"`
	updataGameCofigParams
}

type aggregateData struct {
	Id        string   `bson:"_id" json:"Id"`
	GamesItem []string `json:"GamesItem" bson:"GamesItem"`
}

func syncOperatorBet(ctx *Context, ps getGameCofigParams, ret *respGameCofig) (err error) {

	if ctx.Username != "admin_BW" {
		return errors.New("只有admin_BW可以同步")
	}

	models := []mongo.WriteModel{}

	for _, i2 := range MapGameBet {

		betList := strings.Split(i2.Bet, ",")

		var onlineUpNum float64
		var onlineDownNum float64
		if len(betList) > 0 {

			onlineUpNum, err = strconv.ParseFloat(betList[0], 64)
			if err != nil {
				logger.Err("同步onlineUpNum值：" + err.Error())
			}
			onlineDownNum, err = strconv.ParseFloat(betList[len(betList)-1], 64)
			if err != nil {
				logger.Err("同步onlineDownNum值：" + err.Error())
			}

			setModel := bson.M{
				"BetBase":       i2.Bet,
				"OnlineUpNum":   onlineUpNum,
				"OnlineDownNum": onlineDownNum,
			}
			if strings.HasPrefix(i2.GameID, "pg_") {
				setModel["DefaultCs"] = onlineUpNum
				setModel["DefaultBetLevel"] = 1
			}
			if strings.HasPrefix(i2.GameID, "pp_") {
				setModel["DefaultCs"] = onlineUpNum
			}

			model := &mongo.UpdateManyModel{
				Filter: bson.M{
					"GameId": i2.GameID,
				},
				Update: bson.M{"$set": setModel},
			}
			models = append(models, model)
		}

	}
	_, err = db.Collection2("GameAdmin", "GameConfig").BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(true))

	if err != nil {
		logger.Err("批量处理错误" + err.Error())
		return errors.New("批量处理错误" + err.Error())
	}

	var gameConfigList []*syncGameCofigParams

	err = CollGameConfig.FindAll(bson.M{}, &gameConfigList)
	if err != nil {
		logger.Err("获取全部商户失败" + err.Error())
		return err
	}
	var updataGameCofig updataGameCofigParams
	for _, i2 := range gameConfigList {
		floatBet := convertBetFloat(i2.BetBase)
		updataGameCofig = updataGameCofigParams{
			AppID:          i2.AppID,
			GameId:         i2.GameId,
			GameOn:         i2.GameOn,
			GamePattern:    i2.GamePattern,
			MaxMultiple:    strconv.FormatInt(i2.MaxMultiple, 10),
			MaxWinPoints:   i2.MaxWinPoints,
			RewardPercent:  i2.RewardPercent,
			NoAwardPercent: i2.NoAwardPercent,
			ShowExitBtnOff: i2.ShowExitBtnOff,
		}

		PushMsgStoreSetInfo(updataGameCofig, i2.RewardPercent, i2.NoAwardPercent, floatBet)
	}

	return err
}

func operatorSyncGames() (err error) {

	logger.Infof("启动商户开始执行游戏脚本")
	// ---------------------------------------------------获取商户游戏关联列表-----------------------------------------------------
	var aggregateList []*aggregateData

	GameConfigColl := db.Collection2("GameAdmin", "GameConfig")
	aggregate, err := GameConfigColl.Aggregate(context.TODO(), []bson.M{
		{
			"$group": bson.M{
				"_id": "$AppID",
				"GamesItem": bson.M{
					"$addToSet": "$GameId",
				},
			},
		},
	})
	if err != nil {
		logger.Err("商户游戏关联查询失败：" + err.Error())
		return err
	}

	err = aggregate.All(context.TODO(), &aggregateList)
	if err != nil {
		logger.Err("商户游戏映射结果失败：" + err.Error())
		return err
	}

	// ---------------------------------------------------End-----------------------------------------------------

	// ---------------------------------------------------获取缺失游戏的商户-----------------------------------------------------

	// 缺失游戏的商户
	lackGameMap := map[string][]string{}
	OperatorGameConfigMap := map[string]*comm.Operator_V2{}
	for _, item := range aggregateList {
		lackGameMap[item.Id] = item.GamesItem
	}
	// ---------------------------------------------------End-----------------------------------------------------

	var operator []*comm.Operator_V2

	NewOtherDB("GameAdmin").Collection("AdminOperator").FindAll(bson.M{"OperatorType": 2}, &operator)

	for _, i2 := range operator {
		OperatorGameConfigMap[i2.AppID] = i2
	}

	// ---------------------------------------------------获取缺失游戏的商户-----------------------------------------------------

	for AppID, i := range lackGameMap {
		syncGames(AppID, i, OperatorGameConfigMap[AppID])
	}

	//hml12121707
	return err
}

// 商户同步游戏脚本
func syncGames(AppID string, GameIds []string, operator *comm.Operator_V2) {

	insertList := []bson.M{}

	gamesColl := NewOtherDB("game").Collection("Games")

	var game []*operatorGameConfig
	err := gamesColl.FindAll(bson.M{
		"_id": bson.M{"$nin": GameIds},
	}, &game)

	if err != nil {
		logger.Err("商户" + AppID + "同步游戏时获取商户缺失的所有游戏失败：" + err.Error())
		return
	}

	gameOn := 1
	for _, config := range game {
		gameItem := operatorSetGameConfig(AppID, AppID, *config)

		if config.GameStatus != 0 {
			gameOn = 2
		}
		gameItem["GameOn"] = gameOn

		if strings.ToUpper(config.GameManufacturer) == strings.ToUpper(comm.PP) {

			gameItem["ShowNameAndTimeOff"] = operator.PPConfig.ShowNameAndTimeOff
		}

		if strings.ToUpper(config.GameManufacturer) == strings.ToUpper(comm.PG) {
			gameItem["CarouselOff"] = operator.PGConfig.CarouselOff
			gameItem["ShowNameAndTimeOff"] = operator.PGConfig.ShowNameAndTimeOff
			gameItem["StopLoss"] = operator.PGConfig.CarouselOff
			gameItem["ExitBtnOff"] = operator.PGConfig.ExitBtnOff
			gameItem["ExitLink"] = operator.PGConfig.ExitLink
			gameItem["OfficialVerify"] = operator.PGConfig.OfficialVerify
		}

		if strings.ToUpper(config.GameManufacturer) == strings.ToUpper(comm.JILI) {
			gameItem["BackPackOff"] = operator.JILIConfig.BackPackOff
			gameItem["OpenScreenOff"] = operator.JILIConfig.OpenScreenOff
			gameItem["SidebarOff"] = operator.JILIConfig.SidebarOff
		}
		if strings.ToUpper(config.GameManufacturer) == strings.ToUpper(comm.TADA) {

			gameItem["BackPackOff"] = operator.TADAConfig.BackPackOff
			gameItem["OpenScreenOff"] = operator.TADAConfig.OpenScreenOff
			gameItem["SidebarOff"] = operator.TADAConfig.SidebarOff
		}

		insertList = append(insertList, gameItem)
	}
	if len(insertList) > 0 {

		if err := CollGameConfig.InsertMany(lo.ToAnySlice(insertList)); err != nil {

			logger.Err("商户" + AppID + "批量同步游戏失败：" + err.Error())
		}
	}
}

// 添加游戏时  所有商户同步游戏
func syncOperatorGameConfig(game *comm.Game2) {

	var operList []*comm.Operator_V2

	err := CollAdminOperator.FindAll(bson.M{"OperatorType": 2}, &operList)
	if err != nil {
		logger.Err("获取商户失败未同步游戏：" + err.Error())
		return
	}

	config := operatorGameConfig{}
	config.GameIdd = game.ID

	var insertList []any
	for _, v2 := range operList {
		insertList = append(insertList, operatorSetGameConfig(v2.AppID, v2.UserName, config))
	}

	err = CollGameConfig.InsertMany(insertList)
	if err != nil {
		logger.Err("商户同步游戏失败：" + err.Error())
		return
	}
}
func convertBetFloat(betBase string) []float64 {
	var floatBet []float64

	betList := strings.Split(betBase, ",")
	for _, i := range betList {
		float, err := strconv.ParseFloat(i, 64)
		if err != nil {
			logger.Err("批量处理错误" + err.Error())
		}
		floatBet = append(floatBet, float)
	}

	return floatBet
}
