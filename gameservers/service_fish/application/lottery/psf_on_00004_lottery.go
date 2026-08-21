package lottery

import (
	"serve/fish_comm/flux"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	PSF_ON_00004_1 "serve/service_fish/domain/probability/PSF-ON-00004-1"
	PSFM_00005_95_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-95-1"
	PSFM_00005_96_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-96-1"
	PSFM_00005_97_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-97-1"
	PSFM_00005_98_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-98-1"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"serve/service_fish/models"
)

const (
	PSF_ON_00004_RedEnvelope_TypeId     = 100
	PSF_ON_00004_Slot_TypeId            = 101
	PSF_ON_00004_MachineGun_TypeId      = 201
	PSF_ON_00004_SuperMachineGun_TypeId = 202
)

func psf_on_00004_Check(bulletTypeId, fishTypeId int32) bool {
	if bulletTypeId == PSF_ON_00004_MachineGun_TypeId && fishTypeId == PSF_ON_00004_SuperMachineGun_TypeId {
		return false
	}

	if bulletTypeId == PSF_ON_00004_SuperMachineGun_TypeId && fishTypeId == PSF_ON_00004_MachineGun_TypeId {
		return false
	}

	return true
}

func psf_on_00004_Fish(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, accountingSn uint64) (isDied bool) {
	switch hitResult.FishTypeId {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 300, 301, 302:
		if hitResult.Pay > 0 {
			return true
		}

	case 400: // do nothing
	case PSF_ON_00004_MachineGun_TypeId, PSF_ON_00004_SuperMachineGun_TypeId:
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
				int32(hitResult.FishTypeId),
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
			bonusBullet.Bet = hitBullet.Bet
			bonusBullet.Rate = hitBullet.Rate

			rtpId := "0"
			switch hitFish.MathModuleId {
			case models.PSFM_00005_95_1:
				rtpId = psf_on_00004_getWeaponRtp(hitResult.FishTypeId, PSFM_00005_95_1.RTP95BS.PSF_ON_00004_1_BsMath)

			case models.PSFM_00005_96_1:
				rtpId = psf_on_00004_getWeaponRtp(hitResult.FishTypeId, PSFM_00005_96_1.RTP96BS.PSF_ON_00004_1_BsMath)

			case models.PSFM_00005_97_1:
				rtpId = psf_on_00004_getWeaponRtp(hitResult.FishTypeId, PSFM_00005_97_1.RTP97BS.PSF_ON_00004_1_BsMath)

			case models.PSFM_00005_98_1:
				rtpId = psf_on_00004_getWeaponRtp(hitResult.FishTypeId, PSFM_00005_98_1.RTP98BS.PSF_ON_00004_1_BsMath)
			}
			bonusBullet.RtpId = rtpId

			flux.Send(bullet.ActionBulletBonusCall, Service.HashId(hitBullet.SecWebSocketKey), bullet.Service.HashId(hitBullet.SecWebSocketKey), bonusBullet, accountingSn)
		}

	case PSF_ON_00004_RedEnvelope_TypeId:
		if hitResult.TriggerIconId == PSF_ON_00004_RedEnvelope_TypeId && len(hitResult.BonusPayload.([]int)) > 0 {
			bonusRedEnvelope := redenvelope.New(
				hitResult.BonusUuid,
				hitBullet.SecWebSocketKey,
				hitFish.MathModuleId,
				hitFish.RoomUuid,
				hitFish.SeatId,
				hitFish.TypeId,
				hitResult.BonusPayload.([]int),
				hitBullet.BetIndex,
				hitBullet.BetLevelIndex,
				hitBullet.LineIndex,
				hitBullet.RateIndex,
			)

			bonusRedEnvelope.ExtraData[0] = hitFish
			bonusRedEnvelope.ExtraData[1] = hitBullet
			bonusRedEnvelope.ExtraData[2] = hitResult
			bonusRedEnvelope.Bet = hitBullet.Bet
			bonusRedEnvelope.Rate = hitBullet.Rate

			isOk := make(chan bool, 1024)
			flux.Send(redenvelope.ActionRedEnvelopeBonus, Service.HashId(hitBullet.SecWebSocketKey), redenvelope.Service.Id,
				bonusRedEnvelope,
				isOk,
			)
			<-isOk

			return true
		}

	case PSF_ON_00004_Slot_TypeId:
		// no hit
		if len(hitResult.ExtraData) != 2 {
			return
		}

		pay := hitResult.ExtraData[0].(int)
		reels := hitResult.ExtraData[1].([]int32)
		allPay := hitResult.BonusPayload.([]int64)

		if hitResult.Pay == 0 && pay > 0 && len(reels) == 9 && len(allPay) == 3 {

			bonusSlot := slot.New(
				hitResult.BonusUuid,
				hitBullet.SecWebSocketKey,
				hitFish.MathModuleId,
				hitFish.RoomUuid,
				hitFish.SeatId,
				PSF_ON_00004_Slot_TypeId,
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
			bonusSlot.Bet = hitBullet.Bet
			bonusSlot.Rate = hitBullet.Rate

			isOk := make(chan bool, 1024)
			flux.Send(slot.ActionSlotBonus, Service.HashId(hitBullet.SecWebSocketKey), slot.Service.Id, bonusSlot, isOk)

			<-isOk

			return true
		}
	}

	return false
}

func psf_on_00004_TriggerBonus(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
}

func psf_on_00004_getWeaponRtp(fishId int, math *PSF_ON_00004_1.BsMath) (rtpId string) {
	if fishId == PSF_ON_00004_MachineGun_TypeId {
		rtpId = math.Icons.MachineGun.UseRTP
	} else {
		rtpId = math.Icons.SuperMachineGun.UseRTP
	}

	return rtpId
}
