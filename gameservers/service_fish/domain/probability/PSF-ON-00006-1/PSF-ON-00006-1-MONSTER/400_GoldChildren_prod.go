//go:build prod
// +build prod

package PSF_ON_00006_1_MONSTER

func (o *goldChildren) rtpHit(rtp *struct{ HitBsGoldChildrenMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = o.pick(rtp)
	}

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	return iconPay, triggerIconId, bonusTypeId
}
