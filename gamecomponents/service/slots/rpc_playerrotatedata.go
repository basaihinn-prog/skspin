package main

import (
	"context"
	"fmt"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/slots"

	"go.mongodb.org/mongo-driver/bson"
)

type PlayerRotateRpcServer struct {
}

func (PlayerRotateRpcServer) getId(pid int64, game string) string {
	return fmt.Sprintf("%v-%v", game, pid)
}

func (sf PlayerRotateRpcServer) IncRotate(ctx context.Context, req *slots.SlotsPlayerRotateData_IncRotateReq) (*pb.Empty, error) {
	id := sf.getId(req.Pid, req.Game)
	err := CollPlayerRotateData.UpsertId(id, bson.M{"$inc": bson.M{"RotateCount": 1}})
	return &pb.Empty{}, err
}

func (sf PlayerRotateRpcServer) GetRotate(ctx context.Context, req *slots.SlotsPlayerRotateData_GetRotateReq) (*slots.SlotsPlayerRotateData_GetRotateResp, error) {
	ret := &slots.SlotsPlayerRotateData_GetRotateResp{}

	var doc *slots.SlotsPlayerRotateData_Doc
	id := sf.getId(req.Pid, req.Game)
	err := CollPlayerRotateData.FindId(id, &doc)
	if err != nil {
		return nil, err
	}

	if doc != nil {
		ret.Count = doc.RotateCount
	}

	return ret, nil
}
