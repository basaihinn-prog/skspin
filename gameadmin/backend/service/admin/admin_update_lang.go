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
	RegMsgProc("/AdminInfo/LangUpdate", "修改多语言", "AdminInfo", updateLang, &updateLangParams{
		Lang: lang.Lang{
			ZH: "",
		},
	})
}

type updateLangParams struct {
	lang.Lang
}

func updateLang(ctx *Context, ps updateLangParams, ret *comm.Empty) (err error) {

	if _, ok := IsAdminUser(ctx); !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	coll := db.Collection2("game", "Lang")
	setdata := []any{
		"permission", ps.Permission,
	}
	for k, v := range ps.LangMap {
		coll.UpdateMany(context.TODO(), bson.M{k: bson.M{"$exists": false}}, bson.M{"$set": bson.M{k: ""}})
		setdata = append(setdata, k, v)
	}
	update := db.D(
		"$set", db.D(setdata...),
	)
	_, err = coll.UpdateByID(context.TODO(), ps.ZH, update)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "数据库错误"))
	}

	return err
}
