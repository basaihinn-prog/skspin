package main

import (
	"context"
	"fmt"
	"game/comm"
	"game/comm/db"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetAppGameInfo", "运营商游戏报表", "AdminInfo", getAppGameInfo, getAppGameInfoParams{
		Page:     1,
		PageSize: 20,
	})
}

type getAppGameInfoParams struct {
	Operator int64
	Date     string // "202402"
	Group    int64  // "202402"
	Page     int
	PageSize int
}

type getAppGameInfoResult struct {
	Count int
	List  []*getAppTypeInfoData
}

type getAppTypeInfoData struct {
	AppId           string
	Type            int
	Bet             float32
	Win             float32
	Date            string
	Profit          float32
	PlantRate       float64
	TurnoverPay     float64
	CooperationType int
	GGR             float64
	OperatorData    []*getAppGameInfoData `json:"children"`
}

type getAppGameInfoData struct {
	AppId           string  `bson:"appid"`
	OperatorType    int     `bson:"operator_type" json:"Type"`
	Bet             int64   `bson:"betamount"`
	PlantRate       float64 `bson:"plantrate"`
	TurnoverPay     float64 `bson:"-"`
	CooperationType int     `bson:"-"`
	Profit          float64 `bson:"-"`
	Win             int64   `bson:"winamount"`
	Date            string  `bson:"date"`
	GGR             float64 `bson:"GGR"`
	InCome          float64 `bson:"InCome"`
}

type LotteryMonthData struct {
	AppID     string `bson:"appID"`
	Month     string `bson:"month"`
	Bet       int64  `bson:"bet"`
	PlayerWin int64  `bson:"playerWin"`
}

type statisticsWinBetResult struct {
	Win float32
	Bet float32
}

func getAppGameInfo(ctx *Context, ps getAppGameInfoParams, ret *getAppGameInfoResult) (err error) {
	now := time.Now().Format("200601")
	if ps.Date == "" {
		ps.Date = now
	}
	//startTime := fmt.Sprintf("%s00", ps.Date)
	//endTime := fmt.Sprintf("%s31", ps.Date)
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	ids, _ := GetOperatopAppID(ctx)
	if ps.Operator != 0 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.Operator, &operator)
		if err != nil {
			return err
		}
		ids = ids[:0]
		ids = append(ids, operator.AppID)
	}

	var gameInfoData []*getAppGameInfoData

	psdate := ps.Date
	if len(ps.Date) > 6 {
		psdate = ps.Date[:6]
	}

	if psdate != now {
		gameInfoData, err = getHistoryMonthlyReport(psdate, ids)
		if err != nil {
			return
		}
	} else {
		gameInfoData, err = getNowMonthlyReport(ids)
		if err != nil {
			return
		}
	}

	// 查询所有商户
	var Operators []*comm.Operator_V2
	CollAdminOperator.FindAll(bson.M{}, &Operators)

	OperatorMap := map[string]string{}
	OperatorPlantRateMap := map[string]float64{}
	OperatorInfoMap := map[string]*comm.Operator_V2{}
	for _, item := range Operators {
		OperatorMap[item.AppID] = item.Name
		OperatorPlantRateMap[item.AppID] = item.PlatformPay
		OperatorInfoMap[item.AppID] = item

	}

	selfMapData, onlineMapBetData, onlineMapWinData := generatorOnlineStatistics(gameInfoData, OperatorMap, OperatorInfoMap)

	var resultData []*getAppTypeInfoData

	for key, item := range onlineMapBetData {

		resultItem := &getAppTypeInfoData{
			AppId:           key,
			Type:            1,
			Date:            item[0].Date,
			Win:             onlineMapWinData[key].Win,
			Bet:             onlineMapWinData[key].Bet,
			Profit:          onlineMapWinData[key].Bet - onlineMapWinData[key].Win,
			PlantRate:       OperatorPlantRateMap[key],
			TurnoverPay:     OperatorInfoMap[key].TurnoverPay,
			CooperationType: OperatorInfoMap[key].CooperationType,
			OperatorData:    item,
		}
		resultData = append(resultData, resultItem)
	}

	resultData = append(resultData, selfMapData...)

	ret.Count = len(resultData)
	dataStart := (ps.Page - 1) * ps.PageSize
	dataEndNum := dataStart + ps.PageSize
	if ps.Page == 1 {
		dataStart = ps.Page - 1
		dataEndNum = ps.PageSize
	} else {
		dataStart = (ps.Page - 1) * ps.PageSize
		dataEndNum = dataStart + ps.PageSize
	}

	if ret.Count < dataEndNum {
		dataEndNum = ret.Count
	}

	ret.List = resultData[dataStart:dataEndNum]

	return
}

