package handlers

import (
	"game/comm/define"
	"game/comm/mux"
	"game/service/gamecenter/internal/operator"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/create",
		Handler:      v1_player_create,
		Desc:         "创建玩家帐号",
		Kind:         "api/v1",
		ParamsSample: v1PlayerCreatePs{"abc"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v1PlayerCreatePs struct {
	UserID string
	// Language string
}

type v1PlayerCreateRet struct {
	Pid int64
}

func v1_player_create(app *operator.MemApp, ps v1PlayerCreatePs, ret *v1PlayerCreateRet) (err error) {
	if ps.UserID == "" {
		return define.NewErrCode("UserID is empty.", 1008)
	}
	memplr, err := operator.AppMgr.EnsureUserExists(app, ps.UserID)
	if err != nil {
		return err
	}
	ret.Pid = memplr.Pid

	return
}
