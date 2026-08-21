package PSF_ON_00004_1_MONSTER

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
	RTP1       struct{ HitBsSlotMath }
	RTP2       struct{ HitBsSlotMath }
	RTP3       struct{ HitBsSlotMath }
	RTP4       struct{ HitBsSlotMath }
	RTP5       struct{ HitBsSlotMath }
}

type HitBsSlotMath struct {
	ID        string
	AvgPay    []float32
	HitWeight []int
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

func (s *slot) Hit(rtpId string,
	bsMath *BsSlotMath,
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	faceOffChangerMath1 *FsFaceOffChanger1Math,
	faceOffChangerMath2 *FsFaceOffChanger2Math,
	clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
	snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef *FsFishMath,
) (iconPay int, iconIds []int32, iconPays []int64, avgPay int) {
	switch rtpId {
	case bsMath.RTP1.ID:
		avgPay = int(bsMath.RTP1.AvgPay[0])
		v1, v2, v3 := s.rtpHit(bsMath.RTP1.HitWeight, fsMath, fsStripMath, faceOffChangerMath1, faceOffChangerMath2,
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
		)
		return v1, v2, v3, avgPay

	case bsMath.RTP2.ID:
		avgPay = int(bsMath.RTP2.AvgPay[0])
		v1, v2, v3 := s.rtpHit(bsMath.RTP2.HitWeight, fsMath, fsStripMath, faceOffChangerMath1, faceOffChangerMath2,
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
		)
		return v1, v2, v3, avgPay

	case bsMath.RTP3.ID:
		avgPay = int(bsMath.RTP3.AvgPay[0])
		v1, v2, v3 := s.rtpHit(bsMath.RTP3.HitWeight, fsMath, fsStripMath, faceOffChangerMath1, faceOffChangerMath2,
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
		)
		return v1, v2, v3, avgPay

	case bsMath.RTP4.ID:
		avgPay = int(bsMath.RTP4.AvgPay[0])
		v1, v2, v3 := s.rtpHit(bsMath.RTP4.HitWeight, fsMath, fsStripMath, faceOffChangerMath1, faceOffChangerMath2,
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
		)
		return v1, v2, v3, avgPay

	case bsMath.RTP5.ID:
		avgPay = int(bsMath.RTP5.AvgPay[0])
		v1, v2, v3 := s.rtpHit(bsMath.RTP5.HitWeight, fsMath, fsStripMath, faceOffChangerMath1, faceOffChangerMath2,
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
		)
		return v1, v2, v3, avgPay

	default:
		return 0, nil, nil, 0
	}
}

func (s *slot) hit(
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	faceOffChangerMath1 *FsFaceOffChanger1Math,
	faceOffChangerMath2 *FsFaceOffChanger2Math,
	clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
	snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef *FsFishMath,
) (iconPay int, iconIds []int32, iconPays []int64) {
	blockIndex := s.rngBlockReelIndex()

	iconIds = make([]int32, 0)
	iconPays = make([]int64, 3)

	for i := 0; i < fsMath.MixAwardId; i++ {
		if i == blockIndex {
			v1, v2 := s.rngReel(fsMath, &fsStripMath.Block,
				&faceOffChangerMath1.ChangeCandidateBlock,
				&faceOffChangerMath2.ChangeCandidateBlock,
			)

			for _, v := range v1 {
				iconIds = append(iconIds, v)
			}
			iconPays[i] = int64(v2)
		} else {
			v1, v2 := s.rngReel(fsMath, &fsStripMath.Normal,
				&faceOffChangerMath1.ChangeCandidateNormal,
				&faceOffChangerMath2.ChangeCandidateNormal,
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
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
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
	for i := 0; i < 3; i++ {
		indexOptions = append(indexOptions, rng.Option{1, i})
	}

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

func (s *slot) rngCandidate(iconId int, change1Math *struct{ FsChangeCandidate1Math }, change2Math *struct{ FsChangeCandidate2Math }) int {
	newIconId := iconId

	for {
		iconId = newIconId

		if iconId == 28 {
			newIconId = FaceOffChanger.rngCandidate1(change1Math)
		}
		if iconId == 29 {
			newIconId = FaceOffChanger.rngCandidate2(change2Math)
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

func (s *slot) bonusPay(iconId int32,
	clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
	snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef *FsFishMath,
) (pay int) {
	switch iconId {
	case 0:
		return Clam.HitFs(clam)

	case 4:
		return FlyingFish.HitFs(flyingFish)

	case 9:
		return MahiMahi.HitFs(mahiMahi)

	case 11:
		return LobsterWaiter.HitFs(lobsterWaiter)

	case 13:
		return PlatypusSeniorChef.HitFs(platypusSeniorChef)

	case 14:
		return SeaLionChef.HitFs(seaLionChef)

	case 15:
		return HairtCrab.HitFs(hairtCrab)

	case 16:
		return SnaggletoothShark.HitFs(snaggletoothShark)

	case 17:
		return SwordFish.HitFs(swordfish)

	case 300:
		return LobsterDash.HitFs(lobsterDish)

	case 301:
		return FruitDash.HitFs(fruitDish)

	case 400:
		return WhiteTigerChef.HitFs(whiteTigerChef)

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

func (s *slot) Rtp(order int, math *BsSlotMath) (rtpId string) {
	switch order {
	case 1:
		return math.RTP1.ID

	case 2:
		return math.RTP2.ID

	case 3:
		return math.RTP3.ID

	case 4:
		return math.RTP4.ID

	case 5:
		return math.RTP5.ID

	default:
		return ""
	}
}

func (s *slot) AvgPay(math *BsSlotMath) int {
	return int(math.RTP1.AvgPay[0])
}

func (s *slot) rtpHit(
	weights []int,
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	faceOffChangerMath1 *FsFaceOffChanger1Math,
	faceOffChangerMath2 *FsFaceOffChanger2Math,
	clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
	snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef *FsFishMath,
) (iconPay int, iconIds []int32, iconPays []int64) {
	if MONSTER.isHit(weights) {
		v1, v2, v3 := s.hit(
			fsMath,
			fsStripMath,
			faceOffChangerMath1,
			faceOffChangerMath2,
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
		)

		return v1, v2, v3
	}
	return 0, nil, nil
}
