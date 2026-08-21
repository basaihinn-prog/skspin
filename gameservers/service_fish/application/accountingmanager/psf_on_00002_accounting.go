package accountingmanager

import (
	"fmt"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"serve/service_fish/domain/receipt"
	"serve/service_fish/domain/redenvelope"

	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

func psf_on_00002_HitBullet(action *flux.Action, am *AccountingManager) {
	psf_on_00001_HitBullet(action, am)
}

func psf_on_00002_Bonus(action *flux.Action, am *AccountingManager) {
	if am.guestEnable {
		return
	}

	secWebSocketKey := action.Payload()[0].(string)
	accountingSn := am.accountingSn

	switch action.Payload()[1].(type) {
	case *redenvelope.RedEnvelope:
		re := action.Payload()[1].(*redenvelope.RedEnvelope)

		hitFish := re.ExtraData[0].(*fish.Fish)
		hitBullet := re.ExtraData[1].(*bullet.Bullet)
		hitResult := re.ExtraData[2].(*probability.Probability)

		gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			hitBullet.TypeId,
			hitBullet.Bullets,
			hitBullet.UsedBullets,
			hitFish.TypeId,
			hitFish.Scene,
			hitResult.Multiplier,
		)

		bet := re.Bet * re.Rate

		// hit by bullet 26,27
		if hitResult.TriggerIconId == 28 || hitResult.TriggerIconId == 31 || hitBullet.TypeId == 26 || hitBullet.TypeId == 27 {
			bet = 0
		}

		// hit 26,27 and trigger 31 at the same time
		if (hitResult.FishTypeId == 26 || hitResult.FishTypeId == 27) && hitResult.Pay > 0 {
			bet = 0
			hitResult.Pay = 0 // reset to 0, see f0dc137f
		}

		gameData := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			re.PlayerOptIndex, re.AllPay[0], re.AllPay[1], re.AllPay[2], re.AllPay[3], re.AllPay[4])

		go receipt.Service.New(
			re.Uuid,
			secWebSocketKey,
			am.hostExtId,
			gameResult,
			gameData,
			accountingSn,
			bet,
			re.Pay*re.Bet*re.Rate,
		)

		logger.Service.Zap.Infow("RedEnvelop Accounting",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
		)

	default:
		psf_on_00001_Bonus(action, am)
	}
}

func psf_on_00002_RefundBonusBullet(action *flux.Action, am *AccountingManager) {
	psf_on_00001_RefundBonusBullet(action, am)
}

func psf_on_00002_CheckBulletTypeId(bulletId, restBullets int32,
	am *AccountingManager,
	b *bullet.Bullet,
	drillUuid, machineGunUuid map[string]map[uint64]int32,
) {
	psf_on_00001_CheckBulletTypeId(bulletId, restBullets, am, b, drillUuid, machineGunUuid)
}

func psf_on_00002_BulletHit(bulletId int32,
	bulletUuid map[string]uint64,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	tempReceipts map[uint64]*tempReceipt,
	drillUuid, machineGunUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	isBonusBulletShooted bool,
) {
	psf_on_00001_BulletHit(bulletId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, machineGunUuid, bonusBulletShooted, isBonusBulletShooted)
}

func psf_on_00002_BulletHit_TriggerId(triggerIconId int,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	hitResult *probability.Probability,
	tempReceipts map[uint64]*tempReceipt,
	drillUuid, machineGunUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	redEnvelopeUuid, slotUuid map[string]uint64,
) bool {
	return psf_on_00001_BulletHit_TriggerId(triggerIconId, am, hitBullet, hitResult, tempReceipts,
		drillUuid, machineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
	)
}

func psf_on_00002_BulletBonus(bonusBulletId int32, bonusBullet *bullet.Bullet, machineGunUuid map[string]map[uint64]int32, bonusBulletShooted map[string]map[string]map[uint64]int32) {
	psf_on_00001_BulletBonus(bonusBulletId, bonusBullet, machineGunUuid, bonusBulletShooted)
}
