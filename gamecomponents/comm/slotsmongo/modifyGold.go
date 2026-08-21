package slotsmongo

import (
	"game/comm/mq"
	"game/duck/lazy"
	"game/service/gamecenter/msgdef"

	"github.com/samber/lo"
)

// func ModifyGold(pid int64, change int64, comment string) (balance int64, err error) {

// 	var modifyGoldResp struct {
// 		Balance int64
// 	}

// 	err = mq.Invoke("/gamecenter/player/modifyGold", map[string]any{
// 		"Pid":     pid,
// 		"Change":  change,
// 		"Comment": lazy.ServiceName + ":" + comment,
// 	}, &modifyGoldResp)
// 	if err != nil {
// 		return
// 	}

// 	balance = modifyGoldResp.Balance
// 	return
// }

const (
	ReasonWin    = "win"
	ReasonBet    = "bet"
	ReasonRefund = "refund"
)

// type ModifyGoldPs struct {
// 	GameID  string
// 	Pid     int64
// 	Change  int64
// 	Comment string
// 	RoundID string
// 	Reason  string
// }

type ModifyGoldPs = msgdef.ModifyGoldPs

func ModifyGold(ps *ModifyGoldPs) (balance int64, err error) {
	if ps.GameID == "" {
		ps.GameID = lazy.ServiceName
	}

	if ps.Comment != "" {
		ps.Comment = lazy.ServiceName + ":" + ps.Comment
	}

	// lo.Must0(ps.SpinID != "")
	lo.Must0(ps.Pid != 0)
	// if ps.RoundID == "" {
	// 	ps.RoundID = ps.SpinID
	// }

	if ps.Reason == "" {
		if ps.Change < 0 {
			ps.Reason = ReasonBet
		} else {
			ps.Reason = ReasonWin
		}
	}

	var ret msgdef.ModifyGoldRet
	err = mq.Invoke("/gamecenter/player/modifyGold", ps, &ret)

	if err != nil {
		return
	}

	balance = ret.Balance
	return
}

func GetBalance(pid int64) (balance int64, err error) {

	var modifyGoldResp struct {
		Balance int64
	}

	err = mq.Invoke("/gamecenter/player/getBalance", map[string]any{
		"Pid": pid,
	}, &modifyGoldResp)
	if err != nil {
		return
	}

	balance = modifyGoldResp.Balance
	return
}
