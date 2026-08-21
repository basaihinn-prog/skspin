package gameuser

import (
	"fmt"
	"os"
	"serve/service_fish/application/accountingmanager"
	"serve/service_fish/application/activity"
	"serve/service_fish/application/bonuslottery"
	"serve/service_fish/application/filter"
	"serve/service_fish/application/gamesetting"
	game_setting_proto "serve/service_fish/application/gamesetting/proto"
	"serve/service_fish/application/lobbysetting"
	lobby_setting_proto "serve/service_fish/application/lobbysetting/proto"
	"serve/service_fish/application/lottery"
	"serve/service_fish/application/refund"
	"serve/service_fish/application/rtp"
	"serve/service_fish/domain/auth"
	auth_proto "serve/service_fish/domain/auth/proto"
	"serve/service_fish/domain/bet"
	"serve/service_fish/domain/bullet"
	bullet_proto "serve/service_fish/domain/bullet/proto"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	probability_proto "serve/service_fish/domain/probability/proto"
	common_fish_proto "serve/service_fish/models/proto"
	"strconv"
	"sync"
	"time"

	"serve/fish_comm/blacklist"
	"serve/fish_comm/broadcaster"
	"serve/fish_comm/common"
	common_proto "serve/fish_comm/common/proto"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	heartbeat "serve/fish_comm/heart-beat"
	"serve/fish_comm/idle"
	"serve/fish_comm/jackpot-client/jackpot"
	"serve/fish_comm/logger"
	"serve/fish_comm/maintain"
	"serve/fish_comm/mediator"
	"serve/fish_comm/mysql"
	"serve/fish_comm/session"
	sessionmanager "serve/fish_comm/session-manager"
	"serve/fish_comm/vip"
	"serve/fish_comm/wallet"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 1) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024

	ActionJoinGameRoomCall     = "ActionJoinGameRoomCall"
	ActionJoinGameRoomRecall   = "ActionJoinGameRoomRecall"
	ActionLeaveGameRoomCall    = "ActionLeaveGameRoomCall"
	ActionLeaveGameRoomRecall  = "ActionLeaveGameRoomRecall"
	ActionJoinLobbyRoomCall    = "ActionJoinLobbyRoomCall"
	ActionJoinLobbyRoomRecall  = "ActionJoinLobbyRoomRecall"
	ActionLeaveLobbyRoomCall   = "ActionLeaveLobbyRoomCall"
	ActionLeaveLobbyRoomRecall = "ActionLeaveLobbyRoomRecall"
	ActionBotStart             = "ActionBotStart"
	ActionBotStop              = "ActionBotStop"
	guestOff                   = 0
	guestOn                    = 1
)

type GameUser struct {
	controllerId    string
	lobbyRoomUuid   string
	gameRoomUuid    string
	HostExtId       string
	SecWebSocketKey string
	botUuid         string
	poolSize        string
	seatId          int
	nextScene       int
	conn            *websocket.Conn
	in              chan *flux.Action
	profile         *auth.MemberInfo
	status          *auth.MemberStatus
	isDestroy       bool
	mutex           *sync.Mutex
	maintainChan    chan bool
	banPlayerChan   chan bool
	blacklistValue  uint32
	JackpotGroupId  int32
}

func newGameUser(controllerId, hostExtId, secWebSocketKey string, conn *websocket.Conn) *GameUser {
	u := &GameUser{
		controllerId:    controllerId,
		HostExtId:       hostExtId,
		SecWebSocketKey: secWebSocketKey,
		poolSize:        os.Getenv("POOL_SIZE"),
		seatId:          -1,
		nextScene:       -1,
		conn:            conn,
		in:              make(chan *flux.Action, common.Service.ChanSize),
		status:          &auth.MemberStatus{Alive: true},
		isDestroy:       false,
		mutex:           &sync.Mutex{},
		maintainChan:    make(chan bool),
		banPlayerChan:   make(chan bool),
		blacklistValue:  0,
		JackpotGroupId:  -1,
	}

	isDone := make(chan bool)
	flux.Send(auth.ActionUserHostExtIdSave, secWebSocketKey, auth.Service.Id, hostExtId, isDone)
	<-isDone

	go u.runWriteMessage()
	go u.runReadMessage()

	logger.Service.Zap.Infow("Game User Created",
		"GameUser", u.SecWebSocketKey,
		"Chan", fmt.Sprintf("%p", u.in),
	)
	return u
}

