//go:build dev || staging
// +build dev staging

package PSF_ON_00003_1_MONSTER

func (b *bonusZombie) rtpHit(math *struct{ HitBsBonusZombieMath }) (iconBullets, triggerIconId int) {
	if MONSTER.isHit(math.HitWeight) {
		iconBullets = b.pick(math)
		triggerIconId = math.TriggerIconID
	}

	return iconBullets, triggerIconId
}
