//go:build prod
// +build prod

package PSF_ON_00004_1_MONSTER

func (l *lobsterDash) rtpHit(math *struct{ HitBsDashMath }) (avgPay, multiplier int) {
	avgPay = int(math.AvgPay[0])
	multiplier = 0

	if MONSTER.isHit(math.HitWeight) {
		multiplier = l.pick(math)
	}

	return avgPay, multiplier
}
