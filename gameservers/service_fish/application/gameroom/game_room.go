package gameroom

import (
	"fmt"
	"os"
	"serve/fish_comm/broadcaster"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/format"
	"serve/fish_comm/jackpot-client/jackpot"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"
	"serve/service_fish/domain/auth"
	auth_proto "serve/service_fish/domain/auth/proto"
	"serve/service_fish/domain/bet"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"strconv"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
)

const (
	actionCheckVacancy     = "actionCheckVacancy"
	actionAddUserCall      = "actionAddUserCall"
	actionAddUserRecall    = "actionAddUserRecall"
	actionDeleteUserCall   = "actionDeleteUserCall"
	actionDeleteUserRecall = "actionDeleteUserRecall"
	actionVIPRoomLive      = "actionVIPRoomLive"
)

type gameRoom struct {
	gameId          string
	betList         string
	rateList        string
	mathModuleId    string
	rate            uint64
	uuid            string
	poolSize        string
	roomSize        int
	isVacancy       chan bool
	isReady         chan bool
	isJoined        chan bool
	in              chan *flux.Action
	mutex           *sync.Mutex
	nextScene       int
	isLive          chan bool
	secWebSocketKey string
}

func newGameRoom(gameId, betList, rateList, mathModuleId string, rate uint64, roomSize int, secWebSocketKey string) *gameRoom {
	r := &gameRoom{
		gameId:       gameId,
		betList:      betList,
		rateList:     rateList,
		mathModuleId: mathModuleId,
		rate:         rate,
		uuid:         uuid.New().String(),
		poolSize:     os.Getenv("POOL_SIZE"),
		roomSize:     roomSize,
		isVacancy:    make(chan bool, 1024),
		isReady:      make(chan bool, 1024),
		isJoined:     make(chan bool, 1024),
		in:           make(chan *flux.Action, 1024),
		//in:           make(chan *flux.Action, common.Service.ChanSize),
		mutex:           &sync.Mutex{},
		nextScene:       -1,
		isLive:          make(chan bool, 1024),
		secWebSocketKey: secWebSocketKey,
	}
	r.run()
	logger.Service.Zap.Infow("Game Room Started",
		"GameRoomUuid", r.uuid,
		"Chan", fmt.Sprintf("%p", r.in),
	)
	return r
}

func (r *gameRoom) run() {
	flux.Register(r.uuid, r.in)

	players := make(map[int]string, r.roomSize)

	poolSize := 10

	if r.poolSize != "" {
		if poolSize, _ = strconv.Atoi(r.poolSize); poolSize == 0 {
			poolSize = 10
		}
	}

	for i := 0; i < poolSize; i++ {
		go func() {
			for action := range r.in {
				r.handleAction(action, players)
			}
		}()
	}
}

