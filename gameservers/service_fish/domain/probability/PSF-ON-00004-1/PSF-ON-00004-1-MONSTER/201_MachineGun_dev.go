//go:build dev
// +build dev

package PSF_ON_00004_1_MONSTER

func (m *machineGun) rtpHit(iconPays []int, rtp *struct{ HitBsMachineGunMath }) (iconPay, bonusTimes, avgPay int) {
	avgPay = int(rtp.AvgPay[0])

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = iconPays[0]
		bonusTimes = m.pick(rtp)
	}

	return iconPay, bonusTimes, avgPay
}
