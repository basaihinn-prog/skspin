package db

import (
	"context"
	"fmt"
	"game/comm/ut"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WSReqContext struct {
	Pid      int64
	Language string
}

var (
	WSReqContextType = reflect.TypeOf(&WSReqContext{})
)

// func play(pid int64, ps playPs, ret *playRet) (err error)
// func play(player , ps playPs, ret *playRet) (err error)
func WrapPlayer[PLR, PS, RET any](callee func(*PLR, PS, *RET) error) func(*WSReqContext, PS, *RET) error {
	return func(ctx *WSReqContext, ps PS, ret *RET) error {
		// player := &models.Player{
		// 	PID: pid,
		// }
		pid := ctx.Pid

		var player = new(PLR)

		coll := Collection("players")
		err := coll.FindOne(context.TODO(), ID(pid)).Decode(player)
		if err == mongo.ErrNoDocuments {
			err = nil
			doc, err := getDocPlayer(pid)
			if err != nil {
				err = fmt.Errorf("cannot found player[%v] data, err [%s]", pid, err.Error())
				return err
			}
			doc.Language = ctx.Language
			ut.SetField(player, "DocPlayer", *doc)

			// ut.SetField(player, "PID", pid)
			// ut.SetField(player, "Uid", doc.Uid)
			// ut.SetField(player, "AppID", doc.AppID)
		}

		if err != nil {
			return err
		}

		ut.InitNilFields(player)

		// player.Language = ctx.Language
		ut.SetField(player, "Language", ctx.Language)
		err = callee(player, ps, ret)
		if err != nil {
			return err
		}
		_, err = coll.ReplaceOne(context.TODO(), ID(pid), player, options.Replace().SetUpsert(true))
		return err
	}
}
