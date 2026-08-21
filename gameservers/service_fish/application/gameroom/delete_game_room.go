package gameroom

type deleteGameRoom struct {
	gameRoomUuid string
	gameId       string
	betList      string
	rateList     string
	mathModuleId string
	rate         uint64
}

func newDeleteGameRoom(gameRoomUuid, gameId, betList, rateList, mathModuleId string, rate uint64) *deleteGameRoom {
	return &deleteGameRoom{
		gameRoomUuid: gameRoomUuid,
		gameId:       gameId,
		betList:      betList,
		rateList:     rateList,
		mathModuleId: mathModuleId,
		rate:         rate,
	}
}

func (d *deleteGameRoom) hashId() string {
	return Service.hashId(
		d.gameId,
		d.betList,
		d.rateList,
		d.mathModuleId,
		d.rate,
	)
}
