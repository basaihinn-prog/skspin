package lobbyroom

import (
	"fmt"
	"serve/fish_comm/common"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

const (
	ActionLobbyRoomAutoJoinCall   = "ActionLobbyRoomAutoJoinCall"
	ActionLobbyRoomAutoJoinRecall = "ActionLobbyRoomAutoJoinRecall"
	ActionLobbyRoomLeaveCall      = "ActionLobbyRoomLeaveCall"
	ActionLobbyRoomLeaveRecall    = "ActionLobbyRoomLeaveRecall"
)

var Service = &service{
	Id: "LobbyRoomService",
	in: make(chan *flux.Action, common.Service.ChanSize),
}

type service struct {
	Id string
	in chan *flux.Action
}

func init() {
	go Service.run()
	logger.Service.Zap.Infow("Service created.",
		"Service", Service.Id,
		"Chan", fmt.Sprintf("%p", Service.in),
	)
}

func (s *service) run() {
	flux.Register(s.Id, s.in)
	lobbyRooms := make(map[string]*lobbyRoom)

	// only one lobby room
	//l := newLobbyRoom()
	//lobbyRoom[l.uuid] = l

	for {
		select {
		case action, ok := <-s.in:
			if !ok {
				return
			}

			switch action.Key().Name() {
			case ActionLobbyRoomAutoJoinCall:
				secWebSocketKey := action.Payload()[0].(string)
				gameId := action.Payload()[1].(string)
				controllerId := action.Key().From()

				var l *lobbyRoom
				if lobby, ok := lobbyRooms[gameId]; ok {
					l = lobby
				} else {
					l = newLobbyRoom()
					lobbyRooms[gameId] = l
				}

				flux.Send(actionAddUserCall, s.Id, l.uuid, secWebSocketKey, controllerId, gameId)

			case actionAddUserRecall:
				secWebSocketKey := action.Payload()[0].(string)
				controllerId := action.Payload()[1].(string)
				gameId := action.Payload()[2].(string)
				flux.Send(ActionLobbyRoomAutoJoinRecall, s.Id, controllerId, secWebSocketKey, lobbyRooms[gameId].uuid)

			case ActionLobbyRoomLeaveCall:
				l := action.Payload()[0].(*LeaveLobbyRoom)
				flux.Send(actionDeleteUserCall, s.Id, l.LobbyRoomUuid, l)

			case actionDeleteUserRecall:
				l := action.Payload()[0].(*LeaveLobbyRoom)
				flux.Send(ActionLobbyRoomLeaveRecall, s.Id, l.controllerId, l)
			}
		}
	}
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
