package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/mongodb"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	RegMsgProc("/AdminInfo/getSlotPoolHistory", "奖池修改记录", "AdminInfo", getSlotPoolHistory, getSlotPoolHistoryParams{
		Pid: 0,
	})

	RegMsgProc("/AdminInfo/ensureSlotPoolHistoryList", "奖池修改审核记录(列表)", "AdminInfo", ensureSlotPoolHistoryList, ensureSlotPoolHistoryListParams{
		OperatorId:      0,
		StartTime:       0,
		EndTime:         0,
		OpName:          "",
		Change:          0,
		AnimUserPid:     1,
		Type:            1,
		EnsureStatus:    1,
		EnsureStartTime: 0,
		EnsureEndTime:   0,
		pageSetting: pageSetting{
			PageSize:   10,
			PageNumber: 1,
		},
	})
	//todo:奖池记录
	RegMsgProc("/AdminInfo/slotPoolHistoryType", "奖池修改记录审核的类型", "AdminInfo", slotPoolHistoryType, comm.Empty{})

	RegMsgProc("/AdminInfo/doEnsureslotPoolHistory", "审核奖池修改的记录", "AdminInfo", doEnsureslotPoolHistory, doEnsureslotPoolHistoryParams{
		ID:     "",
		Remark: "",
	})
	RegMsgProc("/AdminInfo/getAlertDataHistory", "获取报警记录", "AdminInfo", getAlertDataHistory, AlertDataListRequest{
		PageIndex: 1,
		PageSize:  20,
	})

	//amount
	RegMsgProc("/AdminInfo/setAlertDataHistory", "设置报警记录", "AdminInfo", setAlertDataHistory, AlertDataListRequest{})

}

type getSlotPoolHistoryParams struct {
	Pid  int64
	Type int
}
type slotPoolHistoryResult struct {
	AnimUserName string             `json:"AnimUserName" bson:"AnimUserName"`
	AnimUserPid  int                `json:"AnimUserPid" bson:"AnimUserPid"`
	AppID        string             `json:"AppID" bson:"AppID"`
	Change       int64              `json:"Change" bson:"Change"`
	Currency     string             `json:"Currency" bson:"Currency"`
	NewGold      int64              `json:"NewGold" bson:"NewGold"`
	OldGold      int64              `json:"OldGold" bson:"OldGold"`
	OpName       string             `json:"OpName" bson:"OpName"`
	OpPid        int                `json:"OpPid" bson:"OpPid"`
	Type         int                `json:"Type" bson:"Type"`
	Time         *mongodb.TimeStamp `json:"time" bson:"time"`
	ID           string             `json:"_id" bson:"_id"`
}

type getSlotPoolHistoryResults struct {
	List []*slotPoolHistoryResult
}

func getSlotPoolHistory(ctx *Context, ps getSlotPoolHistoryParams, ret *getSlotPoolHistoryResults) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	query := bson.M{}
	if ps.Pid > 0 {
		query["AnimUserPid"] = ps.Pid
	}
	if len(query) <= 0 {
		return lang.Error(ctx.Lang, "参数错误")
	}
	query["Type"] = ps.Type

	op := options.Find()
	op.SetSort(bson.M{"_id": -1})
	coll := db.Collection2("GameAdmin", "SlotsPoolHistory")
	cursor, _ := coll.Find(context.TODO(), query, op)
	err = cursor.All(context.TODO(), &ret.List)
	return err
}

type ensureSlotPoolHistoryListParams struct {
	OperatorId      int64  //运营商
	StartTime       int64  //变化时间
	EndTime         int64  //变化时间
	OpName          string //操作人
	Change          int64  //变化值
	AnimUserPid     int64  //玩家id
	Type            int64  //1:Slot奖池,2:百人奖池,3:未转移金额,4:玩家状态
	EnsureStatus    int64  //0:未审核,1已审核
	EnsureStartTime int64  //审核时间
	EnsureEndTime   int64  //审核时间
	pageSetting
}

type ensureSlotPoolHistoryListResults struct {
	List []*comm.SlotsPoolHistory
	pageSetting
}

func ensureSlotPoolHistoryList(ctx *Context, ps ensureSlotPoolHistoryListParams, ret *ensureSlotPoolHistoryListResults) (err error) {

	id, err := GetOperatopAppID(ctx)
	if err != nil {
		return err
	}

	if ps.PageNumber == 0 {
		ps.PageNumber = 1
		ret.PageNumber = 1
	} else {
		ret.PageNumber = ps.PageNumber
	}
	if ps.PageSize == 0 {
		ps.PageSize = 50
		ret.PageSize = 50
	} else {
		ret.PageSize = ps.PageSize
	}
	filter := bson.M{}
	err = filter4ensureSlotPoolHistory(filter, ctx.Lang, ps, id)
	if err != nil {
		return
	}
	findOptions := options.Find()
	skip := (ps.PageNumber - 1) * ps.PageSize
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"time": -1})
	findOptions.SetLimit(int64(ps.PageSize))

	coll := DB.Collection("SlotsPoolHistory")
	err = coll.FindAllOpt(filter, &ret.List, findOptions)
	//for _, v := range ret.List {
	//	changeTime, errTime := time.Parse("2006-01-02 15:04:05", v.Time)
	//	if errTime != nil {
	//		return errTime
	//	}
	//	v.Time = changeTime.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
	//	if ps.EnsureStatus == 1 {
	//		ensureTime, errTime := time.Parse("2006-01-02 15:04:05", v.EnsureTime)
	//		if errTime != nil {
	//			return errTime
	//		}
	//		v.EnsureTime = ensureTime.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
	//	}
	//}
	if err != nil {
		return
	}
	count, _ := coll.CountDocuments(filter)
	ret.Count = int(count)
	return err
}

