package main

import (
	"game/comm"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/getRPTGear", "获取Slots", "AdminInfo", getRPTGear, comm.Empty{})
}

type studnet struct {
	Name string `json:"Name"`
}
type getRPTGearres struct {
	List  []*comm.RTPGear
	Count int
}

// func getOperatorListV2(ctx *Context, ps getOperatorListParamsV2, ret *getOperatorListResultV2) (err error) {
func getRPTGear(ctx *Context, ps comm.Empty, ret *getRPTGearres) (err error) {

	coll := NewOtherDB("GameAdmin").Collection("RTPGear")

	err = coll.FindAll(bson.M{}, &ret.List)
	ret.Count = len(ret.List)

	return err
}
