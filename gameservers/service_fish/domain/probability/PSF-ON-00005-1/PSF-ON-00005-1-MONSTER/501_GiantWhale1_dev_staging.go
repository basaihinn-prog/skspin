//go:build dev || staging
// +build dev staging

package PSF_ON_00005_1_MONSTER

func (g *giantWhale1) rtpHit(math *BsGiantWhaleMath, rtp *struct{ HitBsGiantWhaleMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = g.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId
}
