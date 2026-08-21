package accountingmanager

import (
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/probability"

	"serve/fish_comm/flux"
)

func rkf_h5_00001_HitBullet(action *flux.Action, am *AccountingManager) {
	psf_on_00004_HitBullet(action, am)
}

func rkf_h5_00001_Bonus(action *flux.Action, am *AccountingManager) {
	psf_on_00004_Bonus(action, am)
}

func rkf_h5_00001_RefundBonusBullet(action *flux.Action, am *AccountingManager) {
	psf_on_00004_RefundBonusBullet(action, am)
}

func rkf_h5_00001_CheckBulletTypeId(bulletId, restBullets int32,
	am *AccountingManager,
	b *bullet.Bullet,
	machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
) {
	psf_on_00004_CheckBulletTypeId(bulletId, restBullets, am, b, machineGunUuid, superMachineGunUuid)
}

func rkf_h5_00001_BulletHit(bulletId int32,
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

func rkf_h5_00001_BulletHit_TriggerId(triggerIconId int,
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

func rkf_h5_00001_BulletBonus(bonusBulletId int32, bonusBullet *bullet.Bullet, machineGunUuid, superMachineGunUuid map[string]map[uint64]int32, bonusBulletShooted map[string]map[string]map[uint64]int32) {
	psf_on_00004_BulletBonus(bonusBulletId, bonusBullet, machineGunUuid, superMachineGunUuid, bonusBulletShooted)
}