func (u *GameUser) runReadMessage() {
	defer func() {
		logger.Service.Zap.Infow("Game User Read Message Closed",
			"GameUser", u.SecWebSocketKey,
			"Chan", fmt.Sprintf("%p", u.in),
		)
		flux.Send(auth.ActionUserDisconnect, u.SecWebSocketKey, Service.Id, u.SecWebSocketKey)
		u.conn.Close()
	}()

	u.conn.SetReadLimit(maxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(data string) error {
		flux.Send(heartbeat.ActionHeartBeatQuality, u.SecWebSocketKey, heartbeat.Service.Id, data)
		u.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := u.conn.ReadMessage()

		if err != nil {
			logger.Service.Zap.Infow("Game User Read Message Closed",
				"GameUser", u.SecWebSocketKey,
				"Error", err,
			)
			return
		}

		if v, err := filter.FindActionName(msg); err != nil {
			logger.Service.Zap.Errorw(GameUser_MSGID_PROTO_INVALID,
				"GameUser", u.SecWebSocketKey,
				"Error", err,
			)
			errorcode.Service.Fatal(u.SecWebSocketKey, GameUser_MSGID_PROTO_INVALID)
		} else {
			switch v {
			case auth.EMSGID_eLoginCall:
				flux.Send(auth.EMSGID_eLoginCall, u.SecWebSocketKey, auth.Service.Id, msg)

			case probability.EMSGID_eResultCall:
				if u.gameRoomUuid == "" {
					logger.Service.Zap.Warnw(GameUser_GAME_ROOM_UUID_EMPTY,
						"GameUser", u.SecWebSocketKey,
						"Event", probability.EMSGID_eResultCall,
					)
				} else {
					data := &probability_proto.ResultCall{}

					if err := proto.Unmarshal(msg, data); err != nil {
						logger.Service.Zap.Errorw(GameUser_RESULT_CALL_PROTO_INVALID,
							"GameUser", u.SecWebSocketKey,
							"GameRoomUuid", u.gameRoomUuid,
						)
						errorcode.Service.Fatal(u.SecWebSocketKey, GameUser_RESULT_CALL_PROTO_INVALID)
					}

					if int(data.SeatId) == u.seatId {
						flux.Send(probability.EMSGID_eResultCall, u.SecWebSocketKey, lottery.Service.HashId(u.SecWebSocketKey),
							data,
							u.blacklistValue,
							u.profile.DepositMultiple,
							u.status,
						)

						flux.Send(idle.ActionIdleTimeoutReset, u.SecWebSocketKey, idle.Service.HashId(u.SecWebSocketKey),
							u.SecWebSocketKey,
							gamesetting.Service.Idle(u.SecWebSocketKey),
						)
					}

					u.mutex.Lock()
					if u.botUuid != "" {
						flux.Send(probability.EMSGID_eResultCall, u.SecWebSocketKey, u.botUuid, data)
					}
					u.mutex.Unlock()
				}

			case probability.EMSGID_eMultiResultCall:
				if u.gameRoomUuid == "" {
					logger.Service.Zap.Infow(GameUser_GAME_ROOM_UUID_EMPTY,
						"GameUser", u.SecWebSocketKey,
						"Event", probability.EMSGID_eMultiResultCall,
					)
				} else {
					data := &probability_proto.MultiResultCall{}

					if err := proto.Unmarshal(msg, data); err != nil {
						logger.Service.Zap.Infow(GameUser_RESULT_CALL_PROTO_INVALID,
							"GameUser", u.SecWebSocketKey,
							"GameRoomUuid", u.gameRoomUuid,
						)
						errorcode.Service.Fatal(u.SecWebSocketKey, GameUser_RESULT_CALL_PROTO_INVALID)
					}

					if int(data.SeatId) == u.seatId {
						flux.Send(probability.EMSGID_eMultiResultCall, u.SecWebSocketKey, lottery.Service.HashId(u.SecWebSocketKey),
							data,
							u.blacklistValue,
							u.profile.DepositMultiple,
							u.status,
						)

						flux.Send(idle.ActionIdleTimeoutReset, u.SecWebSocketKey, idle.Service.HashId(u.SecWebSocketKey),
							u.SecWebSocketKey,
							gamesetting.Service.Idle(u.SecWebSocketKey),
						)
					}
				}

			case bullet.EMSGID_eShoot:
				u.mutex.Lock()

				if u.gameRoomUuid == "" {
					logger.Service.Zap.Warnw(GameUser_GAME_ROOM_UUID_EMPTY,
						"GameUser", u.SecWebSocketKey,
						"Event", bullet.EMSGID_eShoot,
					)
				} else {
					data := &bullet_proto.Shoot{}

					if err := proto.Unmarshal(msg, data); err != nil {
						logger.Service.Zap.Errorw(bullet.Bullet_SHOOT_PROTO_INVALID, "GameUser", u.SecWebSocketKey)
						errorcode.Service.Fatal(u.SecWebSocketKey, bullet.Bullet_SHOOT_PROTO_INVALID)
					}
					gameId := gamesetting.Service.GameId(u.SecWebSocketKey)
					bet := gamesetting.Service.Bet(u.SecWebSocketKey, gameId, data.BetRateLine.BetLevelIndex, data.BetRateLine.BetIndex)
					rate := gamesetting.Service.Rate(u.SecWebSocketKey)

					if bet == 0 {
						logger.Service.Zap.Errorw(gamesetting.GameSetting_BET_ZERO, "GameUser", u.SecWebSocketKey)
						errorcode.Service.Fatal(u.SecWebSocketKey, gamesetting.GameSetting_BET_ZERO)
					}

					if rate == 0 {
						logger.Service.Zap.Errorw(gamesetting.GameSetting_RATE_ZERO, "GameUser", u.SecWebSocketKey)
						errorcode.Service.Fatal(u.SecWebSocketKey, gamesetting.GameSetting_RATE_ZERO)
					}

					balance, err := wallet.Service.BalanceCache(u.SecWebSocketKey)
					if wallet.Service.WalletCategory(u.SecWebSocketKey) == wallet.CategoryHost {
						cent, _ := wallet.Service.CentCache(u.SecWebSocketKey)
						balance += cent
					}

					if err != nil {
						logger.Service.Zap.Errorw(GameUser_GET_BALANCE_FAILED,
							"GameUser", u.SecWebSocketKey,
							"Error", err,
						)
						errorcode.Service.Fatal(u.SecWebSocketKey, GameUser_GET_BALANCE_FAILED)
					}

					if int(data.SeatId) == u.seatId {
						mediator.Service.New(u.SecWebSocketKey, data.BulletUuid) // 1. 29b250c5
						mediator.Service.Add(u.SecWebSocketKey, data.BulletUuid) // 2. 29b250c5

						flux.Send(bullet.EMSGID_eShoot, u.SecWebSocketKey, bullet.Service.HashId(u.SecWebSocketKey),
							data,
							bet,
							rate,
							balance,
						)

						flux.Send(idle.ActionIdleTimeoutReset, u.SecWebSocketKey, idle.Service.HashId(u.SecWebSocketKey),
							u.SecWebSocketKey,
							gamesetting.Service.Idle(u.SecWebSocketKey),
						)
					}

					if u.botUuid != "" && int(data.SeatId) != u.seatId {
						flux.Send(bullet.EMSGID_eShoot, u.SecWebSocketKey, u.botUuid, data, bet, rate)
					}
				}
				u.mutex.Unlock()

			case bullet.EMSGID_eBetChange:
				u.mutex.Lock()

				if u.gameRoomUuid == "" {
					logger.Service.Zap.Warnw(GameUser_GAME_ROOM_UUID_EMPTY,
						"GameUser", u.SecWebSocketKey,
						"Event", bullet.EMSGID_eBetChange,
					)
				} else {
					data := &bullet_proto.BetChange{}

					if err := proto.Unmarshal(msg, data); err != nil {
						logger.Service.Zap.Errorw(bullet.Bullet_BET_CHANGE_PROTO_INVALID, "GameUser", u.SecWebSocketKey)
						errorcode.Service.Fatal(u.SecWebSocketKey, bullet.Bullet_BET_CHANGE_PROTO_INVALID)
					}

					if int(data.SeatId) == u.seatId {
						bet.Service.Set(u.SecWebSocketKey, data.BetRateLine)

						mediator.Service.New(u.SecWebSocketKey, bullet.EMSGID_eBetChange) // 1. ac7a2c41
						mediator.Service.Add(u.SecWebSocketKey, bullet.EMSGID_eBetChange) // 2. ac7a2c41

						flux.Send(bullet.EMSGID_eBetChange, u.SecWebSocketKey, bullet.Service.HashId(u.SecWebSocketKey), data)
						flux.Send(idle.ActionIdleTimeoutReset, u.SecWebSocketKey, idle.Service.HashId(u.SecWebSocketKey),
							u.SecWebSocketKey,
							gamesetting.Service.Idle(u.SecWebSocketKey),
						)

						mediator.Service.Wait(u.SecWebSocketKey, bullet.EMSGID_eBetChange) // 4. ac7a2c41
					}

					if u.botUuid != "" && int(data.SeatId) != u.seatId {
						mediator.Service.New(u.botUuid, bullet.EMSGID_eBetChange) // 1. d3f63a43
						mediator.Service.Add(u.botUuid, bullet.EMSGID_eBetChange) // 2. d3f63a43

						flux.Send(bullet.EMSGID_eBetChange, u.SecWebSocketKey, u.botUuid, data)

						mediator.Service.Wait(u.botUuid, bullet.EMSGID_eBetChange) // 4. d3f63a43
					}
				}
				u.mutex.Unlock()

			case fish.EMSGID_eAllFishCall:
				u.mutex.Lock()

				if u.gameRoomUuid == "" {
					logger.Service.Zap.Warnw(GameUser_GAME_ROOM_UUID_EMPTY,
						"GameUser", u.SecWebSocketKey,
						"Event", fish.EMSGID_eAllFishCall,
					)
				} else {
					flux.Send(fish.EMSGID_eAllFishCall, u.SecWebSocketKey, fish.Service.HashId(u.gameRoomUuid), msg)
				}
				u.mutex.Unlock()

			case lobbysetting.EMSGID_eLobbyConfigCall:
				// Process Maintain Component with time 30 Second
				data := &lobby_setting_proto.ConfigCall{}
				if err := proto.Unmarshal(msg, data); err != nil {
					logger.Service.Zap.Errorw(lobbysetting.LobbySetting_CONFIG_CALL_PROTO_INVALID,
						"GameUser", u.SecWebSocketKey,
						"MemberId", u.profile.MemberId,
					)
					errorcode.Service.Fatal(u.SecWebSocketKey, lobbysetting.LobbySetting_CONFIG_CALL_PROTO_INVALID)
					return
				}

				if u.gameRoomUuid != "" {
					flux.Send(auth.EMSGID_eBroadcastPlayerOut, u.SecWebSocketKey, auth.Service.Id,
						u.gameRoomUuid,
						uint32(u.seatId),
						uint64(0),
					)
					flux.Send(ActionLeaveGameRoomCall, u.SecWebSocketKey, u.controllerId, u.gameRoomUuid, false)
					if am := accountingmanager.Service.Get(u.SecWebSocketKey); am != nil {
						done := make(chan bool)

						am.LeaveGameRoom(done)

						for {
							time.Sleep(1 * time.Second)

							if <-done {
								break
							}

							am.LeaveGameRoom(done)
						}

						close(done)
					}
				} else {
					if u.profile == nil {
						logger.Service.Zap.Errorw(GameUser_PROFILE_NIL,
							"GameUser", u.SecWebSocketKey,
						)
						errorcode.Service.Fatal(u.SecWebSocketKey, GameUser_PROFILE_NIL)
						return
					}
					go maintain.Service.CheckMaintainByTime(u.SecWebSocketKey, data.GameId, u.HostExtId, maintain.Fish, u.maintainChan, 30)
					go blacklist.Service.CheckBanByTime(u.SecWebSocketKey, u.HostExtId, u.profile.HostId, u.profile.MemberId, u.banPlayerChan, 30)
					hasSessionMap := checkSessionMap(u.HostExtId, u.profile.HostId, u.profile.MemberId, data.GameId)

					switch wallet.Service.WalletCategory(u.SecWebSocketKey) {
					case wallet.CategoryHost:
						gameId := data.GameId
						// Update Cent Wallet Purpose
						wallet.Service.SetPurpose(u.SecWebSocketKey, lobbysetting.Service.GetPurposeByGameId(gameId))

						if !hasSessionMap {
							if cent, err := wallet.Service.Cent(u.SecWebSocketKey); err == nil {
								if cent > 0 {
									wallet.Service.UpdateCent(u.SecWebSocketKey, gameId, 0)
								}
							} else {
								logger.Service.Zap.Errorw("BUG",
									"GameUser", u.SecWebSocketKey,
									"Error", err,
								)
							}
						}
					}
				}

				bet.Service.Set(u.SecWebSocketKey, &common_fish_proto.BetRateLine{
					BetLevelIndex: 0,
					BetIndex:      0,
					LineIndex:     0,
					RateIndex:     int32(gamesetting.Service.RateIndex(u.SecWebSocketKey)),
				})

				if u.profile != nil {
					lobbysetting.Service.EMSGID_eLobbyConfigCall(
						u.SecWebSocketKey,
						u.profile.MemberId,
						u.profile.GuestEnable,
						msg,
					)
				} else {
					logger.Service.Zap.Warnw(GameUser_PROFILE_NIL,
						"GameUser", u.SecWebSocketKey,
						"HostExtId", u.HostExtId,
					)
				}

			case gamesetting.EMSGID_eGameConfigCall:
				if u.profile != nil {
					gamesetting.Service.EMSGID_eGameConfigCall(
						u.SecWebSocketKey,
						u.profile.MemberId,
						u.profile.GuestEnable,
						msg,
					)
				} else {
					logger.Service.Zap.Warnw(GameUser_PROFILE_NIL,
						"GameUser", u.SecWebSocketKey,
						"HostExtId", u.HostExtId,
					)
				}

			case gamesetting.EMSGID_eGameStripsCall:
				gamesetting.Service.EMSGID_eGameStripsCall(u.SecWebSocketKey, msg)

			case probability.EMSGID_eOptionCall:
				data := &probability_proto.OptionCall{}

				if err := proto.Unmarshal(msg, data); err != nil {
					logger.Service.Zap.Errorw(GameUser_OPTION_CALL_PROTO_INVALID,
						"GameUser", u.SecWebSocketKey,
						"GameRoomUuid", u.gameRoomUuid,
					)
					errorcode.Service.Fatal(u.SecWebSocketKey, GameUser_OPTION_CALL_PROTO_INVALID)
				}

				if int(data.SeatId) == u.seatId {
					flux.Send(probability.EMSGID_eOptionCall, u.SecWebSocketKey, bonuslottery.Service.Id, data, u.profile.DepositMultiple)
				}

				u.mutex.Lock()
				if u.botUuid != "" && int(data.SeatId) != u.seatId {
					flux.Send(probability.EMSGID_eOptionCall, u.SecWebSocketKey, u.botUuid, data, u.profile.DepositMultiple)
				}
				u.mutex.Unlock()

			case heartbeat.ActionHeartBeatCall:
				flux.Send(heartbeat.ActionHeartBeatCall, u.SecWebSocketKey, heartbeat.Service.Id, msg)

			case probability.EMSGID_eMercenaryOpenCall:
				u.mutex.Lock()

				if u.gameRoomUuid == "" {
					logger.Service.Zap.Warnw(GameUser_GAME_ROOM_UUID_EMPTY,
						"GameUser", u.SecWebSocketKey,
						"Event", probability.EMSGID_eMercenaryOpenCall,
					)
				} else {
					data := &probability_proto.MercenaryOpenCall{}

					if err := proto.Unmarshal(msg, data); err != nil {
						logger.Service.Zap.Errorw(bullet.Bullet_MERCENARY_OPEN_PROTO_INVALID, "GameUser", u.SecWebSocketKey)
						errorcode.Service.Fatal(u.SecWebSocketKey, bullet.Bullet_MERCENARY_OPEN_PROTO_INVALID)
					}

					gameId := gamesetting.Service.GameId(u.SecWebSocketKey)
					subgameId := gamesetting.Service.SubgameId(u.SecWebSocketKey)

					openResult := rtp.Service.MercenaryOpen(u.SecWebSocketKey, gameId, subgameId, data.MercenaryType, data.Bullet)
					rtp.Service.Save(u.SecWebSocketKey, gameId, subgameId)
					bulletCollection, mercenaryType, mercenaryBullet := rtp.Service.GetMercenary(u.SecWebSocketKey, gameId, subgameId, data.MercenaryType)

					flux.Send(probability.EMSGID_eMercenaryOpenCall, u.SecWebSocketKey, lottery.Service.HashId(u.SecWebSocketKey),
						int32(u.seatId), bulletCollection, mercenaryType, mercenaryBullet, openResult,
					)
					flux.Send(idle.ActionIdleTimeoutReset, u.SecWebSocketKey, idle.Service.HashId(u.SecWebSocketKey),
						u.SecWebSocketKey,
						gamesetting.Service.Idle(u.SecWebSocketKey),
					)
				}

				u.mutex.Unlock()

			case jackpot.EMSGID_eCentInReask:
				// do nothing now
			}
		}
	}
}

func (u *GameUser) runWriteMessage() {
	flux.Register(u.SecWebSocketKey, u.in)
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		logger.Service.Zap.Infow("Game User Write Message Closed",
			"GameUser", u.SecWebSocketKey,
			"Chan", fmt.Sprintf("%p", u.in),
		)
		u.leaveRooms()
		u.conn.WriteMessage(websocket.CloseMessage, []byte{})
		ticker.Stop()
		u.conn.Close()
	}()

	for {
		select {
		case action, ok := <-u.in:
			if !ok {
				return
			}
			u.handleAction(action, u.in)

		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))

			data := strconv.FormatInt(time.Now().UnixNano(), 10)

			if err := u.conn.WriteMessage(websocket.PingMessage, []byte(data)); err != nil {
				logger.Service.Zap.Warnw("Game User Ping Failed",
					"GameUser", u.SecWebSocketKey,
					"Error", err,
				)
				u.destroy()
				errorcode.Service.Fatal(u.SecWebSocketKey, GameUser_PING_FAILED)
				return
			}
			//TODO User Channel Block (Have Problem)
			//default:
			//   u.destroy()
		}
	}
}

