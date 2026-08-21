package RKF_H5_00001_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var HalonSlot = &halonSlot{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "101"),
}

type halonSlot struct {
	bonusMode bool
}

type BsHalonSlotMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	IconPays   []int
	RTP1       struct{ HitBsHalonSlotMath } `json:"RTP1"`
	RTP2       struct{ HitBsHalonSlotMath } `json:"RTP2"`
	RTP3       struct{ HitBsHalonSlotMath } `json:"RTP3"`
	RTP4       struct{ HitBsHalonSlotMath } `json:"RTP4"`
	RTP5       struct{ HitBsHalonSlotMath } `json:"RTP5"`
	RTP6       struct{ HitBsHalonSlotMath } `json:"RTP6"`
	RTP7       struct{ HitBsHalonSlotMath } `json:"RTP7"`
	RTP8       struct{ HitBsHalonSlotMath } `json:"RTP8"`
	RTP9       struct{ HitBsHalonSlotMath } `json:"RTP9"`
}

type HitBsHalonSlotMath struct {
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

func (h *halonSlot) Hit(
	rtpId string,
	bsMath *BsHalonSlotMath, fsMath *FsSlotMath,
	fsStripTableMath *FsStripTableMath,
	changeLink1Math *FsChangeLink1Math, changeLink2Math *FsChangeLink2Math,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	faceOff1, faceOff3, faceOff5 *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case bsMath.RTP1.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP1, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP2.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP2, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP3.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP3, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP4.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP4, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP5.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP5, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP6.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP6, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP7.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP7, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP8.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP8, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP9.ID:
		v1, v2, v3, v4, v5 := h.rtpHit(&bsMath.RTP9, fsMath, fsStripTableMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, v2, v3, v4, v5

	default:
		return 0, -1, -1, nil, nil
	}
}

func (h *halonSlot) hit(
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	changeLink1Math *FsChangeLink1Math, changeLink2Math *FsChangeLink2Math,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	faceOff1, faceOff3, faceOff5 *FsFishMath,
) (iconPay int, iconIds []int32, iconPays []int64) {
	blockIndex := h.rngBlockReelIndex()

	iconIds = make([]int32, 0)
	iconPays = make([]int64, 3)

	for i := 0; i < fsMath.MixAwardId; i++ {
		if i == blockIndex {
			v1, v2 := h.rngReel(fsMath, &fsStripMath.Block,
				&changeLink1Math.ChangeCandidateBlock,
				&changeLink2Math.ChangeCandidateBlock,
			)

			for _, v := range v1 {
				iconIds = append(iconIds, v)
			}
			iconPays[i] = int64(v2)
		} else {
			v1, v2 := h.rngReel(fsMath, &fsStripMath.Normal,
				&changeLink1Math.ChangeCandidateNormal,
				&changeLink2Math.ChangeCandidateNormal,
			)

			for _, v := range v1 {
				iconIds = append(iconIds, v)
			}
			iconPays[i] = int64(v2)
		}
	}

	// big win
	if iconIds[1] == iconIds[4] && iconIds[1] == iconIds[7] {
		bonusPay := h.bonusPay(iconIds[1],
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)

		if h.contain(fsMath, bonusPay) {
			return bonusPay, iconIds, iconPays
		}
		return 0, nil, nil
	}

	if h.contain(fsMath, int(iconPays[0]+iconPays[1]+iconPays[2])) {
		return int(iconPays[0] + iconPays[1] + iconPays[2]), iconIds, iconPays
	}

	return 0, nil, nil
}

func (h *halonSlot) rngReel(
	fsMath *FsSlotMath,
	fsStripMath *struct{ HitFsStripTableMath },
	change1Math *struct{ FsChangeCandidate1Math },
	change2Math *struct{ FsChangeCandidate2Math },
) (iconIds []int32, pay int) {
	r := h.rngReels(fsStripMath)

	r.previous.iconId = h.rngCandidate(r.previous.iconId, change1Math, change2Math)
	r.iconId = h.rngCandidate(r.iconId, change1Math, change2Math)
	r.next.iconId = h.rngCandidate(r.next.iconId, change1Math, change2Math)

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

func (h *halonSlot) rngCandidate(
	iconId int,
	change1Math *struct{ FsChangeCandidate1Math },
	change2Math *struct{ FsChangeCandidate2Math },
) int {
	newIconId := iconId

	for {
		iconId = newIconId

		if iconId == 28 {
			newIconId = PoseidonChanger.rngCandidate1(change1Math)
		}
		if iconId == 29 {
			newIconId = PoseidonChanger.rngCandidate2(change2Math)
		}

		if newIconId != 28 && newIconId != 29 {
			return newIconId
		}
	}

	return iconId
}

func (h *halonSlot) rngReels(reels *struct{ HitFsStripTableMath }) *reel {
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

func (h *halonSlot) rngBlockReelIndex() int {
	var indexOptions []rng.Option
	indexOptions = append(indexOptions, rng.Option{1, 0})
	indexOptions = append(indexOptions, rng.Option{1, 1})
	indexOptions = append(indexOptions, rng.Option{1, 2})
	return MONSTER.rng(indexOptions).(int)
}

func (h *halonSlot) contain(fsMath *FsSlotMath, pay int) bool {
	for _, v := range fsMath.Pay {
		if v == pay {
			return true
		}
	}

	return false
}

func (h *halonSlot) bonusPay(iconId int32,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	faceOff1, faceOff3, faceOff5 *FsFishMath,
) (pay int) {
	switch iconId {
	case 0:
		return Fish0.HitFs(fish0)

	case 4:
		return Fish4.HitFs(fish4)

	case 9:
		return Fish9.HitFs(fish9)

	case 11:
		return Fish11.HitFs(fish11)

	case 13:
		return Fish13.HitFs(fish13)

	case 14:
		return Fish14.HitFs(fish14)

	case 15:
		return Fish15.HitFs(fish15)

	case 16:
		return Fish16.HitFs(fish16)

	case 17:
		return Fish17.HitFs(fish17)

	case 300:
		return DeepSeaTreasure.HitFs(faceOff1)

	case 301:
		return CasketOfAncientWinters.HitFs(faceOff3)

	case 500:
		return PoseidonChanger.HitFs(faceOff5)

	default:
		return 0
	}
}

func (h *halonSlot) rtpHit(
	rtp *struct{ HitBsHalonSlotMath },
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	changeLink1Math *FsChangeLink1Math, changeLink2Math *FsChangeLink2Math,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	faceOff1, faceOff3, faceOff5 *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		v1, v2, v3 := h.hit(
			fsMath, fsStripMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, triggerIconId, bonusTypeId, v2, v3
	}

	return 0, triggerIconId, bonusTypeId, nil, nil
}
