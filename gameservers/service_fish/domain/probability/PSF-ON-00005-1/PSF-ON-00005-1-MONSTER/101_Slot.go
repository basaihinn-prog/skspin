package PSF_ON_00005_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var Slot = &slot{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "101"),
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
}

type HitBsSlotMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
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

func (s *slot) Hit(
	rtpId string,
	bsMath *BsSlotMath,
	fsMath *FsSlotMath,
	fsStripTableMath *FsStripTableMath,
	giantWhaleChanger1Math *FsGiantWhaleChanger1Math,
	giantWhaleChanger2Math *FsGiantWhaleChanger2Math,
	starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
	fruitDish, aPackOfBeer, giantWhale *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case bsMath.RTP1.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP1, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP2.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP2, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP3.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP3, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP4.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP4, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP5.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP5, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP6.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP6, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP7.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP7, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP8.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP8, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP9.ID:
		v1, v2, v3, v4, v5 := s.rtpHit(&bsMath.RTP9, fsMath, fsStripTableMath,
			giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, v2, v3, v4, v5

	default:
		return 0, -1, -1, nil, nil
	}
}

func (s *slot) hit(
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	giantWhaleMath1 *FsGiantWhaleChanger1Math,
	giantWhaleMath2 *FsGiantWhaleChanger2Math,
	starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab,
	goldGemTurtle, goldCoinToad, goldCrab,
	fruitDish, aPackOfBeer, GiantWhale *FsFishMath,
) (iconPay int, iconIds []int32, iconPays []int64) {
	blockIndex := s.rngBlockReelIndex()

	iconIds = make([]int32, 0)
	iconPays = make([]int64, 3)

	for i := 0; i < fsMath.MixAwardId; i++ {
		if i == blockIndex {
			v1, v2 := s.rngReel(fsMath, &fsStripMath.Block,
				&giantWhaleMath1.ChangeCandidateBlock,
				&giantWhaleMath2.ChangeCandidateBlock,
			)

			for _, v := range v1 {
				iconIds = append(iconIds, v)
			}
			iconPays[i] = int64(v2)
		} else {
			v1, v2 := s.rngReel(fsMath, &fsStripMath.Normal,
				&giantWhaleMath1.ChangeCandidateNormal,
				&giantWhaleMath2.ChangeCandidateNormal,
			)

			for _, v := range v1 {
				iconIds = append(iconIds, v)
			}
			iconPays[i] = int64(v2)
		}
	}

	// big win
	if iconIds[1] == iconIds[4] && iconIds[1] == iconIds[7] {
		bonusPay := s.bonusPay(iconIds[1],
			starFish, clownFish, lionFish, giantJellyFish, lobster,
			hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, GiantWhale,
		)

		if s.contain(fsMath, bonusPay) {
			return bonusPay, iconIds, iconPays
		}
		return 0, nil, nil
	}

	if s.contain(fsMath, int(iconPays[0]+iconPays[1]+iconPays[2])) {
		return int(iconPays[0] + iconPays[1] + iconPays[2]), iconIds, iconPays
	}

	return 0, nil, nil
}

func (s *slot) rngBlockReelIndex() int {
	var indexOptions []rng.Option
	indexOptions = append(indexOptions, rng.Option{1, 0})
	indexOptions = append(indexOptions, rng.Option{1, 1})
	indexOptions = append(indexOptions, rng.Option{1, 2})
	return MONSTER.rng(indexOptions).(int)
}

func (s *slot) rngReel(
	fsMath *FsSlotMath,
	fsStripMath *struct{ HitFsStripTableMath },
	change1Math *struct{ FsChangeCandidate1Math },
	change2Math *struct{ FsChangeCandidate2Math },
) (iconIds []int32, pay int) {
	r := s.rngReels(fsStripMath)

	r.previous.iconId = s.rngCandidate(r.previous.iconId, change1Math, change2Math)
	r.iconId = s.rngCandidate(r.iconId, change1Math, change2Math)
	r.next.iconId = s.rngCandidate(r.next.iconId, change1Math, change2Math)

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

func (s *slot) rngCandidate(
	iconId int,
	change1Math *struct{ FsChangeCandidate1Math },
	change2Math *struct{ FsChangeCandidate2Math },
) int {
	newIconId := iconId

	for {
		iconId = newIconId

		if iconId == 28 {
			newIconId = GiantWhaleChanger.rngCandidate1(change1Math)
		}
		if iconId == 29 {
			newIconId = GiantWhaleChanger.rngCandidate2(change2Math)
		}

		if newIconId != 28 && newIconId != 29 {
			return newIconId
		}
	}

	return iconId
}

func (s *slot) rngReels(reels *struct{ HitFsStripTableMath }) *reel {
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

func (s *slot) contain(fsMath *FsSlotMath, pay int) bool {
	for _, v := range fsMath.Pay {
		if v == pay {
			return true
		}
	}

	return false
}

func (s *slot) bonusPay(iconId int32,
	starFish, clownFish, lionFish, durianFish, muayThaiLobster,
	hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
	fruitDish, aPackOfBeer, giantWhale *FsFishMath,
) (pay int) {
	switch iconId {
	case 0:
		return StarFish.HitFs(starFish)

	case 4:
		return ClownFish.HitFs(clownFish)

	case 9:
		return LionFish.HitFs(lionFish)

	case 11:
		return DurianFish.HitFs(durianFish)

	case 13:
		return MuayThaiLobster.HitFs(muayThaiLobster)

	case 14:
		return HermitCrab.HitFs(hermitCrab)

	case 15:
		return GoldGemTurtle.HitFs(goldGemTurtle)

	case 16:
		return GoldCoinToad.HitFs(goldCoinToad)

	case 17:
		return GoldCrab.HitFs(goldCrab)

	case 300:
		return FruitDish.HitFs(fruitDish)

	case 301:
		return APackOfBeer.HitFs(aPackOfBeer)

	case 500:
		return GiantWhaleChanger.HitFs(giantWhale)

	default:
		return 0
	}
}

func (s *slot) rtpHit(
	rtp *struct{ HitBsSlotMath },
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	giantWhaleChanger1Math *FsGiantWhaleChanger1Math,
	giantWhaleChanger2Math *FsGiantWhaleChanger2Math,
	starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab,
	goldGemTurtle, goldCoinToad, goldCrab,
	fruitDish, aPackOfBeer, giantWhale *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		v1, v2, v3 := s.hit(fsMath, fsStripMath, giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster,
			hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, triggerIconId, bonusTypeId, v2, v3
	}

	return 0, triggerIconId, bonusTypeId, nil, nil
}