func (u *GameUser) handleAction(action *flux.Action, in chan *flux.Action) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.SecWebSocketKey != action.Key().To() {
		return
	}

	logger.Service.Zap.Infow(action.Key().Name(),
		"GameUser", u.SecWebSocketKey,
		"From", action.Key().From(),
		"In", fmt.Sprintf("%p", in),
	)

	switch action.Key().Name() {
	case fish.ActionChangeNextScene:
		u.nextScene = action.Payload()[0].(int)

	case auth.ActionLoginFailed:
		logger.Service.Zap.Errorw("GameUser Login Failed",
			"GameUser", u.SecWebSocketKey,
			"Chan", fmt.Sprintf("%p", u.in),
		)
		errorcode.Service.Fatal(u.SecWebSocketKey, auth.Auth_LOGIN_FAILED)

	case auth.ActionLoginSuccess:
		u.profile = action.Payload()[0].(*auth.MemberInfo)

		logger.Service.Zap.Infow("GameUser Login Success",
			"GameUser", u.SecWebSocketKey,
			"HostId", u.profile.HostId,
			"MemberId", u.profile.MemberId,
			"MemberName", u.profile.MemberName,
		)

		bet.Service.Set(u.SecWebSocketKey, &common_fish_proto.BetRateLine{
			BetLevelIndex: 0,
			BetIndex:      0,
			LineIndex:     0,
			RateIndex:     int32(gamesetting.Service.RateIndex(u.SecWebSocketKey)),
		})

		wallet.Service.Open(
			wallet.Builder().
				SetSecWebSocketKey(u.SecWebSocketKey).
				SetHostId(u.profile.HostId).
				SetHostExtId(u.HostExtId).
				SetMemberName(u.profile.MemberName).
				SetMemberId(u.profile.MemberId).
				SetToken(u.profile.RealToken).
				SetBalance(u.profile.Balance).
				SetCategory(u.profile.WalletType).
				SetPurposeMajor(wallet.WalletSharedFish). // balance
				SetPurposeMinor(wallet.WalletSharedFish). // cent, not used
				SetGuest(u.profile.GuestEnable).
				SetDB(u.HostExtId).
				SetHandler(
					accountingmanager.Service.New(
						accountingmanager.Builder().
							SetSecWebSocketKey(u.SecWebSocketKey).
							SetMemberId(u.profile.MemberId).
							SetHostId(u.profile.HostId).
							SetHostExtId(u.HostExtId).
							SetGuestEnable(u.profile.GuestEnable == guestOn).
							SetCategory(u.profile.WalletType).
							SetBillCurrency(u.profile.BillCurrency).
							SetDepositMultiple(u.profile.DepositMultiple).
							Build(),
					)).
				Build(),
		)

	case wallet.ActionWalletOpened:
		flux.Send(auth.EMSGID_eLoginRecall, u.SecWebSocketKey, auth.Service.Id, true)

	case probability.EMSGID_eResultRecall:
		u.sendAction(action)

	case probability.EMSGID_eMultiResultRecall:
		u.sendAction(action)

	case probability.EMSGID_eBroadcastResult:
		u.sendAction(action)

	case probability.EMSGID_eOptionRecall:
		u.sendAction(action)

	case probability.EMSGID_eBroadcastOption:
		u.sendAction(action)

	case bullet.EMSGID_eBroadcastShoot:
		msg := action.Payload()[0].([]byte)
		data := &bullet_proto.BroadcastShoot{}

		if err := proto.Unmarshal(msg, data); err != nil {
			logger.Service.Zap.Errorw(bullet.Bullet_BROADCAST_SHOOT_PROTO_INVALID,
				"GameUser", u.SecWebSocketKey,
				"GameRoomUuid", u.gameRoomUuid,
			)
			errorcode.Service.Fatal(u.SecWebSocketKey, bullet.Bullet_BROADCAST_SHOOT_PROTO_INVALID)
		}

		u.sendAction(action)

		if int(data.Shoot.SeatId) == u.seatId {
			mediator.Service.Done(u.SecWebSocketKey, data.Shoot.BulletUuid) // 3. 29b250c5
		}

	case bullet.EMSGID_eBroadcastBetChange:
		u.sendAction(action)

	case auth.EMSGID_eLoginRecall:
		u.sendAction(action)

		if u.profile != nil {
			//activity.Service.Check(
			//    u.SecWebSocketKey,
			//    u.HostExtId,
			//    u.profile.HostId,
			//    u.profile.MemberId,
			//    u.profile.GuestEnable == guestOn,
			//)
		} else {
			logger.Service.Zap.Warnw(GameUser_PROFILE_NIL,
				"GameUser", u.SecWebSocketKey,
				"HostExtId", u.HostExtId,
			)
		}

	case lobbysetting.EMSGID_eLobbyConfigRecall:
		data := action.Payload()[0].(*lobby_setting_proto.ConfigRecall)
		jackpotGroup := action.Payload()[1].(int32)
		u.JackpotGroupId = jackpotGroup

		if data.StatusCode == common_proto.Status_kSuccess {
			flux.Send(ActionJoinLobbyRoomCall, u.SecWebSocketKey, u.controllerId, data.RoomInfo[0].GameId)

			balance, _ := wallet.Service.BalanceCache(u.SecWebSocketKey)
			data.PlayerInfo.PlayerCent = balance

			value := vip.Service.CheckVip(
				false,
				u.SecWebSocketKey,
				u.HostExtId,
				u.profile.HostId,
				u.profile.MemberId,
				data.RoomInfo[0].GameId,
				u.profile.BillCurrency,
				1,
				0,
				0,
				0,
			)
			u.checkBlacklistValue(value, data.RoomInfo[0].GameId)

			// Get Jackpot Info
			//flux.Send(jackpot.ActionJackpotUpdate, u.SecWebSocketKey, jackpot.Service.Id, u.SecWebSocketKey, u.JackpotGroupId)
		}
		// Wait Refund Balance Update
		time.Sleep(1000 * time.Millisecond)
		dataBuffer, _ := proto.Marshal(data)
		u.sendByte(dataBuffer)

	case gamesetting.EMSGID_eGameConfigRecall:
		data := action.Payload()[0].(*game_setting_proto.ConfigRecall)

		if u.lobbyRoomUuid != "" {
			flux.Send(ActionLeaveLobbyRoomCall, u.SecWebSocketKey, u.controllerId, u.lobbyRoomUuid, false)
		}

		if data.StatusCode == common_proto.Status_kSuccess {
			gameId := gamesetting.Service.GameId(u.SecWebSocketKey)
			subgameId := gamesetting.Service.SubgameId(u.SecWebSocketKey)

			rtp.Service.Open(
				u.SecWebSocketKey,
				gameId,
				subgameId,
				u.profile.HostId,
				u.HostExtId,
				u.profile.MemberId,
			)

			// 加上傭兵的數據
			bulletCollection, mercenaryInfo := rtp.Service.MercenaryInfo(u.SecWebSocketKey, gameId, subgameId)
			data.MercenaryInfo.BulletCollection = bulletCollection
			data.MercenaryInfo.MercenaryData = mercenaryInfo

			flux.Send(ActionJoinGameRoomCall, u.SecWebSocketKey, u.controllerId, data, u.HostExtId)
		}

	case ActionJoinLobbyRoomRecall:
		u.nextScene = -1
		u.seatId = -1
		u.gameRoomUuid = ""
		u.lobbyRoomUuid = action.Payload()[0].(string)

		flux.Send(idle.ActionIdleTimeoutReset, u.SecWebSocketKey, idle.Service.HashId(u.SecWebSocketKey),
			u.SecWebSocketKey,
			gamesetting.Service.Idle(u.SecWebSocketKey),
		)

	case ActionLeaveLobbyRoomRecall:
		u.lobbyRoomUuid = ""

	case ActionJoinGameRoomRecall:
		u.lobbyRoomUuid = ""
		data := action.Payload()[0].(*game_setting_proto.ConfigRecall)
		u.gameRoomUuid = action.Payload()[1].(string)
		u.seatId = action.Payload()[2].(int)
		u.nextScene = action.Payload()[3].(int)

		lotteryKind := lottery.PLAYER
		bulletKind := bullet.PLAYER

		if u.profile.GuestEnable == guestOn {
			lotteryKind = lottery.GUEST
			bulletKind = bullet.GUEST
		}

		flux.Send(lottery.ActionLotteryStart, u.SecWebSocketKey, lottery.Service.Id,
			gamesetting.Service.GameId(u.SecWebSocketKey),
			gamesetting.Service.MathModuleId(u.SecWebSocketKey),
			lotteryKind,
			accountingmanager.Service.Get(u.SecWebSocketKey),
		)

		flux.Send(bullet.ActionBulletStart, u.SecWebSocketKey, bullet.Service.Id,
			gamesetting.Service.GameId(u.SecWebSocketKey),
			gamesetting.Service.MathModuleId(u.SecWebSocketKey),
			bulletKind,
		)

		flux.Send(idle.ActionIdleTimeoutReset, u.SecWebSocketKey, idle.Service.HashId(u.SecWebSocketKey),
			u.SecWebSocketKey,
			gamesetting.Service.Idle(u.SecWebSocketKey),
		)

		data.RoomUuid = u.gameRoomUuid
		data.NextScene = uint32(u.nextScene)
		data.PlayerInfo.SeatId = uint32(u.seatId)
		dataBuffer, _ := proto.Marshal(data)
		u.sendByte(dataBuffer)

	case ActionLeaveGameRoomRecall:
		u.gameRoomUuid = ""
		gameId := gamesetting.Service.GameId(u.SecWebSocketKey)
		subgameId := gamesetting.Service.SubgameId(u.SecWebSocketKey)
		rtp.Service.Save(u.SecWebSocketKey, gameId, subgameId)

	case gamesetting.EMSGID_eGameStripsRecall:
		data := action.Payload()[0].(*game_setting_proto.StripsRecall)
		data.RoomUuid = u.gameRoomUuid
		dataBuffer, _ := proto.Marshal(data)
		u.sendByte(dataBuffer)

		if data.StatusCode == common_proto.Status_kSuccess {
			flux.Send(auth.EMSGID_eBroadcastPlayerIn, u.SecWebSocketKey, u.gameRoomUuid, auth_proto.PLAYER_TYPE_GAME)
		}

	case auth.EMSGID_eBroadcastPlayerIn:
		u.sendAction(action)

	case auth.EMSGID_eBroadcastPlayerOut:
		u.sendAction(action)

	case fish.EMSGID_eBroadcastFishIn:
		u.sendAction(action)

	case fish.EMSGID_eBroadcastFishOut:
		u.sendAction(action)

	case fish.EMSGID_eBroadcastGroupIn:
		u.sendAction(action)

	case fish.EMSGID_eBroadcastRoundStart:
		u.sendAction(action)

	case fish.EMSGID_eBroadcastRoundEnd:
		u.sendAction(action)

	case fish.EMSGID_eAllFishRecall:
		u.sendAction(action)

	case fish.EMSGID_eBroadcastFishChange:
		u.sendAction(action)

	case broadcaster.EMSGID_eBroadcaster:
		u.sendAction(action)

	case heartbeat.ActionHeartBeatRecall:
		u.sendAction(action)

	case errorcode.ActionErrorCode:
		u.sendAction(action)
		flux.Send(auth.ActionUserDisconnect, u.SecWebSocketKey, Service.Id, u.SecWebSocketKey)

	case auth.ActionAuthUserInfoAsk:
		u.sendAction(action)

	case ActionBotStart:
		u.botUuid = action.Payload()[0].(string)

	case ActionBotStop:
		u.botUuid = ""

	case activity.EMSGID_eActivity:
		u.sendAction(action)

	case probability.EMSGID_eMercenaryOpenRecall:
		u.sendAction(action)

	case jackpot.EMSGID_eCentInAsk:
		jackpotWin := action.Payload()[1].(uint64)

		if am := accountingmanager.Service.Get(u.SecWebSocketKey); am != nil {
			am.JackpotHit(u.SecWebSocketKey, jackpotWin)
		}

		u.sendAction(action)

	case jackpot.EMSGID_eJackpotInfo:
		u.sendAction(action)

	case jackpot.EMSGID_eJackpotNotify:
		u.sendAction(action)

	case accountingmanager.ActionUpdateBlacklist:
		value := action.Payload()[0].(uint32)

		u.checkBlacklistValue(value, gamesetting.Service.GameId(u.SecWebSocketKey))
	}
}

