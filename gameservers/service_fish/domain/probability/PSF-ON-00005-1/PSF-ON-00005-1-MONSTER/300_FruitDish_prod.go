//go:build prod
// +build prod

package PSF_ON_00005_1_MONSTER

func (f *fruitDish) rtpHit(math *BsFruitDish, rtp *struct{ HitBsFruitDish }) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
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
		multiplier = f.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}
