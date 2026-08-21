package main

import (
	"context"
	"encoding/json"
	"game/comm/slotscall"
	"game/comm/ut"
	"game/duck/logger"
	"game/pb/_gen/pb/adminslots"
	"game/pb/_gen/pb/slots"
	"game/pb/_gen/pb/stats"
	"game/service/stats/common"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type StatsRpcServer struct {
}

func (StatsRpcServer) GetLast10DayRunTimeData(ctx context.Context, req *stats.GetLast10DayRunTimeDataReq) (*stats.JsonData, error) {
	now := time.Now()
	endDate := ut.Time2YMD(now)
	lengh := 10
	startDate := ut.Time2YMD(now.AddDate(0, 0, -lengh))
	result := &adminslots.GetLast10DayRunTimeDataResp{
		General: &slots.ThreeMonkeyRunTimeResp{
			Arr: make([]*slots.ThreeMonkeyRunTimeData, lengh),
		},
	}

	subgamearr, err := slotscall.GameDayStatGetDocs(context.TODO(), &slots.SlotsGameDayStatDocsReq{
		StartDate: startDate,
		EndDate:   endDate,
		Game:      req.Game,
		Flags:     []string{"All"},
	})
	if err != nil {
		return nil, err
	}
	docsM := map[string]*slots.SlotsGameDayStatRecord{}
	for _, v := range subgamearr {
		docsM[v.Date] = v
	}

	var arr []*slots.ThreeMonkeyRunTimeData
	days := []string{}

	for i := 0; i < 10; i++ {
		var day time.Time
		if i == 0 {
			day = now
		} else {
			day = now.AddDate(0, 0, -i)
		}
		dateStr := ut.Time2YMD(day)
		days = append(days, dateStr)

		d, _ := slotscall.FixRunTimeDataByGame(ctx, dateStr, req.Game) // slotsBase.FixRunTimeData(client, channel, dateStr)
		if v, ok := docsM[dateStr]; ok {
			d.EnterCount = v.EnterCount
		}
		d.WinPlrCount, _ = slotscall.GetGameDayWinCount(ctx, req.Game, dateStr)
		arr = append(arr, d)
	}
	r := &slots.ThreeMonkeyRunTimeResp{}
	r.Arr = arr
	(r).TotalIncome, (r).TotalExpense, _ = slotscall.RuntimeGetTotalByGame(ctx, req.Game)
	slotscall.ThreeMonkeyRunTimeResp_Add(result.General, r)
	//TODO
	result.BuyGame, _ = slotscall.StatsBuyGameQueryByGame(ctx, days, req.Game)

	bytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return &stats.JsonData{Data: string(bytes)}, nil
}

func (StatsRpcServer) GetTotalExpense(ctx context.Context, req *stats.GetTotalExpenseReq) (*stats.GetTotalExpenseResp, error) {
	resp := &stats.GetTotalExpenseResp{}

	startDay := ut.Time2YMD(time.Now().AddDate(0, 0, -int(req.Days)))

	ret, err := common.Aggregate[stats.GetTotalExpenseResp](CollRuntimeDay.Coll(),
		bson.D{{Key: "$match", Value: bson.M{"Game": req.Game, "Day": bson.M{"$gt": startDay}}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$Game", "TotalExpense": bson.M{"$sum": "$Expense"}, "TotalIncome": bson.M{"$sum": "$Income"}}}},
	)

	logger.Info("GetTotalExpense", req, ret)

	if err != nil {
		return nil, err
	}

	if len(ret) > 0 {
		resp = ret[0]
	}

	return resp, nil
}
