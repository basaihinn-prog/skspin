package main

import (
	"game/comm"
	"game/comm/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	mux.RegRpc("/AdminInfo/Interior/GetOperatorGameMent", "查询用户按钮状态", "AdminInfo", getOperatorGameMent, OperatorGameMentParams{
		AppID: "",
	})
}

type OperatorGameMentParams struct {
	AppID string `json:"AppID"`
}

type OperatorGameMentRes struct {
	ButStatus int64  `json:"ButStatus"`
	ButLink   string `json:"ButLink"`
}

func getOperatorGameMent(ps OperatorGameMentParams, ret *OperatorGameMentRes) (err error) {
	var oper *comm.Operator_V2
	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &oper)
	if err != nil {
		return
	}
	ret.ButStatus = oper.ShowExitBtnOff
	ret.ButLink = oper.ExitLink
	return

}
