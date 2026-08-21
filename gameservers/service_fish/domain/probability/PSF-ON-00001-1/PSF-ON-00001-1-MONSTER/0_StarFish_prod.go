//go:build prod
// +build prod

package PSF_ON_00001_1_MONSTER

func (s *starFish) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	if MONSTER.isHit(math.HitWeight) {
		if MONSTER.isHit(math.TriggerWeight) {
			return math.IconPays[0], math.TriggerIconID
		}
		return math.IconPays[0], -1
	}
	return 0, -1
}

func (s *starFish) HitFs(math *FsFishMath) (iconPays int) {
	return math.IconPays[0]
}
