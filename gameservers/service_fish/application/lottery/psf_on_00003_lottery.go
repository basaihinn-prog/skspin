package lottery

import (
	"serve/fish_comm/flux"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	PSFM_00004_95_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-95-1"
	PSFM_00004_96_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-96-1"
	PSFM_00004_97_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-97-1"
	PSFM_00004_98_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-98-1"
	"serve/service_fish/domain/slot"
	"serve/service_fish/models"
)

const (
	slot_bigWinPay = 80
)

func psf_on_00003_Check(bulletTypeId, fishTypeId int32) bool {
	return true
}

func psf_on_00003_TriggerBonus(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {

}

func psf_on_00003_Fish(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, checkMulti bool, mercenaryType int32, fishAmount int, accountingSn uint64) (isDied bool, bonusBullet *bullet.Bullet) {
	switch hitResult.FishTypeId {
	case 0, 1, 2, 3, 4, 5:
		if hitResult.Pay > 0 {
			return true, nil
		}

	case 6, 7, 8, 9, 10, 11:
		if len(hitResult.ExtraData) != 2 {
			return
		}

		pay := hitResult.ExtraData[0].(int)
		reels := hitResult.ExtraData[1].([]int32)
		allPay := hitResult.BonusPayload.([]int64)

		bet := psf_on_00003_processWeapon(hitBullet.TypeId, mercenaryType, fishAmount, hitBullet.Bet)

		if hitResult.Pay != 0 && pay > 0 && len(reels) == 9 && len(allPay) == 3 {
			bonusSlot := slot.New(
				hitResult.BonusUuid,
				hitBullet.SecWebSocketKey,
				hitFish.MathModuleId,
				hitFish.RoomUuid,
				hitFish.SeatId,
				int32(hitResult.FishTypeId),
				hitBullet.BetIndex,
				hitBullet.BetLevelIndex,
				hitBullet.LineIndex,
				hitBullet.RateIndex,
			)

			bonusSlot.Pay = uint64(pay)
			bonusSlot.AllPay = allPay
			bonusSlot.Reels = reels
			bonusSlot.ExtraData[0] = hitFish
			bonusSlot.ExtraData[1] = hitBullet
			bonusSlot.ExtraData[2] = hitResult
			bonusSlot.Bet = bet
			bonusSlot.Rate = hitBullet.Rate

			isOk := make(chan bool, 1024)
			flux.Send(slot.ActionSlotBonus, Service.HashId(hitBullet.SecWebSocketKey), slot.Service.Id, bonusSlot, isOk)

			<-isOk
			return true, nil
		}

	case 12, 13, 14, 15, 16, 17:
		bet := psf_on_00003_processWeapon(hitBullet.TypeId, mercenaryType, fishAmount, hitBullet.Bet)

		if hitResult.Pay > 0 && hitResult.BonusPayload.(int) > 0 {
			bonusBullet := bullet.New(
				hitResult.BonusUuid,
				hitBullet.Uuid,
				hitBullet.SecWebSocketKey,
				hitFish.RoomUuid,
				hitBullet.BetIndex,
				hitBullet.LineIndex,
				hitBullet.RateIndex,
				hitBullet.BetLevelIndex,
				200,
				//int32(hitResult.FishTypeId),
				int32(hitResult.BonusPayload.(int)),
				0,
				0,
				1,
				hitBullet.Target,
			)

			bonusBullet.ExtraData[0] = hitFish
			bonusBullet.ExtraData[1] = hitBullet
			bonusBullet.ExtraData[2] = hitResult
			bonusBullet.ExtraData[3] = uint64(0)
			bonusBullet.Bet = bet
			bonusBullet.Rate = hitBullet.Rate

			rtpId := "0"
			switch hitFish.MathModuleId {
			case models.PSFM_00004_95_1:
				rtpId = PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
			case models.PSFM_00004_96_1:
				rtpId = PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
			case models.PSFM_00004_97_1:
				rtpId = PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
			case models.PSFM_00004_98_1:
				rtpId = PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
			}
			bonusBullet.RtpId = rtpId
			if !checkMulti {
				flux.Send(bullet.ActionBulletBonusCall, Service.HashId(hitBullet.SecWebSocketKey), bullet.Service.HashId(hitBullet.SecWebSocketKey), bonusBullet, mercenaryType, accountingSn)
			}

			return true, bonusBullet
		}

	case 203:
		if hitResult.Bullet > 0 {
			return true, nil
		}

	case 303:
		if hitResult.Pay > 0 || hitResult.Bullet > 0 {
			return true, nil
		}

	case 401, 402, 403: // do nothing
	}

	return false, nil
}

func checkTriggerIconId(triggerIconId int) bool {
	switch triggerIconId {
	case 12, 13, 14, 15, 16, 17:
		return false
	default:
		return true
	}
}

func psf_on_00003_processWeapon(bulletType, mercenaryType int32, fishAmount int, bet uint64) uint64 {
	if bulletType == 999 {
		return 1
	}

	// Process 散彈槍 and 迫擊砲
	switch {
	case bulletType == 0 && fishAmount == 3, mercenaryType == 1:
		return 1
	case bulletType == 0 && fishAmount == 5, mercenaryType == 2:
		return 1
	case bulletType == 0 && fishAmount == 10:
		return 1
	}

	return bet
}

func psf_on_00003_broadcasterCheck(typeId int32, isDied bool) (fishId int32, bigWinShow bool) {
	fishId = 0
	bigWinShow = false

	switch typeId {
	case 4, 10, 16: // x80
		fishId = int32(4)
		if isDied {
			bigWinShow = true
		}
	case 5, 11, 17: // x100
		fishId = int32(5)
		if isDied {
			bigWinShow = true
		}
	default:
		fishId = typeId
	}

	return fishId, bigWinShow
}
