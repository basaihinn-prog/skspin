package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/comm/ut"
	"game/duck/lang"
	"game/duck/mongodb"
	"log/slog"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	//玩家管理 / 玩家信息
	RegMsgProc("/AdminInfo/GetPlayerList", "获取玩家列表", "AdminInfo", getPlayerList, getPlayerListParams{
		PageIndex:  1,
		PageSize:   10,
		Uid:        "",
		OperatorId: 0,
		Pid:        0,
		SortBet:    0,
		SortWin:    0,
	})

	RegMsgProc("/AdminInfo/UpdatePlayerStatus", "玩家启用/禁用", "AdminInfo", UpdatePlayerStatus, UpdatePlayerStatusParams{
		Pid:    0,
		Status: 1,
	})

	RegMsgProc("/AdminInfo/GetPlayerInfo", "获取玩家详情", "AdminInfo", getPlayerInfo, getPlayerListParams{
		Uid:        "",
		OperatorId: 0,
		Pid:        0,
	})
}

type getPlayerListParams struct {
	PageIndex  int64
	PageSize   int64
	Uid        string
	OperatorId int64
	Pid        int64
	SortBet    int64
	SortWin    int64
	//Status         *int  // 0:启用,1:禁用 (当0时查询所有，1时查询禁用)(兼容玩家信息页面和冻结玩家信息页面)

	OnLineStatus int64

	Status         *int64 // 1:启用,0:禁用 (当0时查询所有，1时查询禁用)(兼容玩家信息页面和冻结玩家信息页面)  2024 11 1 改动
	CloseTimeStart int64  //禁用(冻结)时间开始
	CloseTimeEnd   int64  //禁用(冻结)时间结束
	CurrencyCode   string //币种
}

type getPlayerListResult struct {
	List     []*player
	AllCount int64
}

type player struct {
	Id                   int64                  `json:"Pid" bson:"_id"`
	Uid                  string                 `json:"Uid" bson:"Uid"`
	AppID                string                 `json:"AppID" bson:"AppID" md:"运营商"`
	CurrencyKey          string                 `json:"CurrencyKey" bson:"CurrencyKey"`
	CurrencyName         string                 `json:"CurrencyName" bson:"CurrencyName"`
	LoginAt              *mongodb.TimeStamp     `json:"LoginAt" bson:"LoginAt"`
	Status               int64                  `json:"Status" bson:"Status"`
	UName                string                 `json:"UName" bson:"UName"`
	CloseTime            *mongodb.TimeStamp     `json:"CloseTime" bson:"CloseTime"`
	CreateAt             *mongodb.TimeStamp     `json:"CreateAt" bson:"CreateAt"`
	TypeInfo             map[int]*comm.TypeInfo `json:"TypeInfo" bson:"TypeInfo"`
	Bet                  int64                  `json:"Bet" bson:"Bet"`
	Win                  int64                  `json:"Win" bson:"Win"`
	Multi                float64                `json:"Multi" bson:"Multi" description:"玩家当前赢取倍数"`
	LastSpinTime         string                 `json:"LastSpinTime" bson:"LastSpinTime"`
	OnLineStatus         int64                  `json:"OnLineStatus"`
	Balance              int64                  `json:"Balance" bson:"Balance"`
	RestrictionsStatus   int64                  `json:"RestrictionsStatus" bson:"RestrictionsStatus" description:"玩家约束状态"`
	RestrictionsMaxWin   int64                  `json:"RestrictionsMaxWin" bson:"RestrictionsMaxWin" description:"玩家约束最大赢取金额"`
	RestrictionsMaxMulti float64                `json:"RestrictionsMaxMulti" bson:"RestrictionsMaxMulti" description:"玩家约束最大赢取倍数"`
	RestrictionsTime     *mongodb.TimeStamp     `json:"RestrictionsTime" bson:"RestrictionsTime" description:"玩家约束设置时间"`
}

type playAndMerchantInfo struct {
	Operator             *comm.Operator
	OnlineOperator       *comm.Operator
	OnlineStatus         bool
	RTPControlStatus     bool
	PlayerRTPControllist []*comm.PlayerRTPControlModel
	player
}

type getPlayerListResultV2 struct {
	List     []*comm.Player
	AllCount int64
}

type curr struct {
	Name string `bson:"CurrencyName"`
	Code string `bson:"CurrencyCode"`
}

