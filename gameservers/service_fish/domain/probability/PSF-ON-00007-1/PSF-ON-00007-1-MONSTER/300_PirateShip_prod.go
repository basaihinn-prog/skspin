//go:build prod
// +build prod

package PSF_ON_00007_1_MONSTER

func (p *pirateShip) rtpHit(math *BsPirateShip, rtp *struct{ HitBsPirateShip }) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
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
		multiplier = p.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}
