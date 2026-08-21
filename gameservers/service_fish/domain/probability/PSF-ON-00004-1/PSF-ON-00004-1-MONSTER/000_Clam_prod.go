//go:build prod
// +build prod

package PSF_ON_00004_1_MONSTER

func (c *clam) rtpHit(iconPays int, bsFish *struct{ HitBsFishMath }) (iconPay int) {
	if MONSTER.isHit(bsFish.HitWeight) {
		iconPay = iconPays
	}

	return iconPay
}
