package bullet

import (
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	bullet_proto "serve/service_fish/domain/bullet/proto"
)

const (
	PSF_ON_00001_MAX_27_MachineGunBullets = 999
)

func psf_on_00001_Shoot(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	secWebSocketKey := action.Key().From()

	shootProto := action.Payload()[0].(*bullet_proto.Shoot)
	bet := action.Payload()[1].(uint64)
	rate := action.Payload()[2].(uint64)

	switch shootProto.BulletTypeId {
	case 0:
		if _, ok := bullets[shootProto.BulletUuid]; ok {
			logger.Service.Zap.Errorw("Duplicate Regular Bullet",
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"BulletType", shootProto.BulletTypeId,
				"BulletUuid", shootProto.BulletUuid,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BULLET_FOUND)
			return
		}

		shootBullet := New(
			"",
			shootProto.BulletUuid,
			secWebSocketKey,
			shootProto.RoomUuid,
			shootProto.BetRateLine.BetIndex,
			shootProto.BetRateLine.LineIndex,
			shootProto.BetRateLine.RateIndex,
			shootProto.BetRateLine.BetLevelIndex,
			0,
			-1,
			-1,
			-1,
			shootProto.ExtraData.ShootMode,
			shootProto.ExtraData.Target,
		)
		shootBullet.Bet = bet
		shootBullet.Rate = rate
		bullets[shootProto.BulletUuid] = shootBullet

		logger.Service.Zap.Debugw("Got New Bullet",
			"GameUser", secWebSocketKey,
			"GameRoomUuid", shootProto.RoomUuid,
			"BulletType", shootProto.BulletTypeId,
			"BulletUuid", shootProto.BulletUuid,
		)

	case 26:
		if len(shootProto.Drill) <= 0 {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"Drill", shootProto.Drill,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND)
			return
		}

		if shootProto.Drill[0] == nil {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"Drill", shootProto.Drill,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND)
			return
		}

		if b, ok := bonusBullets[shootProto.Drill[0].Uuid]; ok {
			logger.Service.Zap.Infow("Shoot Drill Bullet",
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"BulletType", b.TypeId,
				"BulletUuid", b.Uuid,
				"Bullets", b.Bullets,
				"ShootBullets", b.ShootBullets,
				"UsedBullets", b.UsedBullets,
			)
		} else {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND)
		}

	case 27:
		if len(shootProto.MachineGun) <= 0 {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_27_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"MachineGun", shootProto.MachineGun,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_27_NOT_FOUND)
			return
		}

		if shootProto.MachineGun[0] == nil {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_27_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"MachineGun", shootProto.MachineGun,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_27_NOT_FOUND)
			return
		}

		if b, ok := bonusBullets[shootProto.MachineGun[0].Uuid]; ok {
			logger.Service.Zap.Infow("Shoot MachineGun Bullet",
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"BulletType", b.TypeId,
				"BulletUuid", b.Uuid,
				"Bullets", b.Bullets,
				"ShootBullets", b.ShootBullets,
				"UsedBullets", b.UsedBullets,
			)
			b.ShootBullets++
		} else {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_27_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_27_NOT_FOUND)
		}

	default:
		logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BULLET_TYPE_INVALID,
			"GameUser", action.Key().From(),
			"GameRoomUuid", shootProto.RoomUuid,
			"BulletType", shootProto.BulletTypeId,
			"BulletUuid", shootProto.BulletUuid,
		)
		errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BULLET_TYPE_INVALID)
	}
}

