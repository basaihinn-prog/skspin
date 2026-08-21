package main

import (
	"errors"
	"fmt"
	"game/comm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	RegMsgProc("/AdminInfo/GetSlotsPool", "获取Slots", "AdminInfo", getSlotsPool, getSlotsPoolParams{
		Pid:  0,
		Uid:  "",
		Type: 0,
		Gold: 0,
	})
}

type getSlotsPoolParams struct {
	Pid   int64
	Uid   string
	AppId string
	Type  int64
	Gold  int64
}

type getSlotsPoolResults struct {
	Gold int64
}

func getSlotsPool(ctx *Context, ps setSlotsPoolParams, ret *getSlotsPoolResults) (err error) {
	filter := bson.M{}
	if ps.Uid != "" {
		filter["Uid"] = ps.Uid
	}

	if ps.Pid > 0 {
		filter["_id"] = ps.Pid
	}
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}
	if user.AppID == "admin" {
		if ps.AppId != "" {
			filter["AppID"] = ps.AppId
		}
	} else {
		filter["AppID"] = user.AppID
	}

	var player comm.Player
	err = NewOtherDB("game").Collection("Players").FindOne(filter, &player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}

	if ps.Type != 1 && ps.Type != 2 {
		return errors.New(fmt.Sprintf("error type, type:%d", ps.Type))
	}

	id := fmt.Sprintf("%d-%d", player.Id, ps.Type)
	var slotsPool comm.SlotsPool
	NewOtherDB("slots").Collection("Pool").FindId(id, &slotsPool)
	ret.Gold = slotsPool.Gold
	return err
}