func getPlayerList(ctx *Context, ps getPlayerListParams, ret *getPlayerListResult) (err error) {
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	query := bson.M{}

	appidlist := []string{}
	var operator *comm.Operator
	if user.GroupId != 3 {
		ctx.CurrencyKey = ps.CurrencyCode
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		err = CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &operator)
		if err != nil {
			return err
		}
		if (ps.CurrencyCode != "" && operator.CurrencyKey == ps.CurrencyCode) || ps.CurrencyCode == "" {
			appidlist = append(appidlist, user.AppID)
		}
	}
	if ps.OperatorId > 0 {
		err = CollAdminOperator.FindId(ps.OperatorId, &operator)
		if err != nil {
			return err
		}
		temp := false
		for _, appid := range appidlist {
			if operator.AppID == appid {
				temp = true
			}
		}
		if temp == false {
			return errors.New(lang.GetLang(ctx.Lang, "未填写商户"))
		}
		appidlist = appidlist[:0]
		if (ps.CurrencyCode != "" && operator.CurrencyKey == ps.CurrencyCode) || ps.CurrencyCode == "" {
			appidlist = append(appidlist, operator.AppID)
		}
	}

	query["AppID"] = bson.M{"$in": appidlist}
	if ps.Uid != "" {
		query["Uid"] = ps.Uid
	}
	if ps.Pid > 0 {
		query["_id"] = ps.Pid
	}

	if ps.Status != nil {
		query["Status"] = ps.Status
	}

	//玩家登录状态
	t := time.Now().Add(-5 * time.Minute).UTC()
	if ps.OnLineStatus != 0 && ps.OnLineStatus != -1 { // 0 全部   1 在线   2 离线
		if ps.OnLineStatus == 1 {

			query["$or"] = []bson.M{
				{"LastSpinTime": bson.M{"$gte": t.Format("2006-01-02T15:04:05")}},
				{"LoginAt": bson.M{"$gte": t}},
			}
			//query["LastSpinTime"] = bson.M{"$gte": t}
			//query["LoginAt"] = bson.M{"$gte": time.Now().Add(-5 * time.Minute).Unix()}
		} else {

			//query["LastSpinTime"] = bson.M{"$lte": t}
			query["$and"] = []bson.M{
				{"LastSpinTime": bson.M{"$lte": t.Format("2006-01-02T15:04:05")}},
				{"LoginAt": bson.M{"$lte": t}},
			}
		}
	}
	sort := bson.M{"_id": -1}
	filter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Query:    query,
		Sort:     sort,
	}

	var players []*player
	count, err := NewOtherDB("game").Collection("Players").FindPage(filter, &players)
	if err != nil {
		return
	}

	//var currency []*curr
	//err = NewOtherDB("GameAdmin").Collection("CurrencyType").FindAll(bson.M{}, &currency)
	//if err != nil {
	//	return
	//}
	//币种
	currency := []*comm.CurrencyType{}
	err = NewOtherDB("GameAdmin").Collection("CurrencyType").FindAll(bson.M{}, &currency)
	if err != nil {
		return
	}
	mapCurrency := make(map[string]*comm.CurrencyType)
	for _, currencyType := range currency {
		mapCurrency[currencyType.CurrencyCode] = currencyType
	}

	appidlist = appidlist[:0]
	for _, re := range players {
		appidlist = append(appidlist, re.AppID)
	}
	operlist := []*comm.Operator{}
	err = CollAdminOperator.FindAll(bson.M{"AppID": bson.M{"$in": appidlist}}, &operlist)
	if err != nil {
		return
	}
	mapOper := make(map[string]*comm.Operator)
	for _, operator := range operlist {
		mapOper[operator.AppID] = operator
	}

	for k, re := range players {
		players[k].CurrencyKey = mapOper[re.AppID].CurrencyKey
		if players[k].CurrencyKey == "" {
			players[k].CurrencyName = ""
			continue
		}
		players[k].CurrencyName = mapCurrency[players[k].CurrencyKey].CurrencyName
		if ctx.Lang != "zh" {
			players[k].CurrencyName = mapCurrency[players[k].CurrencyKey].CurrencyCode
		}
	}
	agoTime := time.Now().Add(-5 * time.Minute).UTC().Unix()
	for i, play := range players {

		players[i].OnLineStatus = 0

		if play.LastSpinTime != "" {
			layout := "2006-01-02T15:04:05.000Z"
			playerSpinTime, err := time.Parse(layout, play.LastSpinTime)
			if err != nil {
				slog.Error("Time Format field: ", play.Id, err)
			}
			if playerSpinTime.Unix() > agoTime {
				players[i].OnLineStatus = 1
			}
		}
		loginTime := play.LoginAt.AsTime().UTC().Unix()
		if loginTime > agoTime {
			players[i].OnLineStatus = 1
		}

	}

	ret.AllCount = count
	ret.List = players
	for i := range ret.List {
		for _, info := range ret.List[i].TypeInfo {
			info.AvgBet = comm.DIV(info.Bet, info.SpinCnt)
		}
	}

	return nil
}

type UpdatePlayerStatusParams struct {
	UId    string
	Pid    int64
	Status int
	AppID  string
}

