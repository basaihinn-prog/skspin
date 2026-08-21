package operator

import (
	"context"
	"game/comm"
	"game/pb/_gen/pb/gamepb"
	"net/url"
)

type HttpApi struct {
	ApiUrl string
}

func NewHttpApi(apiurl string) HttpApi {
	return HttpApi{ApiUrl: apiurl}
}

func (h HttpApi) ModifyGold( /*req *gamepb.ModifyGoldUidReq*/ memplr *MemPlr, change int64, comment string) (int64, error) {
	apiUrl, err := url.JoinPath(h.ApiUrl, "ModifyGold")
	if err != nil {
		return 0, err
	}

	req := &gamepb.ModifyGoldUidReq{
		UserID:  memplr.Uid,
		Change:  change,
		Comment: comment,
	}

	resp := gamepb.GetBalanceUidResp{}
	err = comm.PostJsonCode(context.TODO(), apiUrl, req, &resp, nil)
	if err != nil {
		return 0, err
	}

	return resp.Balance, nil
}

func (h HttpApi) GetBalance( /*req *gamepb.GetBalanceUidReq*/ memplr *MemPlr) (int64, error) {
	apiUrl, err := url.JoinPath(h.ApiUrl, "GetBalance")
	if err != nil {
		return 0, err
	}
	resp := gamepb.GetBalanceUidResp{}
	req := &gamepb.GetBalanceUidReq{UserID: memplr.Uid}
	err = comm.PostJsonCode(context.TODO(), apiUrl, req, &resp, nil)
	if err != nil {
		return 0, err
	}

	return resp.Balance, nil
}
