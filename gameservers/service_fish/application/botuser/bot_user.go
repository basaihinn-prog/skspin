package botuser

import (
	"fmt"
	"math/rand"
	"serve/service_fish/application/bonuslottery"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/application/lottery"
	"serve/service_fish/application/refund"
	"serve/service_fish/domain/auth"
	auth_proto "serve/service_fish/domain/auth/proto"
	"serve/service_fish/domain/bet"
	"serve/service_fish/domain/bullet"
	bullet_proto "serve/service_fish/domain/bullet/proto"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	probability_proto "serve/service_fish/domain/probability/proto"
	common_proto "serve/service_fish/models/proto"
	"sync"
	"time"

	"serve/fish_comm/broadcaster"
	"serve/fish_comm/common"
	Common_Function "serve/fish_comm/common-function"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/idle"
	"serve/fish_comm/logger"
	"serve/fish_comm/mediator"
	"serve/fish_comm/wallet"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
)

const (
	ActionBotUserJoinGameRoomCall   = "ActionBotUserJoinGameRoomCall"
	ActionBotUserJoinGameRoomRecall = "ActionBotUserJoinGameRoomRecall"
	ActionBotUserLeaveGameRoom      = "ActionBotUserLeaveGameRoom"
	ActionBotUserReady              = "ActionBotUserReady"
)

type botUser struct {
	uuid            string
	controllerId    string
	secWebSocketKey string
	gameRoomUuid    string
	memberId        string
	memberName      string
	seatId          int
	bet             uint64
	maxBet          uint64
	rate            uint64
	rateIndex       uint32
	status          *auth.MemberStatus
	in              chan *flux.Action
	isDestroy       bool
	rwMutex         *sync.RWMutex
}

func newBotUser(controllerId, secWebSocketKey, gameRoomUuid string) *botUser {
	b := &botUser{
		uuid:            uuid.New().String(),
		controllerId:    controllerId,
		secWebSocketKey: secWebSocketKey,
		gameRoomUuid:    gameRoomUuid,
		seatId:          -1,
		status:          &auth.MemberStatus{Alive: true},
		in:              make(chan *flux.Action, common.Service.ChanSize),
		isDestroy:       false,
		rwMutex:         &sync.RWMutex{},
	}

	if len(secWebSocketKey) == len(b.uuid) {
		return nil
	}

	if _, err := wallet.Service.BalanceCache(secWebSocketKey); err != nil {
		logger.Service.Zap.Warnw("Check Players Alive",
			"GameUser", secWebSocketKey,
			"Bot", b.uuid,
			"SeatId", b.seatId,
			"GameRoomUuid", b.gameRoomUuid,
			"Chan", fmt.Sprintf("%p", b.in),
		)
		return nil
	}

	gamesetting.Service.Clone(secWebSocketKey, b.uuid)

	b.rate = gamesetting.Service.Rate(b.uuid)
	b.maxBet = gamesetting.Service.MaxBet(b.uuid)
	b.rateIndex = gamesetting.Service.RateIndex(b.uuid)

	b.memberId = b.rngName()
	b.memberName = b.memberId

	go b.run()
	return b
}

func (b *botUser) run() {
	flux.Register(b.uuid, b.in)

	t := time.NewTimer(time.Duration(5) * time.Second)

	defer t.Stop()

	for {
		select {
		case action, ok := <-b.in:
			if !ok {
				t.Stop()
				return
			}
			b.handleAction(action, b.in)

		case <-t.C:
			flux.Send(ActionBotUserJoinGameRoomCall, b.uuid, b.controllerId, b.gameRoomUuid)
			t.Stop()
		}
	}
}

