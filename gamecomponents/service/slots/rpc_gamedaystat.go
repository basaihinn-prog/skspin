package main

import (
	"context"
	"fmt"
	"game/duck/logger"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/slots"

	"go.mongodb.org/mongo-driver/bson"
)

type GameDayStatRpc struct {
}

func (sf *GameDayStatRpc) IncAll(ctx context.Context, req *slots.SlotsGameDayStatRecord) (*pb.Empty, error) {
	req.Flag = "All"
	err := sf.IncFlag(req)
	if err != nil {
		logger.Err(err)
	}
	return &pb.Empty{}, err
}

func (sf *GameDayStatRpc) IncBoth(ctx context.Context, req *slots.SlotsGameDayStatRecord) (*pb.Empty, error) {
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

func (sf *GameDayStatRpc) IncFlag(req *slots.SlotsGameDayStatRecord) (err error) {
	if req == nil {
		return
	}
	if req.Game == "" {
		err = fmt.Errorf("game is empty")
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

	err = CollGameDayStat.UpsertOne(
		bson.M{"Game": req.Game, "Date": req.Date, "Flag": req.Flag},
		bson.M{"$inc": bson.M{
			"EnterCount":     req.EnterCount,
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

func (sf *GameDayStatRpc) DocGames(channel string, date string) ([]*slots.SlotsGameDayStatRecord, error) {

	var err error
	if channel == "" {
		err = fmt.Errorf("channel is empty")
		return nil, err
	}
	if date == "" {
		err = fmt.Errorf("date is empty")
		return nil, err
	}
	filter := bson.M{
		"Channel": channel,
		"Date":    date,
		"Flag":    "All",
	}
	var docs = make([]*slots.SlotsGameDayStatRecord, 0, 50)
	err = CollGameDayStat.FindAll(filter, &docs)
	return docs, err
}

func (sf *GameDayStatRpc) Doc(date string, game string, flags ...string) ([]*slots.SlotsGameDayStatRecord, error) {
	var err error
	if date == "" {
		err = fmt.Errorf("date is empty")
		return nil, err
	}
	if game == "" {
		err = fmt.Errorf("game is empty")
		return nil, err
	}

	filter := bson.M{"Date": date, "Game": game}

	flagsLength := len(flags)
	switch flagsLength {
	case 0:
		break
	case 1:
		filter["Flag"] = flags[0]
	default:
		filter["Flag"] = bson.M{"$in": flags}
	}
	var docs = make([]*slots.SlotsGameDayStatRecord, 0)
	err = CollGameDayStat.FindAll(filter, &docs)
	return docs, err
}

func (sf *GameDayStatRpc) Docs(ctx context.Context, req *slots.SlotsGameDayStatDocsReq) (*slots.SlotsGameDayStatDocsResp, error) {
	var startdate = req.StartDate
	var enddate = req.EndDate
	var game = req.Game
	var flags = req.Flags

	var err error

	if startdate == "" {
		err = fmt.Errorf("startdate is empty")
		return nil, err
	}
	if enddate == "" {
		err = fmt.Errorf("enddate is empty")
		return nil, err
	}
	if game == "" {
		err = fmt.Errorf("game is empty")
		return nil, err
	}
	filter := bson.M{}
	if startdate == enddate {
		filter["Date"] = startdate
	} else if startdate > enddate {
		return nil, fmt.Errorf("startdate greater enddate")
	} else if startdate < enddate {
		filter["Date"] = bson.M{
			"$lte": enddate,
			"$gt":  startdate,
		}
	}
	filter["Game"] = game
	flagsLength := len(flags)
	switch flagsLength {
	case 0:
		break
	case 1:
		filter["Flag"] = flags[0]
	default:
		filter["Flag"] = bson.M{"$in": flags}
	}
	var docs = make([]*slots.SlotsGameDayStatRecord, 0, 10)
	err = CollGameDayStat.FindAll(filter, &docs)
	return &slots.SlotsGameDayStatDocsResp{List: docs}, err
}
