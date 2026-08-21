package bullet

import (
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	bullet_proto "serve/service_fish/domain/bullet/proto"
)

const (
	PSF_ON_00004_MAX_201_MachineGunBullets      = 999
	PSF_ON_00004_MAX_202_SuperMachineGunBullets = 999

	PSF_ON_00004_MachineGun_Name      = "MachineGun"
	PSF_ON_00004_SuperMachineGun_Name = "SuperMachineGun"

	PSF_ON_00004_MachineGun_TypeId      = 201
	PSF_ON_00004_SuperMachineGun_TypeId = 202
)

func psf_on_00004_Shoot(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
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

			errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00004_BULLET_FOUND)
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

	case 201:
		machineGunShootCheck(action, secWebSocketKey, PSF_ON_00004_MachineGun_TypeId, shootProto, bonusBullets)

	case 202:
		machineGunShootCheck(action, secWebSocketKey, PSF_ON_00004_SuperMachineGun_TypeId, shootProto, bonusBullets)

	default:
		logger.Service.Zap.Errorw(Bullet_PSF_ON_00004_BULLET_TYPE_INVALID,
			"GameUser", action.Key().From(),
			"GameRoomUuid", shootProto.RoomUuid,
			"BulletType", shootProto.BulletTypeId,
			"BulletUuid", shootProto.BulletUuid,
		)

		errorcode.Service.Fatal(secWebSocketKey, Bullet_PSF_ON_00004_BULLET_TYPE_INVALID)
	}
}

func psf_on_00004_Bonus(action *flux.Action, bonusBullets map[string]*Bullet) {
	bonusBullet := action.Payload()[0].(*Bullet)
	accountingSn := action.Payload()[1].(uint64)

	switch bonusBullet.TypeId {
	case 201:
		machineGunBonusCheck(action, bonusBullet, bonusBullets, 201, accountingSn)

	case 202:
		machineGunBonusCheck(action, bonusBullet, bonusBullets, 202, accountingSn)

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

		errorcode.Service.Fatal(bonusBullet.SecWebSocketKey, Bullet_PSF_ON_00004_BONUS_BULLET_TYPE_INVALID)
	}
}

func psf_on_00004_Hit(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	hitFish := action.Payload()[0]
	hitBullet := action.Payload()[1].(*Bullet)
	hitResult := action.Payload()[2]
	pay := action.Payload()[3].(int)
	accountingSn := action.Payload()[6].(uint64)

	// AccumulatedPay
	hitBullet.ExtraData[3] = uint64(0)

	switch int(hitBullet.TypeId) {
	case 201:
		machineGunHitCheck(action, bonusBullets, PSF_ON_00004_MachineGun_TypeId, accountingSn)

	case 202:
		machineGunHitCheck(action, bonusBullets, PSF_ON_00004_SuperMachineGun_TypeId, accountingSn)

	case 0:
		if b, ok := bullets[hitBullet.Uuid]; ok {
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

			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BULLET_NOT_FOUND)
		}

	default:
		logger.Service.Zap.Errorw("Incorrect Bullet Type",
			"GameUser", hitBullet.SecWebSocketKey,
			"BulletType", hitBullet.TypeId,
			"BulletUuid", hitBullet.Uuid,
			"Pay", pay,
		)

		errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BULLET_TYPE_INVALID)
	}
}

func machineGunHitCheck(action *flux.Action, bonusBullets map[string]*Bullet, typeId int32, accountingSn uint64) {
	hitFish := action.Payload()[0]
	hitBullet := action.Payload()[1].(*Bullet)
	hitResult := action.Payload()[2]
	pay := action.Payload()[3].(int)

	var bulletName = ""
	switch typeId {
	case 201:
		bulletName = PSF_ON_00004_MachineGun_Name

	case 202:
		bulletName = PSF_ON_00004_SuperMachineGun_Name
	}

	if b, ok := bonusBullets[hitBullet.BonusUuid]; ok {
		if b.UsedBullets < 0 {
			b.UsedBullets = 0
		}
		b.UsedBullets++

		if b.ShootBullets > 0 {
			b.ShootBullets--
		}

		logger.Service.Zap.Infow(bulletName+" Bullet Hit",
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

		flux.Send(ActionBulletHitRecall, Service.HashId(b.SecWebSocketKey), action.Key().From(), hitFish, hitBullet, hitResult, accountingSn)

		if b.UsedBullets == b.Bullets {
			logger.Service.Zap.Infow("End "+bulletName+" Bullet",
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
				logger.Service.Zap.Errorw("Incorrect "+bulletName+" Bullet Count",
					"GameUser", b.SecWebSocketKey,
					"BulletType", b.TypeId,
					"BonusUuid", b.BonusUuid,
					"Bullets", b.Bullets,
					"ShootBullets", b.ShootBullets,
					"UsedBullets", b.UsedBullets,
					"Pay", pay,
					"AccumulatedPay", b.ExtraData[3],
				)

				switch typeId {
				case 201:
					errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BONUS_BULLET_201_INVALID)
				case 202:
					errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BONUS_BULLET_202_INVALID)
				}
			}
			return
		}

		if b.UsedBullets > b.Bullets {
			logger.Service.Zap.Errorw("Incorrect "+bulletName+" Bullet Count",
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

			switch typeId {
			case 201:
				errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BONUS_BULLET_201_INVALID)
			case 202:
				errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BONUS_BULLET_202_INVALID)
			}
		}
	} else {
		logger.Service.Zap.Errorw(bulletName+" Bullet Not Found",
			"GameUser", hitBullet.SecWebSocketKey,
			"BulletType", hitBullet.TypeId,
			"BonusUuid", hitBullet.BonusUuid,
			"Bullets", hitBullet.Bullets,
			"ShootBullets", hitBullet.ShootBullets,
			"UsedBullets", hitBullet.UsedBullets,
			"Pay", pay,
			"AccumulatedPay", hitBullet.ExtraData[3],
		)

		switch typeId {
		case 201:
			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BONUS_BULLET_201_NOT_FOUND)
		case 202:
			errorcode.Service.Fatal(hitBullet.SecWebSocketKey, Bullet_PSF_ON_00004_HIT_BONUS_BULLET_202_NOT_FOUND)
		}
	}
}

