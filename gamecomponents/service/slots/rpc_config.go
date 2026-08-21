package main

import (
	"context"
	"game/pb/_gen/pb/slots"
)

type ConfigRpcServer struct {
}

func (ConfigRpcServer) GetYeJiRatio(ctx context.Context, req *slots.GetYeJiRatioReq) (*slots.GetYeJiRatioResp, error) {
	resp := &slots.GetYeJiRatioResp{}

	switch req.Game {
	case slots.GetYeJiRatioReq_Slots:
		resp.Ratio = FlowTransferYejiConfig.Slots
	}

	return resp, nil
}