func (u *GameUser) sendAction(action *flux.Action) {
	u.sendByte(action.Payload()[0].([]byte))
}

func (u *GameUser) sendByte(data []byte) {
	u.conn.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := u.conn.NextWriter(websocket.BinaryMessage)

	defer func() {
		if w != nil {
			w.Close()
		}
	}()

	if err == nil && w != nil {
		w.Write(data)
	}
}

func (u *GameUser) leaveRooms() {
	gameRoomUuid := u.gameRoomUuid
	lobbyRoomUuid := u.lobbyRoomUuid

	if gameRoomUuid != "" {
		flux.Send(auth.EMSGID_eBroadcastPlayerOut, u.SecWebSocketKey, auth.Service.Id,
			gameRoomUuid,
			uint32(u.seatId),
			uint64(0),
		)
		flux.Send(ActionLeaveGameRoomCall, u.SecWebSocketKey, u.controllerId, gameRoomUuid, true)
	} else {
		flux.Send(ActionGameUserDelete, u.SecWebSocketKey, Service.Id, u.SecWebSocketKey)
	}

	if lobbyRoomUuid != "" {
		flux.Send(ActionLeaveLobbyRoomCall, u.SecWebSocketKey, u.controllerId, lobbyRoomUuid, true)
	} else {
		flux.Send(ActionGameUserDelete, u.SecWebSocketKey, Service.Id, u.SecWebSocketKey)
	}
}

