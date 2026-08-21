//go:build prod
// +build prod

package PSF_ON_00004_1_MONSTER

func (r *redEnvelope) rtpHit(hitWeight []int, rtp *struct{ HitBsRedEnvelopeMath }) (iconPays []int, avgPay int) {
	avgPay = int(rtp.AvgPay[0])

	if MONSTER.isHit(hitWeight) {
		return r.pick(rtp), avgPay
	}

	return nil, avgPay
}
