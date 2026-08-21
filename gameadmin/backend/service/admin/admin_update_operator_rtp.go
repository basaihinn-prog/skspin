package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/UpdateOperatorRtp", "设置运营商rtp", "AdminInfo", updateOperatorRtp, UpdateOperatorRtpParams{
		OperatorId: 0,
		Rtp:        0,
	})
}

type UpdateOperatorRtpParams struct {
	OperatorId int64
	Rtp        int
}

func updateOperatorRtp(ctx *Context, ps UpdateOperatorRtpParams, ret *comm.Empty) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	err = CollAdminOperator.UpdateId(ps.OperatorId, bson.M{
		"$set": bson.M{
			"CurrentRtp": ps.Rtp,
		},
	})

	return err
}
