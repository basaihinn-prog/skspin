package PSF_ON_00003_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var Zombie1Slot = &zombie1Slot{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "006"),
}

type zombie1Slot struct {
	bonusMode bool
}

type BsZombieSlotMath struct {
	IconID     int                           `json:"IconID"`
	BonusID    string                        `json:"BonusID"`
	BonusTimes []int                         `json:"BonusTimes"`
	IconPays   []int                         `json:"IconPays"`
	RTP1       struct{ HitBsZombieSlotMath } `json:"RTP1"`
	RTP2       struct{ HitBsZombieSlotMath } `json:"RTP2"`
	RTP3       struct{ HitBsZombieSlotMath } `json:"RTP3"`
	RTP4       struct{ HitBsZombieSlotMath } `json:"RTP4"`
	RTP5       struct{ HitBsZombieSlotMath } `json:"RTP5"`
	RTP6       struct{ HitBsZombieSlotMath } `json:"RTP6"`
	RTP7       struct{ HitBsZombieSlotMath } `json:"RTP7"`
	RTP8       struct{ HitBsZombieSlotMath } `json:"RTP8"`
	RTP9       struct{ HitBsZombieSlotMath } `json:"RTP9"`
	RTP10      struct{ HitBsZombieSlotMath } `json:"RTP10"`
	RTP11      struct{ HitBsZombieSlotMath } `json:"RTP11"`
}

type HitBsZombieSlotMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
}

type FsStripTableMath struct {
	Normal struct{ HitFsStripTableMath }
	Block  struct{ HitFsStripTableMath }
}

type HitFsStripTableMath struct {
	UseIcon   []int
	HitWeight []int
}

type FsSlotMath struct {
	MixAwardId        int
	IconLinkThreshold int
	Pay               []int
	MixIconGroup      []struct {
		IconId   int
		Quantity int
	}
	WildGroup []interface{}
}

type reel struct {
	previous *reel
	next     *reel
	index    int
	iconId   int
}

func (z *zombie1Slot) Hit(rtpId string,
	bsMath *BsZombieSlotMath,
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	slotChangerMath1 *FsSlotChanger1Math,
	littleZombie1, littleZombie2, littleZombie3, littleZombie4,
	zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
	demon31, demon32, demon33 *FsZombieMath,
	bullet *BsBullet,
) (fishPay, iconPay int, iconIds []int32, iconPays []int64, bullets int) {
	bullets = -1

	switch rtpId {
	case bsMath.RTP1.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP1, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP2.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP2, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP3.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP3, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP4.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP4, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP5.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP5, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP6.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP6, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP7.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP7, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP8.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP8, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP9.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP9, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP10.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP10, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP11.ID:
		v1, v2, v3, v4, v5 := z.rtpHit(
			bsMath, &bsMath.RTP11, fsMath, fsStripMath,
			slotChangerMath1,
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
			bullet,
		)
		return v1, v2, v3, v4, v5

	default:
		return 0, 0, nil, nil, bullets
	}
}

func (z *zombie1Slot) hit(
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	slotChangerMath1 *FsSlotChanger1Math,
	littleZombie1, littleZombie2, littleZombie3, littleZombie4,
	zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
	demon31, demon32, demon33 *FsZombieMath,
) (iconPay int, iconIds []int32, iconPays []int64) {
	blockIndex := z.rngBlockReelIndex()

	iconIds = make([]int32, 0)
	iconPays = make([]int64, 3)

	for i := 0; i < fsMath.MixAwardId; i++ {
		if i == blockIndex {
			v1, v2 := z.rngReel(fsMath, &fsStripMath.Block, &slotChangerMath1.ChangeCandidateBlock)

			for _, v := range v1 {
				iconIds = append(iconIds, v)
			}
			iconPays[i] = int64(v2)
		} else {
			v1, v2 := z.rngReel(fsMath, &fsStripMath.Normal, &slotChangerMath1.ChangeCandidateNormal)

			for _, v := range v1 {
				iconIds = append(iconIds, v)
			}
			iconPays[i] = int64(v2)
		}
	}

	// Big Win
	if iconIds[1] == iconIds[4] && iconIds[1] == iconIds[7] {
		bonusPay := z.bonusPay(iconIds[1],
			littleZombie1, littleZombie2, littleZombie3, littleZombie4,
			zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
			demon31, demon32, demon33,
		)

		if z.contain(fsMath, bonusPay) {
			return bonusPay, iconIds, iconPays
		}

		return 0, nil, nil
	}

	if z.contain(fsMath, int(iconPays[0]+iconPays[1]+iconPays[2])) {
		return int(iconPays[0] + iconPays[1] + iconPays[2]), iconIds, iconPays
	}

	return 0, nil, nil
}

