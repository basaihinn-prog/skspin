package lottery

import (
	"serve/fish_comm/flux"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	PSFM_00003_93_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-93-1"
	PSFM_00003_94_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-94-1"
	PSFM_00003_95_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-95-1"
	PSFM_00003_96_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-96-1"
	PSFM_00003_97_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-97-1"
	PSFM_00003_98_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-98-1"
	"serve/service_fish/models"
)

func psf_on_00002_Check(bulletTypeId, fishTypeId int32) bool {
	return psf_on_00001_Check(bulletTypeId, fishTypeId)
}

func psf_on_00002_Fish(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, accountingSn uint64) (isDied bool) {
	switch hitResult.FishTypeId {
	case 26, 27:
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
			case models.PSFM_00003_93_1:
				rtpId = PSFM_00003_93_1.RTP93BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00003_94_1:
				rtpId = PSFM_00003_94_1.RTP94BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00003_95_1:
				rtpId = PSFM_00003_95_1.RTP95BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00003_96_1:
				rtpId = PSFM_00003_96_1.RTP96BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00003_97_1:
				rtpId = PSFM_00003_97_1.RTP97BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00003_98_1:
				rtpId = PSFM_00003_98_1.RTP98BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID
			}

			bonusBullet.RtpId = rtpId

			flux.Send(bullet.ActionBulletBonusCall, Service.HashId(hitBullet.SecWebSocketKey), bullet.Service.HashId(hitBullet.SecWebSocketKey), bonusBullet, accountingSn)
		}
		return false
	}

	return psf_on_00001_Fish(hitFish, hitBullet, hitResult, accountingSn)
}

func psf_on_00002_TriggerBonus(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
	psf_on_00001_TriggerBonus(hitFish, hitBullet, hitResult)
}

func psf_on_00002_ExtraTriggerBonus(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
	if hitResult.ExtraTriggerBonus != nil {

		switch hitResult.TriggerIconId {
		case 20, 26, 27:
			for _, v := range hitResult.ExtraTriggerBonus {
				v.Pay = hitResult.Pay // this is for accounting purpose, see f0dc137f
				psf_on_00001_TriggerBonus(hitFish, hitBullet, v)
			}
		}
	}
}
