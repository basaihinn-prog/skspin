//go:build prod
// +build prod

package PSF_ON_00004_1_MONSTER

func (x *xiaoLongBao) rtpHit(rtp *struct{ HitBsXiaoLongBaoMath }) (totalPay int, iconPays, iconPaysPick []int, avgPay int) {
	totalPay = 0
	avgPay = rtp.AvgPay[0]

	if MONSTER.isHit(rtp.HitWeight) {
		gainAmount := x.rngGain(rtp)
		iconPays = x.rngIconPay(gainAmount, rtp)
		iconPaysPick = x.rngPick(gainAmount)
	}

	if len(iconPays) > 0 {
		for i := 0; i < len(iconPays); i++ {
			if iconPaysPick[i] == 1 {
				totalPay += iconPays[i]
			}
		}
	}

	return totalPay, iconPays, iconPaysPick, avgPay
}
