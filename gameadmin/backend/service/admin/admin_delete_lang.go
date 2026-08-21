package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/LangDelete", "删除多语言", "AdminInfo", deleteLang, &deleteLangParams{
		Id: "",
	})
}

type deleteLangParams struct {
	Id string
}

func deleteLang(ctx *Context, ps deleteLangParams, ret *comm.Empty) (err error) {
	if _, ok := IsAdminUser(ctx); !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	coll := db.Collection2("game", "Lang")
	hex, err := mongodb.ObjectIDFromHex(ps.Id)
	if err != nil {
		return err
	}
	_, err = coll.DeleteOne(context.TODO(), bson.M{"_id": hex})
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "数据库错误"))
	}
	return err
}
