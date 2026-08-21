//go:build dev
// +build dev

package RKF_H5_00001_1_MONSTER

func (t *thunderHammer) rtpHit(iconPays []int,
	math *BsThunderHammerMath, rtp *struct{ HitBsThunderHammerMath }) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
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
		bonusTimes = t.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, bonusTimes
}
