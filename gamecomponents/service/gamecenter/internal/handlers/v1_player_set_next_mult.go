package handlers

import (
	"game/comm"
	"game/comm/define"
	"game/comm/mux"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/setNextMult",
		Handler:      v1_player_set_next_mult,
		Desc:         "设置玩家下次游戏中奖倍数",
		Kind:         "api/v1",
		ParamsSample: v1PlayerSetNextMultPs{"abc", "pg_89", 90, 95},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v1PlayerSetNextMultPs struct {
	UserId  string  `json:"UserId"`
	GameId  string  `json:"GameId"`
	MinMult float64 `json:"MinMult"`
	MaxMult float64 `json:"MaxMult"`
}

func v1_player_set_next_mult(app *operator.MemApp, ps v1PlayerSetNextMultPs, ret *emptypb.Empty) (err error) {
	if ps.UserId == "" {
		err = define.NewErrCode("UserID is empty", 1008)
		return
	}
	if ps.GameId == "" {
		err = define.NewErrCode("GameId is empty", 1004)
		return
	}
	if ps.MinMult >= ps.MaxMult {
		err = define.NewErrCode("The minimum value of the multiple cannot be greater than the maximum value", 1026)
		return
	}
	memplr, err := operator.AppMgr.GetPlr2(app, ps.UserId)
	if err != nil {
		return
	}
	nextMultKey := "pnd_" + strconv.FormatInt(memplr.Pid, 10) + "_" + ps.GameId
	nextMultValue := &comm.NextMult{MinMult: ps.MinMult, MaxMult: ps.MaxMult}
	err = gcdb.StoreNextMult(nextMultKey, nextMultValue)
	if err != nil {
		return
	}
	return
}
