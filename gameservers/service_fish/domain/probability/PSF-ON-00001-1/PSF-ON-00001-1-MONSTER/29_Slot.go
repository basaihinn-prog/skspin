package PSF_ON_00001_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var Slot = &slot{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "29"),
}

type slot struct {
	bonusMode bool
}

type BsSlotMath struct {
	IconID           int
	IsWildSubstitute bool
	BonusID          string
	BonusTimes       []int
	HitWeight        []int
	IconPays         []int
}

type FsReelMath struct {
	Normal struct {
		Reel []int
	} `json:"Normal"`
	Block struct {
		Reel []int
	} `json:"Block"`
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

func (s *slot) rngBlockReelIndex() int {
	var indexOptions []rng.Option
	indexOptions = append(indexOptions, rng.Option{1, 0})
	indexOptions = append(indexOptions, rng.Option{1, 1})
	indexOptions = append(indexOptions, rng.Option{1, 2})
	return MONSTER.rng(indexOptions).(int)
}

func (s *slot) rngNormalReel(fsMath *FsSlotMath, reelMath *FsReelMath, faceOffChangerMath *FsFaceOffChangerMath) (iconIds []int32, pay int) {
	r := s.rngReels(reelMath.Normal.Reel)

	if r.previous.iconId == 30 {
		r.previous.iconId = FaceOffChanger.rngNormalReel(faceOffChangerMath)
	}
	iconIds = append(iconIds, int32(r.previous.iconId))

	if r.iconId == 30 {
		r.iconId = FaceOffChanger.rngNormalReel(faceOffChangerMath)
	}
	iconIds = append(iconIds, int32(r.iconId))

	if r.next.iconId == 30 {
		r.next.iconId = FaceOffChanger.rngNormalReel(faceOffChangerMath)
	}
	iconIds = append(iconIds, int32(r.next.iconId))

	for _, v := range fsMath.MixIconGroup {
		if v.IconId == r.iconId {
			return iconIds, v.Quantity
		}
	}
	return nil, 0
}

func (s *slot) rngBlockReel(fsMath *FsSlotMath, reelMath *FsReelMath, faceOffChangerMath *FsFaceOffChangerMath) (iconIds []int32, pay int) {
	r := s.rngReels(reelMath.Block.Reel)

	if r.previous.iconId == 30 {
		r.previous.iconId = FaceOffChanger.rngBlockReel(faceOffChangerMath)
	}
	iconIds = append(iconIds, int32(r.previous.iconId))

	if r.iconId == 30 {
		r.iconId = FaceOffChanger.rngBlockReel(faceOffChangerMath)
	}
	iconIds = append(iconIds, int32(r.iconId))

	if r.next.iconId == 30 {
		r.next.iconId = FaceOffChanger.rngBlockReel(faceOffChangerMath)
	}
	iconIds = append(iconIds, int32(r.next.iconId))

	for _, v := range fsMath.MixIconGroup {
		if v.IconId == r.iconId {
			return iconIds, v.Quantity
		}
	}
	return nil, 0
}

func (s *slot) rngReels(reels []int) *reel {
	var reelOptions []rng.Option

	size := len(reels)

	for i := 0; i < size; i++ {
		p := &reel{
			previous: nil,
			next:     nil,
		}

		if i-1 < 0 {
			p.index = i - 1
			p.iconId = reels[size-1]
		} else {
			p.index = i - 1
			p.iconId = reels[i-1]
		}

		n := &reel{
			previous: nil,
			next:     nil,
		}

		if i+1 >= size {
			n.index = 0
			n.iconId = reels[0]
		} else {
			n.index = i + 1
			n.iconId = reels[i+1]
		}

		r := &reel{
			previous: p,
			next:     n,
			index:    i,
			iconId:   reels[i],
		}
		reelOptions = append(reelOptions, rng.Option{1, r})
	}
	return MONSTER.rng(reelOptions).(*reel)
}

func (s *slot) bonusPay(
	iconId int32,
	starFish, dory, lionFish, lobster, platypus, manatee, dolphin, koi, hermitCrab *FsFishMath,
	faceOff1, faceOff3, faceOff5 *FsFaceOffMath,
) (pay int) {
	switch iconId {
	case 0:
		return StarFish.HitFs(starFish)

	case 4:
		return Dory.HitFs(dory)

	case 9:
		return LionFish.HitFs(lionFish)

	case 11:
		return Lobster.HitFs(lobster)

	case 13:
		return Platypus.HitFs(platypus)

	case 14:
		return Manatee.HitFs(manatee)

	case 15:
		return Dolphin.HitFs(dolphin)

	case 16:
		return Koi.HitFs(koi)

	case 17:
		return HermitCrab.HitFs(hermitCrab)

	case 21:
		return FaceOff1.HitFs(faceOff1)

	case 23:
		return FaceOff3.HitFs(faceOff3)

	case 25:
		return FaceOff5.HitFs(faceOff5)

	default:
		return 0
	}
}

func (s *slot) contain(fsMath *FsSlotMath, pay int) bool {
	for _, v := range fsMath.Pay {
		if v == pay {
			return true
		}
	}
	return false
}

func (s *slot) Hit(
	bsMath *BsSlotMath,
	fsMath *FsSlotMath,
	reelMath *FsReelMath,
	faceOffChangerMath *FsFaceOffChangerMath,
	starFish, dory, lionFish, lobster, platypus, manatee, dolphin, koi, hermitCrab *FsFishMath,
	faceOff1, faceOff3, faceOff5 *FsFaceOffMath,
) (pay int, iconIds []int32, iconPays []int64) {
	if MONSTER.isHit(bsMath.HitWeight) {
		blockReelIndex := s.rngBlockReelIndex()

		iconIds := make([]int32, 0)
		iconPays := make([]int64, 3)

		for i := 0; i < fsMath.MixAwardId; i++ {
			if i == blockReelIndex {
				v1, v2 := s.rngBlockReel(fsMath, reelMath, faceOffChangerMath)

				for _, v := range v1 {
					iconIds = append(iconIds, v)
				}

				iconPays[i] = int64(v2)
			} else {
				v1, v2 := s.rngNormalReel(fsMath, reelMath, faceOffChangerMath)

				for _, v := range v1 {
					iconIds = append(iconIds, v)
				}

				iconPays[i] = int64(v2)
			}
		}

		// big win
		if iconIds[1] == iconIds[4] && iconIds[1] == iconIds[7] {
			bonusPay := s.bonusPay(
				iconIds[1],
				starFish, dory, lionFish, lobster, platypus, manatee, dolphin, koi, hermitCrab,
				faceOff1, faceOff3, faceOff5,
			)

			if s.contain(fsMath, bonusPay) {
				return bonusPay, iconIds, iconPays
			}
			return 0, nil, nil
		}

		if s.contain(fsMath, int(iconPays[0]+iconPays[1]+iconPays[2])) {
			return int(iconPays[0] + iconPays[1] + iconPays[2]), iconIds, iconPays
		}
	}
	return 0, nil, nil
}