func psf_on_00001_Bonus(action *flux.Action, bonusBullets map[string]*Bullet) {
	bonusBullet := action.Payload()[0].(*Bullet)
	accountingSn := action.Payload()[1].(uint64)

	switch bonusBullet.TypeId {
	case 26:
		bonusBullets[bonusBullet.BonusUuid] = bonusBullet

		logger.Service.Zap.Infow("Got Drill Bullets",
			"GameUser", bonusBullet.SecWebSocketKey,
			"BulletType", bonusBullet.TypeId,
			"BulletUuid", bonusBullet.Uuid,
			"BonusUuid", bonusBullet.BonusUuid,
			"Bet", bonusBullet.Bet,
			"Rate", bonusBullet.Rate,
			"Bullets", bonusBullet.Bullets,
			"UsedBullets", bonusBullet.UsedBullets,
		)

		flux.Send(ActionBulletBonusRecall, action.Key().To(), action.Key().From(), bonusBullet, accountingSn)

	case 27:
		var remainBullets int32 = 0
		var shootBullets int32 = 0

		for _, v := range bonusBullets {
			if v.TypeId == 27 && v.SecWebSocketKey == bonusBullet.SecWebSocketKey {
				shootBullets += v.ShootBullets
				remain := v.Bullets - v.UsedBullets
				remainBullets += remain
			}
		}

		quota := PSF_ON_00001_MAX_27_MachineGunBullets + shootBullets - remainBullets

		if quota < 0 {
			logger.Service.Zap.Errorw("Got MachineGun Bullets",
				"GameUser", bonusBullet.SecWebSocketKey,
				"BulletType", bonusBullet.TypeId,
				"BulletUuid", bonusBullet.Uuid,
				"BonusUuid", bonusBullet.BonusUuid,
				"Bet", bonusBullet.Bet,
				"Rate", bonusBullet.Rate,
				"Bullets", bonusBullet.Bullets,
				"ShootBullets", shootBullets,
				"UsedBullets", bonusBullet.UsedBullets,
				"RemainBullets", remainBullets,
				"Quota", quota,
			)

			errorcode.Service.Fatal(bonusBullet.SecWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_27_INVALID)
			return
		}

		if quota == 0 {
			bonusBullet.Bullets = 0

			logger.Service.Zap.Infow("Got MachineGun Bullets",
				"GameUser", bonusBullet.SecWebSocketKey,
				"BulletType", bonusBullet.TypeId,
				"BulletUuid", bonusBullet.Uuid,
				"BonusUuid", bonusBullet.BonusUuid,
				"Bet", bonusBullet.Bet,
				"Rate", bonusBullet.Rate,
				"Bullets", bonusBullet.Bullets,
				"ShootBullets", shootBullets,
				"UsedBullets", bonusBullet.UsedBullets,
				"RemainBullets", remainBullets,
				"Quota", quota,
			)

			flux.Send(ActionBulletBonusRecall, action.Key().To(), action.Key().From(), bonusBullet, accountingSn)
			return
		}

		if quota >= bonusBullet.Bullets {
			bonusBullets[bonusBullet.BonusUuid] = bonusBullet

			logger.Service.Zap.Infow("Got MachineGun Bullets",
				"GameUser", bonusBullet.SecWebSocketKey,
				"BulletType", bonusBullet.TypeId,
				"BulletUuid", bonusBullet.Uuid,
				"BonusUuid", bonusBullet.BonusUuid,
				"Bet", bonusBullet.Bet,
				"Rate", bonusBullet.Rate,
				"Bullets", bonusBullet.Bullets,
				"ShootBullets", shootBullets,
				"UsedBullets", bonusBullet.UsedBullets,
				"RemainBullets", remainBullets,
				"Quota", quota,
			)

			flux.Send(ActionBulletBonusRecall, action.Key().To(), action.Key().From(), bonusBullet, accountingSn)
			return
		}

		bonusBullet.Bullets = quota
		bonusBullets[bonusBullet.BonusUuid] = bonusBullet

		logger.Service.Zap.Infow("Got MachineGun Bullets",
			"GameUser", bonusBullet.SecWebSocketKey,
			"BulletType", bonusBullet.TypeId,
			"BulletUuid", bonusBullet.Uuid,
			"BonusUuid", bonusBullet.BonusUuid,
			"Bet", bonusBullet.Bet,
			"Rate", bonusBullet.Rate,
			"Bullets", bonusBullet.Bullets,
			"ShootBullets", shootBullets,
			"UsedBullets", bonusBullet.UsedBullets,
			"RemainBullets", remainBullets,
			"Quota", quota,
		)

		flux.Send(ActionBulletBonusRecall, action.Key().To(), action.Key().From(), bonusBullet, accountingSn)

	default:
		logger.Service.Zap.Errorw("Incorrect Bonus Bullet Type",
			"GameUser", bonusBullet.SecWebSocketKey,
			"BulletType", bonusBullet.TypeId,
			"BulletUuid", bonusBullet.Uuid,
			"BonusUuid", bonusBullet.BonusUuid,
			"Bet", bonusBullet.Bet,
			"Rate", bonusBullet.Rate,
			"Bullets", bonusBullet.Bullets,
			"UsedBullets", bonusBullet.UsedBullets,
		)
		errorcode.Service.Fatal(bonusBullet.SecWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_TYPE_INVALID)
	}
}