func generatorOnlineStatistics(gameInfoData []*getAppGameInfoData, OperatorMap map[string]string, operatorInfoMap map[string]*comm.Operator_V2) ([]*getAppTypeInfoData, map[string][]*getAppGameInfoData, map[string]statisticsWinBetResult) {
	// 没有线路商
	var SelfMapBetData []*getAppTypeInfoData
	// 有线路商
	onlineMapBetData := map[string][]*getAppGameInfoData{}
	// 每个线路商所对应的投注总数与盈利总数
	onlineMapWinData := map[string]statisticsWinBetResult{}

	for _, item := range gameInfoData {

		onlineName := OperatorMap[item.AppId]

		if onlineName == item.AppId || onlineName == "admin" || onlineName == "" {

			SelfMapBetData = append(SelfMapBetData, &getAppTypeInfoData{
				AppId:           item.AppId,
				Type:            2,
				Date:            item.Date,
				Win:             float32(item.Win),
				Bet:             float32(item.Bet),
				Profit:          float32(item.Bet) - float32(item.Win),
				GGR:             item.GGR,
				PlantRate:       item.PlantRate,
				TurnoverPay:     operatorInfoMap[item.AppId].TurnoverPay,
				CooperationType: operatorInfoMap[item.AppId].CooperationType,
			})

		} else {

			result := getAppGameInfoData{
				AppId:           item.AppId,
				OperatorType:    2,
				Date:            item.Date,
				Win:             item.Win,
				Profit:          float64(item.Bet - item.Win),
				Bet:             item.Bet,
				GGR:             item.GGR,
				PlantRate:       item.PlantRate,
				TurnoverPay:     operatorInfoMap[item.AppId].TurnoverPay,
				CooperationType: operatorInfoMap[item.AppId].CooperationType,
			}

			onlineMapBetData[onlineName] = append(onlineMapBetData[onlineName], &result)
		}

		onlineMapWinData[onlineName] = statisticsWinBetResult{
			Win: onlineMapWinData[onlineName].Win + float32(item.Win),
			Bet: onlineMapWinData[onlineName].Bet + float32(item.Bet),
		}
	}

	return SelfMapBetData, onlineMapBetData, onlineMapWinData

}

// 查其他月份
func getHistoryMonthlyReport(date string, ids []string) ([]*getAppGameInfoData, error) {

	var gameInfoData []*getAppGameInfoData
	filter := bson.M{
		"date": date,
		"appid": bson.M{
			"$in": ids,
		}}
	err := NewOtherDB("reports").Collection("OperatorMonthlyGGR").FindAll(filter, &gameInfoData)
	if err != nil {
		return nil, err
	}
	return gameInfoData, err
}

// 查询本月的
func getNowMonthlyReport(ids []string) ([]*getAppGameInfoData, error) {

	var gameInfoData []*getAppGameInfoData
	date := time.Now().UTC().Format("200601")
	filter := bson.M{
		"date": date,
		"appid": bson.M{
			"$in": ids,
		}}
	err := NewOtherDB("reports").Collection("OperatorMonthlyGGR").FindAll(filter, &gameInfoData)
	if err != nil {
		return nil, err
	}

	mapNow, err := getNowReport(ids)
	if err != nil {
		return nil, err
	}
	for i, datum := range gameInfoData {
		if _, ok := mapNow[datum.AppId]; !ok {
			continue
		}
		gameInfoData[i].Win += mapNow[datum.AppId].Win
		gameInfoData[i].Bet += mapNow[datum.AppId].Bet
		gameInfoData[i].GGR += mapNow[datum.AppId].GGR
		gameInfoData[i].InCome += mapNow[datum.AppId].InCome
		gameInfoData[i].PlantRate = mapNow[datum.AppId].PlantRate
		delete(mapNow, datum.AppId)
	}

	for _, data := range mapNow {
		gameInfoData = append(gameInfoData, data)
	}

	return gameInfoData, err
}

// 查询今天的 费率 和 费用
func getNowReport(ids []string) (map[string]*getAppGameInfoData, error) {
	match := bson.M{
		"$match": bson.M{
			"date":  time.Now().Format("20060102"),
			"appid": bson.M{"$in": ids},
		}}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"appid": "$appid",
			},
			"winamount": bson.M{"$sum": "$winamount"},
			"betamount": bson.M{"$sum": "$betamount"},
			"appid":     bson.M{"$first": "$appid"},
			"date":      bson.M{"$first": "$date"},
		}}

	sort := bson.M{
		"$sort": bson.M{
			"date": -1,
		}}

	addfields := bson.M{
		"$addFields": bson.M{
			"operator_type": 3,
		}}

	aggregate := []bson.M{}

	aggregate = append(aggregate, match, group, sort, addfields)

	coll := db.Collection2("reports", "BetDailyReport")

	cursor, err := coll.Aggregate(context.TODO(), aggregate)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var gameInfoData []*getAppGameInfoData
	err = cursor.All(context.TODO(), &gameInfoData)

	if err != nil {
		return nil, err
	}

	var Operators []*comm.Operator_V2
	CollAdminOperator.FindAll(bson.M{}, &Operators)

	PlantRateMap := map[string]float64{}
	mapParent := make(map[string]string)
	mapCooperstion := make(map[string]int)
	mapTurnoverRale := make(map[string]float64)
	for _, item := range Operators {
		PlantRateMap[item.AppID] = item.PlatformPay
		mapParent[item.AppID] = item.Name
		mapCooperstion[item.AppID] = item.CooperationType
		mapTurnoverRale[item.AppID] = item.TurnoverPay
	}

	result := make(map[string]*getAppGameInfoData)
	for i, datum := range gameInfoData {
		Rate := PlantRateMap[datum.AppId]
		if mapCooperstion[datum.AppId] == 2 {
			Rate = mapTurnoverRale[datum.AppId]
		}
		bet := gameInfoData[i].Bet
		win := gameInfoData[i].Win
		InCome := float64(0.0)
		GGR := float64(0.0)
		if mapParent[datum.AppId] != "admin" {
			if mapCooperstion[datum.AppId] == 1 {
				parentRate := PlantRateMap[mapParent[datum.AppId]]
				InCome = float64(bet-win) * ((Rate - parentRate) * 0.01)
				GGR = float64(bet-win) * (Rate * 0.01)
			} else {
				parentRate := mapTurnoverRale[mapParent[datum.AppId]]
				InCome = float64(bet) * ((Rate - parentRate) * 0.01)
				GGR = float64(bet) * (Rate * 0.01)
			}
		}

		gameInfoData[i].PlantRate = Rate
		gameInfoData[i].GGR = GGR
		gameInfoData[i].InCome = InCome

		result[datum.AppId] = datum
	}

	return result, err
}
