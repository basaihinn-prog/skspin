package main

import (
	"context"
	"fmt"
	"game/duck/logger"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/slots"

	"go.mongodb.org/mongo-driver/bson"
)

type PlayerGameDayStatRpc struct {
}

func (sf *PlayerGameDayStatRpc) IncAll(ctx context.Context, req *slots.SlotsPlayerGameDayStatRecord) (*pb.Empty, error) {
	req.Flag = "All"
	err := sf.IncFlag(req)
	if err != nil {
		logger.Err(err)
	}
	return &pb.Empty{}, err
}

func (sf *PlayerGameDayStatRpc) IncBoth(ctx context.Context, req *slots.SlotsPlayerGameDayStatRecord) (*pb.Empty, error) {
	err := sf.IncFlag(req)
	if err != nil {
		return &pb.Empty{}, err
	}
	if req.Flag == "All" {
		return &pb.Empty{}, nil
	}

	req.Flag = "All"
	return sf.IncAll(ctx, req)
}

func (sf *PlayerGameDayStatRpc) IncFlag(req *slots.SlotsPlayerGameDayStatRecord) (err error) {
	if req == nil {
		return
	}
	if req.Game == "" {
		err = fmt.Errorf("game is empty")
		return
	}
	if req.Pid == 0 {
		err = fmt.Errorf("Pid is zero")
		return
	}
	if req.Date == "" {
		err = fmt.Errorf("date is empty")
		return
	}
	if req.Flag == "" {
		err = fmt.Errorf("flag is empty")
		return
	}

	CollPlayerGameDayStat.UpsertOne(
		bson.M{"Game": req.Game, "Pid": req.Pid, "Date": req.Date, "Flag": req.Flag},
		bson.M{"$inc": bson.M{
			"Flow":           req.Flow,
			"Win":            req.Win,
			"YeJi":           req.YeJi,
			"NetProfit":      req.NetProfit,
			"Fee":            req.Fee,
			"SelfPoolReward": req.SelfPoolReward,
		}},
	)

	return
}

func (sf *PlayerGameDayStatRpc) Docs(game string, Pid int, startdate, enddate string, flags ...string) ([]*slots.SlotsPlayerGameDayStatRecord, error) {
	var err error
	if game == "" {
		err = fmt.Errorf("game is empty")
		return nil, err
	}
	filter := bson.M{"Game": game}
	if Pid != 0 {
		filter["Pid"] = Pid
	}
	if startdate == "" {
		err = fmt.Errorf("startdate is empty")
		return nil, err
	}
	if enddate == "" {
		err = fmt.Errorf("enddate is empty")
		return nil, err
	}
	if startdate == enddate {
		filter["Date"] = startdate
	} else if startdate > enddate {
		return nil, fmt.Errorf("startdate greater enddate")
	} else if startdate < enddate {
		filter["Date"] = bson.M{
			"$lte": enddate,
			"$gte": startdate,
		}
	}
	flagsLength := len(flags)
	switch flagsLength {
	case 0:
		break
	case 1:
		filter["Flag"] = flags[0]
	default:
		filter["Flag"] = bson.M{"$in": flags}
	}
	var docs = make([]*slots.SlotsPlayerGameDayStatRecord, 0)

	err = CollPlayerGameDayStat.FindAll(filter, &docs)

	return docs, err
}

func (sf *PlayerGameDayStatRpc) GetGameDayWinCount(ctx context.Context, req *slots.GetGameDayWinCountReq) (*slots.GetGameDayWinCountResp, error) {
	var filer = bson.M{
		"Game": req.Game,
		"Date": req.Day,
		"NetProfit": bson.M{
			"$gt": 0,
		},
	}

	c, err := CollPlayerGameDayStat.CountDocuments(filer)

	return &slots.GetGameDayWinCountResp{Count: c}, err
}

func (sf *PlayerGameDayStatRpc) GetAllDayInfo(ctx context.Context, req *slots.SlotsDayBetReq) (*slots.SlotsGetAllDayInfoResp, error) {
	// filter := bson.D{{Key: "Game", Value: string(req.T)}, {Key: "Pid", Value: req.Pid}}
	filter := bson.M{"Game": req.Game, "Pid": req.Pid}

	var docs = make([]*slots.SlotsPlayerGameDayStatRecord, 0, 10)
	err := CollPlayerGameDayStat.FindAll(filter, &docs)

	return &slots.SlotsGetAllDayInfoResp{List: docs}, err
}

func (sf *PlayerGameDayStatRpc) GetPlrDayInfoByGame(ctx context.Context, req *slots.SlotsDayBetReq) (*slots.SlotsPlayerGameDayStatRecord, error) {
	ret := &slots.SlotsPlayerGameDayStatRecord{}
	filter := bson.M{"Game": req.Game, "Pid": req.Pid}
	docs := make([]*slots.SlotsPlayerGameDayStatRecord, 0, 10)
	err := CollPlayerGameDayStat.FindAll(filter, &docs)
	if err != nil {
		return ret, err
	}

	for _, doc := range docs {
		if !PlayerGameDayStatRecord_IsAll(doc) {
			continue
		}

		PlayerGameDayStatRecord_Add(ret, doc)
	}

	return ret, nil
}

func PlayerGameDayStatRecord_Add(this *slots.SlotsPlayerGameDayStatRecord, delta *slots.SlotsPlayerGameDayStatRecord) {
	if this == nil || delta == nil {
		return
	}
	this.Flow += delta.Flow
	this.Win += delta.Win
	this.YeJi += delta.YeJi
	this.NetProfit += delta.NetProfit
	this.Fee += delta.Fee
	this.SelfPoolReward += delta.SelfPoolReward
}

func PlayerGameDayStatRecord_IsAll(this *slots.SlotsPlayerGameDayStatRecord) bool {
	return this.Flag == "All"
}
