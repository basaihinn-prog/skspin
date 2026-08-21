package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/UpdateGame", "修改游戏", "AdminInfo", updateGame, updateGameParams{
		ID:     "游戏Id",
		Name:   "游戏名称",
		Type:   0,
		Status: 0,
		Icon:   "游戏图片",
	})
}

type updateGameParams struct {
	ID          string
	Name        string
	Type        int
	Status      int
	Icon        string
	BuyBetMulti int
}

func updateGame(ctx *Context, ps addGameParams, ret *comm.Empty) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	coll := NewOtherDB("game").Collection("Games")

	exists, err := coll.Exists(bson.M{"_id": ps.ID, "Status": ps.Status})
	if err != nil {
		return err
	}

	settingInfo := bson.M{}
	settingInfo["_id"] = ps.ID
	settingInfo["GameId"] = ps.GameId
	settingInfo["LineNum"] = ps.LineNum
	settingInfo["ChangeBetOff"] = ps.ChangeBetOff
	settingInfo["ManufacturerName"] = ps.ManufacturerName
	settingInfo["BuyBetMulti"] = ps.Type
	settingInfo["Type"] = ps.Type
	settingInfo["Status"] = ps.Status
	settingInfo["BetList"] = ps.BetList
	settingInfo["NoAwardPercent"] = ps.NoAwardPercent
	settingInfo["DefaultBet"] = ps.DefaultBet
	settingInfo["DefaultBetLevel"] = ps.DefaultBetLevel
	settingInfo["GameNameConfig"] = ps.GameNameConfig

	updateInfo := bson.M{
		"$set": settingInfo,
	}
	err = coll.UpdateId(ps.ID, updateInfo)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	if !exists {
		updateGameOn := 1
		if ps.Status != 0 {
			updateGameOn = 2
		}

		gameConfigColl := DB.Collection("GameConfig")
		gameConfigColl.Update(bson.M{"GameId": ps.ID}, bson.M{"$set": bson.M{"GameOn": updateGameOn}})
	}

	go initGameConfig()
	go resetGameManufacturer()
	return err
}
