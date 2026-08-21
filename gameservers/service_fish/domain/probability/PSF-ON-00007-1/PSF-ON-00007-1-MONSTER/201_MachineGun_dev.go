//go:build dev
// +build dev

package PSF_ON_00007_1_MONSTER

func (m *machineGun) rtpHit(iconPays []int,
	math *BsMachineGunMath, rtp *struct{ HitBsMachineGunMath }) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	bonusTimes = 0

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = iconPays[0]
		bonusTimes = m.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, bonusTimes
}
