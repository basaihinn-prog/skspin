//go:build prod
// +build prod

package PSF_ON_00006_1_MONSTER

func (t *treasureHermitCrab) rtpHit(rtp *struct{ HitBsTreasureHermitCrabMath }) (totalPay int, iconPays, iconPaysPick []int, triggerIconId, bonusTypeId int) {
	totalPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.HitWeight) {
		gainAmount := t.rngGain(rtp)
		iconPays = t.rngIconPay(gainAmount, rtp)
		iconPaysPick = t.rngPick(gainAmount)
	}

	if len(iconPays) > 0 {
		for i := 0; i < len(iconPays); i++ {
			if iconPaysPick[i] == 1 {
				totalPay += iconPays[i]
			}
		}
	}

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	return totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId
}
