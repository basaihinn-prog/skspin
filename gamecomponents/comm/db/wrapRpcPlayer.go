package db

import (
	"context"
	"fmt"
	"game/comm/define"
	"game/comm/ut"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//	type RpcReqContext struct {
//		Pid      int64
//		Language string
//	}
type RpcReqContext = WSReqContext

// func play(pid int64, ps playPs, ret *playRet) (err error)
// func play(player , ps playPs, ret *playRet) (err error)
func WrapRpcPlayer[PLR, RET any](callee func(*PLR, define.PGParams, *RET) error) func(define.PGParams, *RET) error {
	return func(ps define.PGParams, ret *RET) error {
		// player := &models.Player{
		// 	PID: pid,
		// }
		//pgps := define.PGParams{}
		//fmt.Println(pgps)
		pid := ps.Pid

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

			// doc.Language = ctx.Language
			ut.SetField(player, "DocPlayer", *doc)
		}

		if err != nil {
			return err
		}

		ut.InitNilFields(player)

		// player.Language = ctx.Language
		// ut.SetField(player, "Language", ctx.Language)
		err = callee(player, ps, ret)
		if err != nil {
			return err
		}
		_, err = coll.ReplaceOne(context.TODO(), ID(pid), player, options.Replace().SetUpsert(true))
		return err
	}
}
