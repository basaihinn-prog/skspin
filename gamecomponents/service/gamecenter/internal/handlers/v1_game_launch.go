package handlers

import (
	"game/comm/mux"
	"game/service/gamecenter/internal/operator"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/game/launch",
		Handler:      v1_game_launch,
		Desc:         "启动游戏",
		Kind:         "api/v1",
		ParamsSample: v1GameLaunchPs{"operator_user_abcd", "XingYunXiang", "th", "mobile"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v1GameLaunchPs struct {
	UserID   string // 玩家id
	GameID   string // 游戏ID
	Language string // 游戏语言
	Platform string // mobile (default), desktop
}

type v1GameLaunchRet struct {
	Url string // 游戏启动链接
}

func v1_game_launch(app *operator.MemApp, ps v1GameLaunchPs, ret *v1GameLaunchRet) (err error) {
	// var loginResp LoginResp
	return login(app, ps, ret)
}
