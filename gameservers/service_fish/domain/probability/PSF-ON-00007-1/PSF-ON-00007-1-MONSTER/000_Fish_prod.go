//go:build prod
// +build prod

package PSF_ON_00007_1_MONSTER

func (f *fish0) rtpHit(iconPays int, bsFish *struct{ HitBsFishMath }) (iconPay, triggerIconId, bonusTypeId int) {
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
