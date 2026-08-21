package gameroom

type LeaveGameRoom struct {
	controllerId    string
	SecWebSocketKey string
	GameRoomUuid    string
	IsDisconnect    bool
	Players         map[int]string
}

func NewLeaveGameRoom(controllerId, secWebSocketKey, gameRoomUuid string, isDisconnect bool) *LeaveGameRoom {
	return &LeaveGameRoom{
		controllerId:    controllerId,
		SecWebSocketKey: secWebSocketKey,
		GameRoomUuid:    gameRoomUuid,
		IsDisconnect:    isDisconnect,
		Players:         make(map[int]string),
	}
}
