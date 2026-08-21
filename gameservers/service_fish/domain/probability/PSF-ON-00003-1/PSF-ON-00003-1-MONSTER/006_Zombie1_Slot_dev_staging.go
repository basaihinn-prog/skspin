//go:build dev || staging
// +build dev staging

package PSF_ON_00003_1_MONSTER

func (z *zombie1Slot) rtpHit(
	bsMath *BsZombieSlotMath,
	rtp *struct{ HitBsZombieSlotMath },
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	slotChangerMath1 *FsSlotChanger1Math,
	littleZombie1, littleZombie2, littleZombie3, littleZombie4,
	zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
	demon31, demon32, demon33 *FsZombieMath,
	bulletMath *BsBullet,
) (fishPay, iconPay int, iconIds []int32, iconPays []int64, bullets int) {
	fishPay = 0
	iconPay = 0
	iconIds = nil
	iconPays = nil

	bullets = 0

	if MONSTER.isHit(rtp.HitWeight) {
		v1, v2, v3 := z.hit(
			fsMath,
			fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
		)
		v0 := bsMath.IconPays[0]

		fishPay = v0
		iconPay = v1
		iconIds = v2
		iconPays = v3
	}

	if MONSTER.isHit(rtp.TriggerWeight) {
		bullets = Bullet.RngBonusBullets(bulletMath)
	}

	return fishPay, iconPay, iconIds, iconPays, bullets
}
