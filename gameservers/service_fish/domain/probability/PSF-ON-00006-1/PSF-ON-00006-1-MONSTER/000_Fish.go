package PSF_ON_00006_1_MONSTER

import (
	"os"
	"strings"
)

var Fish0 = &fish0{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "102"),
}

type fish0 struct {
	bonusMode bool
}

type BsFishMath struct {
	IconID     int                     `json:"IconID"`
	BonusID    string                  `json:"BonusID"`
	BonusTimes []int                   `json:"BonusTimes"`
	IconPays   []int                   `json:"IconPays"`
	RTP1       struct{ HitBsFishMath } `json:"RTP1"`
	RTP2       struct{ HitBsFishMath } `json:"RTP2"`
	RTP3       struct{ HitBsFishMath } `json:"RTP3"`
	RTP4       struct{ HitBsFishMath } `json:"RTP4"`
	RTP5       struct{ HitBsFishMath } `json:"RTP5"`
	RTP6       struct{ HitBsFishMath } `json:"RTP6"`
	RTP7       struct{ HitBsFishMath } `json:"RTP7"`
	RTP8       struct{ HitBsFishMath } `json:"RTP8"`
	RTP9       struct{ HitBsFishMath } `json:"RTP9"`
}

type HitBsFishMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

type FsFishMath struct {
	IconID      int
	IconsAmount []int
	IconPays    []int
}

func (f *fish0) Hit(rtpId string, math *BsFishMath) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP1)

	case math.RTP2.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP2)

	case math.RTP3.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP3)

	case math.RTP4.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP4)

	case math.RTP5.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP5)

	case math.RTP6.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP6)

	case math.RTP7.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP7)

	case math.RTP8.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP8)

	case math.RTP9.ID:
		return f.rtpHit(math.IconPays[0], &math.RTP9)

	default:
		return iconPay, triggerIconId, bonusTypeId
	}

	return iconPay, triggerIconId, bonusTypeId
}

func (f *fish0) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return math.IconPays[payIndex]
}
func (f *fish0) rtpHit(iconPays int, bsFish *struct{ HitBsFishMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(bsFish.HitWeight) {
		iconPay = iconPays
	}

	if MONSTER.isHit(bsFish.TriggerWeight) {
		triggerIconId = bsFish.TriggerIconID
		bonusTypeId = bsFish.Type
	}

	return iconPay, triggerIconId, bonusTypeId
}
