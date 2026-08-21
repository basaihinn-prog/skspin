package main

import (
	"context"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/slots"
)

type AdminRpcServer struct {
}

func (AdminRpcServer) GetSelfSlotsPool(ctx context.Context, req *slots.GetSelfSlotsPoolReq) (*slots.GetSelfSlotsPoolResp, error) {
	return (&SlotsPoolRpcServer{}).GetSelfSlotsPool(ctx, req)
}

func (AdminRpcServer) SetSelfSlotsPool(ctx context.Context, req *slots.SetSelfSlotsPoolReq) (*pb.Empty, error) {
	return (&SlotsPoolRpcServer{}).SetSelfSlotsPool(ctx, req)
}
