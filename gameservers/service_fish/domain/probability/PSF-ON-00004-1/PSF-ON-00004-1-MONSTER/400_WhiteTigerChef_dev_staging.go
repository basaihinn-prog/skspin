//go:build dev || staging
// +build dev staging

package PSF_ON_00004_1_MONSTER

func (w *whiteTigerChef) rtpHit(math *struct{ HitBsWhiteTigerChefMath }) (iconPay, avgPay int) {
	iconPay = 0
	avgPay = int(math.AvgPay[0])

	if MONSTER.isHit(math.HitWeight) {
		iconPay = w.pick(math)
	}

	return iconPay, avgPay
}