func (u *GameUser) checkBlacklistValue(blacklistValue uint32, gameId string) {
	if blacklistValue == 0 {
		u.blacklistValue = uint32(blacklist.Service.Check(
			u.HostExtId,
			u.profile.HostId,
			u.profile.MemberId,
			"fish",
			gameId,
		))
	} else {
		u.blacklistValue = blacklistValue
	}
}

func (u *GameUser) cleanData() {
	mediator.Service.Delete(u.SecWebSocketKey)

	u.maintainChan <- true
	u.banPlayerChan <- true
	if u.profile != nil {
		var rounds uint64 = 0

		if am := accountingmanager.Service.Get(u.SecWebSocketKey); am != nil {
			rounds = am.Rounds()
		}

		session.Service.End(
			u.SecWebSocketKey,
			u.profile.HostId,
			u.HostExtId,
			u.profile.MemberId,
			rounds,
		)
	} else {
		logger.Service.Zap.Warnw(GameUser_PROFILE_NIL,
			"GameUser", u.SecWebSocketKey,
			"HostExtId", u.HostExtId,
		)
	}

	sessionmanager.Service.End(u.SecWebSocketKey, u.HostExtId)

	if am := accountingmanager.Service.Get(u.SecWebSocketKey); am != nil {
		done := make(chan bool)
		for {
			time.Sleep(1 * time.Second)

			logger.Service.Zap.Infow("Accounting Manager Disconnect",
				"GameUser", u.SecWebSocketKey)

			if !u.isBulletsProcessedComplete() {
				continue
			}

			am.Disconnect(done)

			if <-done {
				break
			}
		}
		close(done)
	}

	// for lag reason in some game user, we clean data later
	time.Sleep(10 * time.Second)

	gameId := gamesetting.Service.GameId(u.SecWebSocketKey)
	subgameId := gamesetting.Service.SubgameId(u.SecWebSocketKey)

	rtp.Service.Close(u.SecWebSocketKey, gameId, subgameId)
	flux.Send(refund.ActionRefundDelete, u.SecWebSocketKey, refund.Service.Id, u.SecWebSocketKey)
	flux.Send(bullet.ActionBulletDelete, u.SecWebSocketKey, bullet.Service.Id, refund.Service.Id)
	flux.Send(lottery.ActionLotteryDelete, u.SecWebSocketKey, lottery.Service.Id, "")
	flux.Send(auth.ActionUserDisconnect, u.SecWebSocketKey, auth.Service.Id, u.SecWebSocketKey)
	lobbysetting.Service.Delete(u.SecWebSocketKey)
	gamesetting.Service.Delete(u.SecWebSocketKey)
	wallet.Service.Close(u.SecWebSocketKey)
	flux.Send(idle.ActionIdleTimeoutStop, u.SecWebSocketKey, idle.Service.Id, u.SecWebSocketKey)
	mediator.Service.Delete(u.SecWebSocketKey)
	accountingmanager.Service.Delete(u.SecWebSocketKey)
	bet.Service.Delete(u.SecWebSocketKey)
	activity.Service.Delete(u.SecWebSocketKey)
	//flux.Send(jackpot.ActionJackpotDeleteUser, u.SecWebSocketKey, jackpot.Service.Id, u.SecWebSocketKey)
}

func (u *GameUser) destroy() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if !u.isDestroy {
		logger.Service.Zap.Infow("Game User Unregister",
			"GameUser", u.SecWebSocketKey,
			"SeatId", u.seatId,
			"Chan", fmt.Sprintf("%p", u.in),
		)
		u.isDestroy = true
		u.status.Alive = false
		flux.UnRegister(u.SecWebSocketKey, u.in)
	}
}

func checkSessionMap(hostExtId, host_id, member_id, game_id string) bool {
	if db, err := mysql.Repository.GameDB(hostExtId); err == nil {

		rows, err := db.Table("session_map").
			Where("host_id = ? AND member_id = ? AND game_id = ?", host_id, member_id, game_id).
			Rows()

		if err != nil {
			return false
		} else {
			defer rows.Close()
			result := false

			for rows.Next() {
				result = true
				break
			}
			return result
		}
	}

	return false
}

func (u *GameUser) isBulletsProcessedComplete() bool {
	isDone := make(chan bool)
	flux.Send(lottery.ActionUnprocessedBullet, u.SecWebSocketKey, lottery.Service.HashId(u.SecWebSocketKey), u.SecWebSocketKey, isDone)
	return <-isDone
}
