//go:build prod
// +build prod

package PSF_ON_00007_1_MONSTER

func (m *magicShellFaFaFa) rtpHit(bsMath *BsMagicShellFaFaFaMath, rtp *struct{ HitBsMagicShellFaFaFaMath }) (triggerIconId, bonusTypeId int, iconPays []int) {
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
