//go:build dev || staging
// +build dev staging

package RKF_H5_00001_1_MONSTER

func (y *yggdrasil) rtpHit(bsMath *BsYggdrasilMath, rtp *struct{ HitBsYggdrasilMath }) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1
	iconPays = nil

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPays = y.pick(bsMath)
	}

	return triggerIconId, bonusTypeId, iconPays
}
