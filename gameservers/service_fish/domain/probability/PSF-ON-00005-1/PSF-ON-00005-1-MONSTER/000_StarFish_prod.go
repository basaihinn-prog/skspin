//go:build prod
// +build prod

package PSF_ON_00005_1_MONSTER

func (s *starfish) rtpHit(iconPays int, bsFish *struct{ HitBsFishMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(bsFish.HitWeight) {
		iconPay = iconPays
	}

	if MONSTER.isHit(bsFish.TriggerWeight) {
		triggerIconId = bsFish.TriggerIconID
		bonusTypeId = bsFish.Type
	}

	return iconPay, triggerIconId, bonusTypeId
}
