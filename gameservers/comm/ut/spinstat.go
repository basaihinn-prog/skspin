package ut

type SpinStatPs struct {
	Pid      int64
	BuyBonus bool
	Count    int
	Bet      float64
	// Extra1   bool
	// Extra2   bool
	Extra  int // 0 1 2
	BuyInt int
	BuyInd int
}

type SpinStatRet struct {
	Balance_before int64
	Balance_after  int64
	Bet            int64
	Win            int64
	Count          int
}
