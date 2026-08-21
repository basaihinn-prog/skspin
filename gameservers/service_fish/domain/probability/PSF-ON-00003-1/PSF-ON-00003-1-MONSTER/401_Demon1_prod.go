//go:build prod
// +build prod

package PSF_ON_00003_1_MONSTER

func (d *demon1) rtpHit(math *struct{ HitBsDemonMath }) (iconPay int) {
	if MONSTER.isHit(math.HitWeight) {
		iconPay = d.pick(math)
	}

	return iconPay
}
