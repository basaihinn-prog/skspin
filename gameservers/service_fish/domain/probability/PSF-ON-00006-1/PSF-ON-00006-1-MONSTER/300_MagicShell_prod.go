//go:build prod
// +build prod

package PSF_ON_00006_1_MONSTER

func (p *magicShell) rtpHit(math *BsMagicShell, rtp *struct{ HitBsMagicShell }) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
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
