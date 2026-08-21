package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	RegMsgProc("/AdminInfo/UpdateOperatorPermission", "修改运营商权限", "AdminInfo", updateOperatorPermission, updateOperatorPermissionParams{
		OperatorId: 1,
		MenuIds:    []int64{1, 2},
	})
}

type updateOperatorPermissionParams struct {
	OperatorId     int64
	MenuIds        []int64
	AppSecret      string
	Address        string
	WhiteIps       []string
	CurrencyKey    string
	ExcluedGameId  primitive.ObjectID
	ExcluedGameIds []string
	Robot          string //telegram 机器人
	ChatID         string //telegram 聊天群id
}

func updateOperatorPermission(ctx *Context, ps updateOperatorPermissionParams, ret *comm.Empty) (err error) {
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}
	if user.AppID != "admin" {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	var operator *comm.Operator
	err = CollAdminOperator.FindId(ps.OperatorId, &operator)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	ps.MenuIds = NormalMenu(ps.MenuIds)

	addIds := make([]int64, 0, 1)
	for i := range ps.MenuIds {
		if !lo.Contains(operator.MenuIds, ps.MenuIds[i]) {
			addIds = append(addIds, ps.MenuIds[i])
		}
	}
	delIds := make([]int64, 0, 1)
	for i := range operator.MenuIds {
		if !lo.Contains(ps.MenuIds, operator.MenuIds[i]) {
			delIds = append(delIds, operator.MenuIds[i])
		}
	}
	err = CollAdminOperator.UpdateId(ps.OperatorId, bson.M{
		"$set": bson.M{
			"MenuIds":        ps.MenuIds,
			"AppSecret":      ps.AppSecret,
			"Address":        ps.Address,
			"WhiteIps":       ps.WhiteIps,
			"CurrencyKey":    ps.CurrencyKey,
			"ExcluedGameId":  ps.ExcluedGameId,
			"ExcluedGameIds": ps.ExcluedGameIds,
			"Robot":          ps.Robot,
			"ChatID":         ps.ChatID,
		},
	})
	if err != nil {
		return err
	}

	for i := range delIds {
		CollAdminPermission.Update(
			bson.M{
				"OperatorId": ps.OperatorId,
				"$in": bson.M{
					"MenuIds": delIds[i],
				},
			},
			bson.M{
				"$pull": bson.M{"MenuIds": delIds[i]},
			},
		)
	}
	CollAdminPermission.UpdateId(operator.PermissionId, bson.M{
		"$set": bson.M{
			"MenuIds":    ps.MenuIds,
			"UpdateTime": mongodb.NewTimeStamp(time.Now()),
		},
	})

	return nil
}
