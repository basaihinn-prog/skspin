package gameroom

type JoinGameRoom struct {
	actionType      string
	controllerId    string
	GameRoomUuid    string
	HostExtId       string
	SecWebSocketKey string
	gameId          string
	mathModuleId    string
	betList         string
	rateList        string
	rate            uint64
	SeatId          int
	NextScene       int
	Players         map[int]string
	ExtraData       []interface{}
}

func NewJoinGameRoom(actionType, controllerId, hostExtId, secWebSocketKey, gameId, mathModuleId, betList, rateList string,
	rate uint64,
	extraData []interface{},
) *JoinGameRoom {
	return &JoinGameRoom{
		actionType:      actionType,
		controllerId:    controllerId,
		HostExtId:       hostExtId,
		SecWebSocketKey: secWebSocketKey,
		gameId:          gameId,
		mathModuleId:    mathModuleId,
		betList:         betList,
		rateList:        rateList,
		rate:            rate,
		SeatId:          -1,
		NextScene:       -1,
		Players:         make(map[int]string),
		ExtraData:       extraData,
	}
}

func (j *JoinGameRoom) hashId() string {
	return Service.hashId(
		j.gameId,
		j.betList,
		j.rateList,
		j.mathModuleId,
		j.rate,
	)
}
