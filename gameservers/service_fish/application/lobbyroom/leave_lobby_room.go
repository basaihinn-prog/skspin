package lobbyroom

type LeaveLobbyRoom struct {
	controllerId    string
	SecWebSocketKey string
	LobbyRoomUuid   string
	IsDisconnect    bool
}

func NewLeaveLobbyRoom(controllerId, secWebSocketKey, lobbyRoomUuid string, isDisconnect bool) *LeaveLobbyRoom {
	return &LeaveLobbyRoom{
		controllerId:    controllerId,
		SecWebSocketKey: secWebSocketKey,
		LobbyRoomUuid:   lobbyRoomUuid,
		IsDisconnect:    isDisconnect,
	}
}
