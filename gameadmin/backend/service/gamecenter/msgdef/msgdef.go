package msgdef

type ModifyGoldPs struct {
	GameID  string
	Pid     int64
	Change  int64
	Comment string
	RoundID string
	Reason  string
}

const (
	ReasonWin    = "win"
	ReasonBet    = "bet"
	ReasonRefund = "refund"
)

const (
	min = -1 << 50
	max = 1 << 50
)

func (mp *ModifyGoldPs) Valid() bool {
	if mp.Change < min || max < mp.Change {
		return false
	}
	switch mp.Reason {
	case ReasonBet:
		if 0 < mp.Change {
			return false
		}
	case ReasonWin, ReasonRefund:
		if mp.Change < 0 {
			return false
		}
	}

	return true
}

type ModifyGoldRet struct {
	Balance int64
	Uid     string
}
