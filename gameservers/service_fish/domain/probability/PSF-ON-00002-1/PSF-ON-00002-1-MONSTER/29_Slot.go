package PSF_ON_00002_1_MONSTER

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
	IconID     int
	BonusID    string
	BonusTimes []int
	IconPays   []int
	RTP1       struct{ HitBsSlotMath } `json:"RTP1"`
	RTP2       struct{ HitBsSlotMath } `json:"RTP2"`
	RTP3       struct{ HitBsSlotMath } `json:"RTP3"`
	RTP4       struct{ HitBsSlotMath } `json:"RTP4"`
	RTP5       struct{ HitBsSlotMath } `json:"RTP5"`
	RTP6       struct{ HitBsSlotMath } `json:"RTP6"`
	RTP7       struct{ HitBsSlotMath } `json:"RTP7"`
	RTP8       struct{ HitBsSlotMath } `json:"RTP8"`
	RTP9       struct{ HitBsSlotMath } `json:"RTP9"`
	RTP10      struct{ HitBsSlotMath } `json:"RTP10"`
	RTP11      struct{ HitBsFishMath } `json:"RTP11"`
}

type HitBsSlotMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

type FsReelMath struct {
	Normal reelStrape `json:"Normal"`
	Block  reelStrape `json:"Block"`
}

