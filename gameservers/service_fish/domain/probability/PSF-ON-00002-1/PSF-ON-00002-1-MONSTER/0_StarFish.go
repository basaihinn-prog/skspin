package PSF_ON_00002_1_MONSTER

import (
	"os"
	"strings"
)

var StarFish = &starFish{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "31"),
}

type starFish struct {
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
	RTP10      struct{ HitBsFishMath } `json:"RTP10"`
	RTP11      struct{ HitBsFishMath } `json:"RTP11"`
}

type HitBsFishMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

type FsFishMath struct {
	IconID   int   `json:"IconID"`
	IconPays []int `json:"IconPays"`
}

func (s *starFish) Hit(rtpId string, math *BsFishMath) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case "0":
		if MONSTER.isHit(math.RTP1.TriggerWeight) {
			triggerIconId = math.RTP1.TriggerIconID
			bonusTypeId = math.RTP1.Type
		}

		if MONSTER.isHit(math.RTP1.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "20":
		if MONSTER.isHit(math.RTP2.TriggerWeight) {
			triggerIconId = math.RTP2.TriggerIconID
			bonusTypeId = math.RTP2.Type
		}

		if MONSTER.isHit(math.RTP2.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "40":
		if MONSTER.isHit(math.RTP3.TriggerWeight) {
			triggerIconId = math.RTP3.TriggerIconID
			bonusTypeId = math.RTP3.Type
		}

		if MONSTER.isHit(math.RTP3.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "60":
		if MONSTER.isHit(math.RTP4.TriggerWeight) {
			triggerIconId = math.RTP4.TriggerIconID
			bonusTypeId = math.RTP4.Type
		}

		if MONSTER.isHit(math.RTP4.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "80":
		if MONSTER.isHit(math.RTP5.TriggerWeight) {
			triggerIconId = math.RTP5.TriggerIconID
			bonusTypeId = math.RTP5.Type
		}

		if MONSTER.isHit(math.RTP5.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "901":
		if MONSTER.isHit(math.RTP6.TriggerWeight) {
			triggerIconId = math.RTP6.TriggerIconID
			bonusTypeId = math.RTP6.Type
		}

		if MONSTER.isHit(math.RTP6.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "902":
		if MONSTER.isHit(math.RTP7.TriggerWeight) {
			triggerIconId = math.RTP7.TriggerIconID
			bonusTypeId = math.RTP7.Type
		}

		if MONSTER.isHit(math.RTP7.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "100":
		if MONSTER.isHit(math.RTP8.TriggerWeight) {
			triggerIconId = math.RTP8.TriggerIconID
			bonusTypeId = math.RTP8.Type
		}

		if MONSTER.isHit(math.RTP8.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "150":
		if MONSTER.isHit(math.RTP9.TriggerWeight) {
			triggerIconId = math.RTP9.TriggerIconID
			bonusTypeId = math.RTP9.Type
		}

		if MONSTER.isHit(math.RTP9.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "200":
		if MONSTER.isHit(math.RTP10.TriggerWeight) {
			triggerIconId = math.RTP10.TriggerIconID
			bonusTypeId = math.RTP10.Type
		}

		if MONSTER.isHit(math.RTP10.HitWeight) {
			iconPay = math.IconPays[0]
		}

	case "300":
		if MONSTER.isHit(math.RTP11.TriggerWeight) {
			triggerIconId = math.RTP11.TriggerIconID
			bonusTypeId = math.RTP11.Type
		}

		if MONSTER.isHit(math.RTP11.HitWeight) {
			iconPay = math.IconPays[0]
		}

	default:
		return iconPay, triggerIconId, bonusTypeId
	}

	return iconPay, triggerIconId, bonusTypeId
}

func (s *starFish) HitFs(math *FsFishMath) (iconPays int) {
	return math.IconPays[0]
}
