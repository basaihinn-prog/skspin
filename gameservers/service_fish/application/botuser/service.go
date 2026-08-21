package botuser

import (
	"fmt"
	"os"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"strconv"
	"sync"
)

const (
	ActionBotUserCheck  = "ActionBotUserCheck"
	ActionBotUserCreate = "ActionBotUserCreate"
	ActionBotUserDelete = "ActionBotUserDelete"
)

var Service = &service{
	Id:       "BotUserService",
	poolSize: os.Getenv("POOL_SIZE"),
	in:       make(chan *flux.Action, 1024),
	rwMutex:  &sync.RWMutex{},
}

type service struct {
	Id       string
	poolSize string
	in       chan *flux.Action
	rwMutex  *sync.RWMutex
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
	allBotUsers := make(map[string]*botUser)

	poolSize := 10

	if s.poolSize != "" {
		if poolSize, _ = strconv.Atoi(s.poolSize); poolSize == 0 {
			poolSize = 10
		}
	}

	for i := 0; i < poolSize*10; i++ {
		go func() {
			for action := range s.in {
				s.handleAction(action, allBotUsers)
			}
		}()
	}
}

func (s *service) handleAction(action *flux.Action, allBotUsers map[string]*botUser) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	logger.Service.Zap.Debugw("AllBots",
		"Size", len(allBotUsers),
		"Map", allBotUsers,
	)

	switch action.Key().Name() {
	case ActionBotUserCheck:
		controllerId := action.Key().From()
		gameRoomUuid := action.Payload()[0].(string)
		gameRoomPlayers := action.Payload()[1].(map[int]string)

		switch {
		case len(gameRoomPlayers) == 1:
			for _, v := range gameRoomPlayers {
				if b, ok := allBotUsers[v]; ok {
					b.destroy()
					delete(allBotUsers, v)

					logger.Service.Zap.Infow("Deleted Bot",
						"GameUser", b.secWebSocketKey,
						"Bot", b.uuid,
						"GameRoomUuid", b.gameRoomUuid,
						"AllBotUsers", len(allBotUsers),
						"GameRoomPlayersMap", gameRoomPlayers,
					)
				} else {
					if b := newBotUser(controllerId, v, gameRoomUuid); b != nil {
						allBotUsers[b.uuid] = b

						logger.Service.Zap.Infow("Added Bot",
							"GameUser", b.secWebSocketKey,
							"Bot", b.uuid,
							"Rate", b.rate,
							"MaxBet", b.maxBet,
							"RateIndex", b.rateIndex,
							"GameRoomUuid", b.gameRoomUuid,
							"AllBotUsers", len(allBotUsers),
							"GameRoomPlayersMap", gameRoomPlayers,
						)
					}
				}
			}

		case len(gameRoomPlayers) == 2:
			for k, v := range gameRoomPlayers {
				if b, ok := allBotUsers[v]; ok {
					flux.Send(ActionBotUserJoinGameRoomRecall, controllerId, b.uuid, k)

					logger.Service.Zap.Infow("Running Bot",
						"GameUser", b.secWebSocketKey,
						"Bot", b.uuid,
						"Rate", b.rate,
						"MaxBet", b.maxBet,
						"RateIndex", b.rateIndex,
						"GameRoomUuid", b.gameRoomUuid,
						"AllBotUsers", len(allBotUsers),
						"GameRoomPlayersMap", gameRoomPlayers,
					)
				}
			}

		default:
			for k, v := range allBotUsers {
				if v.gameRoomUuid == gameRoomUuid {
					v.destroy()
					delete(allBotUsers, k)

					logger.Service.Zap.Infow("Deleted Game Room's Bots",
						"GameUser", v.secWebSocketKey,
						"Bot", v.uuid,
						"Rate", v.rate,
						"MaxBet", v.maxBet,
						"RateIndex", v.rateIndex,
						"GameRoomUuid", v.gameRoomUuid,
						"AllBotUsers", len(allBotUsers),
						"GameRoomPlayersMap", gameRoomPlayers,
					)
				}
			}
		}

		// not used now
	case ActionBotUserCreate:
		c := action.Payload()[0].(*createBotUser)

		b := newBotUser(
			action.Key().From(),
			c.secWebSocketKey,
			c.gameRoomUuid,
		)

		allBotUsers[b.uuid] = b

		logger.Service.Zap.Infow("Added Bot",
			"GameUser", b.secWebSocketKey,
			"Bot", b.uuid,
			"Rate", b.rate,
			"MaxBet", b.maxBet,
			"RateIndex", b.rateIndex,
			"GameRoomUuid", b.gameRoomUuid,
			"AllBotUsers", len(allBotUsers),
		)

	case ActionBotUserDelete:
		secWebSocketKey := action.Payload()[0].(string)

		for k, v := range allBotUsers {
			if v.uuid == secWebSocketKey || v.secWebSocketKey == secWebSocketKey {
				v.destroy()
				delete(allBotUsers, k)

				logger.Service.Zap.Infow("Deleted Bot",
					"GameUser", v.secWebSocketKey,
					"Bot", v.uuid,
					"GameRoomUuid", v.gameRoomUuid,
					"AllBotUsers", len(allBotUsers),
				)
			}
		}
	}
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
