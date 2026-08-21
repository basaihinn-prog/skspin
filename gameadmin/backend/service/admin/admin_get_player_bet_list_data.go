package main

/*
import (
	"context"
	"encoding/json"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/duck/mongodb"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func init() {
	mux.RegHttpWithSample("/AdminInfo/GetPlayerBetLog", "获取玩家下注记录", "player", getPlayerBetLogList, &getPlayerBetLogParams{})
	mux.RegHttpWithSample("/AdminInfo/GetPlayerSpinDetails", "获取玩家下注详情", "player", getPlayerSpinDetails, &getPlayerSpinDetailsParams{
		BetID: lo.Must(primitive.ObjectIDFromHex("6548a3a02fe7bff0fd271de9")),
	})
}

type getPlayerBetLogParams struct {
	GameID    string
	Pid       int
	PageIndex int64
	PageSize  int64
}

type GetPlayerBetLogResults struct {
	Count int64
	List  []*slotsmongo.DocBetLog
}

//mx.HandleFunc("/AdminInfo/GetPlayerBetLog", GetPlayerBetLogList)
//	mx.HandleFunc("/AdminInfo/GetPlayerSpinDetails", GetPlayerSpinDetails)

func getPlayerBetLogList(_ *http.Request, ps getPlayerBetLogParams, ret *GetPlayerBetLogResults) (err error) {
	if ps.GameID == "" || ps.Pid <= 0 {
		return
	}
	query := bson.M{}
	query["GameID"] = ps.GameID
	query["Pid"] = ps.Pid

	if ps.PageSize > 20 {
		ps.PageSize = 20
	}

	filter := mongodb.FindPageOpt{
		Page:       ps.PageIndex,
		PageSize:   ps.PageSize,
		Sort:       db.D("_id", -1),
		Query:      query,
		Projection: bson.M{"SpinDetailsJson": 0},
	}
	count, err := NewOtherDB("reports").Collection("BetLog").FindPage(filter, &ret.List)
	if err != nil {
		return
	}
	ret.Count = count
	return
}

type getPlayerSpinDetailsParams struct {
	BetID primitive.ObjectID
}

func getPlayerSpinDetails(_ *http.Request, ps *getPlayerSpinDetailsParams, ret *json.RawMessage) (err error) {
	var doc struct {
		ID              primitive.ObjectID `bson:"_id"`
		SpinDetailsJson string             `bson:"SpinDetailsJson"`
	}
	coll := db.Collection2("reports", "BetLog")
	err = coll.FindOne(context.TODO(), db.ID(ps.BetID)).Decode(&doc)
	if err != nil {
		return
	}
	*ret = json.RawMessage(doc.SpinDetailsJson)
	return
}
*/
