package botuser

type createBotUser struct {
	secWebSocketKey string
	gameRoomUuid    string
	rate            uint64
	maxBet          uint64
	rateIndex       uint32
}

func NewCreateBotUser(secWebSocketKey, gameRoomUuid string, rate, maxBet uint64, rateIndex uint32) *createBotUser {
	return &createBotUser{
		secWebSocketKey: secWebSocketKey,
		gameRoomUuid:    gameRoomUuid,
		rate:            rate,
		maxBet:          maxBet,
		rateIndex:       rateIndex,
	}
}
