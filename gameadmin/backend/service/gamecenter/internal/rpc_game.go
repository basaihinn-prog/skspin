package internal

// import (
// 	"context"
// 	"errors"
// 	"game/duck/ut2/jwtutil"
// 	"game/pb/_gen/pb"
// 	"game/pb/_gen/pb/gamepb"
// 	"game/service/gamecenter/internal/operator"
// )

// type GameRpcServer struct{}

// func (GameRpcServer) ModifyGold(ctx context.Context, req *gamepb.ModifyGoldReq) (*gamepb.GetBalanceResp, error) {
// 	// am := operator.AppMgr
// 	// plr, err := am.GetPlr(req.Pid)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// app := am.GetApp(plr.AppID)
// 	// if app == nil {
// 	// 	return nil, errors.New("app not found.")
// 	// }

// 	// balance, err := app.ModifyGold(plr, req.Change, req.Comment)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// return &gamepb.GetBalanceResp{Balance: balance, Uid: plr.Uid}, nil
// 	return nil, errors.New("Method Unimplemented")
// }

// func (GameRpcServer) GetBalance(ctx context.Context, req *gamepb.GetBalanceReq) (*gamepb.GetBalanceResp, error) {
// 	// am := operator.AppMgr
// 	// plr, err := am.GetPlr(req.Pid)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// app := am.GetApp(plr.AppID)
// 	// if app == nil {
// 	// 	return nil, errors.New("app not found.")
// 	// }
// 	// balance, err := app.GetBalance(plr)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// return &gamepb.GetBalanceResp{
// 	// 	Balance: balance,
// 	// 	Uid:     plr.Uid,
// 	// }, nil

// 	return nil, errors.New("Method Unimplemented")
// }

// func (GameRpcServer) AddLog(ctx context.Context, req *gamepb.AddLogReq) (*pb.Empty, error) {
// 	return nil, errors.New("Method Unimplemented")
// }

// func (GameRpcServer) ValidateToken(ctx context.Context, req *gamepb.TokenReq) (*pb.Empty, error) {
// 	pid, _, err := jwtutil.ParseTokenData(req.Token)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pb.Empty{}, operator.AppMgr.UpdatePlrLoginTime(pid)
// }
