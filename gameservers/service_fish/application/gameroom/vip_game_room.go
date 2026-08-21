package gameroom

func newVipGameRoom(gameId, betList, rateList, mathModuleId string, rate uint64, secWebSocketKey string) *gameRoom {
	return newGameRoom(gameId, betList, rateList, mathModuleId, rate, 1, secWebSocketKey)
}
