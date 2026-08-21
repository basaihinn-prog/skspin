package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/ut"
	"game/duck/lang"
	"game/duck/logger"
	"game/duck/ut2/fileutil"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	RegMsgProc("/AdminInfo/AddGame", "添加游戏", "AdminInfo", addGame, addGameParams{
		ID:     "游戏Id",
		Name:   "游戏名称",
		Type:   0,
		Status: 0,
		Icon:   "游戏图片",
	})
}

type addGameParams struct {
	ID               string
	GameId           string
	BetList          string
	NoAwardPercent   int64
	DefaultBet       float64
	DefaultBetLevel  int64
	Name             string
	Type             int
	Status           int
	Icon             string
	LineNum          float64
	ChangeBetOff     int64
	GameNameConfig   map[string]*comm.GameConfigItem
	ManufacturerName string
}

func addGame(ctx *Context, ps addGameParams, ret *comm.Empty) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	coll := db.Collection2("game", "Games")
	data := []*comm.Game2{}
	find, err := coll.Find(context.TODO(), bson.M{"_id": ps.ID})
	if err != nil {
		return err
	}

	err = find.All(context.TODO(), &data)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		return errors.New(lang.GetLang(ctx.Lang, "游戏ID已存在"))
	}

	insertData := &comm.Game2{
		ID:               ps.ID,
		GameID:           ps.GameId,
		Name:             ps.Name,
		Type:             ps.Type,
		Status:           ps.Status,
		LineNum:          ps.LineNum,
		ChangeBetOff:     ps.ChangeBetOff,
		ManufacturerName: ps.ManufacturerName,
		GameNameConfig:   ps.GameNameConfig,
		BetList:          ps.BetList,
		NoAwardPercent:   ps.NoAwardPercent,
		DefaultBet:       ps.DefaultBet,
		DefaultBetLevel:  ps.DefaultBetLevel,
	}
	if insertData.ID == "" {

		return errors.New("插入失败ID不能为：" + insertData.ID)
	}
	_, err = coll.InsertOne(context.TODO(), insertData)

	if err != nil {
		return err
	}

	go initGameConfig()
	go operatorSyncGames()
	go resetGameManufacturer()
	return nil
}

func uploadGameFile(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")

	excel, err := fileutil.ReadExcel(file)

	if err != nil {
		ut.HttpReturnJson(w, map[string]string{
			"code":  "0",
			"error": "上传错误" + err.Error(),
			"data":  "",
		})
	}

	gameColl := db.Collection2("game", "Games")

	// =================================表中的原数据分类  状态为正常与非正常游戏(关闭   隐藏等)====================================

	var originGameData []*comm.Game2
	originGameStatusSyncConfig := map[string]string{}
	originGameConfigSyncConfig := map[string]map[string]*comm.GameConfigItem{}

	find, err := gameColl.Find(context.TODO(), bson.M{})
	if err != nil {
		ut.HttpReturnJson(w, map[string]string{
			"code":  "1",
			"error": "上传查询全部游戏错误" + err.Error(),
			"data":  "",
		})
		return
	}

	err = find.All(context.TODO(), &originGameData)
	if err != nil {
		ut.HttpReturnJson(w, map[string]string{
			"code":  "1",
			"error": "上传查询映射错误" + err.Error(),
			"data":  "",
		})
		return
	}

	for _, i2 := range originGameData {
		originGameStatusSyncConfig[i2.ID] = strconv.Itoa(i2.Status)
		originGameConfigSyncConfig[i2.ID] = i2.GameNameConfig
	}

	gameIdList := ""
	var gameExitsId []string
	var gameWriteModel []interface{}
	// 通过游戏状态分类   分为正常  与 非正常数据
	gameStatusSyncConfig := map[string][]string{
		"0": nil,
		"1": nil,
		"2": nil,
	}

	for _, m := range excel {
		data := generatorGameData(m, originGameConfigSyncConfig[m["_id"]])
		if data == nil {
			continue
		}

		if m["Status"] == "0" {
			if m["Status"] != originGameStatusSyncConfig[m["_id"]] {

				gameStatusSyncConfig["1"] = append(gameStatusSyncConfig["1"], m["_id"])
			}

		} else {
			gameStatusSyncConfig["2"] = append(gameStatusSyncConfig["2"], m["_id"])
		}

		if originGameStatusSyncConfig[m["_id"]] != "" {
			delete(originGameStatusSyncConfig, m["_id"])
		}

		index := strings.Index(gameIdList, m["_id"]+",")
		if index < 0 {

			gameIdList += m["_id"] + ","
			gameWriteModel = append(gameWriteModel, data)
		} else {
			if gameExitsId == nil {
				gameExitsId = []string{}
			}
			gameExitsId = append(gameExitsId, m["_id"])
		}
	}

	// ================================= 清空表中元数据 ====================================
	_, err = gameColl.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		ut.HttpReturnJson(w, map[string]string{
			"code": "0",
			"msg":  "上传错误" + err.Error(),
			"data": "",
		})
		return
	}

	_, err = gameColl.InsertMany(context.TODO(), gameWriteModel)
	if err != nil {
		ut.HttpReturnJson(w, map[string]string{
			"code":  "1",
			"error": "上传错误" + err.Error(),
			"data":  "",
		})
		return
	}

	gameIds := []string{}
	for gid, _ := range originGameStatusSyncConfig {
		gameIds = append(gameIds, gid)
	}
	db.Collection2("GameAdmin", "GameConfig").DeleteMany(context.TODO(), bson.M{"GameId": bson.M{"$in": gameIds}})

	syncGamesConfigStatus(gameStatusSyncConfig)

	initGameConfig()

	resetGameManufacturer()

	err = operatorSyncGames()
	if err != nil {
		ut.HttpReturnJson(w, map[string]string{
			"code":  "1",
			"error": "上传失败" + err.Error(),
			"data":  "",
		})
		return
	}

	msg := "上传成功"
	if len(gameExitsId) > 0 {
		msg += fmt.Sprintf("，游戏ID重复：%s", strings.Join(gameExitsId, ","))
	}

	ut.HttpReturnJson(w, map[string]string{
		"code":  "0",
		"error": msg,
		"data":  "",
	})
}