type historyType struct {
	Key  int
	Desc string
}

type slotPoolHistoryTypeResult struct {
	Types []historyType
}

func slotPoolHistoryType(ctx *Context, ps comm.Empty, ret *slotPoolHistoryTypeResult) (err error) {
	ret.Types =
		[]historyType{{Key: 1, Desc: lang.GetLang(ctx.Lang, "Slot奖池")},
			{Key: 2, Desc: lang.GetLang(ctx.Lang, "百人奖池")},
			{Key: 3, Desc: lang.GetLang(ctx.Lang, "未转移金额")},
			{Key: 4, Desc: lang.GetLang(ctx.Lang, "玩家状态")},
		}
	return
}

type doEnsureslotPoolHistoryParams struct {
	ID string
	// EnsureStatus int    //0:未审核,1已审核
	Remark string //审核备注
}

func doEnsureslotPoolHistory(ctx *Context, ps doEnsureslotPoolHistoryParams, ret *comm.Empty) (err error) {
	//filter := bson.M{}
	//ensureFilter := bson.M{}
	//if len(ps.ID) > 19 {
	//	id, err := primitive.ObjectIDFromHex(ps.ID)
	//	if err != nil {
	//		return err
	//	}
	//	filter = bson.M{"_id": id}
	//	ensureFilter = bson.M{"_id": id, "Remark": bson.M{"$exists": false}, "EnsureStatus": bson.M{"$exists": false}}
	//} else {
	//	filter = bson.M{"_id": ps.ID}
	//	ensureFilter = bson.M{"_id": ps.ID, "Remark": bson.M{"$exists": false}, "EnsureStatus": bson.M{"$exists": false}}
	//}

	// 已审核不能再次审核
	ensureFilter := bson.M{"_id": ps.ID, "Remark": bson.M{"$exists": false}, "EnsureStatus": bson.M{"$exists": false}}
	coll := db.Collection2("GameAdmin", "SlotsPoolHistory")
	sr := coll.FindOne(context.TODO(), ensureFilter)
	if sr.Err() != nil {
		return errors.New("此记录已被审核")
	}
	filter := bson.M{"_id": ps.ID}
	update :=
		bson.M{"$set": bson.M{"EnsureStatus": 1,
			"Remark":       ps.Remark,
			"EnsureOpPid":  ctx.PID,
			"EnsureOpName": ctx.Username,
			"EnsureTime":   mongodb.NewTimeStamp(time.Now())},
		}
	opts := options.Update().SetUpsert(true)
	_, err = coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return
	}
	return
}

func filter4ensureSlotPoolHistory(match bson.M, Lang string, ps ensureSlotPoolHistoryListParams, AppIDs []string) (err error) {
	if ps.EnsureStatus == 0 {
		match["EnsureStatus"] = bson.M{"$ne": 1}
	} else if ps.EnsureStatus == 1 {
		match["EnsureStatus"] = 1
	} else {
		return errors.New(lang.GetLang(Lang, "参数错误"))
	}

	appId := ""
	if ps.OperatorId != 0 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.OperatorId, &operator)
		appId = operator.AppID
		if err != nil {
			return err
		}
		temp := false
		for _, appid := range AppIDs {
			if operator.AppID == appid {
				temp = true
			}
		}
		if temp == false {
			return errors.New(lang.GetLang(Lang, "未代理该商户"))
		} else {
			AppIDs = AppIDs[:0]
			AppIDs = append(AppIDs, appId)
		}
	}
	match["AppID"] = bson.M{
		"$in": AppIDs,
	}

	if ps.StartTime != 0 || ps.EndTime != 0 {
		timeTemp := bson.M{}
		match["time"] = timeTemp
		if ps.StartTime != 0 {
			timeTemp["$gte"] = mongodb.NewTimeStamp(time.UnixMilli(ps.StartTime))
		}
		if ps.EndTime != 0 {
			timeTemp["$lte"] = mongodb.NewTimeStamp(time.UnixMilli(ps.EndTime))
		}
	}
	if ps.EnsureStartTime != 0 || ps.EnsureEndTime != 0 {
		timeTemp := bson.M{}
		match["EnsureTime"] = timeTemp
		if ps.EnsureStartTime != 0 {
			timeTemp["$gte"] = mongodb.NewTimeStamp(time.UnixMilli(ps.EnsureStartTime))
		}
		if ps.EnsureEndTime != 0 {
			timeTemp["$lte"] = mongodb.NewTimeStamp(time.UnixMilli(ps.EnsureEndTime))
		}

	}
	if len(ps.OpName) != 0 {
		match["OpName"] = ps.OpName
	}
	if ps.Change != 0 {
		match["Change"] = ps.Change
	}
	if ps.AnimUserPid != 0 {
		match["AnimUserPid"] = ps.AnimUserPid
	}
	if ps.Type != 0 {
		match["Type"] = ps.Type
	}
	return nil
}

