package lobbyroom

import (
	"fmt"
	"serve/fish_comm/broadcaster"
	"serve/fish_comm/flux"
	"serve/fish_comm/jackpot-client/jackpot"
	"serve/fish_comm/logger"

	"github.com/google/uuid"
)

const (
	actionAddUserCall      = "actionAddUserCall"
	actionAddUserRecall    = "actionAddUserRecall"
	actionDeleteUserCall   = "actionDeleteUserCall"
	actionDeleteUserRecall = "actionDeleteUserRecall"
)

type lobbyRoom struct {
	uuid string
	in   chan *flux.Action
}

func newLobbyRoom() *lobbyRoom {
	l := &lobbyRoom{
		uuid: uuid.New().String(),
		in:   make(chan *flux.Action, 1024),
	}

	go l.run()

	logger.Service.Zap.Infow("LobbyRoom Created.",
		"LobbyRoomUuid", l.uuid,
		"Chan", fmt.Sprintf("%p", l.in),
	)
	return l
}

func (l *lobbyRoom) run() {
	flux.Register(l.uuid, l.in)
	lobbyUsers := make(map[string]interface{})

	isRegistered := false

	for {
		select {
		case action, ok := <-l.in:
			if !ok {
				return
			}

			switch action.Key().Name() {
			case actionAddUserCall:
				gameId := action.Payload()[2].(string)

				if !isRegistered {
					isRegistered = true
					flux.Send(
						broadcaster.ActionBroadcasterRegister,
						l.uuid,
						broadcaster.Service.Id,
						broadcaster.NewRegisterBroadcaster(l.uuid, gameId, l.uuid, l.uuid, l.uuid, 0),
					)
				}

				secWebSocketKey := action.Payload()[0].(string)
				controllerId := action.Payload()[1].(string)
				lobbyUsers[secWebSocketKey] = secWebSocketKey
				flux.Send(actionAddUserRecall, l.uuid, action.Key().From(), secWebSocketKey, controllerId, gameId)

			case actionDeleteUserCall:
				data := action.Payload()[0].(*LeaveLobbyRoom)

				if data.LobbyRoomUuid == l.uuid {
					delete(lobbyUsers, data.SecWebSocketKey)
					flux.Send(actionDeleteUserRecall, l.uuid, action.Key().From(), data)
				}

			case broadcaster.EMSGID_eBroadcaster:
				if action.Key().To() == l.uuid {
					l.broadcastAction(broadcaster.EMSGID_eBroadcaster, action, lobbyUsers)
				}

			case jackpot.EMSGID_eJackpotInfo:
				if action.Key().To() == l.uuid {
					l.broadcastAction(jackpot.EMSGID_eJackpotInfo, action, lobbyUsers)
				}

			case jackpot.EMSGID_eJackpotNotify:
				if action.Key().To() == l.uuid {
					l.broadcastAction(jackpot.EMSGID_eJackpotNotify, action, lobbyUsers)
				}

			}
		}
	}
}

func (l *lobbyRoom) broadcastAction(actionName string, action *flux.Action, lobbyUsers map[string]interface{}) {
	l.broadcastData(actionName, action.Payload()[0], lobbyUsers)
}

func (l *lobbyRoom) broadcastData(actionName string, data interface{}, lobbyUsers map[string]interface{}) {
	for k := range lobbyUsers {
		flux.Send(actionName, l.uuid, k, data)
	}
}

func (l *lobbyRoom) Destroy() {
	flux.UnRegister(l.uuid, l.in)
}