type reelStrape struct {
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

func (s *slot) hit(
	fsMath *FsSlotMath,
	reelMath *FsReelMath,
	faceOffChangerMath31 *FsFaceOffChangerMath_31,
	faceOffChangerMath *FsFaceOffChangerMath,
	starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab *FsFishMath,
	faceOff1, faceOff3, faceOff5 *FsFaceOffMath,
) (iconPay int, iconIds []int32, iconPays []int64) {

	iconIds = make([]int32, 0)
	iconPays = make([]int64, 3)

	// N, N, B
	for i := 0; i < 2; i++ {
		v1, v2 := s.rngNormalReel(fsMath, reelMath, faceOffChangerMath31, faceOffChangerMath)

		for _, v := range v1 {
			iconIds = append(iconIds, v)
		}
		iconPays[i] = int64(v2)
	}

	v1, v2 := s.rngBlockReel(fsMath, reelMath, faceOffChangerMath31, faceOffChangerMath)

	for _, v := range v1 {
		iconIds = append(iconIds, v)
	}

	iconPays[2] = int64(v2)

	// big win
	if iconIds[1] == iconIds[4] && iconIds[1] == iconIds[7] {
		bonusPay := s.bonusPay(
			iconIds[1],
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			faceOff1, faceOff3, faceOff5,
		)

		if s.contain(fsMath, bonusPay) {
			return bonusPay, iconIds, iconPays
		}
		return 0, nil, nil
	}
	// Normal win
	if s.contain(fsMath, int(iconPays[0]+iconPays[1]+iconPays[2])) {
		return int(iconPays[0] + iconPays[1] + iconPays[2]), iconIds, iconPays
	}
	return 0, nil, nil
}

func (s *slot) rngNormalReel(fsMath *FsSlotMath, reelMath *FsReelMath, faceOffChangerMath31 *FsFaceOffChangerMath_31,
	faceOffChangerMath *FsFaceOffChangerMath) (iconIds []int32, pay int) {
	r := s.rngReels(&reelMath.Normal)

	switch r.previous.iconId {
	case 31:
		r.previous.iconId = FaceOffChanger.rngNormalReel31(faceOffChangerMath31)
	case 30:
		r.previous.iconId = FaceOffChanger.rngNormalReel(faceOffChangerMath)
	}

	iconIds = append(iconIds, int32(r.previous.iconId))

	switch r.iconId {
	case 31:
		r.iconId = FaceOffChanger.rngNormalReel31(faceOffChangerMath31)
	case 30:
		r.iconId = FaceOffChanger.rngNormalReel(faceOffChangerMath)
	}

	iconIds = append(iconIds, int32(r.iconId))

	switch r.next.iconId {
	case 31:
		r.next.iconId = FaceOffChanger.rngNormalReel31(faceOffChangerMath31)
	case 30:
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

func (s *slot) rngBlockReel(fsMath *FsSlotMath, reelMath *FsReelMath, faceOffChangerMath31 *FsFaceOffChangerMath_31,
	faceOffChangerMath *FsFaceOffChangerMath) (iconIds []int32, pay int) {

	r := s.rngReels(&reelMath.Block)

	for r.previous.iconId == 31 || r.previous.iconId == 30 {
		switch r.previous.iconId {
		case 31:
			r.previous.iconId = FaceOffChanger.rngBlockReel31(faceOffChangerMath31)
		case 30:
			r.previous.iconId = FaceOffChanger.rngBlockReel(faceOffChangerMath)
		default:
			break
		}
	}
	iconIds = append(iconIds, int32(r.previous.iconId))

	for r.iconId == 31 || r.iconId == 30 {
		switch r.iconId {
		case 31:
			r.iconId = FaceOffChanger.rngBlockReel31(faceOffChangerMath31)
		case 30:
			r.iconId = FaceOffChanger.rngBlockReel(faceOffChangerMath)
		default:
			break
		}
	}
	iconIds = append(iconIds, int32(r.iconId))

	for r.next.iconId == 31 || r.next.iconId == 30 {
		switch r.next.iconId {
		case 31:
			r.next.iconId = FaceOffChanger.rngBlockReel31(faceOffChangerMath31)
		case 30:
			r.next.iconId = FaceOffChanger.rngBlockReel(faceOffChangerMath)
		}
	}
	iconIds = append(iconIds, int32(r.next.iconId))

	for _, v := range fsMath.MixIconGroup {
		if v.IconId == r.iconId {
			return iconIds, v.Quantity
		}
	}

	return nil, 0
}

func (s *slot) rngReels(reelMath *reelStrape) *reel {

	var reelOptions []rng.Option
	var upDownOptions []rng.Option

	for i := 0; i < len(reelMath.UseIcon); i++ {
		allIcon := &reel{
			previous: nil,
			next:     nil,
			index:    0,
			iconId:   reelMath.UseIcon[i],
		}
		upDownOptions = append(upDownOptions, rng.Option{reelMath.HitWeight[i], allIcon})

		topIcon := MONSTER.rng(upDownOptions).(*reel)
		downIcon := MONSTER.rng(upDownOptions).(*reel)

		r := &reel{
			previous: topIcon,
			next:     downIcon,
			index:    i,
			iconId:   reelMath.UseIcon[i],
		}
		reelOptions = append(reelOptions, rng.Option{reelMath.HitWeight[i], r})
	}

	return MONSTER.rng(reelOptions).(*reel)
}

func (s *slot) bonusPay(
	iconId int32,
	starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab *FsFishMath,
	faceOff1, faceOff3, faceOff5 *FsFaceOffMath,
) (pay int) {
	switch iconId {
	case 0:
		return StarFish.HitFs(starFish)

	case 4:
		return Clownfish.HitFs(clownFish)

	case 9:
		return LionFish.HitFs(lionFish)

	case 11:
		return GiantJellyfish.HitFs(giantJellyFish)

	case 13:
		return Lobster.HitFs(lobster)

	case 14:
		return HermitCrab.HitFs(hermitCrab)

	case 15:
		return GoldGemTurtle.HitFs(goldGemTurtle)

	case 16:
		return GoldCoinToad.HitFs(goldCoinToad)

	case 17:
		return GoldCrab.HitFs(goldCrab)

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
	rtpId string,
	bsMath *BsSlotMath,
	fsMath *FsSlotMath,
	reelMath *FsReelMath,
	faceOffChangerMath31 *FsFaceOffChangerMath_31,
	faceOffChangerMath *FsFaceOffChangerMath,
	starFish, clownfish, lionFish, giantJellyfish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab *FsFishMath,
	faceOff1, faceOff3, faceOff5 *FsFaceOffMath,
) (iconPay int, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case "0":
		if MONSTER.isHit(bsMath.RTP1.TriggerWeight) {
			triggerIconId = bsMath.RTP1.TriggerIconID
			bonusTypeId = bsMath.RTP1.Type
		}

		if MONSTER.isHit(bsMath.RTP1.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "20":
		if MONSTER.isHit(bsMath.RTP2.TriggerWeight) {
			triggerIconId = bsMath.RTP2.TriggerIconID
			bonusTypeId = bsMath.RTP2.Type
		}

		if MONSTER.isHit(bsMath.RTP2.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "40":
		if MONSTER.isHit(bsMath.RTP3.TriggerWeight) {
			triggerIconId = bsMath.RTP3.TriggerIconID
			bonusTypeId = bsMath.RTP3.Type
		}

		if MONSTER.isHit(bsMath.RTP3.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "60":
		if MONSTER.isHit(bsMath.RTP4.TriggerWeight) {
			triggerIconId = bsMath.RTP4.TriggerIconID
			bonusTypeId = bsMath.RTP4.Type
		}

		if MONSTER.isHit(bsMath.RTP4.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "80":
		if MONSTER.isHit(bsMath.RTP5.TriggerWeight) {
			triggerIconId = bsMath.RTP5.TriggerIconID
			bonusTypeId = bsMath.RTP5.Type
		}

		if MONSTER.isHit(bsMath.RTP5.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "901":
		if MONSTER.isHit(bsMath.RTP6.TriggerWeight) {
			triggerIconId = bsMath.RTP6.TriggerIconID
			bonusTypeId = bsMath.RTP6.Type
		}

		if MONSTER.isHit(bsMath.RTP6.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "902":
		if MONSTER.isHit(bsMath.RTP7.TriggerWeight) {
			triggerIconId = bsMath.RTP7.TriggerIconID
			bonusTypeId = bsMath.RTP7.Type
		}

		if MONSTER.isHit(bsMath.RTP7.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "100":
		if MONSTER.isHit(bsMath.RTP8.TriggerWeight) {
			triggerIconId = bsMath.RTP8.TriggerIconID
			bonusTypeId = bsMath.RTP8.Type
		}

		if MONSTER.isHit(bsMath.RTP8.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "150":
		if MONSTER.isHit(bsMath.RTP9.TriggerWeight) {
			triggerIconId = bsMath.RTP9.TriggerIconID
			bonusTypeId = bsMath.RTP9.Type
		}

		if MONSTER.isHit(bsMath.RTP9.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "200":
		if MONSTER.isHit(bsMath.RTP10.TriggerWeight) {
			triggerIconId = bsMath.RTP10.TriggerIconID
			bonusTypeId = bsMath.RTP10.Type
		}

		if MONSTER.isHit(bsMath.RTP10.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "300":
		if MONSTER.isHit(bsMath.RTP11.TriggerWeight) {
			triggerIconId = bsMath.RTP11.TriggerIconID
			bonusTypeId = bsMath.RTP11.Type
		}

		if MONSTER.isHit(bsMath.RTP11.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyfish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	default:
		return 0, -1, -1, nil, nil
	}
}