type EditAlerterData struct {
	Type  string
	Id    string
	Range string
}
type AlertDataListRequest struct {
	Pid            int64
	GameId         string
	AppID          string
	AlarmTimeStart int64
	AlarmTimeEnd   int64
	PageIndex      int64 `json:"Page"`
	PageSize       int64 `json:"PageSize"`
	Type           string
	Manufacturer   string `json:"Manufacturer"`
}

type AlertDataList struct {
	List  []*AlertData
	Count int64
}
type AlertData struct {
	Id              primitive.ObjectID `bson:"_id" bson:"Id"`
	Pid             int64              `bson:"pid"`
	UserId          string             `bson:"userid"`
	GameId          string             `bson:"gameId"`
	WinMoney        string             `bson:"winMoney"`
	OrderId         string             `bson:"orderId"`
	TotalWinLoss    string             `bson:"totalWinLoss"`
	TotalBet        string             `bson:"totalBet"`
	WinRate         string             `bson:"winRate"`
	Currency        string             `bson:"currency"`
	Balance         string             `bson:"balance"`
	AppId           string             `bson:"appId"`
	Amount          string             `bson:"amount"`
	ReadStatus      int64              `bson:"readStatus"`
	CreateTime      *mongodb.TimeStamp `bson:"createTime"`
	CooperationType int                `bson:"cooperationType"`
	WalletType      int                `bson:"walletType"`
	MerchantBalance float64            `bson:"merchantBalance"`
}

func getAlertDataHistory(ctx *Context, ps AlertDataListRequest, ret *AlertDataList) (err error) {

	user, ok := IsAdminUser(ctx)
	query := bson.M{}

	if ps.AlarmTimeStart != 0 && ps.AlarmTimeEnd != 0 {
		startTime := time.Unix(ps.AlarmTimeStart, 0)
		endTime := time.Unix(ps.AlarmTimeEnd, 0)
		query["createTime"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	if ps.GameId != "ALL" && ps.GameId != "" && ps.Type != "balanceAlert" {
		query["gameId"] = ps.GameId
	}

	if ps.Pid != 0 && ps.Type != "balanceAlert" {
		query["pid"] = ps.Pid
	}
	if ok {
		if ps.AppID != "" {
			query["appId"] = ps.AppID
		}
	} else {
		query["appId"] = user.AppID
	}

	if ps.Type == "transfer" {
		query["amount"] = bson.M{"$exists": true}
		query["winMoney"] = bson.M{"$exists": false}
	}
	if ps.Type == "turnRate" {
		query["amount"] = bson.M{"$exists": false}
		query["winMoney"] = bson.M{"$exists": true}
	}
	if ps.Type == "balanceAlert" {
		query["amount"] = bson.M{"$exists": false}
		query["winMoney"] = bson.M{"$exists": false}
		query["merchantBalance"] = bson.M{"$exists": true}
	}

	alertData := NewOtherDB("reports").Collection("alert")

	filter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Query:    query,
		Sort:     bson.M{"createTime": -1},
	}
	ret.Count, err = alertData.FindPage(filter, &ret.List)
	//for _, v := range ret.List {
	//	createTime, errTime := time.Parse("2006-01-02 15:04:05", v.CreateTime)
	//	if errTime != nil {
	//		return
	//	}
	//	v.CreateTime = createTime.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
	//}

	if err != nil {
		return err
	}

	return nil
}

func setAlertDataHistory(ctx *Context, ps EditAlerterData, ret *AlertDataList) (err error) {

	alertColl := NewOtherDB("reports").Collection("alert")

	query := bson.M{}

	if strings.ToLower(ps.Range) == "all" {

		if ps.Type == "transfer" {
			query["amount"] = bson.M{"$exists": true}
			query["winMoney"] = bson.M{"$exists": false}
		}
		if ps.Type == "turnRate" {
			query["amount"] = bson.M{"$exists": false}
			query["winMoney"] = bson.M{"$exists": true}
		}
		if ps.Type == "balanceAlert" {
			query["amount"] = bson.M{"$exists": false}
			query["winMoney"] = bson.M{"$exists": false}
			query["merchantBalance"] = bson.M{"$exists": true}
		}

		query["readStatus"] = bson.M{"$ne": 1}

		err = alertColl.Update(query, bson.M{"$set": bson.M{"readStatus": 1}})

	} else {

		id, err := primitive.ObjectIDFromHex(ps.Id)
		query["_id"] = id

		if err != nil {
			return errors.New("状态调整失败")
		}

		err = alertColl.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"readStatus": 1}})
	}

	return err
}