func syncGamesConfigStatus(gameStatusSyncConfig map[string][]string) {

	var err error

	dbColl := db.Collection2("GameAdmin", "GameConfig")

	// 状态为开启
	_, err = dbColl.UpdateMany(context.TODO(), bson.M{
		"GameId": bson.M{
			"$in": gameStatusSyncConfig["0"],
		},
	}, bson.M{
		"$set": bson.M{
			"GameOn": 0,
		},
	})

	// 状态为开启
	_, err = dbColl.UpdateMany(context.TODO(), bson.M{
		"GameId": bson.M{
			"$in": gameStatusSyncConfig["1"],
		},
	}, bson.M{
		"$set": bson.M{
			"GameOn": 1,
		},
	})

	// 状态为关闭
	_, err = dbColl.UpdateMany(context.TODO(), bson.M{
		"GameId": bson.M{
			"$in": gameStatusSyncConfig["2"],
		},
	}, bson.M{
		"$set": bson.M{
			"GameOn": 2,
		},
	})

	if err != nil {
		return
	}

}

func generatorGameData(m map[string]string, origin map[string]*comm.GameConfigItem) interface{} {

	gameIconAndNameMap := map[string]*comm.GameConfigItem{}
	lineNum, err := strconv.Atoi(m["LineNum"])

	if err != nil {
		logger.Err("线数错误：" + err.Error())
	}

	gameType, err := strconv.Atoi(m["Type"])

	if err != nil {
		logger.Err("游戏类型错误：" + err.Error())
	}

	Status, err := strconv.Atoi(m["Status"])

	if err != nil {
		logger.Err("状态错误：" + err.Error())
	}

	NoAwardPercent, err := strconv.Atoi(m["NoAwardPercent"])
	if err != nil {
		logger.Err("NoAwardPercent数值错误：" + err.Error())
	}
	RewardPercent, err := strconv.Atoi(m["RewardPercent"])
	if err != nil {
		logger.Err("RewardPercent数值错误：" + err.Error())
	}
	DefaultBet, err := strconv.ParseFloat(m["DefaultBet"], 64)
	if err != nil {
		logger.Err("DefaultBet数值错误：" + err.Error())
	}
	DefaultBetLevel, err := strconv.ParseFloat(m["DefaultBetLevel"], 64)
	if err != nil {
		logger.Err("DefaultBetLevel数值错误：" + err.Error())
	}

	ChangeBetOff, err := strconv.Atoi(m["ChangeBetOff"])
	if err != nil {
		logger.Err("改变可修改Bet状态字段错误：" + err.Error())
	}

	BuyType, err := strconv.Atoi(m["BuyType"])
	if err != nil {
		logger.Err("改变可修改BuyType状态字段错误：" + err.Error())
	}

	if err != nil {
		return nil
	}

	// 语言名称与图标配置
	gameIconAndNameMap["zh"] = &comm.GameConfigItem{
		GameName: m["zh"],
	}
	gameIconAndNameMap["en"] = &comm.GameConfigItem{
		GameName: m["en"],
	}
	gameIconAndNameMap["idr"] = &comm.GameConfigItem{
		GameName: m["idr"],
	}
	gameIconAndNameMap["it"] = &comm.GameConfigItem{
		GameName: m["it"],
	}
	gameIconAndNameMap["es"] = &comm.GameConfigItem{
		GameName: m["es"],
	}
	gameIconAndNameMap["th"] = &comm.GameConfigItem{
		GameName: m["th"],
	}

	if origin != nil {

		gameIconAndNameMap["zh"].Icon = origin["zh"].Icon

		gameIconAndNameMap["en"].Icon = origin["en"].Icon

		gameIconAndNameMap["idr"].Icon = origin["idr"].Icon

		gameIconAndNameMap["it"].Icon = origin["it"].Icon

		gameIconAndNameMap["es"].Icon = origin["es"].Icon

		gameIconAndNameMap["th"].Icon = origin["th"].Icon
	}
	return bson.M{
		"_id":              m["_id"],
		"BetList":          m["Bet"],
		"ChangeBetOff":     ChangeBetOff,
		"BuyType":          BuyType,
		"GameId":           m["GameId"],
		"LineNum":          lineNum,
		"ManufacturerName": m["ManufacturerName"],
		"Name":             m["Name"],
		"NoAwardPercent":   NoAwardPercent,
		"RewardPercent":    RewardPercent,
		"BuyBetMulti":      0,
		"Status":           Status,
		"DefaultBet":       DefaultBet,
		"DefaultBetLevel":  DefaultBetLevel,
		"GameNameConfig":   gameIconAndNameMap,
		"Type":             gameType,
	}
}
