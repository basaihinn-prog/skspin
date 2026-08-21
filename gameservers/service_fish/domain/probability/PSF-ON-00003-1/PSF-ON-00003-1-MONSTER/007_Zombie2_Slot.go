package PSF_ON_00003_1_MONSTER

import (
	"os"
	"strings"
)

var Zombie2Slot = &zombie2Slot{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "007"),
}

type zombie2Slot struct {
	bonusMode bool
}

func (z *zombie2Slot) Hit(rtpId string,
	bsMath *BsZombieSlotMath,
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	slotChangerMath1 *FsSlotChanger1Math,
	littleZombie1, littleZombie2, littleZombie3, littleZombie4,
	zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
	demon31, demon32, demon33 *FsZombieMath,
	bullet *BsBullet,
) (fishPay, iconPay int, iconIds []int32, iconPays []int64, bullets int) {
	return Zombie1Slot.Hit(
		rtpId, bsMath, fsMath, fsStripMath, slotChangerMath1,
		littleZombie1, littleZombie2, littleZombie3, littleZombie4,
		zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
		demon31, demon32, demon33,
		bullet,
	)
}
