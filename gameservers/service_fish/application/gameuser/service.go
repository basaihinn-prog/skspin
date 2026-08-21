package gameuser

import (
	"fmt"
	"os"
	"serve/fish_comm/common"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"
	"serve/service_fish/domain/auth"
	"strconv"

	"github.com/gorilla/websocket"
	sync "github.com/sasha-s/go-deadlock"
)

const (
	ActionGameUserCreate = "ActionGameUserCreate"
	ActionGameUserDelete = "ActionGameUserDelete"
	ActionGameUserLeave  = "ActionGameUserLeave"
)

var Service = &service{
	Id:           "GameUserService",
	poolSize:     os.Getenv("POOL_SIZE"),
	in:           make(chan *flux.Action, common.Service.ChanSize),
	mutex:        sync.Mutex{},
	allGameUsers: make(map[string]*GameUser),
}

type service struct {
	Id           string
	poolSize     string
	in           chan *flux.Action
	mutex        sync.Mutex
	allGameUsers map[string]*GameUser
}

func init() {
	Service.run()
	logger.Service.Zap.Infow("Service created.",
		"Service", Service.Id,
		"Chan", fmt.Sprintf("%p", Service.in),
	)
}

func (s *service) run() {
	flux.Register(s.Id, s.in)
	//allGameUsers := make(map[string]*GameUser)

	poolSize := 10

	if s.poolSize != "" {
		if poolSize, _ = strconv.Atoi(s.poolSize); poolSize == 0 {
			poolSize = 10
		}
	}

	for i := 0; i < poolSize*10; i++ {
		go func() {
			for action := range s.in {
				s.handleAction(action, s.allGameUsers)
			}
		}()
	}
}

func (s *service) handleAction(action *flux.Action, allGameUsers map[string]*GameUser) {
	switch action.Key().Name() {
	case ActionGameUserCreate:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		hostExtId := action.Payload()[0].(string)
		secWebSocketKey := action.Payload()[1].(string)
		conn := action.Payload()[2].(*websocket.Conn)

		allGameUsers[secWebSocketKey] = newGameUser(action.Key().From(), hostExtId, secWebSocketKey, conn)

		logger.Service.Zap.Infow("Add Game User",
			"GameUser", secWebSocketKey,
			"HostExtId", hostExtId,
			"AllGameUsers", len(allGameUsers),
		)

	case auth.ActionUserDisconnect:
		s.mutex.Lock()

		secWebSocketKey := action.Payload()[0].(string)

		if v, ok := allGameUsers[secWebSocketKey]; ok {
			s.mutex.Unlock()

			v.destroy()
		} else {
			s.mutex.Unlock()

			logger.Service.Zap.Warnw(GameUser_NOT_FOUND,
				"GameUser", secWebSocketKey,
				"AllGameUsers", len(allGameUsers),
			)
			errorcode.Service.Fatal(secWebSocketKey, GameUser_NOT_FOUND)
		}

	case ActionGameUserDelete:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		secWebSocketKey := action.Payload()[0].(string)

		if v, ok := allGameUsers[secWebSocketKey]; ok {
			go v.cleanData()
			delete(allGameUsers, secWebSocketKey)

			logger.Service.Zap.Infow("Delete Game User",
				"GameUser", secWebSocketKey,
				"AllGameUsers", len(allGameUsers),
			)
		}

	case ActionGameUserLeave:
		s.mutex.Lock()

		secWebSocketKey := action.Payload()[0].(string)

		if _, ok := allGameUsers[secWebSocketKey]; ok {
			s.mutex.Unlock()

			balance, err := wallet.Service.Balance(secWebSocketKey)

			if err != nil {
				logger.Service.Zap.Warnw(GameUser_GET_BALANCE_FAILED,
					"GameUser", secWebSocketKey,
					"Error", err,
				)
			}

			flux.Send(auth.ActionAuthUserInfoAsk, s.Id, auth.Service.Id, secWebSocketKey, balance)

		} else {
			s.mutex.Unlock()

			logger.Service.Zap.Warnw(GameUser_NOT_FOUND,
				"GameUser", secWebSocketKey,
				"AllGameUsers", len(allGameUsers),
			)
			errorcode.Service.Fatal(secWebSocketKey, GameUser_NOT_FOUND)
		}
	}
}

func (s *service) GetAllGameUsers() map[string]*GameUser {
	return s.allGameUsers
}

func (s *service) GetRoomUuid(secWebSocketKey string) string {
	if v, ok := s.allGameUsers[secWebSocketKey]; ok {
		if v.lobbyRoomUuid != "" {
			return v.lobbyRoomUuid
		}

		if v.gameRoomUuid != "" {
			return v.gameRoomUuid
		}
	}

	return ""
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