func (b *botUser) handleAction(action *flux.Action, in chan *flux.Action) {
	if b.uuid != action.Key().To() {
		return
	}

	logger.Service.Zap.Infow(action.Key().Name(),
		"GameUser", b.uuid,
		"From", action.Key().From(),
		"In", fmt.Sprintf("%p", in),
	)

	switch action.Key().Name() {
	case flux.ActionFluxRegisterDone:
		bet.Service.Set(b.uuid, &common_proto.BetRateLine{
			BetLevelIndex: 0,
			BetIndex:      0,
			LineIndex:     0,
			RateIndex:     int32(b.rateIndex),
		})

		wallet.Service.Open(
			wallet.Builder().
				SetSecWebSocketKey(b.uuid).
				SetHostId("BOT").
				SetHostExtId("BOT").
				SetMemberName(b.memberName).
				SetMemberId(b.memberId).
				SetBalance(b.rngCent()).
				SetCategory(wallet.CategoryGuest).
				SetGuest(wallet.GuestOn).
				SetHandler(nil).
				Build(),
		)

		flux.Send(lottery.ActionLotteryStart, b.uuid, lottery.Service.Id,
			gamesetting.Service.GameId(b.uuid),
			gamesetting.Service.MathModuleId(b.uuid),
			lottery.BOT,
			nil,
		)

		flux.Send(bullet.ActionBulletStart, b.uuid, bullet.Service.Id,
			gamesetting.Service.GameId(b.uuid),
			gamesetting.Service.MathModuleId(b.uuid),
			bullet.BOT,
		)

		flux.Send(idle.ActionIdleTimeoutStart, b.uuid, idle.Service.Id, b.uuid, uint32(60))

	case wallet.ActionWalletOpened:
		// do nothing

	case bullet.EMSGID_eBetChange:
		data := action.Payload()[0].(*bullet_proto.BetChange)

		if int(data.SeatId) != b.seatId {
			return
		}

		bet.Service.Set(b.uuid, data.BetRateLine)
		flux.Send(bullet.EMSGID_eBetChange, b.uuid, bullet.Service.HashId(b.uuid), data)
		flux.Send(idle.ActionIdleTimeoutReset, b.uuid, idle.Service.HashId(b.uuid), gamesetting.Service.Idle(b.uuid))

	case bullet.EMSGID_eShoot:
		data := action.Payload()[0].(*bullet_proto.Shoot)
		bet := action.Payload()[1].(uint64)
		rate := action.Payload()[2].(uint64)

		if int(data.SeatId) != b.seatId {
			return
		}

		b.rwMutex.Lock()
		b.bet = bet
		b.rwMutex.Unlock()

		balance, _ := wallet.Service.Balance(b.uuid)

		mediator.Service.New(b.uuid, data.BulletUuid) // 1. 78adf904
		mediator.Service.Add(b.uuid, data.BulletUuid) // 2. 78adf904
		flux.Send(bullet.EMSGID_eShoot, b.uuid, bullet.Service.HashId(b.uuid), data, bet, rate, uint64(1), balance)
		flux.Send(idle.ActionIdleTimeoutReset, b.uuid, idle.Service.HashId(b.uuid), gamesetting.Service.Idle(b.uuid))

	case probability.EMSGID_eResultCall:
		data := action.Payload()[0].(*probability_proto.ResultCall)

		if int(data.SeatId) != b.seatId {
			return
		}

		flux.Send(idle.ActionIdleTimeoutReset, b.uuid, idle.Service.HashId(b.uuid), gamesetting.Service.Idle(b.uuid))

		balance, _ := wallet.Service.Balance(b.uuid)

		// no money
		switch {
		case balance <= b.bet*b.rate:
			fallthrough
		case balance == 0:
			flux.Send(ActionBotUserDelete, b.uuid, Service.Id, b.secWebSocketKey)
		default:
			flux.Send(probability.EMSGID_eResultCall, b.uuid, lottery.Service.HashId(b.uuid), data, uint32(0), uint64(20), b.status)
		}

	case probability.EMSGID_eOptionCall:
		data := action.Payload()[0].(*probability_proto.OptionCall)

		if int(data.SeatId) != b.seatId {
			return
		}
		flux.Send(probability.EMSGID_eOptionCall, b.uuid, bonuslottery.Service.Id, data, uint64(0))
		flux.Send(idle.ActionIdleTimeoutReset, b.uuid, idle.Service.HashId(b.uuid), gamesetting.Service.Idle(b.uuid))

	case bullet.EMSGID_eBroadcastShoot:
		msg := action.Payload()[0].([]byte)
		data := &bullet_proto.BroadcastShoot{}

		if err := proto.Unmarshal(msg, data); err != nil {
			logger.Service.Zap.Errorw(bullet.Bullet_BROADCAST_SHOOT_PROTO_INVALID,
				"GameUser", b.uuid,
				"GameRoomUuid", b.gameRoomUuid,
			)
			errorcode.Service.Fatal(b.uuid, bullet.Bullet_BROADCAST_SHOOT_PROTO_INVALID)
		}

		if int(data.Shoot.SeatId) == b.seatId {
			mediator.Service.Done(b.uuid, data.Shoot.BulletUuid) // 3. 78adf904
		}

	case probability.EMSGID_eBroadcastResult:
		fallthrough
	case bullet.EMSGID_eBroadcastBetChange:
		fallthrough
	case probability.EMSGID_eBroadcastOption:
		fallthrough
	case fish.EMSGID_eBroadcastFishIn:
		fallthrough
	case fish.EMSGID_eBroadcastFishOut:
		fallthrough
	case fish.EMSGID_eBroadcastGroupIn:
		fallthrough
	case fish.EMSGID_eBroadcastRoundStart:
		fallthrough
	case fish.EMSGID_eBroadcastRoundEnd:
		fallthrough
	case fish.EMSGID_eBroadcastFishChange:
		fallthrough
	case auth.EMSGID_eBroadcastPlayerIn:
		fallthrough
	case auth.EMSGID_eBroadcastPlayerOut:
		fallthrough
	case broadcaster.EMSGID_eBroadcaster:
		// do nothing

	case ActionBotUserJoinGameRoomRecall:
		b.rwMutex.Lock()
		b.seatId = action.Payload()[0].(int)
		b.rwMutex.Unlock()

		flux.Send(auth.EMSGID_eBroadcastPlayerIn, b.uuid, b.gameRoomUuid, auth_proto.PLAYER_TYPE_BOT)
		flux.Send(ActionBotUserReady, b.uuid, b.controllerId, b.secWebSocketKey)

	case errorcode.ActionErrorCode:
		flux.Send(ActionBotUserDelete, b.uuid, Service.Id, b.secWebSocketKey)

	case probability.EMSGID_eResultRecall:
		balance, _ := wallet.Service.Balance(b.uuid)

		// no money
		switch {
		case balance <= b.bet*b.rate:
			fallthrough
		case balance == 0:
			flux.Send(ActionBotUserDelete, b.uuid, Service.Id, b.secWebSocketKey)
		}

		fallthrough

	default:
		flux.Send(action.Key().Name(), b.uuid, b.secWebSocketKey, action.Payload()...)
	}
}