func (z *zombie1Slot) rngBlockReelIndex() int {
	var indexOptions []rng.Option
	for i := 0; i < 3; i++ {
		indexOptions = append(indexOptions, rng.Option{1, i})
	}

	return MONSTER.rng(indexOptions).(int)
}

func (z *zombie1Slot) rngReel(
	fsMath *FsSlotMath,
	fsStripMath *struct{ HitFsStripTableMath },
	change1Math *struct{ FsChangeCandidate1Math },
) (iconIds []int32, pay int) {
	r := z.rngReels(fsStripMath)

	r.previous.iconId = z.rngCandidate(r.previous.iconId, change1Math)
	r.iconId = z.rngCandidate(r.iconId, change1Math)
	r.next.iconId = z.rngCandidate(r.next.iconId, change1Math)

	iconIds = append(iconIds, int32(r.previous.iconId))
	iconIds = append(iconIds, int32(r.iconId))
	iconIds = append(iconIds, int32(r.next.iconId))

	for _, v := range fsMath.MixIconGroup {
		if v.IconId == r.iconId {
			return iconIds, v.Quantity
		}
	}

	return nil, 0
}

func (z *zombie1Slot) rngCandidate(iconId int, change1Math *struct{ FsChangeCandidate1Math }) int {
	newIconId := iconId

	for {
		iconId = newIconId

		if iconId == 26 {
			newIconId = SlotChanger.rngCandidate1(change1Math)
		}

		if newIconId != 26 {
			return newIconId
		}
	}

	return iconId
}

func (z *zombie1Slot) rngReels(reels *struct{ HitFsStripTableMath }) *reel {
	var reelOptions []rng.Option
	size := len(reels.UseIcon)

	for i := 0; i < size; i++ {
		p := &reel{
			previous: nil,
			next:     nil,
		}

		if i-1 < 0 {
			p.index = i - 1
			p.iconId = reels.UseIcon[size-1]
		} else {
			p.index = i - 1
			p.iconId = reels.UseIcon[i-1]
		}

		n := &reel{
			previous: nil,
			next:     nil,
		}

		if i+1 >= size {
			n.index = 0
			n.iconId = reels.UseIcon[0]
		} else {
			n.index = i + 1
			n.iconId = reels.UseIcon[i+1]
		}

		r := &reel{
			previous: p,
			next:     n,
			index:    i,
			iconId:   reels.UseIcon[i],
		}

		reelOptions = append(reelOptions, rng.Option{reels.HitWeight[i], r})
	}

	return MONSTER.rng(reelOptions).(*reel)
}

func (z *zombie1Slot) bonusPay(iconId int32,
	littleZombie1, littleZombie2, littleZombie3, littleZombie4,
	zombie1, zombie2, zombie3, zombie4, zombie5, zombie6,
	demon31, demon32, demon33 *FsZombieMath,
) (pay int) {
	switch iconId {
	case 3031:
		return LittleZombie.HitFs(littleZombie1)
	case 3032:
		return LittleZombie.HitFs(littleZombie2)
	case 3033:
		return LittleZombie.HitFs(littleZombie3)
	case 3034:
		return LittleZombie.HitFs(littleZombie4)
	case 0:
		return Zombie1.HitFs(zombie1)
	case 1:
		return Zombie2.HitFs(zombie2)
	case 2:
		return Zombie3.HitFs(zombie3)
	case 3:
		return Zombie4.HitFs(zombie4)
	case 4:
		return Zombie1.HitFs(zombie5)
	case 5:
		return Zombie1.HitFs(zombie6)
	case 4031:
		return Demon3.HitFs(demon31)
	case 4032:
		return Demon3.HitFs(demon32)
	case 4033:
		return Demon3.HitFs(demon33)
	default:
		return 0
	}
}

func (z *zombie1Slot) contain(fsMath *FsSlotMath, pay int) bool {
	for _, v := range fsMath.Pay {
		if v == pay {
			return true
		}
	}

	return false
}

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
