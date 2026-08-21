package bullet

import (
	"fmt"
	"os"
	"serve/service_fish/application/gamesetting"
	bullet_proto "serve/service_fish/domain/bullet/proto"
	"serve/service_fish/models"
	common_proto "serve/service_fish/models/proto"
	"strconv"
	"sync"

	"serve/fish_comm/mediator"

	"serve/fish_comm/common"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"

	"github.com/gogo/protobuf/proto"
)

const (
	EMSGID_eShoot                  = "EMSGID_eShoot"
	EMSGID_eBroadcastShoot         = "EMSGID_eBroadcastShoot"
	EMSGID_eBetChange              = "EMSGID_eBetChange"
	EMSGID_eBroadcastBetChange     = "EMSGID_eBroadcastBetChange"
	ActionBulletExistCheck         = "ActionBulletExistCheck"
	ActionBulletExist              = "ActionBulletExist"
	ActionBulletNotExist           = "ActionBulletNotExist"
	ActionBulletMultiExistCheck    = "ActionBulletMultiExistCheck"
	ActionBulletMultiExist         = "ActionBulletMultiExist"
	ActionBulletMultiNotExist      = "ActionBulletMultiNotExist"
	ActionBulletHitCall            = "ActionBulletHitCall"
	ActionBulletHitRecall          = "ActionBulletHitRecall"
	ActionBulletMultiCall          = "ActionBulletMultiCall"
	ActionBulletMultiRecall        = "ActionBulletMultiRecall"
	ActionBulletBonusCall          = "ActionBulletBonusCall"
	ActionBulletBonusRecall        = "ActionBulletBonusRecall"
	ActionBulletMultiBonusCall     = "ActionBulletMultiBonusCall"
	ActionBulletRefundCall         = "ActionBulletRefundCall"
	ActionBulletRefundRecall       = "ActionBulletRefundRecall"
	ActionBulletRefundDone         = "ActionBulletRefundDone"
	ActionBulletBonusExist         = "ActionBulletBonusExist"
	ActionBulletBonusNotExist      = "ActionBulletBonusNotExist"
	ActionBulletMultiBonusExist    = "ActionBulletMultiBonusExist"
	ActionBulletMultiBonusNotExist = "ActionBulletMultiBonusNotExist"
	actionBulletDestroyCall        = "actionBulletDestroyCall"
	actionBulletDestroyRecall      = "actionBulletDestroyRecall"
	ActionBulletDestroyDone        = "ActionBulletDestroyDone"
	ActionBonusBulletShooted       = "ActionBonusBulletShooted"
)

type bulletPool struct {
	secWebSocketKey string
	gameId          string
	mathModuleId    string // not used now
	kind            string
	poolSize        string
	in              chan *flux.Action
	rwMutex         *sync.RWMutex
}

func newBulletPool(secWebSocketKey, gameId, mathModuleId, kind string) *bulletPool {
	b := &bulletPool{
		secWebSocketKey: secWebSocketKey,
		gameId:          gameId,
		mathModuleId:    mathModuleId,
		kind:            kind,
		poolSize:        os.Getenv("POOL_SIZE"),
		in:              make(chan *flux.Action, common.Service.ChanSize),
		rwMutex:         &sync.RWMutex{},
	}
	b.run()
	logger.Service.Zap.Infow("bullet pool started run",
		"Chan", fmt.Sprintf("%p", b.in),
		"GameId", gameId,
		"MathModuleId", mathModuleId,
	)
	return b
}

func (b *bulletPool) run() {
	flux.Register(b.hashId(), b.in)
	bullets := make(map[string]*Bullet)
	bonusBullets := make(map[string]*Bullet)

	poolSize := 10

	if b.poolSize != "" {
		if poolSize, _ = strconv.Atoi(b.poolSize); poolSize == 0 {
			poolSize = 10
		}
	}

	for i := 0; i < poolSize; i++ {
		go func() {
			for action := range b.in {
				b.handleAction(action, bullets, bonusBullets)
			}
		}()
	}
}