func (r *gameRoom) handleAction(action *flux.Action, players map[int]string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if action.Key().To() != r.uuid {
		return
	}

	switch action.Key().Name() {
	case actionVIPRoomLive:
		r.isLive <- true

	case flux.ActionFluxRegisterDone:
		r.isReady <- true

	case fish.ActionChangeNextScene:
		r.nextScene = action.Payload()[0].(int)
		r.broadcastAction(fish.ActionChangeNextScene, action, players)

	case actionCheckVacancy:
		for k, v := range players {
			if _, err := wallet.Service.BalanceCache(v); err != nil {
				logger.Service.Zap.Warnw("Check Players Alive",
					"GameUser", v,
					"SeatId", k,
					"GameRoomUuid", r.uuid,
					"Chan", fmt.Sprintf("%p", r.in),
					"GameRoomPlayersMap", players,
				)
				delete(players, k)
			}
		}

		if len(players) == 1 {
			for k, v := range players {
				if len(v) == len(r.uuid) {
					delete(players, k)
				}
			}
		}

		switch action.Payload()[0].(string) {
		case ActionGameRoomBotJoinCall:
			if (len(players) >= 2) || (len(players) == 0) {
				r.isVacancy <- false

			} else {
				r.isVacancy <- len(players) < r.roomSize
			}

		case ActionRoomAutoJoinCall:
			r.isVacancy <- len(players) < r.roomSize
		}

	case actionAddUserCall:
		j := action.Payload()[0].(*JoinGameRoom)

		for i := 0; i < r.roomSize; i++ {
			if _, ok := players[i]; !ok {
				j.GameRoomUuid = r.uuid
				j.SeatId = i
				j.NextScene = r.nextScene
				players[i] = j.SecWebSocketKey

				for k, v := range players {
					j.Players[k] = v
				}

				flux.Send(actionAddUserRecall, r.uuid, action.Key().From(), j)
				break
			}
		}

		logger.Service.Zap.Infow("Join Game Room",
			"GameUser", j.SecWebSocketKey,
			"SeatId", j.SeatId,
			"GameRoomUuid", r.uuid,
			"Chan", fmt.Sprintf("%p", r.in),
			"GameRoomPlayers", len(players),
			"GameRoomPlayersMap", players,
		)

		r.isJoined <- true

	case actionDeleteUserCall:
		l := action.Payload()[0].(*LeaveGameRoom)

		for k, v := range players {
			if v == l.SecWebSocketKey {
				delete(players, k)
			}

			if _, err := wallet.Service.BalanceCache(v); err != nil {
				logger.Service.Zap.Warnw("Check Players Alive",
					"GameUser", v,
					"SeatId", k,
					"GameRoomUuid", r.uuid,
					"Chan", fmt.Sprintf("%p", r.in),
					"GameRoomPlayersMap", players,
				)
				delete(players, k)
			}
		}

		if len(players) == 1 {
			for k, v := range players {
				if len(v) == len(r.uuid) {
					delete(players, k)
				}
			}
		}

		for k, v := range players {
			l.Players[k] = v
		}

		logger.Service.Zap.Infow("Leave Game Room",
			"GameUser", l.SecWebSocketKey,
			"GameRoomUuid", r.uuid,
			"Chan", fmt.Sprintf("%p", r.in),
			"GameRoomPlayers", len(players),
			"GameRoomPlayersMap", players,
		)

		flux.Send(actionDeleteUserRecall, r.uuid, action.Key().From(), l)

		if len(players) == 0 {
			flux.Send(actionRoomDelete, r.uuid, action.Key().From(), newDeleteGameRoom(
				r.uuid, r.gameId, r.betList, r.rateList, r.mathModuleId, r.rate),
			)
		}

	case fish.EMSGID_eBroadcastFishIn:
		r.broadcastAction(fish.EMSGID_eBroadcastFishIn, action, players)

	case fish.EMSGID_eBroadcastFishOut:
		r.broadcastAction(fish.EMSGID_eBroadcastFishOut, action, players)

	case fish.EMSGID_eBroadcastGroupIn:
		r.broadcastAction(fish.EMSGID_eBroadcastGroupIn, action, players)

	case fish.EMSGID_eBroadcastRoundStart:
		r.broadcastAction(fish.EMSGID_eBroadcastRoundStart, action, players)

	case fish.EMSGID_eBroadcastRoundEnd:
		r.broadcastAction(fish.EMSGID_eBroadcastRoundEnd, action, players)

	case fish.EMSGID_eBroadcastFishChange:
		r.broadcastAction(fish.EMSGID_eBroadcastFishChange, action, players)

	case bullet.EMSGID_eBroadcastShoot:

		// TODO JOHNNY 這邊拿出來

		r.broadcastAction(bullet.EMSGID_eBroadcastShoot, action, players)

	case bullet.EMSGID_eBroadcastBetChange:
		r.broadcastAction(bullet.EMSGID_eBroadcastBetChange, action, players)

	case broadcaster.EMSGID_eBroadcaster:
		r.broadcastAction(broadcaster.EMSGID_eBroadcaster, action, players)

	case auth.EMSGID_eBroadcastPlayerIn:
		secWebSocketKey := action.Key().From() // include player or bot
		playerType := action.Payload()[0].(auth_proto.PLAYER_TYPE)

		data := auth.Service.BroadcastPlayerIn(r.uuid)

		for k, v := range players {

			balance, err := wallet.Service.BalanceCache(v)

			if err != nil {
				logger.Service.Zap.Warnw("Check Players Alive",
					"GameUser", v,
					"SeatId", k,
					"GameRoomUuid", r.uuid,
					"Chan", fmt.Sprintf("%p", r.in),
					"GameRoomPlayersMap", players,
				)
				delete(players, k)
				return
			}

			u := &auth_proto.UserInfo{
				PlayerId:    wallet.Service.MemberId(v),
				PlayerCent:  balance,
				PlayerName:  format.HideWithStar(wallet.Service.MemberName(v)),
				SeatId:      uint32(k),
				BetRateLine: bet.Service.Get(v),
				Visible:     nil,
			}

			if len(v) == len(secWebSocketKey) {
				u.PlayerType = playerType
			} else {
				switch playerType {
				case auth_proto.PLAYER_TYPE_GAME:
					u.PlayerType = auth_proto.PLAYER_TYPE_BOT

				case auth_proto.PLAYER_TYPE_BOT:
					u.PlayerType = auth_proto.PLAYER_TYPE_GAME
				}
			}

			data.Player = append(data.Player, u)
		}

		dataBuffer, _ := proto.Marshal(data)
		r.broadcastData(auth.EMSGID_eBroadcastPlayerIn, dataBuffer, players)

	case auth.EMSGID_eBroadcastPlayerOut:
		r.broadcastAction(auth.EMSGID_eBroadcastPlayerOut, action, players)

	case probability.EMSGID_eBroadcastResult:
		r.broadcastAction(probability.EMSGID_eBroadcastResult, action, players)

	case probability.EMSGID_eBroadcastOption:
		r.broadcastAction(probability.EMSGID_eBroadcastOption, action, players)

	case errorcode.ActionErrorCode:
		r.broadcastAction(errorcode.ActionErrorCode, action, players)

	case jackpot.EMSGID_eJackpotInfo:
		r.broadcastAction(jackpot.EMSGID_eJackpotInfo, action, players)

	case jackpot.EMSGID_eJackpotNotify:
		r.broadcastAction(jackpot.EMSGID_eJackpotNotify, action, players)

	}
}

func (r *gameRoom) broadcastAction(actionName string, action *flux.Action, players map[int]string) {
	r.broadcastData(actionName, action.Payload()[0], players)
}

func (r *gameRoom) broadcastData(actionName string, data interface{}, players map[int]string) {
	for _, v := range players {
		flux.Send(actionName, r.uuid, v, data)
	}
}

func (r *gameRoom) hashId() string {
	return Service.hashId(
		r.gameId,
		r.betList,
		r.rateList,
		r.mathModuleId,
		r.rate,
	)
}

func (r *gameRoom) destroy() {
	logger.Service.Zap.Infow("Game Room Unregister",
		"GameRoomUuid", r.uuid,
		"Chan", fmt.Sprintf("%p", r.in),
	)

	flux.UnRegister(r.uuid, r.in)
	close(r.isReady)
	close(r.isVacancy)
	close(r.isJoined)
}
