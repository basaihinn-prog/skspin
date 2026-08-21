package accountingmanager

import (
	"fmt"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"serve/service_fish/domain/receipt"
	"serve/service_fish/domain/slot"

	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

func psf_on_00006_HitBullet(action *flux.Action, am *AccountingManager) {
	psf_on_00004_HitBullet(action, am)
}

func psf_on_00006_Bonus(action *flux.Action, am *AccountingManager) {
	secWebSocketKey := action.Payload()[0].(string)
	accountingSn := am.accountingSn

	switch action.Payload()[1].(type) {
	case *slot.Slot:
		s := action.Payload()[1].(*slot.Slot)

		hitFish := s.ExtraData[0].(*fish.Fish)
		hitBullet := s.ExtraData[1].(*bullet.Bullet)
		hitResult := s.ExtraData[2].(*probability.Probability)

		gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			hitBullet.TypeId,
			hitBullet.Bullets,
			hitBullet.UsedBullets,
			hitFish.TypeId,
			hitFish.Scene,
			hitResult.Multiplier,
		)

		bet := s.Bet * s.Rate

		if hitBullet.TypeId == 201 || hitBullet.TypeId == 202 {
			bet = 0
		}

		gameData := fmt.Sprintf("%d_%d/%d_%d/%d_%d/%d_%d",
			s.AllPay[0], s.Reels[1],
			s.AllPay[1], s.Reels[4],
			s.AllPay[2], s.Reels[7],
			s.AllPay[3], s.Reels[10],
		)

		go receipt.Service.New(
			s.Uuid,
			secWebSocketKey,
			am.hostExtId,
			gameResult,
			gameData,
			accountingSn,
			bet,
			s.Pay*s.Bet*s.Rate,
		)

		logger.Service.Zap.Infow("Slot Accounting",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
		)

	default:
		psf_on_00004_Bonus(action, am)
	}
}

func psf_on_00006_RefundBonusBullet(action *flux.Action, am *AccountingManager) {
	psf_on_00004_RefundBonusBullet(action, am)
}

func psf_on_00006_CheckBulletTypeId(bulletId, restBullets int32,
	am *AccountingManager,
	b *bullet.Bullet,
	machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
) {
	psf_on_00004_CheckBulletTypeId(bulletId, restBullets, am, b, machineGunUuid, superMachineGunUuid)
}

func psf_on_00006_BulletHit(bulletId int32,
	bulletUuid map[string]uint64,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	tempReceipts map[uint64]*tempReceipt,
	machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	isBonusBulletShooted bool,
) {
	psf_on_00004_BulletHit(bulletId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, isBonusBulletShooted)
}

func psf_on_00006_BulletHit_TriggerId(triggerIconId int,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	hitResult *probability.Probability,
	tempReceipts map[uint64]*tempReceipt,
	machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	redEnvelopeUuid, slotUuid map[string]uint64,
) bool {
	return psf_on_00004_BulletHit_TriggerId(triggerIconId, am, hitBullet, hitResult, tempReceipts,
		machineGunUuid, superMachineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
	)
}

func psf_on_00006_BulletBonus(bonusBulletId int32, bonusBullet *bullet.Bullet, machineGunUuid, superMachineGunUuid map[string]map[uint64]int32, bonusBulletShooted map[string]map[string]map[uint64]int32) {
	psf_on_00004_BulletBonus(bonusBulletId, bonusBullet, machineGunUuid, superMachineGunUuid, bonusBulletShooted)
}