func (b *bulletPool) handleAction(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	if action.Key().To() != b.hashId() {
		return
	}

	switch action.Key().Name() {
	case EMSGID_eBetChange:
		data := action.Payload()[0].(*bullet_proto.BetChange)

		flux.Send(EMSGID_eBroadcastBetChange, b.hashId(), data.RoomUuid, b.eBroadcastBetChange(data))

		mediator.Service.Done(action.Key().From(), EMSGID_eBetChange) // 3. ac7a2c41, d3f63a43

	case EMSGID_eShoot:
		secWebSocketKey := action.Key().From()
		shootProto := action.Payload()[0].(*bullet_proto.Shoot)

		if dataByte := b.eBroadcastShoot(action); dataByte != nil {
			mediator.Service.Add(secWebSocketKey, shootProto.BulletUuid) // 2. 29b250c5, 78adf904
			flux.Send(EMSGID_eBroadcastShoot, b.hashId(), shootProto.RoomUuid, dataByte)
		} else {
			logger.Service.Zap.Errorw(Bullet_BROADCAST_SHOOT_PROTO_INVALID,
				"GameUser", secWebSocketKey,
				"GameRoomUuid", shootProto.RoomUuid,
				"BulletType", shootProto.BulletTypeId,
				"BulletUuid", shootProto.BulletUuid,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_BROADCAST_SHOOT_PROTO_INVALID)
		}

		switch b.gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_Shoot(action, bullets, bonusBullets)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_Shoot(action, bullets, bonusBullets)

		case models.PSF_ON_00003:
			psf_on_00003_Shoot(action, bullets, bonusBullets)

		case models.PSF_ON_00004:
			psf_on_00004_Shoot(action, bullets, bonusBullets)

		case models.PSF_ON_00005:
			psf_on_00005_Shoot(action, bullets, bonusBullets)

		case models.PSF_ON_00006:
			psf_on_00006_Shoot(action, bullets, bonusBullets)

		case models.PSF_ON_00007:
			psf_on_00007_Shoot(action, bullets, bonusBullets)

		case models.RKF_H5_00001:
			rkf_h5_00001_Shoot(action, bullets, bonusBullets)

		default:
			logger.Service.Zap.Errorw(Bullet_GAME_ID_INVALID,
				"Service", Service.Id,
				"Chan", fmt.Sprintf("%p", b.in),
				"GameUser", secWebSocketKey,
				"GameRoomUuid", shootProto.RoomUuid,
				"GameId", b.gameId,
				"MathModuleId", b.mathModuleId,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_GAME_ID_INVALID)
		}

		mediator.Service.Done(secWebSocketKey, shootProto.BulletUuid) // 3. 29b250c5, 78adf904

	case ActionBulletExistCheck:
		hitBullet := action.Payload()[0].(*Bullet)
		mercenaryType := action.Payload()[1].(int32)
		blacklistValue := action.Payload()[2].(uint32)
		depositMultiple := action.Payload()[3].(uint64)
		isFishExist := action.Payload()[4].(bool)

		gameId := gamesetting.Service.GameId(hitBullet.SecWebSocketKey)

		if isFishExist {
			// regular bullet
			if hitBullet.BonusUuid == "" {
				if v, ok := bullets[hitBullet.Uuid]; ok {
					hitBullet.Rate = v.Rate
					hitBullet.Bet = v.Bet
					hitBullet.BetIndex = v.BetIndex
					hitBullet.BetLevelIndex = v.BetLevelIndex
					hitBullet.RateIndex = v.RateIndex
					hitBullet.LineIndex = v.LineIndex
					v.hitLock = true
					flux.Send(ActionBulletExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType, blacklistValue, depositMultiple)
				} else {
					flux.Send(ActionBulletNotExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType)
				}
			} else {
				// bonus bullet
				if v, ok := bonusBullets[hitBullet.BonusUuid]; ok {

					if b.checkDrill(gameId, hitBullet.TypeId) {
						if v.ShootBullets == v.Bullets {

							flux.Send(ActionBulletBonusNotExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType)

							logger.Service.Zap.Warnw("Last Bullet Double Hit",
								"GameUser", hitBullet.SecWebSocketKey,
								"GameRoomUuid", hitBullet.RoomUuid,
								"BulletType", v.TypeId,
								"BonusUuid", v.BonusUuid,
								"Bullets", v.Bullets,
								"ShootBullets", v.ShootBullets,
								"UsedBullets", v.UsedBullets,
							)
							return
						}
						v.ShootBullets += 1
					}

					hitBullet.Rate = v.Rate
					hitBullet.Bet = v.Bet
					hitBullet.BetIndex = v.BetIndex
					hitBullet.BetLevelIndex = v.BetLevelIndex
					hitBullet.RateIndex = v.RateIndex
					hitBullet.LineIndex = v.LineIndex
					hitBullet.RtpId = v.RtpId
					flux.Send(ActionBulletBonusExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType, blacklistValue, depositMultiple)
				} else {
					flux.Send(ActionBulletBonusNotExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType)
				}
			}
		} else {
			// bonus bullet
			if hitBullet.BonusUuid != "" {
				if v, ok := bonusBullets[hitBullet.BonusUuid]; ok && !b.checkDrill(gameId, hitBullet.TypeId) {
					hitBullet.Rate = v.Rate
					hitBullet.Bet = v.Bet
					hitBullet.BetIndex = v.BetIndex
					hitBullet.BetLevelIndex = v.BetLevelIndex
					hitBullet.RateIndex = v.RateIndex
					hitBullet.LineIndex = v.LineIndex
					hitBullet.RtpId = v.RtpId
					flux.Send(ActionBonusBulletShooted, b.hashId(), action.Key().From(), hitBullet)
				}
			}
		}

	case ActionBulletMultiExistCheck:
		hitBullet := action.Payload()[0].(*Bullet)
		checkFishiesExist := action.Payload()[1].([]bool)
		mercenaryType := action.Payload()[2].(int32)
		blacklistValue := action.Payload()[3].(uint32)
		depositMultiple := action.Payload()[4].(uint64)

		// Regular Bullet
		if hitBullet.BonusUuid == "" {
			if v, ok := bullets[hitBullet.Uuid]; ok {
				hitBullet.Rate = v.Rate
				hitBullet.Bet = v.Bet
				hitBullet.BetIndex = v.BetIndex
				hitBullet.BetLevelIndex = v.BetLevelIndex
				hitBullet.RateIndex = v.RateIndex
				hitBullet.LineIndex = v.LineIndex
				v.hitLock = true
				flux.Send(ActionBulletMultiExist, b.hashId(), action.Key().From(), hitBullet, checkFishiesExist, mercenaryType, blacklistValue, depositMultiple)
			} else {
				flux.Send(ActionBulletMultiNotExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType)
			}
		} else {
			// Bonus Bullet
			if v, ok := bonusBullets[hitBullet.BonusUuid]; ok {
				if hitBullet.TypeId == 200 {
					if v.ShootBullets == v.Bullets {
						flux.Send(ActionBulletMultiBonusNotExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType)

						logger.Service.Zap.Warnw("Last Bullet Double Hit",
							"GameUser", hitBullet.SecWebSocketKey,
							"GameRoomUuid", hitBullet.RoomUuid,
							"BulletType", v.TypeId,
							"BonusUuid", v.BonusUuid,
							"Bullets", v.Bullets,
							"ShootBullets", v.ShootBullets,
							"UsedBullets", v.UsedBullets,
						)
						return
					}
					v.ShootBullets += 1
				}

				hitBullet.Rate = v.Rate
				hitBullet.Bet = v.Bet
				hitBullet.BetIndex = v.BetIndex
				hitBullet.BetLevelIndex = v.BetLevelIndex
				hitBullet.RateIndex = v.RateIndex
				hitBullet.LineIndex = v.LineIndex
				hitBullet.RtpId = v.RtpId
				flux.Send(ActionBulletMultiBonusExist, b.hashId(), action.Key().From(), hitBullet, checkFishiesExist, mercenaryType, blacklistValue)
			} else {
				flux.Send(ActionBulletMultiBonusNotExist, b.hashId(), action.Key().From(), hitBullet, mercenaryType)
			}
		}

	case ActionBulletBonusCall:
		bonusBullet := action.Payload()[0].(*Bullet)

		switch b.gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_Bonus(action, bonusBullets)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_Bonus(action, bonusBullets)

		case models.PSF_ON_00003:
			psf_on_00003_Bonus(action, bonusBullets)

		case models.PSF_ON_00004:
			psf_on_00004_Bonus(action, bonusBullets)

		case models.PSF_ON_00005:
			psf_on_00005_Bonus(action, bonusBullets)

		case models.PSF_ON_00006:
			psf_on_00006_Bonus(action, bonusBullets)

		case models.PSF_ON_00007:
			psf_on_00007_Bonus(action, bonusBullets)

		case models.RKF_H5_00001:
			rkf_h5_00001_Bonus(action, bonusBullets)

		default:
			logger.Service.Zap.Errorw(Bullet_GAME_ID_INVALID,
				"Service", Service.Id,
				"Chan", fmt.Sprintf("%p", b.in),
				"GameUser", bonusBullet.SecWebSocketKey,
				"GameRoomUuid", bonusBullet.RoomUuid,
				"GameId", b.gameId,
				"MathModuleId", b.mathModuleId,
			)
			errorcode.Service.Fatal(bonusBullet.SecWebSocketKey, Bullet_GAME_ID_INVALID)
		}

	case ActionBulletMultiBonusCall:
		bonusBulletList := action.Payload()[0].([]*Bullet)

		switch b.gameId {
		case models.PSF_ON_00003:
			psf_on_00003_multi_Bonus(action, bonusBullets)

		default:
			logger.Service.Zap.Errorw(Bullet_GAME_ID_INVALID,
				"Service", Service.Id,
				"Chan", fmt.Sprintf("%p", b.in),
				"GameUser", bonusBulletList[0].SecWebSocketKey,
				"GameRoomUuid", bonusBulletList[0].RoomUuid,
				"GameId", b.gameId,
				"MathModuleId", b.mathModuleId,
			)
			errorcode.Service.Fatal(bonusBulletList[0].SecWebSocketKey, Bullet_GAME_ID_INVALID)
		}

	case ActionBulletHitCall:
		switch b.gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_Hit(action, bullets, bonusBullets)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_Hit(action, bullets, bonusBullets)

		case models.PSF_ON_00003:
			psf_on_00003_Hit(action, bullets, bonusBullets)

		case models.PSF_ON_00004:
			psf_on_00004_Hit(action, bullets, bonusBullets)

		case models.PSF_ON_00005:
			psf_on_00005_Hit(action, bullets, bonusBullets)

		case models.PSF_ON_00006:
			psf_on_00006_Hit(action, bullets, bonusBullets)

		case models.PSF_ON_00007:
			psf_on_00007_Hit(action, bullets, bonusBullets)

		case models.RKF_H5_00001:
			rkf_h5_00001_Hit(action, bullets, bonusBullets)

		default:
			hitBullet := action.Payload()[1].(*Bullet)

			logger.Service.Zap.Errorw(Bullet_GAME_ID_INVALID,
				"Service", Service.Id,
				"Chan", fmt.Sprintf("%p", b.in),
				"GameUser", hitBullet.SecWebSocketKey,
				"GameRoomUuid", hitBullet.RoomUuid,
				"GameId", b.gameId,
				"MathModuleId", b.mathModuleId,
			)
			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_GAME_ID_INVALID)
		}

	case ActionBulletMultiCall:
		switch b.gameId {
		case models.PSF_ON_00003:
			psf_on_00003_Multi_Hit(action, bullets, bonusBullets)

		default:
			hitBullet := action.Payload()[1].(*Bullet)

			logger.Service.Zap.Errorw(Bullet_GAME_ID_INVALID,
				"Service", Service.Id,
				"Chan", fmt.Sprintf("%p", b.in),
				"GameUser", hitBullet.SecWebSocketKey,
				"GameRoomUuid", hitBullet.RoomUuid,
				"GameId", b.gameId,
				"MathModuleId", b.mathModuleId,
			)
			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_GAME_ID_INVALID)
		}

	case ActionBulletRefundCall:
		secWebSocketKey := action.Payload()[0].(string)
		accountingSn := action.Payload()[1].(uint64)

		for k, v := range bullets {
			if !v.hitLock {
				delete(bullets, k)
			}
		}

		for k, v := range bonusBullets {
			logger.Service.Zap.Infow("Refund Bonus Bullet",
				"GameUser", secWebSocketKey,
				"GameRoomUuid", v.RoomUuid,
				"BonusBulletUuid", v.BonusUuid,
				"Bullets", v.Bullets,
				"ShootBullets", v.ShootBullets,
				"UsedBullets", v.UsedBullets,
				"RemainBullets", v.Bullets-v.UsedBullets,
				"Bet", v.Bet,
				"Rate", v.Rate,
			)

			delete(bonusBullets, k)

			if b.kind == PLAYER {
				flux.Send(ActionBulletRefundRecall, b.hashId(), action.Key().From(), secWebSocketKey, v, accountingSn)
			}
		}
		flux.Send(ActionBulletRefundDone, b.hashId(), action.Key().From(), secWebSocketKey)

	case actionBulletDestroyCall:
		refundServiceId := action.Payload()[0].(string)

		for k := range bullets {
			delete(bullets, k)
		}

		for k, v := range bonusBullets {
			if b.kind != BOT {
				logger.Service.Zap.Infow("Refund Bonus Bullet",
					"GameUser", v.SecWebSocketKey,
					"GameRoomUuid", v.RoomUuid,
					"Bullets", v.Bullets,
					"ShootBullets", v.ShootBullets,
					"UsedBullets", v.UsedBullets,
					"RestBullets", v.Bullets-v.UsedBullets,
					"Bet", v.Bet,
					"Rate", v.Rate,
				)
			}
			delete(bonusBullets, k)

			if b.kind == PLAYER {
				flux.Send(ActionBulletRefundRecall, b.hashId(), refundServiceId, v.SecWebSocketKey, v, 0)
			}
		}
		flux.Send(actionBulletDestroyRecall, b.hashId(), action.Key().From(), refundServiceId, b.secWebSocketKey)
	}
}

