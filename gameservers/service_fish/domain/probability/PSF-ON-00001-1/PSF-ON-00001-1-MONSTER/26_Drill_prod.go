//go:build prod
// +build prod

package PSF_ON_00001_1_MONSTER

func (d *drill) Hit(math *BsDrillMath) (iconPay, bonusTimes int) {
	if MONSTER.isHit(math.HitWeight) {
		return math.IconPays[0], d.pick(math)
	}
	return 0, 0
}
