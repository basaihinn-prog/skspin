package main

import (
	"context"
	"errors"
	"fmt"
	"game/duck/logger"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/slots"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocSlotsPool struct {
	ID   string `bson:"_id"`
	Pid  int
	Type slots.PoolType
	Gold int
}

type SlotsPoolRpcServer struct {
}

func SlotsPoolRpcServer_GetID(pid int, typ slots.PoolType) string {
	return fmt.Sprintf("%v-%v", pid, int(typ)) // 注意int(typ)，否则会是字符串。
}

func SlotsPoolRpcServer_EnsureExist(pid int, typ slots.PoolType) error {
	id := SlotsPoolRpcServer_GetID(pid, typ)

	var doc *DocSlotsPool
	err := CollPool.FindId(id, &doc)
	if err != nil && err != mongo.ErrNoDocuments {
		logger.Err("IncSelfSlotsPool err:", err, pid, typ)
		return err
	}

	if err == mongo.ErrNoDocuments {
		err = CollPool.InsertOne(&DocSlotsPool{
			ID:   id,
			Pid:  pid,
			Type: typ,
		})
		if err != nil {
			logger.Err("IncSelfSlotsPool err:", err, pid, typ)
			return err
		}
	}

	return nil
}

func (SlotsPoolRpcServer) IncSelfSlotsPool(ctx context.Context, req *slots.IncSelfSlotsPoolReq) (*pb.Empty, error) {
	if req.Type == slots.PoolType_Invalid {
		logger.Err("req.Type是无效值", req)
		return nil, errors.New("req.Type是无效值")
	}

	err := SlotsPoolRpcServer_EnsureExist(int(req.Pid), req.Type)
	if err != nil {
		return nil, err
	}

	id := SlotsPoolRpcServer_GetID(int(req.Pid), req.Type)

	err = CollPool.UpdateId(id, bson.M{"$inc": bson.M{"Gold": int(req.Gold)}})
	if err != nil {
		logger.Err("IncSelfSlotsPool err:", err, req)
		return nil, err
	}

	return &pb.Empty{}, err
}

func (SlotsPoolRpcServer) GetSelfSlotsPool(ctx context.Context, req *slots.GetSelfSlotsPoolReq) (*slots.GetSelfSlotsPoolResp, error) {
	id := SlotsPoolRpcServer_GetID(int(req.Pid), req.Type)

	var doc *DocSlotsPool
	err := CollPool.FindId(id, &doc)

	if err != nil && err != mongo.ErrNoDocuments {
		logger.Err("GetSelfSlotsPool err:", err, req)
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return &slots.GetSelfSlotsPoolResp{}, nil
	}

	return &slots.GetSelfSlotsPoolResp{Gold: int64(doc.Gold)}, nil
}

func (SlotsPoolRpcServer) SetSelfSlotsPool(ctx context.Context, req *slots.SetSelfSlotsPoolReq) (*pb.Empty, error) {
	err := SlotsPoolRpcServer_EnsureExist(int(req.Pid), req.Type)
	if err != nil {
		return nil, err
	}

	id := SlotsPoolRpcServer_GetID(int(req.Pid), req.Type)

	err = CollPool.UpdateId(id, bson.M{"$set": bson.M{"Gold": req.Gold}})
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
