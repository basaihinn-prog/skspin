//go:build prod
// +build prod

package PSF_ON_00007_1_MONSTER

func (o *oceanKing) rtpHit(rtp *struct{ HitBsOceanKingMath }) (iconPay, triggerIconId, bonusTypeId int) {
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
