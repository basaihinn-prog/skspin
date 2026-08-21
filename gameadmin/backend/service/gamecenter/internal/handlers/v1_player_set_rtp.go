package handlers

import (
	"errors"
	"game/comm/mux"
	"game/service/gamecenter/internal/operator"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/setRtp",
		Handler:      v1_player_set_rtp,
		Desc:         "设置玩家RTP",
		Kind:         "api/v1",
		ParamsSample: v1PlayerSetRtpPs{"abc", 80},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v1PlayerSetRtpPs struct {
	UserID string
	Rtp    int
}

type v1PlayerSetRtpRet struct {
	Pid int64
	Rtp int
}

func v1_player_set_rtp(app *operator.MemApp, ps v1PlayerSetRtpPs, ret *v1PlayerSetRtpRet) (err error) {
	memplr, err := operator.AppMgr.GetPlr2(app, ps.UserID)
	if err != nil {
		return
	}

	if ps.Rtp < -1 || ps.Rtp > 200 {
		err = errors.New("incorrect rtp value")
		return
	}

	old_rtp, err := memplr.SetRtp(ps.Rtp)
	if err != nil {
		return
	}
	ret.Pid = memplr.Pid
	ret.Rtp = old_rtp
	return
}
