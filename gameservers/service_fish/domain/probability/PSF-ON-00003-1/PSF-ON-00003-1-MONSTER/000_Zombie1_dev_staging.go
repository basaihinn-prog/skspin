//go:build dev || staging
// +build dev staging

package PSF_ON_00003_1_MONSTER

func (z *zombie1) rtpHit(iconPays int, bsZombie *struct{ HitBsZombieMath }, bulletMath *BsBullet) (iconPay, bullets int) {
	if MONSTER.isHit(bsZombie.HitWeight) {
		iconPay = iconPays
	}

	if MONSTER.isHit(bsZombie.TriggerWeight) {
		bullets = Bullet.RngBonusBullets(bulletMath)
	}

	return iconPay, bullets
}