func machineGunShootCheck(action *flux.Action, secWebSocketKey string, bulletType int32,
	shootProto *bullet_proto.Shoot, bonusBullets map[string]*Bullet) {
	var bulletName = ""
	var errorMsg = ""
	switch shootProto.BulletTypeId {
	case 201:
		bulletName = PSF_ON_00004_MachineGun_Name
		errorMsg = Bullet_PSF_ON_00004_BONUS_BULLET_201_NOT_FOUND

	case 202:
		bulletName = PSF_ON_00004_SuperMachineGun_Name
		errorMsg = Bullet_PSF_ON_00004_BONUS_BULLET_202_NOT_FOUND
	}

	if len(shootProto.BonusBullet) <= 0 {
		logger.Service.Zap.Errorw(errorMsg,
			"GameUser", action.Key().From(),
			"GameRoomUuid", shootProto.RoomUuid,
			"BonusBullet", shootProto.BonusBullet,
		)
		errorcode.Service.Fatal(secWebSocketKey, errorMsg)
		return
	}

	if shootProto.BonusBullet[0] == nil {
		logger.Service.Zap.Errorw(errorMsg,
			"GameUser", action.Key().From(),
			"GameRoomUuid", shootProto.RoomUuid,
			"BonusBullet", shootProto.BonusBullet,
		)
		errorcode.Service.Fatal(secWebSocketKey, errorMsg)
		return
	}

	if shootProto.BonusBullet[0].Type == bulletType {
		if b, ok := bonusBullets[shootProto.BonusBullet[0].Uuid]; ok {
			logger.Service.Zap.Infow("Shoot "+bulletName+" Bullet",
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
				"BulletType", shootProto.BonusBullet[0].Type,
				"BulletUuid", shootProto.BulletUuid,
				"Bullets", b.Bullets,
				"ShootBullets", b.ShootBullets,
				"UsedBullets", b.UsedBullets,
			)
			b.ShootBullets++
		} else {
			logger.Service.Zap.Errorw(errorMsg,
				"GameUser", action.Key().From(),
				"GameRoomUuid", shootProto.RoomUuid,
			)

			errorcode.Service.Fatal(secWebSocketKey, errorMsg)
		}
	}
}

func machineGunBonusCheck(action *flux.Action, bonusBullet *Bullet, bonusBullets map[string]*Bullet, typeId int32, accountingSn uint64) {
	var remainBullets int32 = 0
	var shootBullets int32 = 0

	for _, v := range bonusBullets {
		if v.TypeId == typeId && v.SecWebSocketKey == bonusBullet.SecWebSocketKey {
			shootBullets += v.ShootBullets
			remain := v.Bullets - v.UsedBullets
			remainBullets += remain
		}
	}

	var quota int32 = 0
	var bulletName = ""
	var errorMsg = ""
	switch typeId {
	case 201:
		quota = PSF_ON_00004_MAX_201_MachineGunBullets + shootBullets - remainBullets
		bulletName = PSF_ON_00004_MachineGun_Name
		errorMsg = Bullet_PSF_ON_00004_BONUS_BULLET_201_NOT_INVALID

	case 202:
		quota = PSF_ON_00004_MAX_202_SuperMachineGunBullets + shootBullets - remainBullets
		bulletName = PSF_ON_00004_SuperMachineGun_Name
		errorMsg = Bullet_PSF_ON_00004_BONUS_BULLET_202_NOT_INVALID
	}

	if quota < 0 {
		logger.Service.Zap.Errorw("Got "+bulletName+" Bullets",
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

		errorcode.Service.Fatal(bonusBullet.SecWebSocketKey, errorMsg)
		return
	}

	if quota >= bonusBullet.Bullets {
		bonusBullets[bonusBullet.BonusUuid] = bonusBullet

		logger.Service.Zap.Infow("Got "+bulletName+" Bullets",
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

	logger.Service.Zap.Infow("Got "+bulletName+" Bullets",
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
}
