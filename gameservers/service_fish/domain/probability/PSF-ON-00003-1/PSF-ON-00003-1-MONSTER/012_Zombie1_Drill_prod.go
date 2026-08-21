//go:build prod
// +build prod

package PSF_ON_00003_1_MONSTER

func (z *zombie1Drill) rtpHit(math *BsZombieDrillMath, rtp *struct{ HitBsZombieDrillMath }, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = math.IconPays[0]
		bonusTimes = z.pick(math)
	}

	if MONSTER.isHit(rtp.TriggerWeight) {
		bullets = Bullet.RngBonusBullets(bulletMath)
	}

	return iconPay, bonusTimes, bullets
}