func (b *botUser) rngCent() (cent uint64) {
	rand.Seed(time.Now().UnixNano())
	min := int(b.maxBet) * 20
	max := int(b.maxBet) * 100
	return uint64(rand.Intn(max-min+1)+min) * b.rate
}

func (b *botUser) rngName() (name string) {
	return Common_Function.Service.RngBotName()
}

func (b *botUser) destroy() {
	mediator.Service.Delete(b.uuid)

	logger.Service.Zap.Infow("Bot Unregister",
		"GameUser", b.secWebSocketKey,
		"Bot", b.uuid,
		"SeatId", b.seatId,
		"GameRoomUuid", b.gameRoomUuid,
		"Chan", fmt.Sprintf("%p", b.in),
	)

	b.rwMutex.Lock()
	if !b.isDestroy {
		b.isDestroy = true
		b.status.Alive = false
		flux.UnRegister(b.uuid, b.in)
	}
	b.rwMutex.Unlock()

	if b.seatId > -1 {
		flux.Send(auth.EMSGID_eBroadcastPlayerOut, b.uuid, b.gameRoomUuid, auth.Service.BroadcastPlayerOut(
			b.gameRoomUuid,
			b.memberId,
			b.memberName,
			uint32(b.seatId),
			uint64(0),
		))
	}

	flux.Send(ActionBotUserLeaveGameRoom, b.uuid, b.controllerId, b.secWebSocketKey, b.gameRoomUuid)

	// this delay is necessary which use to waiting bullet/lottery/accounting started
	time.Sleep(5 * time.Second)

	flux.Send(idle.ActionIdleTimeoutStop, b.uuid, idle.Service.Id, b.uuid)
	flux.Send(bullet.ActionBulletDelete, b.uuid, bullet.Service.Id, refund.Service.Id)
	flux.Send(lottery.ActionLotteryDelete, b.uuid, lottery.Service.Id, "")
	wallet.Service.Close(b.uuid)
	gamesetting.Service.Delete(b.uuid)
	bet.Service.Delete(b.uuid)
}