func (b *bulletPool) eBroadcastShoot(action *flux.Action) []byte {
	shootProto := action.Payload()[0].(*bullet_proto.Shoot)
	balance := action.Payload()[3].(uint64)

	data := &bullet_proto.BroadcastShoot{
		Msgid:      common_proto.EMSGID_eBroadcastShoot,
		Shoot:      shootProto,
		PlayerCent: balance,
	}
	dataBuffer, _ := proto.Marshal(data)
	return dataBuffer
}

func (b *bulletPool) eBroadcastBetChange(betChangeProto *bullet_proto.BetChange) []byte {
	data := &bullet_proto.BroadcastBetChange{
		Msgid:     common_proto.EMSGID_eBroadcastBetChange,
		BetChange: betChangeProto,
	}
	dataBuffer, _ := proto.Marshal(data)
	return dataBuffer
}

func (b *bulletPool) destroy() {
	logger.Service.Zap.Infow("Bullet Pool Unregister",
		"GameUser", b.secWebSocketKey,
		"Chan", fmt.Sprintf("%p", b.in),
		"GameId", b.gameId,
		"MathModuleId", b.mathModuleId,
	)
	flux.UnRegister(b.hashId(), b.in)
}

func (b *bulletPool) hashId() string {
	return Service.HashId(b.secWebSocketKey)
}

func (b *bulletPool) checkDrill(gameId string, bulletType int32) bool {
	switch gameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		return psf_on_00001_CheckDrill(bulletType)
	case models.PSF_ON_00002, models.PSF_ON_20002:
		return psf_on_00002_CheckDrill(bulletType)
	case models.PSF_ON_00003:
		return psf_on_00003_CheckDrill(bulletType)
	default:
		return false
	}
}