func psf_on_00001_Hit(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	hitFish := action.Payload()[0]
	hitBullet := action.Payload()[1].(*Bullet)
	hitResult := action.Payload()[2]
	pay := action.Payload()[3].(int)
	accountingSn := action.Payload()[6].(uint64)

	hitBullet.ExtraData[3] = uint64(0) // AccumulatedPay

	switch int(hitBullet.TypeId) {
	case 26:
		if b, ok := bonusBullets[hitBullet.BonusUuid]; ok {
			if b.UsedBullets < 0 {
				b.UsedBullets = 0
			}

			b.UsedBullets++
			b.ExtraData[3] = b.ExtraData[3].(uint64) + uint64(pay) // AccumulatedPay

			logger.Service.Zap.Infow("Drill Bullet Hit",
				"GameUser", b.SecWebSocketKey,
				"BulletType", b.TypeId,
				"BonusUuid", b.BonusUuid,
				"Bullets", b.Bullets,
				"ShootBullets", b.ShootBullets,
				"UsedBullets", b.UsedBullets,
				"Pay", pay,
				"AccumulatedPay", b.ExtraData[3],
			)

			hitBullet.Bullets = b.Bullets
			hitBullet.UsedBullets = b.UsedBullets
			//hitBullet.ExtraData[3] = uint64(0) // AccumulatedPay

			if b.UsedBullets == b.Bullets {
				hitBullet.ExtraData[3] = b.ExtraData[3] // AccumulatedPay

				logger.Service.Zap.Infow("End Drill Bullet",
					"GameUser", b.SecWebSocketKey,
					"BulletType", b.TypeId,
					"BonusUuid", b.BonusUuid,
					"Bullets", b.Bullets,
					"ShootBullets", b.ShootBullets,
					"UsedBullets", b.UsedBullets,
					"Pay", pay,
					"AccumulatedPay", b.ExtraData[3],
				)
				delete(bonusBullets, hitBullet.BonusUuid)
			}

			if b.UsedBullets > b.Bullets {
				logger.Service.Zap.Errorw("Incorrect Drill Bullet Count",
					"GameUser", b.SecWebSocketKey,
					"BulletType", b.TypeId,
					"BonusUuid", b.BonusUuid,
					"Bullets", b.Bullets,
					"ShootBullets", b.ShootBullets,
					"UsedBullets", b.UsedBullets,
					"Pay", pay,
					"AccumulatedPay", b.ExtraData[3],
				)
				delete(bonusBullets, hitBullet.BonusUuid)
				errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00001_HIT_BONUS_BULLET_26_INVALID)
			}

			flux.Send(ActionBulletHitRecall, Service.HashId(b.SecWebSocketKey), action.Key().From(), hitFish, hitBullet, hitResult, accountingSn)

		} else {
			logger.Service.Zap.Errorw("Drill Bullet Not Found",
				"GameUser", hitBullet.SecWebSocketKey,
				"BulletType", hitBullet.TypeId,
				"BonusUuid", hitBullet.BonusUuid,
				"Bullets", hitBullet.Bullets,
				"ShootBullets", hitBullet.ShootBullets,
				"UsedBullets", hitBullet.UsedBullets,
				"Pay", pay,
				"AccumulatedPay", hitBullet.ExtraData[3],
			)
			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00001_HIT_BONUS_BULLET_26_NOT_FOUND)
		}

	case 27:
		if b, ok := bonusBullets[hitBullet.BonusUuid]; ok {
			if b.UsedBullets < 0 {
				b.UsedBullets = 0
			}

			b.UsedBullets++

			if b.ShootBullets > 0 {
				b.ShootBullets--
			}

			logger.Service.Zap.Infow("MachineGun Bullet Hit",
				"GameUser", b.SecWebSocketKey,
				"BulletType", b.TypeId,
				"BonusUuid", b.BonusUuid,
				"Bullets", b.Bullets,
				"ShootBullets", b.ShootBullets,
				"UsedBullets", b.UsedBullets,
				"Pay", pay,
				"AccumulatedPay", b.ExtraData[3],
			)

			hitBullet.Bullets = b.Bullets
			hitBullet.UsedBullets = b.UsedBullets
			hitBullet.ShootBullets = b.ShootBullets
			//hitBullet.ExtraData[3] = uint64(0) // AccumulatedPay

			flux.Send(ActionBulletHitRecall, Service.HashId(b.SecWebSocketKey), action.Key().From(), hitFish, hitBullet, hitResult, accountingSn)

			if b.UsedBullets == b.Bullets {
				logger.Service.Zap.Infow("End MachineGun Bullet",
					"GameUser", b.SecWebSocketKey,
					"BulletType", b.TypeId,
					"BonusUuid", b.BonusUuid,
					"Bullets", b.Bullets,
					"ShootBullets", b.ShootBullets,
					"UsedBullets", b.UsedBullets,
					"Pay", pay,
					"AccumulatedPay", b.ExtraData[3],
				)
				delete(bonusBullets, hitBullet.BonusUuid)

				if b.ShootBullets != 0 {
					logger.Service.Zap.Errorw("Incorrect MachineGun Bullet Count",
						"GameUser", b.SecWebSocketKey,
						"BulletType", b.TypeId,
						"BonusUuid", b.BonusUuid,
						"Bullets", b.Bullets,
						"ShootBullets", b.ShootBullets,
						"UsedBullets", b.UsedBullets,
						"Pay", pay,
						"AccumulatedPay", b.ExtraData[3],
					)
					errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00001_HIT_BONUS_BULLET_27_INVALID)
				}
				return
			}

			if b.UsedBullets > b.Bullets {
				logger.Service.Zap.Errorw("Incorrect MachineGun Bullet Count",
					"GameUser", b.SecWebSocketKey,
					"BulletType", b.TypeId,
					"BonusUuid", b.BonusUuid,
					"Bullets", b.Bullets,
					"ShootBullets", b.ShootBullets,
					"UsedBullets", b.UsedBullets,
					"Pay", pay,
					"AccumulatedPay", b.ExtraData[3],
				)
				delete(bonusBullets, hitBullet.BonusUuid)
				errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00001_HIT_BONUS_BULLET_27_INVALID)
			}
		} else {
			logger.Service.Zap.Errorw("MachineGun Bullet Not Found",
				"GameUser", hitBullet.SecWebSocketKey,
				"BulletType", hitBullet.TypeId,
				"BonusUuid", hitBullet.BonusUuid,
				"Bullets", hitBullet.Bullets,
				"ShootBullets", hitBullet.ShootBullets,
				"UsedBullets", hitBullet.UsedBullets,
				"Pay", pay,
				"AccumulatedPay", hitBullet.ExtraData[3],
			)
			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00001_HIT_BONUS_BULLET_27_NOT_FOUND)
		}

	case 0:
		if b, ok := bullets[hitBullet.Uuid]; ok {
			//hitBullet.ExtraData[3] = uint64(0)

			delete(bullets, hitBullet.Uuid)

			flux.Send(ActionBulletHitRecall, Service.HashId(b.SecWebSocketKey), action.Key().From(), hitFish, hitBullet, hitResult, accountingSn)

			logger.Service.Zap.Infow("Delete Hit Bullet",
				"GameUser", b.SecWebSocketKey,
				"BulletType", b.TypeId,
				"BulletUuid", b.Uuid,
				"Bet", b.Bet,
				"Rate", b.Rate,
				"RegularBulletsCount", len(bullets),
				"Pay", pay,
			)
		} else {
			logger.Service.Zap.Errorw("Hit Bullet Not Found",
				"GameUser", hitBullet.SecWebSocketKey,
				"BulletType", hitBullet.TypeId,
				"BulletUuid", hitBullet.Uuid,
				"Pay", pay,
			)
			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00001_HIT_BULLET_NOT_FOUND)
		}

	default:
		logger.Service.Zap.Errorw("Incorrect Bullet Type",
			"GameUser", hitBullet.SecWebSocketKey,
			"BulletType", hitBullet.TypeId,
			"BulletUuid", hitBullet.Uuid,
			"Pay", pay,
		)
		errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00001_HIT_BULLET_TYPE_INVALID)
	}
}

func psf_on_00001_CheckDrill(bulletType int32) bool {
	if bulletType == 26 {
		return true
	}

	return false
}
