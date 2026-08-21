//go:build dev || staging
// +build dev staging

package PSF_ON_00003_1_MONSTER

func (l *littleZombie) rtpHit(rtp *struct{ HitBsLittleZombieMath }) (iconPay, iconBullet int) {
	if MONSTER.isHit(rtp.HitWeight) {
		payOrBullet := l.pick(rtp)

		switch payOrBullet {
		case 0: // IconPay
			iconPay = l.pickPay(rtp)
		case 1: // Bullets
			iconBullet = l.pickBullet(rtp)
		}
	}

	return iconPay, iconBullet
}
