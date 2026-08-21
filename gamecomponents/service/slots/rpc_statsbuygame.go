package main

import (
	"context"
	"fmt"
	"game/duck/logger"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/slots"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SlotsStatsBuyGameRpcServer struct {
}

func StatsBuyGame_GetID(day string, game string) string {
	return fmt.Sprintf("%v-%v", day, game)
}

func (SlotsStatsBuyGameRpcServer) Add(ctx context.Context, req *slots.SlotsStatsBuyGame_Doc) (*pb.Empty, error) {
	id := StatsBuyGame_GetID(req.Day, req.Game)

	var doc *slots.SlotsStatsBuyGame_Doc
	err := CollStatsBuyGame.FindId(id, &doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			CollStatsBuyGame.InsertOne(&slots.SlotsStatsBuyGame_Doc{ID: id, Day: req.Day, Game: req.Game})
		} else {
			logger.Err("CollStatsBuyGame.FindId err:", id, err)
			return &pb.Empty{}, err
		}
	}

	err = CollStatsBuyGame.UpdateId(id, bson.M{
		"$inc": bson.M{
			"Bet": req.Bet,
			"Win": req.Win,
		},
	})

	return &pb.Empty{}, err
}

func (SlotsStatsBuyGameRpcServer) Query(ctx context.Context, req *slots.SlotsStatsBuyGame_QueryReq) (*slots.SlotsStatsBuyGame_QueryResp, error) {
	var all = []*slots.SlotsStatsBuyGame_Doc{}

	err := CollStatsBuyGame.FindAll(bson.M{"Day": bson.M{"$in": req.Days}, "Game": req.Game}, &all)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Err("CollGameDayStat.FindAll err:", req.Days, err)
			return nil, err
		}
	}

	return &slots.SlotsStatsBuyGame_QueryResp{List: all}, nil
}
