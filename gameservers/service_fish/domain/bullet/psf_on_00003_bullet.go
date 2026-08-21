package bullet

import (
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	bullet_proto "serve/service_fish/domain/bullet/proto"
)

func psf_on_00003_Shoot(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	secWebSocketKey := action.Key().From()

	shootProto := action.Payload()[0].(*bullet_proto.Shoot)
	bet := action.Payload()[1].(uint64)
	rate := action.Payload()[2].(uint64)

	switch shootProto.BulletTypeId {
	case 0, 999:
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
			shootProto.BulletTypeId,
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

	// Drill
	case 200:
		if len(shootProto.BonusBullet) <= 0 {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"BonusBullet", shootProto.BonusBullet,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND)
			return
		}

		if shootProto.BonusBullet[0] == nil {
			logger.Service.Zap.Errorw(Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"BonusBullet", shootProto.BonusBullet,
			)
			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_26_NOT_FOUND)
			return
		}

		if b, ok := bonusBullets[shootProto.BonusBullet[0].Uuid]; ok {
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

func psf_on_00003_Multi_Hit(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	hitFishies := action.Payload()[0]
	hitBullet := action.Payload()[1].(*Bullet)
	hitResults := action.Payload()[2]
	pay := action.Payload()[3].(int)
	mercenaryType := action.Payload()[4].(int32)
	accountingSn := action.Payload()[5].(uint64)

	hitBullet.ExtraData[3] = uint64(0)

	switch int(hitBullet.TypeId) {
	case 0, 999:
		if b, ok := bullets[hitBullet.Uuid]; ok {
			delete(bullets, hitBullet.Uuid)
			flux.Send(ActionBulletMultiRecall, Service.HashId(b.SecWebSocketKey), action.Key().From(), hitFishies, hitBullet, hitResults, accountingSn, mercenaryType)

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

func psf_on_00003_Hit(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	hitFish := action.Payload()[0]
	hitBullet := action.Payload()[1].(*Bullet)
	hitResult := action.Payload()[2]
	pay := action.Payload()[3].(int)
	mercenaryType := action.Payload()[4].(int32)
	accountingSn := action.Payload()[6].(uint64)

	hitBullet.ExtraData[3] = uint64(0) // AccumulatedPay

	switch int(hitBullet.TypeId) {
	// Drill
	case 200:
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

			if b.UsedBullets == b.Bullets {
				hitBullet.ExtraData[3] = b.ExtraData[3]

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

			flux.Send(ActionBulletHitRecall, Service.HashId(b.SecWebSocketKey), action.Key().From(), hitFish, hitBullet, hitResult, accountingSn, mercenaryType)

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
	case 0, 999:
		if b, ok := bullets[hitBullet.Uuid]; ok {
			delete(bullets, hitBullet.Uuid)
			flux.Send(ActionBulletHitRecall, Service.HashId(b.SecWebSocketKey), action.Key().From(), hitFish, hitBullet, hitResult, accountingSn, mercenaryType)

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

func psf_on_00003_multi_Bonus(action *flux.Action, bonusBullets map[string]*Bullet) {
	bonusBulletList := action.Payload()[0].([]*Bullet)

	for i := 0; i < len(bonusBulletList); i++ {
		if bonusBulletList[i] != nil {
			switch bonusBulletList[i].TypeId {
			case 12, 13, 14, 15, 16, 17, 200:
				bonusBullets[bonusBulletList[i].BonusUuid] = bonusBulletList[i]

				logger.Service.Zap.Infow("Got Drill Bullets",
					"GameUser", bonusBulletList[i].SecWebSocketKey,
					"BulletType", bonusBulletList[i].TypeId,
					"BulletUuid", bonusBulletList[i].Uuid,
					"BonusUuid", bonusBulletList[i].BonusUuid,
					"Bet", bonusBulletList[i].Bet,
					"Rate", bonusBulletList[i].Rate,
					"Bullets", bonusBulletList[i].Bullets,
					"UsedBullets", bonusBulletList[i].UsedBullets,
				)

			default:
				logger.Service.Zap.Errorw("Incorrect Bonus Bullet Type",
					"GameUser", bonusBulletList[i].SecWebSocketKey,
					"BulletType", bonusBulletList[i].TypeId,
					"BulletUuid", bonusBulletList[i].Uuid,
					"BonusUuid", bonusBulletList[i].BonusUuid,
					"Bet", bonusBulletList[i].Bet,
					"Rate", bonusBulletList[i].Rate,
					"Bullets", bonusBulletList[i].Bullets,
					"UsedBullets", bonusBulletList[i].UsedBullets,
				)
				errorcode.Service.Fatal(bonusBulletList[i].SecWebSocketKey, Bullet_PSF_ON_00001_BONUS_BULLET_TYPE_INVALID)
			}
		}
	}
}

func psf_on_00003_Bonus(action *flux.Action, bonusBullets map[string]*Bullet) {
	bonusBullet := action.Payload()[0].(*Bullet)
	mercenaryType := action.Payload()[1].(int32)
	accountingSn := action.Payload()[2].(uint64)

	switch bonusBullet.TypeId {
	case 12, 13, 14, 15, 16, 17, 200:
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
		flux.Send(ActionBulletBonusRecall, action.Key().To(), action.Key().From(), bonusBullet, accountingSn, mercenaryType)

	default:
		logger.Service.Zap.Errorw("Incorrect Bonus Bullet Type",
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

func psf_on_00003_CheckDrill(bulletType int32) bool {
	switch bulletType {
	case 12, 13, 14, 15, 16, 17, 200:
		return true

	default:
		return false
	}
}
