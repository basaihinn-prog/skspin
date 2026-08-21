package PSF_ON_00006_1_MONSTER

import (
	"os"
	"strings"

	"serve/fish_comm/rng"
)

var LuckSlot = &luckSlot{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "101"),
}

type luckSlot struct {
	bonusMode bool
}

type BsLuckSlotMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	IconPays   []int
	RTP1       struct{ HitBsLuckSlotMath } `json:"RTP1"`
	RTP2       struct{ HitBsLuckSlotMath } `json:"RTP2"`
	RTP3       struct{ HitBsLuckSlotMath } `json:"RTP3"`
	RTP4       struct{ HitBsLuckSlotMath } `json:"RTP4"`
	RTP5       struct{ HitBsLuckSlotMath } `json:"RTP5"`
	RTP6       struct{ HitBsLuckSlotMath } `json:"RTP6"`
	RTP7       struct{ HitBsLuckSlotMath } `json:"RTP7"`
	RTP8       struct{ HitBsLuckSlotMath } `json:"RTP8"`
	RTP9       struct{ HitBsLuckSlotMath } `json:"RTP9"`
}

type HitBsLuckSlotMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

type FsStripTableMath struct {
	Reel1 struct{ HitFsStripTableMath } `json:"Reel_1"`
	Reel2 struct{ HitFsStripTableMath } `json:"Reel_2"`
	Reel3 struct{ HitFsStripTableMath } `json:"Reel_3"`
	Reel4 struct{ HitFsStripTableMath } `json:"Reel_4"`
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

type FsChangeLink1 struct {
	ChangeIconId         []int
	ChangeCandidateReel1 struct{ RandomFsChange } `json:"ChangeCandidate_Reel_1"`
	ChangeCandidateReel2 struct{ RandomFsChange } `json:"ChangeCandidate_Reel_2"`
	ChangeCandidateReel3 struct{ RandomFsChange } `json:"ChangeCandidate_Reel_3"`
	ChangeCandidateReel4 struct{ RandomFsChange } `json:"ChangeCandidate_Reel_4"`
}

type RandomFsChange struct {
	ChangeCandidate1  struct{ RandomFsChangeLink } `json:"ChangeCandidate_1"`
	ChangeCandidate2  struct{ RandomFsChangeLink } `json:"ChangeCandidate_2"`
	ChangeCandidate3  struct{ RandomFsChangeLink } `json:"ChangeCandidate_3"`
	ChangeCandidate4  struct{ RandomFsChangeLink } `json:"ChangeCandidate_4"`
	ChangeCandidate5  struct{ RandomFsChangeLink } `json:"ChangeCandidate_5"`
	ChangeCandidate6  struct{ RandomFsChangeLink } `json:"ChangeCandidate_6"`
	ChangeCandidate7  struct{ RandomFsChangeLink } `json:"ChangeCandidate_7"`
	ChangeCandidate8  struct{ RandomFsChangeLink } `json:"ChangeCandidate_8"`
	ChangeCandidate9  struct{ RandomFsChangeLink } `json:"ChangeCandidate_9"`
	ChangeCandidate10 struct{ RandomFsChangeLink } `json:"ChangeCandidate_10"`
	ChangeCandidate11 struct{ RandomFsChangeLink } `json:"ChangeCandidate_11"`
	ChangeCandidate12 struct{ RandomFsChangeLink } `json:"ChangeCandidate_12"`
}

type RandomFsChangeLink struct {
	IconId int `json:"IconID"`
	Weight int `json:"Weight"`
}

type reel struct {
	previous *reel
	next     *reel
	index    int
	iconId   int
}

func (l *luckSlot) Hit(
	rtpId string,
	bsMath *BsLuckSlotMath,
	fsMath *FsSlotMath,
	fsStripTableMath *FsStripTableMath,
	changeLink *FsChangeLink1,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	magicShell, fuChe, giantWhale *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case bsMath.RTP1.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP1, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP2.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP2, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP3.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP3, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP4.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP4, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP5.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP5, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP6.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP6, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP7.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP7, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP8.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP8, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	case bsMath.RTP9.ID:
		v1, v2, v3, v4, v5 := l.rtpHit(&bsMath.RTP9, fsMath, fsStripTableMath,
			changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, v2, v3, v4, v5

	default:
		return 0, -1, -1, nil, nil
	}
}

func (l *luckSlot) hit(
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	changeLink *FsChangeLink1,
	fish0, fish4, fish9, fish11, fish13,
	fish14, fish15, fish16, fish17,
	magicShell, fuChe, giantWhale *FsFishMath,
) (iconPay int, iconIds []int32, iconPays []int64) {
	iconIds = make([]int32, 0)
	iconPays = make([]int64, 4)

	for i := 0; i < fsMath.MixAwardId; i++ {
		reels := [4]*struct{ HitFsStripTableMath }{
			&fsStripMath.Reel1, &fsStripMath.Reel2,
			&fsStripMath.Reel3, &fsStripMath.Reel4,
		}

		changes := [4]*struct{ RandomFsChange }{
			&changeLink.ChangeCandidateReel1, &changeLink.ChangeCandidateReel2,
			&changeLink.ChangeCandidateReel3, &changeLink.ChangeCandidateReel4,
		}

		v1, v2 := l.rngReel(fsMath, reels[i], changes[i])

		for _, v := range v1 {
			iconIds = append(iconIds, v)
		}
		iconPays[i] = int64(v2)
	}

	// big win
	if l.bigWin3(iconIds) {
		payIndex := 0
		if l.bigWin4(iconIds) {
			payIndex = 1
		}

		bonusPay := l.bonusPay(iconIds[1], payIndex,
			fish0, fish4, fish9, fish11, fish13,
			fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)

		if l.contain(fsMath, bonusPay) {
			return bonusPay, iconIds, iconPays
		}

		return 0, nil, nil
	}

	totalPay := 0
	for _, v := range iconPays {
		totalPay += int(v)
	}

	if l.contain(fsMath, totalPay) {
		return totalPay, iconIds, iconPays
	}

	return 0, nil, nil
}

func (l *luckSlot) rngReel(
	fsMath *FsSlotMath,
	fsStripMath *struct{ HitFsStripTableMath },
	changeLinkMath *struct{ RandomFsChange },
) (iconIds []int32, pay int) {
	r := l.rngReels(fsStripMath)

	r.previous.iconId = l.rngCandidate(r.previous.iconId, changeLinkMath)
	r.iconId = l.rngCandidate(r.iconId, changeLinkMath)
	r.next.iconId = l.rngCandidate(r.next.iconId, changeLinkMath)

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

func (l *luckSlot) rngCandidate(
	iconId int,
	change1Math *struct{ RandomFsChange },
) int {
	newIconId := iconId

	for {
		iconId = newIconId

		if iconId == 28 {
			newIconId = l.rngCandidateMath(change1Math)
		}

		if newIconId != 28 {
			return newIconId
		}
	}
}

func (l *luckSlot) rngReels(reels *struct{ HitFsStripTableMath }) *reel {
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

func (l *luckSlot) contain(fsMath *FsSlotMath, pay int) bool {
	for _, v := range fsMath.Pay {
		if v == pay {
			return true
		}
	}

	return false
}

func (l *luckSlot) bonusPay(iconId int32, payIndex int,
	fish0, fish4, fish9, fish11, fish13,
	fish14, fish15, fish16, fish17,
	magicShell, fuChe, giantWhale *FsFishMath,
) (pay int) {
	switch iconId {
	case 1:
		return Fish0.HitFs(fish0, payIndex)

	case 2:
		return Fish4.HitFs(fish4, payIndex)

	case 4:
		return Fish9.HitFs(fish9, payIndex)

	case 7:
		return Fish11.HitFs(fish11, payIndex)

	case 9:
		return Fish13.HitFs(fish13, payIndex)

	case 11:
		return Fish14.HitFs(fish14, payIndex)

	case 14:
		return Fish15.HitFs(fish15, payIndex)

	case 15:
		return Fish16.HitFs(fish16, payIndex)

	case 17:
		return Fish17.HitFs(fish17, payIndex)

	case 300:
		return MagicShell.HitFs(magicShell, payIndex)

	case 301:
		return FuChe.HitFs(fuChe, payIndex)

	case 400:
		return Fish0.HitFs(giantWhale, payIndex)

	default:
		return 0
	}
}

func (l *luckSlot) bigWin4(iconIds []int32) bool {
	if iconIds[1] == iconIds[4] && iconIds[4] == iconIds[7] && iconIds[7] == iconIds[10] {
		return true
	}

	return false
}

func (l *luckSlot) bigWin3(iconIds []int32) bool {
	if (iconIds[1] == iconIds[4] && iconIds[4] == iconIds[7]) ||
		(iconIds[1] == iconIds[4] && iconIds[4] == iconIds[10]) ||
		(iconIds[1] == iconIds[7] && iconIds[7] == iconIds[10]) ||
		(iconIds[4] == iconIds[7] && iconIds[7] == iconIds[10]) {
		return true
	}

	return false
}

func (l *luckSlot) rngCandidateMath(math *struct{ RandomFsChange }) (iconId int) {
	options := make([]rng.Option, 0)

	options = append(options, rng.Option{math.ChangeCandidate1.Weight, math.ChangeCandidate1.IconId})
	options = append(options, rng.Option{math.ChangeCandidate2.Weight, math.ChangeCandidate2.IconId})
	options = append(options, rng.Option{math.ChangeCandidate3.Weight, math.ChangeCandidate3.IconId})
	options = append(options, rng.Option{math.ChangeCandidate4.Weight, math.ChangeCandidate4.IconId})
	options = append(options, rng.Option{math.ChangeCandidate5.Weight, math.ChangeCandidate5.IconId})
	options = append(options, rng.Option{math.ChangeCandidate6.Weight, math.ChangeCandidate6.IconId})
	options = append(options, rng.Option{math.ChangeCandidate7.Weight, math.ChangeCandidate7.IconId})
	options = append(options, rng.Option{math.ChangeCandidate8.Weight, math.ChangeCandidate8.IconId})
	options = append(options, rng.Option{math.ChangeCandidate9.Weight, math.ChangeCandidate9.IconId})
	options = append(options, rng.Option{math.ChangeCandidate10.Weight, math.ChangeCandidate10.IconId})
	options = append(options, rng.Option{math.ChangeCandidate11.Weight, math.ChangeCandidate11.IconId})
	options = append(options, rng.Option{math.ChangeCandidate12.Weight, math.ChangeCandidate12.IconId})

	return MONSTER.rng(options).(int)
}

func (l *luckSlot) rtpHit(
	rtp *struct{ HitBsLuckSlotMath },
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	changeLink *FsChangeLink1,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	magicShell, fuChe, giantWhale *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		v1, v2, v3 := l.hit(fsMath, fsStripMath, changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			magicShell, fuChe, giantWhale,
		)
		return v1, triggerIconId, bonusTypeId, v2, v3
	}

	return 0, triggerIconId, bonusTypeId, nil, nil
}
