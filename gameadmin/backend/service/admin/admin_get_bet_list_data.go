package main

import (
	"context"
	"encoding/json"
	"game/comm"
	"game/comm/db"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/mongodb"
	"game/service/pggateway/pgcomm"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	RegMsgProc("/AdminInfo/GetBetLogList", "获取下注列表", "AdminInfo", getBetLogList, getBetLogListParams{
		PageIndex: 1,
		PageSize:  10,
	})

	RegMsgProc("/AdminInfo/GetSpinDetails", "获取下注详情", "AdminInfo", getSpinDetails, getSpinDetailsPs{
		BetID: "",
	})
}

type getBetLogListParams struct {
	OperatorId int64
	OrderId    string
	Pid        int64
	UserID     string
	StartTime  int64
	EndTime    int64
	GameID     string
	Bet        int64
	Win        int64
	WinLose    int64
	Balance    int64
	PageIndex  int64
	PageSize   int64
}

type getBetLogListResults struct {
	Count int64
	List  []*slotsmongo.DocBetLog
}

func getBetLogList(ctx *Context, ps getBetLogListParams, ret *getBetLogListResults) (err error) {
	query := bson.M{}
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}
	if user.AppID == "admin" {
		if ps.OperatorId > 0 {
			var operator *comm.Operator
			err = CollAdminOperator.FindId(ps.OperatorId, &operator)
			if err != nil {
				return err
			}
			query["AppID"] = operator.AppID
		}
	} else {
		query["AppID"] = user.AppID
	}
	if ps.OrderId != "" {
		query["_id"], _ = primitive.ObjectIDFromHex(ps.OrderId)
	}
	if ps.Pid > 0 {
		query["Pid"] = ps.Pid
	}
	if ps.UserID != "" {
		query["UserID"] = ps.UserID
	}
	var startT, endT time.Time
	if ps.StartTime > 0 && ps.EndTime > 0 {
		startT = time.Unix(ps.StartTime, 0)
		endT = time.Unix(ps.EndTime, 0)
		query["InsertTime"] = bson.M{
			"$gte": startT,
			"$lte": endT,
		}
	}
	//else {
	//	endT = time.Now()
	//	startT = endT.Add(-24 * time.Hour)
	//}
	//query["InsertTime"] = bson.M{
	//	"$gte": startT,
	//	"$lte": endT,
	//}
	if ps.GameID != "" {
		query["GameID"] = ps.GameID
	}
	if ps.Bet > 0 {
		query["Bet"] = ps.Bet
	}
	if ps.Win > 0 {
		query["Win"] = ps.Win
	}
	if ps.WinLose > 0 {
		query["WinLose"] = ps.WinLose
	}
	if ps.Balance > 0 {
		query["Balance"] = ps.Balance
	}

	query["LogType"] = db.D("$not", db.D("$gt", 0))

	sort := bson.D{}
	//if ps.SortBet != 0 {
	//	sort = append(sort, bson.E{Key: "Bet", Value: ps.SortBet})
	//}
	//if ps.SortWin != 0 {
	//	sort = append(sort, bson.E{Key: "Win", Value: ps.SortWin})
	//}
	//if ps.SortWinLose != 0 {
	//	sort = append(sort, bson.E{Key: "WinLose", Value: ps.SortWinLose})
	//}
	//if ps.SortBalance != 0 {
	//	sort = append(sort, bson.E{Key: "Balance", Value: ps.SortBalance})
	//}

	if len(sort) == 0 {
		sort = append(sort, bson.E{Key: "_id", Value: -1})
	}

	filter := mongodb.FindPageOpt{
		Page:       ps.PageIndex,
		PageSize:   ps.PageSize,
		Sort:       sort,
		Query:      query,
		Projection: bson.M{"SpinDetailsJson": 0},
	}
	// {SpinDetails: 0}
	count, err := NewOtherDB("reports").Collection("BetLog").FindPage(filter, &ret.List, false)
	if err != nil {
		return err
	}
	count = 1000
	ret.Count = count

	return nil
}

type getSpinDetailsPs struct {
	BetID string
	AppID string
}

func getSpinDetails(ctx *Context, ps getSpinDetailsPs, ret *json.RawMessage) (err error) {
	coll := db.Collection2("reports", "BetLog")

	if ps.AppID != "" {
		coll = slotsmongo.GetBetLogColl(ps.AppID)
	}

	var doc struct {
		ID              string `bson:"_id"`
		GameID          string `bson:"GameID"`
		SpinDetailsJson string `bson:"SpinDetailsJson"`
	}
	err = coll.FindOne(context.TODO(), db.ID(ps.BetID)).Decode(&doc)
	if err != nil {
		return
	}
	if strings.HasPrefix(doc.GameID, "pg_") {
		coll2 := db.Collection2(doc.GameID, "BetHistory")
		var objId primitive.ObjectID
		err = json.Unmarshal([]byte(doc.SpinDetailsJson), &objId)
		if err != nil {
			return err
		}

		var data pgcomm.BHItem
		err = coll2.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&data)
		if err != nil {
			return err
		}
		*ret, err = ut.GetJsonRaw(data)
		return err
	}

	*ret = json.RawMessage(doc.SpinDetailsJson)
	return err
}
