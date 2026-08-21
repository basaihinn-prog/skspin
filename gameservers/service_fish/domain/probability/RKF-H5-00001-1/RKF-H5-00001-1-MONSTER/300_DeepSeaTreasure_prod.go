//go:build prod
// +build prod

package RKF_H5_00001_1_MONSTER

func (d *deepSeaTreasure) rtpHit(math *BsDeepSeaTreasureMath, rtp *struct{ HitBsDeepSeaTreasureMath }) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	multiplier = 0

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = 1
		multiplier = d.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}
