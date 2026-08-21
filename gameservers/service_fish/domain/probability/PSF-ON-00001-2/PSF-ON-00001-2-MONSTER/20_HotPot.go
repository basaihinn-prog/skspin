package PSF_ON_00001_2_MONSTER

import "serve/fish_comm/rng"

var HotPot = &hotPot{}

type hotPot struct {
}

type BsHotPotMath struct {
	IconID         int                       `json:"IconID"`
	BonusID        string                    `json:"BonusID"`
	BonusTimes     []int                     `json:"BonusTimes"`
	IconPays       []int                     `json:"IconPays"`
	IconPaysWeight []int                     `json:"IconPaysWeight"`
	RTP1           struct{ HitBsHotPotMath } `json:"RTP1"`
	RTP2           struct{ HitBsHotPotMath } `json:"RTP2"`
	RTP3           struct{ HitBsHotPotMath } `json:"RTP3"`
	RTP4           struct{ HitBsHotPotMath } `json:"RTP4"`
	RTP5           struct{ HitBsHotPotMath } `json:"RTP5"`
	RTP6           struct{ HitBsHotPotMath } `json:"RTP6"`
	RTP7           struct{ HitBsHotPotMath } `json:"RTP7"`
	RTP8           struct{ HitBsHotPotMath } `json:"RTP8"`
	RTP9           struct{ HitBsHotPotMath } `json:"RTP9"`
	RTP10          struct{ HitBsHotPotMath } `json:"RTP10"`
	RTP11          struct{ HitBsHotPotMath } `json:"RTP11"`
}

type HitBsHotPotMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (h *hotPot) Hit(rtpId string, math *BsHotPotMath) (iconPay int, triggerIconId int, bonusTypeId int) {
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
			iconPay = h.pick(math)
		}

	case "20":
		if MONSTER.isHit(math.RTP2.TriggerWeight) {
			triggerIconId = math.RTP2.TriggerIconID
			bonusTypeId = math.RTP2.Type
		}

		if MONSTER.isHit(math.RTP2.HitWeight) {
			iconPay = h.pick(math)
		}

	case "40":
		if MONSTER.isHit(math.RTP3.TriggerWeight) {
			triggerIconId = math.RTP3.TriggerIconID
			bonusTypeId = math.RTP3.Type
		}

		if MONSTER.isHit(math.RTP3.HitWeight) {
			iconPay = h.pick(math)
		}

	case "60":
		if MONSTER.isHit(math.RTP4.TriggerWeight) {
			triggerIconId = math.RTP4.TriggerIconID
			bonusTypeId = math.RTP4.Type
		}

		if MONSTER.isHit(math.RTP4.HitWeight) {
			iconPay = h.pick(math)
		}

	case "80":
		if MONSTER.isHit(math.RTP5.TriggerWeight) {
			triggerIconId = math.RTP5.TriggerIconID
			bonusTypeId = math.RTP5.Type
		}

		if MONSTER.isHit(math.RTP5.HitWeight) {
			iconPay = h.pick(math)
		}

	case "901":
		if MONSTER.isHit(math.RTP6.TriggerWeight) {
			triggerIconId = math.RTP6.TriggerIconID
			bonusTypeId = math.RTP6.Type
		}

		if MONSTER.isHit(math.RTP6.HitWeight) {
			iconPay = h.pick(math)
		}

	case "902":
		if MONSTER.isHit(math.RTP7.TriggerWeight) {
			triggerIconId = math.RTP7.TriggerIconID
			bonusTypeId = math.RTP7.Type
		}

		if MONSTER.isHit(math.RTP7.HitWeight) {
			iconPay = h.pick(math)
		}

	case "100":
		if MONSTER.isHit(math.RTP8.TriggerWeight) {
			triggerIconId = math.RTP8.TriggerIconID
			bonusTypeId = math.RTP8.Type
		}

		if MONSTER.isHit(math.RTP8.HitWeight) {
			iconPay = h.pick(math)
		}

	case "150":
		if MONSTER.isHit(math.RTP9.TriggerWeight) {
			triggerIconId = math.RTP9.TriggerIconID
			bonusTypeId = math.RTP9.Type
		}

		if MONSTER.isHit(math.RTP9.HitWeight) {
			iconPay = h.pick(math)
		}

	case "200":
		if MONSTER.isHit(math.RTP10.TriggerWeight) {
			triggerIconId = math.RTP10.TriggerIconID
			bonusTypeId = math.RTP10.Type
		}

		if MONSTER.isHit(math.RTP10.HitWeight) {
			iconPay = h.pick(math)
		}

	case "300":
		if MONSTER.isHit(math.RTP11.TriggerWeight) {
			triggerIconId = math.RTP11.TriggerIconID
			bonusTypeId = math.RTP11.Type
		}

		if MONSTER.isHit(math.RTP11.HitWeight) {
			iconPay = h.pick(math)
		}

	default:
		return iconPay, triggerIconId, bonusTypeId
	}

	return iconPay, triggerIconId, bonusTypeId
}

func (h *hotPot) pick(math *BsHotPotMath) int {
	iconPays := make([]rng.Option, 0, 9)
	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPays = append(iconPays, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPays).(int)
}
