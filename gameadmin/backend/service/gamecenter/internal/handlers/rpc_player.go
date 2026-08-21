package handlers

import (
	"errors"
	"fmt"
	"game/comm/mux"
	"game/service/gamecenter/internal/operator"
	"game/service/gamecenter/msgdef"
)

func init() {
	mux.RegRpc("/gamecenter/player/getBalance", "获取玩家余额", "player", getBalance, getBalancePs{100001})
	mux.RegRpc("/gamecenter/player/modifyGold", "修改玩家余额", "player", modifyGold, modifyGoldReq{
		Pid:    100001,
		Change: 10000,
		// Comment: "http 测试修改",
	})
}

// type GameRpcPlayer struct{}

type getBalancePs struct {
	Pid int64
}

type getBalanceRet struct {
	Balance int64
	Uid     string
}

func getBalance(req getBalancePs, ret *getBalanceRet) (err error) {
	am := operator.AppMgr
	plr, err := am.GetPlr(req.Pid)
	if err != nil {
		return
	}
	app := am.GetApp(plr.AppID)
	if app == nil {
		err = errors.New("app not found.")
		return
	}

	// if lazy.CommCfg().IsDev {
	// 	time.Sleep(50 * time.Millisecond)
	// }

	balance, err := app.GetBalance(plr)
	if err != nil {
		return
	}

	ret.Balance = balance
	ret.Uid = plr.Uid

	return
}

type modifyGoldReq = msgdef.ModifyGoldPs

// type modifyGoldReq struct {
// 	Pid     int64 // 玩家ID
// 	Change  int64 // 金币变化：>0 加钱， <0 扣钱
// 	SpinID  string
// 	RoundID string
// 	Reason  string
// 	Comment string // 注释
// }

func modifyGold(ps *msgdef.ModifyGoldPs, ret *msgdef.ModifyGoldRet) (err error) {
	if !ps.Valid() {
		err = fmt.Errorf("error ps[%+v]", ps)
		return
	}

	am := operator.AppMgr

	plr, err := am.GetPlr(ps.Pid)
	if err != nil {
		return
	}
	app := am.GetApp(plr.AppID)
	if app == nil {
		err = errors.New("app not found.")
		return
	}

	// balance, err := app.ModifyGold(plr, ps.Change, ps.Comment)
	balance, err := app.ModifyGold(plr, ps)
	if err != nil {
		return
	}
	ret.Balance = balance
	ret.Uid = plr.Uid

	return
}
