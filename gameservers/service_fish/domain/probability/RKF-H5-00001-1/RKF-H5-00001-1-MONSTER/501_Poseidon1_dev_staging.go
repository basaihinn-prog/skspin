//go:build dev || staging
// +build dev staging

package RKF_H5_00001_1_MONSTER

func (p *poseidon1) rtpHit(math *BsPoseidonMath, rtp *struct{ HitBsPoseidonMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = p.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId
}
