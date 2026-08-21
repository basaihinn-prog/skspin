package receipt

type receipt struct {
	secWebSocketKey string
	hostExtId       string
	fish            *fish
}

type fish struct {
	Sn           string
	AccountingSn uint64
	Bet          uint64
	Win          uint64
	GameData     string
	GameResult   string
}

func newReceipt(bulletUuid, secWebSocketKey, hostExtId, gameResult, gameData string,
	accountingSn, bet, pay uint64) *receipt {
	return &receipt{
		secWebSocketKey: secWebSocketKey,
		hostExtId:       hostExtId,
		fish: &fish{
			Sn:           bulletUuid,
			GameData:     gameData,
			GameResult:   gameResult,
			AccountingSn: accountingSn,
			Bet:          bet,
			Win:          pay,
		},
	}
}
