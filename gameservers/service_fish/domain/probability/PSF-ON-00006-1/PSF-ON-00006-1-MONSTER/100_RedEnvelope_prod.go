//go:build prod
// +build prod

package PSF_ON_00006_1_MONSTER

func (m *redEnvelope) rtpHit(bsMath *BsRedEnvelopeMath, rtp *struct{ HitBsRedEnvelopeMath }) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1
	iconPays = nil

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPays = m.pick(bsMath)
	}

	return triggerIconId, bonusTypeId, iconPays
}
