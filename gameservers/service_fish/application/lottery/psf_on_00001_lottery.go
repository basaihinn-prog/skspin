package lottery

import (
	"serve/fish_comm/flux"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	PSFM_00002_95_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-95-1"
	PSFM_00002_96_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-96-1"
	PSFM_00002_97_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-97-1"
	PSFM_00002_98_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-98-1"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"serve/service_fish/models"
	"strconv"
)

func psf_on_00001_Check(bulletTypeId, fishTypeId int32) bool {
	if bulletTypeId == 26 && fishTypeId == 26 {
		return false
	}

	if bulletTypeId == 26 && fishTypeId == 27 {
		return false
	}

	if bulletTypeId == 27 && fishTypeId == 26 {
		return false
	}
	return true
}

func psf_on_00001_Fish(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, accountingSn uint64) (isDied bool) {
	switch hitResult.FishTypeId {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20:
		if hitResult.Pay > 0 {
			return true
		}

	case 21, 22, 23, 24, 25:
		// do nothing

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
			case models.PSFM_00001_98_1:
				rtpId = hitBullet.RtpId

			case models.PSFM_00002_95_1:
				rtpId = PSFM_00002_95_1.RTP95BS.PSF_ON_00001_2_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00002_96_1:
				rtpId = PSFM_00002_96_1.RTP96BS.PSF_ON_00001_2_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00002_97_1:
				rtpId = PSFM_00002_97_1.RTP97BS.PSF_ON_00001_2_BsMath.Icons.MachineGun.RTP8.ID

			case models.PSFM_00002_98_1:
				rtpId = PSFM_00002_98_1.RTP98BS.PSF_ON_00001_2_BsMath.Icons.MachineGun.RTP8.ID
			}

			bonusBullet.RtpId = rtpId

			flux.Send(bullet.ActionBulletBonusCall, Service.HashId(hitBullet.SecWebSocketKey), bullet.Service.HashId(hitBullet.SecWebSocketKey), bonusBullet, accountingSn)
		}

	case 28:
		if hitResult.TriggerIconId == 28 && len(hitResult.BonusPayload.([]int)) > 0 {
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

	case 29:
		if len(hitResult.ExtraData) != 2 {
			// no hit
			return
		}

		// client don't want to get hitResult.Pay when hit slot
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
				29,
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

func psf_on_00001_TriggerBonus(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
	switch hitResult.TriggerIconId {
	case 31:
		redEnvelopeResult := probability.Service.Calc(
			hitFish.GameId,
			hitFish.MathModuleId,
			strconv.Itoa(hitResult.BonusTypeId), // decided to pick high/low RedEnvelope group
			31,
			-1,
			0,
			0,
		)

		bonusRedEnvelope := redenvelope.New(
			hitResult.BonusUuid,
			hitBullet.SecWebSocketKey,
			hitFish.MathModuleId,
			hitFish.RoomUuid,
			hitFish.SeatId,
			31,
			redEnvelopeResult.BonusPayload.([]int),
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
	}
}
