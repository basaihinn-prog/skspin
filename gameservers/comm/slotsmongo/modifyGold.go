package slotsmongo

import (
	"log/slog"

	"serve/comm/lazy"
	"serve/comm/mq"

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

type ModifyGoldPs struct {
	GameID  string
	Pid     int64
	Change  int64
	Comment string

	// 关联 BetLog.RoundID
	RoundID string
	Reason  string
	IsEnd   bool // 结束标识
}

type ModifyGoldRet struct {
	Balance int64
	Uid     string
}

const (
	ReasonWin    = "win"
	ReasonBet    = "bet"
	ReasonFRBBet = "frbbet"
	ReasonRefund = "refund"
)

func ModifyGold(ps *ModifyGoldPs) (balance int64, err error) {
	ps.GameID = lazy.ServiceName

	if ps.Comment != "" {
		ps.Comment = lazy.ServiceName + ":" + ps.Comment
	}

	if ps.RoundID == "" {
		slog.Error("RoundID is empty", "ps", ps)
	}

	lo.Must0(ps.Pid != 0)

	if ps.Reason == "" {
		if ps.Change < 0 {
			ps.Reason = ReasonBet
		} else {
			ps.Reason = ReasonWin
		}
	}

	var ret ModifyGoldRet
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

func GetGain(pid int64) (gain int64, err error) {

	var gainResp struct {
		Gain int64
	}

	err = mq.Invoke("/gamecenter/player/getGain", map[string]any{
		"Pid": pid,
	}, &gainResp)
	if err != nil {
		return
	}

	gain = gainResp.Gain
	return
}
