package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	RegMsgProc("/AdminInfo/excludeGameAdd", "运营商不显示游戏配置(添加)", "AdminInfo", addExclude, ExcludeGames{
		Name:           "",
		Remark:         "",
		ExcluedGameIds: []string{},
	})
	RegMsgProc("/AdminInfo/excludeGameUp", "运营商不显示游戏配置(修改)", "AdminInfo", upExclude, ExcludeGames{
		ID:             primitive.ObjectID{},
		Name:           "",
		Remark:         "",
		ExcluedGameIds: []string{},
	})
	RegMsgProc("/AdminInfo/excludeGameAll", "运营商不显示游戏配置(列表)", "AdminInfo", allExclude, ExcludeGames{})
	RegMsgProc("/AdminInfo/excludeGameDel", "运营商不显示游戏配置(删除)", "AdminInfo", delExclude, ExcludeGames{
		ID: primitive.ObjectID{},
	})
}

var excludeGameColl *mongodb.Collection = DB.Collection("ExcludeGame")

type ExcludeGames struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `json:Name,bson:Name`
	Remark         string             `json:Remark,bson:Remark`
	ExcluedGameIds []string           `json:ExcluedGameIds,bson:ExcluedGameIds`
}

func addExclude(ctx *Context, ps ExcludeGames, ret *comm.Empty) (err error) {
	ps.Name = strings.TrimSpace(ps.Name)
	if containName(ps.Name, primitive.NilObjectID) {
		return errors.New(lang.GetLang(ctx.Lang, "名字重复"))
	}
	excludeG := &ExcludeGames{
		ID:             primitive.NewObjectID(),
		Name:           ps.Name,
		Remark:         ps.Remark,
		ExcluedGameIds: ps.ExcluedGameIds,
	}
	err = excludeGameColl.InsertOne(excludeG)
	if err != nil {
		return err
	}
	return nil
}

func upExclude(ctx *Context, ps ExcludeGames, ret *comm.Empty) (err error) {
	ps.Name = strings.TrimSpace(ps.Name)
	if containName(ps.Name, ps.ID) {
		return errors.New(lang.GetLang(ctx.Lang, "名字重复"))
	}
	update := bson.M{
		"$set": bson.M{
			"Name":           ps.Name,
			"Remark":         ps.Remark,
			"ExcluedGameIds": ps.ExcluedGameIds,
		},
	}
	err = excludeGameColl.UpdateOne(bson.M{"_id": ps.ID}, update)
	if err != nil {
		panic(err)
	}
	return nil
}

func allExclude(ctx *Context, _ interface{}, ret *[](*ExcludeGames)) (err error) {
	err = excludeGameColl.FindAll(bson.M{}, ret)
	if err != nil {
		panic(err)
	}
	return
}

func delExclude(ctx *Context, ps ExcludeGames, ret *comm.Empty) (err error) {
	err = excludeGameColl.DeleteId(ps.ID)
	if err != nil {
		panic(err)
	}
	return
}

func containName(Name string, ID primitive.ObjectID) (ret bool) {
	e := &ExcludeGames{}
	ret = false
	excludeGameColl.FindOne(bson.M{"Name": Name}, e)
	if ID.IsZero() {
		if !e.ID.IsZero() {
			ret = true
		}
	} else {
		if !e.ID.IsZero() && ID != e.ID {
			ret = true
		}
	}
	return
}
