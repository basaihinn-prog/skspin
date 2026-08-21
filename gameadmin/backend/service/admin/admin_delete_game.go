package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/DeleteGame", "删除游戏", "AdminInfo", deleteGame, deleteGameParams{
		ID: "游戏Id",
	})
}

type deleteGameParams struct {
	ID string
}

func deleteGame(ctx *Context, ps addGameParams, ret *comm.Empty) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	coll := NewOtherDB("game").Collection("Games")
	err = coll.DeleteOne(bson.M{"_id": ps.ID, "Status": comm.GameStatus_Stop})
	if err != nil {
		return err
	}

	configColl := db.Collection2("GameAdmin", "GameConfig")
	_, err = configColl.DeleteMany(context.TODO(), bson.M{"GameId": ps.ID})
	if err != nil {
		return err
	}
	go initGameConfig()
	go resetGameManufacturer()
	return nil
}