func UpdatePlayerStatus(ctx *Context, ps UpdatePlayerStatusParams, ret *mux.EmptyResult) (err error) {

	if ps.Pid == 0 {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	update := bson.M{}
	if ps.Status == comm.Player_Status_Normal {
		update["$set"] = bson.M{
			"Status":    comm.Player_Status_Normal,
			"CloseTime": time.Now().UTC(),
		}
	} else if ps.Status == comm.Player_Status_Stop {
		update["$set"] = bson.M{
			"Status":    comm.Player_Status_Stop,
			"CloseTime": time.Now().UTC(),
		}
	} else {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	coll := db.Collection2("game", "Players")
	old := &comm.Player{}
	err = coll.FindOne(context.TODO(), bson.M{"_id": ps.Pid}).Decode(&old)
	if err != nil {
		return
	}
	_, err = coll.UpdateByID(context.TODO(), ps.Pid, update, options.Update().SetUpsert(true))

	count := 0
	PublishPlayStatus("/player/status", PlayerStatus{
		Status: int64(ps.Status),
		UserID: ps.UId,
		AppID:  ps.AppID,
	}, count)

	coll = db.Collection2("GameAdmin", "SlotsPoolHistory")
	sf := ut.NewSnowflake()
	sfId := strconv.Itoa(int(sf.NextID()))
	_, err = coll.InsertOne(context.TODO(), bson.M{
		"_id":          sfId,
		"OpName":       ctx.Username,
		"OpPid":        ctx.PID,
		"AnimUserPid":  ps.Pid,
		"AnimUserName": old.Uid,
		"Type":         4, //4:玩家状态
		"OldGold":      old.Status,
		"NewGold":      ps.Status,
		"time":         mongodb.NewTimeStamp(time.Now()),
		"Currency":     "",
		"AppID":        old.AppID,
	})
	if err != nil {
		return
	}

	return nil
}
func getUpoperator(Lineoper string, AppID string, OperatorId ...string) {

	filter := bson.M{}

	if len(OperatorId) > 0 {
		filter = bson.M{"operatorId": OperatorId[0]}
	} else {
		filter = bson.M{"AppID": AppID}
	}
	oper := comm.Operator_V2{}
	err := CollAdminOperator.FindOne(filter, &oper)
	if err != nil {
		return
	}
	filter["AppID"] = oper.Name

	err = CollAdminOperator.FindOne(filter, &oper)
	if err != nil {
		return
	}

	if oper.OperatorType == 1 {
		Lineoper = oper.AppID
	}
}

func getPlayerInfo(ctx *Context, ps getPlayerListParams, ret *playAndMerchantInfo) (err error) {
	var playerInfo player
	filter := bson.M{
		"_id": ps.Pid,
	}

	var operator comm.Operator
	if ps.OperatorId != 0 {
		err = CollAdminOperator.FindId(ps.OperatorId, &operator)
		if err != nil {
			return
		}

		filter["AppId"] = operator.AppID
	}

	err = NewOtherDB("game").Collection("Players").FindOne(filter, &playerInfo)

	ret.player = playerInfo

	req, err := Getrrrk(ret.player)
	if err != nil {
		return err
	}

	ret.OnlineOperator = req.Line
	ret.Operator = req.oper

	loginTime := playerInfo.LoginAt.AsTime()
	now := time.Now()
	timeDifference := now.Sub(loginTime)
	ret.OnlineStatus = false
	if timeDifference < 5*time.Minute || len(req.betlog) > 0 {
		ret.OnlineStatus = true //todo：改成查player表中 LastSpinTime 获取状态
	}

	// 查询是否控制RTP和具体RTP配置的值是什么
	var updata_rTPControl []*comm.PlayerRTPControlModel

	err = CollPlayerRTPControl.FindAll(bson.M{"Pid": ps.Pid, "Status": 1, "FromPlant": "admin"}, &updata_rTPControl)

	if err != nil {
		return
	}

	if len(updata_rTPControl) > 0 {
		ret.RTPControlStatus = true

		gameNameMap := map[string]string{}
		_, games := GetGameListFormatName(ctx.Lang, nil)

		for _, game := range games {
			gameNameMap[game.ID] = game.Name
		}

		for i, play := range updata_rTPControl {
			_, ok := gameNameMap[play.GameID]
			if ok == false {
				gameNameMap[play.GameID] = "/"
			}
			updata_rTPControl[i].GameName = gameNameMap[play.GameID]
		}
	}

	ret.PlayerRTPControllist = updata_rTPControl
	return
}

type PlayerStatus struct {
	UserID string
	Status int64
	AppID  string
}

func PublishPlayStatus(url string, ps interface{}, count int) {

	if count == 3 {
		count = 0
		return
	}

	err := mq.JsonNC.Publish(url, ps)
	if err != nil {
		count++
		PublishPlayStatus(url, ps, count)
	}
	return
}
