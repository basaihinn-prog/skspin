package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/ut"
	"game/duck/lang"
	"game/duck/mongodb"
	"log/slog"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	RegMsgProc("/AdminInfo/SetSlotsPool", "设置Slots", "AdminInfo", setSlotsPool, setSlotsPoolParams{
		Pid:  0,
		Uid:  "",
		Type: 0,
		Gold: 0,
	})
}

type setSlotsPoolParams struct {
	Pid   int64
	Uid   string
	AppId string
	Type  int64
	Gold  int64
}

func setSlotsPool(ctx *Context, ps setSlotsPoolParams, ret *comm.Empty) (err error) {
	currentUser, _ := IsAdminUser(ctx)
	if currentUser.GroupId > 1 {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	filter := bson.M{}
	if ps.Uid != "" {
		filter["Uid"] = ps.Uid
	}

	if ps.Pid > 0 {
		filter["_id"] = ps.Pid
	}
	var user comm.User
	err = CollAdminUser.FindId(ctx.PID, &user)
	if err != nil {
		return err
	}
	if user.AppID == "admin" {
		filter["AppID"] = ps.AppId
	} else {
		filter["AppID"] = user.AppID
	}

	var animUser comm.Player
	err = NewOtherDB("game").Collection("Players").FindOne(filter, &animUser)
	if err != nil {
		return err
	}

	if ps.Type != 1 && ps.Type != 2 {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	id := fmt.Sprintf("%d-%d", animUser.Id, ps.Type)
	var slotsPool *comm.SlotsPool
	err = NewOtherDB("slots").Collection("Pool").FindId(id, &slotsPool)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = NewOtherDB("slots").Collection("Pool").InsertOne(comm.SlotsPool{
				ID:   id,
				Pid:  animUser.Id,
				Type: ps.Type,
				Gold: 0,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err = NewOtherDB("slots").Collection("Pool").UpdateId(id, bson.M{
		"$set": bson.M{"Gold": ps.Gold},
	})
	if err != nil {
		return err
	}

	if slotsPool == nil {
		slotsPool = &comm.SlotsPool{}
	}
	slog.Info("update player slots", "UserName", ctx.Username, "Pid", ctx.PID, "animUserPid", animUser.Id, "TypeInfo", ps.Type, "OldGold", slotsPool.Gold, "NewGold", ps.Gold)

	coll := db.Collection2("GameAdmin", "AdminOperator")
	operator := &comm.Operator{}
	err = coll.FindOne(context.TODO(), bson.M{"AppID": animUser.AppID}).Decode(&operator)
	if err != nil {
		return err
	}

	coll = db.Collection2("GameAdmin", "SlotsPoolHistory")
	sf := ut.NewSnowflake()
	sfId := strconv.Itoa(int(sf.NextID()))
	coll.InsertOne(context.TODO(), bson.M{
		"_id":          sfId,
		"OpName":       ctx.Username,
		"OpPid":        ctx.PID,
		"AnimUserPid":  animUser.Id,
		"AnimUserName": animUser.Uid,
		"Type":         ps.Type,
		"OldGold":      slotsPool.Gold,
		"NewGold":      ps.Gold,
		"Change":       ps.Gold - slotsPool.Gold,
		"time":         mongodb.NewTimeStamp(time.Now()),
		"Currency":     operator.CurrencyKey,
		"AppID":        operator.AppID,
	})
	return err
}
