package lottery

import (
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	PSFM_00013_93_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-93-1"
	PSFM_00013_94_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-94-1"
	PSFM_00013_95_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-95-1"
	PSFM_00013_96_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-96-1"
	PSFM_00013_97_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-97-1"
	PSFM_00013_98_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-98-1"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"serve/service_fish/models"
	"strconv"

	"serve/fish_comm/flux"
)

const (
	RKF_H5_00001_Yggdrasil_TypeId              = 100
	RKF_H5_00001_HalonSlot_TypeId              = 101
	RKF_H5_00001_RandomTriggerYggdrasil_TypeId = 102
	RKF_H5_00001_ThunderHarmmer_TypeId         = 201
	RKF_H5_00001_StormBreaker_TypeId           = 202
)

func rkf_h5_00001_Check(bulletTypeId, fishTypeId int32) bool {
	if bulletTypeId == RKF_H5_00001_ThunderHarmmer_TypeId && fishTypeId == RKF_H5_00001_StormBreaker_TypeId {
		return false
	}

	if bulletTypeId == RKF_H5_00001_StormBreaker_TypeId && fishTypeId == RKF_H5_00001_ThunderHarmmer_TypeId {
		return false
	}

	return true
}

func rkf_h5_00001_Fish(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, accountingSn uint64) (isDied bool) {
	switch hitResult.FishTypeId {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 300, 301:
		if hitResult.Pay > 0 {
			return true
		}

	case 501, 502, 503, 504, 505:
		// do nothing

	case RKF_H5_00001_ThunderHarmmer_TypeId, RKF_H5_00001_StormBreaker_TypeId:
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
			case models.PSFM_00013_93_1:
				if hitResult.FishTypeId == RKF_H5_00001_ThunderHarmmer_TypeId {
					rtpId = PSFM_00013_93_1.RTP93BS.RKF_H5_00001_1_BsMath.Icons.ThunderHammer.UseRTP
				} else {
					rtpId = PSFM_00013_93_1.RTP93BS.RKF_H5_00001_1_BsMath.Icons.Stormbreaker.UseRTP
				}

			case models.PSFM_00013_94_1:
				if hitResult.FishTypeId == RKF_H5_00001_ThunderHarmmer_TypeId {
					rtpId = PSFM_00013_94_1.RTP94BS.RKF_H5_00001_1_BsMath.Icons.ThunderHammer.UseRTP
				} else {
					rtpId = PSFM_00013_94_1.RTP94BS.RKF_H5_00001_1_BsMath.Icons.Stormbreaker.UseRTP
				}

			case models.PSFM_00013_95_1:
				if hitResult.FishTypeId == RKF_H5_00001_ThunderHarmmer_TypeId {
					rtpId = PSFM_00013_95_1.RTP95BS.RKF_H5_00001_1_BsMath.Icons.ThunderHammer.UseRTP
				} else {
					rtpId = PSFM_00013_95_1.RTP95BS.RKF_H5_00001_1_BsMath.Icons.Stormbreaker.UseRTP
				}

			case models.PSFM_00013_96_1:
				if hitResult.FishTypeId == RKF_H5_00001_ThunderHarmmer_TypeId {
					rtpId = PSFM_00013_96_1.RTP96BS.RKF_H5_00001_1_BsMath.Icons.ThunderHammer.UseRTP
				} else {
					rtpId = PSFM_00013_96_1.RTP96BS.RKF_H5_00001_1_BsMath.Icons.Stormbreaker.UseRTP
				}

			case models.PSFM_00013_97_1:
				if hitResult.FishTypeId == RKF_H5_00001_ThunderHarmmer_TypeId {
					rtpId = PSFM_00013_97_1.RTP97BS.RKF_H5_00001_1_BsMath.Icons.ThunderHammer.UseRTP
				} else {
					rtpId = PSFM_00013_97_1.RTP97BS.RKF_H5_00001_1_BsMath.Icons.Stormbreaker.UseRTP
				}

			case models.PSFM_00013_98_1:
				if hitResult.FishTypeId == RKF_H5_00001_ThunderHarmmer_TypeId {
					rtpId = PSFM_00013_98_1.RTP98BS.RKF_H5_00001_1_BsMath.Icons.ThunderHammer.UseRTP
				} else {
					rtpId = PSFM_00013_98_1.RTP98BS.RKF_H5_00001_1_BsMath.Icons.Stormbreaker.UseRTP
				}
			}

			bonusBullet.RtpId = rtpId

			flux.Send(bullet.ActionBulletBonusCall, Service.HashId(hitBullet.SecWebSocketKey), bullet.Service.HashId(hitBullet.SecWebSocketKey), bonusBullet, accountingSn)
		}

	case RKF_H5_00001_Yggdrasil_TypeId:
		if hitResult.TriggerIconId == RKF_H5_00001_Yggdrasil_TypeId && len(hitResult.BonusPayload.([]int)) > 0 {
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

	case RKF_H5_00001_HalonSlot_TypeId:
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
				RKF_H5_00001_HalonSlot_TypeId,
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

func rkf_h5_00001_TriggerBonus(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
	switch hitResult.TriggerIconId {
	case RKF_H5_00001_RandomTriggerYggdrasil_TypeId:
		redEnvelopeResult := probability.Service.Calc(
			hitFish.GameId,
			hitFish.MathModuleId,
			strconv.Itoa(hitResult.BonusTypeId),
			RKF_H5_00001_RandomTriggerYggdrasil_TypeId,
			-1,
			0,
			hitBullet.Bet*hitBullet.Rate,
		)

		bonusRedEnvelope := redenvelope.New(
			hitResult.BonusUuid,
			hitBullet.SecWebSocketKey,
			hitFish.MathModuleId,
			hitFish.RoomUuid,
			hitFish.SeatId,
			RKF_H5_00001_RandomTriggerYggdrasil_TypeId,
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

func rkf_h5_00001_ExtraTriggerBonus(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
	if hitResult.ExtraTriggerBonus != nil {
		switch hitResult.TriggerIconId {
		case 201, 202:
			for _, v := range hitResult.ExtraTriggerBonus {
				v.Pay = hitResult.Pay
				rkf_h5_00001_TriggerBonus(hitFish, hitBullet, v)
			}
		}
	}
}
