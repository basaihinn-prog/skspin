//go:build prod
// +build prod

package PSF_ON_00001_1_MONSTER

func (m *machineGun) Hit(math *BsMachineGunMath) (iconPay, bonusTimes int) {
	if MONSTER.isHit(math.HitWeight) {
		return math.IconPays[0], m.pick(math)
	}
	return 0, 0
}
