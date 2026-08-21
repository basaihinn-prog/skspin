package gameroom

import (
	"crypto/sha1"
	"fmt"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/domain/fish"
	"serve/service_fish/models"
	"strconv"
	"time"

	"serve/fish_comm/broadcaster"
	"serve/fish_comm/common"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

const (
	actionRoomDelete            = "actionRoomDelete"
	actionVIPRoomDelete         = "actionVIPRoomDelete"
	ActionRoomAutoJoinCall      = "ActionRoomAutoJoinCall"
	ActionRoomAutoJoinRecall    = "ActionRoomAutoJoinRecall"
	ActionRoomLeaveCall         = "ActionRoomLeaveCall"
	ActionRoomLeaveRecall       = "ActionRoomLeaveRecall"
	ActionGameRoomBotJoinCall   = "ActionGameRoomBotJoinCall"
	ActionGameRoomBotJoinRecall = "ActionGameRoomBotJoinRecall"
)

var Service = &service{
	Id:            "GameRoomService",
	reservedRooms: 1,
	roomSize:      4,
	in:            make(chan *flux.Action, common.Service.ChanSize),
}

type service struct {
	Id            string
	reservedRooms int
	roomSize      int
	in            chan *flux.Action
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

	allGames := make(map[string]*gameRoom)

	for {
		select {
		case action, ok := <-s.in:
			if !ok {
				return
			}

			switch action.Key().Name() {
			case ActionGameRoomBotJoinCall:
				j := action.Payload()[0].(*JoinGameRoom)
				isJoined := false

				if r, ok := allGames[j.GameRoomUuid]; ok {
					flux.Send(actionCheckVacancy, s.Id, r.uuid, ActionGameRoomBotJoinCall)
					if <-r.isVacancy {
						flux.Send(actionAddUserCall, s.Id, r.uuid, j)
						if <-r.isJoined {
							isJoined = true
							logger.Service.Zap.Infow("bot user join game room",
								"GameUser", j.SecWebSocketKey,
								"GameRoomUuid", r.uuid,
								"GameId", j.gameId,
								"SeatId", j.SeatId,
								"BetList", j.betList,
								"RateList", j.rateList,
								"Rate", j.rate,
								"AllGameRooms", len(allGames),
							)
						}
					}
				}

				if !isJoined {
					flux.Send(ActionGameRoomBotJoinRecall, s.Id, j.controllerId, j)
				}

			case ActionRoomAutoJoinCall:
				j := action.Payload()[0].(*JoinGameRoom)

				count := 0

				for uuid, r := range allGames {
					count++

					// Process VIP Room One Room only have one Player with SecWebSocketKey
					if s.checkVipRoom(j.gameId, j.rate) {
						if r.hashId() == j.hashId() && r.secWebSocketKey == j.SecWebSocketKey {
							flux.Send(actionCheckVacancy, s.Id, uuid, ActionRoomAutoJoinCall)

							if <-r.isVacancy {
								flux.Send(actionAddUserCall, s.Id, uuid, j)

								if <-r.isJoined {
									logger.Service.Zap.Infow("GameUser is joined the existed room.",
										"GameUser", j.SecWebSocketKey,
										"GameRoomUuid", uuid,
										"GameId", j.gameId,
										"SeatId", j.SeatId,
										"BetList", j.betList,
										"RateList", j.rateList,
										"Rate", j.rate,
										"AllGameRooms", len(allGames),
									)
									flux.Send(actionVIPRoomLive, s.Id, j.GameRoomUuid)

									break
								}
							}
						}
					} else {
						if r.hashId() == j.hashId() {
							flux.Send(actionCheckVacancy, s.Id, uuid, ActionRoomAutoJoinCall)

							if <-r.isVacancy {
								flux.Send(actionAddUserCall, s.Id, uuid, j)

								if <-r.isJoined {
									logger.Service.Zap.Infow("GameUser is joined the existed room.",
										"GameUser", j.SecWebSocketKey,
										"GameRoomUuid", uuid,
										"GameId", j.gameId,
										"SeatId", j.SeatId,
										"BetList", j.betList,
										"RateList", j.rateList,
										"Rate", j.rate,
										"AllGameRooms", len(allGames),
									)
									break
								}
							}
						}
					}

					// all rooms full
					if count >= len(allGames) {
						r := s.createRoom(j)
						allGames[r.uuid] = r
						break
					}
				}

				// first game room
				if len(allGames) == 0 {
					r := s.createRoom(j)
					allGames[r.uuid] = r
				}

				logger.Service.Zap.Infow("All Game Rooms",
					"GameUser", j.SecWebSocketKey,
					"AllGameRooms", len(allGames),
				)

			case actionAddUserRecall:
				j := action.Payload()[0].(*JoinGameRoom)

				switch j.actionType {
				case ActionGameRoomBotJoinCall:
					flux.Send(ActionGameRoomBotJoinRecall, s.Id, j.controllerId, j)

				case ActionRoomAutoJoinCall:
					flux.Send(ActionRoomAutoJoinRecall, s.Id, j.controllerId, j)
				}

			case ActionRoomLeaveCall:
				l := action.Payload()[0].(*LeaveGameRoom)
				flux.Send(actionDeleteUserCall, s.Id, l.GameRoomUuid, l)

			case actionDeleteUserRecall:
				l := action.Payload()[0].(*LeaveGameRoom)
				flux.Send(ActionRoomLeaveRecall, s.Id, l.controllerId, l)

			case actionRoomDelete:
				d := action.Payload()[0].(*deleteGameRoom)

				count := 0

				for _, r := range allGames {
					if r.hashId() == d.hashId() {
						count++
					}
				}

				if s.checkVipRoom(d.gameId, d.rate) {
					r := allGames[d.gameRoomUuid]
					go s.vipRoomDelete(r, d)
				} else {
					if count > s.reservedRooms {
						if r, ok := allGames[d.gameRoomUuid]; ok {
							flux.Send(fish.ActionFishStop, d.gameRoomUuid, fish.Service.Id, "")
							flux.Send(
								broadcaster.ActionBroadcasterUnRegister,
								s.Id,
								broadcaster.Service.Id,
								broadcaster.NewUnRegisterBroadcaster(d.gameRoomUuid, d.gameId, d.betList, d.rateList, d.mathModuleId, d.rate),
							)

							r.destroy()
							delete(allGames, d.gameRoomUuid)
						}
					}
				}

			case actionVIPRoomDelete:
				d := action.Payload()[0].(*deleteGameRoom)

				if r, ok := allGames[d.gameRoomUuid]; ok {
					flux.Send(fish.ActionFishStop, d.gameRoomUuid, fish.Service.Id, "")
					flux.Send(
						broadcaster.ActionBroadcasterUnRegister,
						s.Id,
						broadcaster.Service.Id,
						broadcaster.NewUnRegisterBroadcaster(d.gameRoomUuid, d.gameId, d.betList, d.rateList, d.mathModuleId, d.rate),
					)

					r.destroy()
					delete(allGames, d.gameRoomUuid)
				}
			}
		}
	}
}

func (s *service) createRoom(j *JoinGameRoom) *gameRoom {
	// VIP廳只有單獨玩家
	var r = &gameRoom{}
	subGameId := gamesetting.Service.SubgameId(j.SecWebSocketKey)
	switch {
	case j.gameId == models.PSF_ON_00003:
		r = newVipGameRoom(j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate, j.SecWebSocketKey)
	case j.gameId == models.PSF_ON_00004 && subGameId == 2:
		r = newVipGameRoom(j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate, j.SecWebSocketKey)
	case j.gameId == models.PSF_ON_00005 && subGameId == 2:
		r = newVipGameRoom(j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate, j.SecWebSocketKey)
	case j.gameId == models.PSF_ON_00006 && subGameId == 2:
		r = newVipGameRoom(j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate, j.SecWebSocketKey)
	case j.gameId == models.PSF_ON_00007 && subGameId == 2:
		r = newVipGameRoom(j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate, j.SecWebSocketKey)
	case j.gameId == models.RKF_H5_00001 && subGameId == 2:
		r = newVipGameRoom(j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate, j.SecWebSocketKey)
	default:
		r = newGameRoom(j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate, s.roomSize, j.SecWebSocketKey)
	}

	if <-r.isReady {
		flux.Send(fish.ActionFishStart, r.uuid, fish.Service.Id, j.gameId, j.mathModuleId)

		logger.Service.Zap.Infow("New GameRoom is created.",
			"GameUser", j.SecWebSocketKey,
			"GameRoomUuid", r.uuid,
			"GameId", j.gameId,
			"BetList", j.betList,
			"RateList", j.rateList,
			"Rate", j.rate,
			"MathModuleId", j.mathModuleId,
		)

		flux.Send(
			broadcaster.ActionBroadcasterRegister,
			s.Id,
			broadcaster.Service.Id,
			broadcaster.NewRegisterBroadcaster(r.uuid, j.gameId, j.betList, j.rateList, j.mathModuleId, j.rate),
		)

		flux.Send(actionAddUserCall, s.Id, r.uuid, j)

		if <-r.isJoined {
			logger.Service.Zap.Infow("GameUser is joined the new room.",
				"GameUser", j.SecWebSocketKey,
				"GameRoomUuid", r.uuid,
				"GameId", j.gameId,
				"BetList", j.betList,
				"RateList", j.rateList,
				"Rate", j.rate,
				"MathModuleId", j.mathModuleId,
			)
		}
	}
	return r
}

func (s *service) vipRoomDelete(r *gameRoom, d *deleteGameRoom) {
	t := time.NewTimer(time.Duration(5) * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-r.isLive:
			return

		case <-t.C:
			flux.Send(actionVIPRoomDelete, s.Id, s.Id, d)
		}
	}
}

func (s *service) checkVipRoom(gameId string, rate uint64) bool {
	switch {
	case gameId == models.PSF_ON_00003:
		fallthrough
	case gameId == models.PSF_ON_00004 && rate == 1000:
		fallthrough
	case gameId == models.PSF_ON_00005 && rate == 1000:
		fallthrough
	case gameId == models.PSF_ON_00006 && rate == 1000:
		fallthrough
	case gameId == models.PSF_ON_00007 && rate == 1000:
		fallthrough
	case gameId == models.RKF_H5_00001 && rate == 1000:
		return true
	}

	return false
}

func (s *service) hashId(gameId, betList, rateList, mathModuleId string, rate uint64) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(gameId+betList+rateList+mathModuleId+strconv.FormatUint(rate, 10))))
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
