//go:build prod
// +build prod

package PSF_ON_00005_1_MONSTER

func (r *redEnvelopeTree) rtpHit(bsMath *BsRedEnvelopeTreeMath, rtp *struct{ HitBsRedEnvelopeTreeMath }) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1
	iconPays = nil

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPays = r.pick(bsMath)
	}

	return triggerIconId, bonusTypeId, iconPays
}
