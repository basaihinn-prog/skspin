package PSF_ON_00002_1_MONSTER

import "serve/fish_comm/rng"

var HotPot = &hotPot{}

type hotPot struct {
}

type BsHotPotMath struct {
	IconID         int                       `json:"IconID"`
	BonusID        string                    `json:"BonusID"`
	BonusTimes     []int                     `json:"BonusTimes"`
	IconPays       []int                     `json:"IconPays"`
	Multiplier     []int                     `json:"Multiplier"`
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

func (h *hotPot) Hit(rtpId string, math *BsHotPotMath) (iconPay int, triggerIconId int, bonusTypeId int, multiplier int) {
	iconPay = 0
	multiplier = 1
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case "0":
		if MONSTER.isHit(math.RTP1.TriggerWeight) {
			triggerIconId = math.RTP1.TriggerIconID
			bonusTypeId = math.RTP1.Type
		}

		if MONSTER.isHit(math.RTP1.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "20":
		if MONSTER.isHit(math.RTP2.TriggerWeight) {
			triggerIconId = math.RTP2.TriggerIconID
			bonusTypeId = math.RTP2.Type
		}

		if MONSTER.isHit(math.RTP2.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "40":
		if MONSTER.isHit(math.RTP3.TriggerWeight) {
			triggerIconId = math.RTP3.TriggerIconID
			bonusTypeId = math.RTP3.Type
		}

		if MONSTER.isHit(math.RTP3.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "60":
		if MONSTER.isHit(math.RTP4.TriggerWeight) {
			triggerIconId = math.RTP4.TriggerIconID
			bonusTypeId = math.RTP4.Type
		}

		if MONSTER.isHit(math.RTP4.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "80":
		if MONSTER.isHit(math.RTP5.TriggerWeight) {
			triggerIconId = math.RTP5.TriggerIconID
			bonusTypeId = math.RTP5.Type
		}

		if MONSTER.isHit(math.RTP5.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "901":
		if MONSTER.isHit(math.RTP6.TriggerWeight) {
			triggerIconId = math.RTP6.TriggerIconID
			bonusTypeId = math.RTP6.Type
		}

		if MONSTER.isHit(math.RTP6.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "902":
		if MONSTER.isHit(math.RTP7.TriggerWeight) {
			triggerIconId = math.RTP7.TriggerIconID
			bonusTypeId = math.RTP7.Type
		}

		if MONSTER.isHit(math.RTP7.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "100":
		if MONSTER.isHit(math.RTP8.TriggerWeight) {
			triggerIconId = math.RTP8.TriggerIconID
			bonusTypeId = math.RTP8.Type
		}

		if MONSTER.isHit(math.RTP8.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "150":
		if MONSTER.isHit(math.RTP9.TriggerWeight) {
			triggerIconId = math.RTP9.TriggerIconID
			bonusTypeId = math.RTP9.Type
		}

		if MONSTER.isHit(math.RTP9.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "200":
		if MONSTER.isHit(math.RTP10.TriggerWeight) {
			triggerIconId = math.RTP10.TriggerIconID
			bonusTypeId = math.RTP10.Type
		}

		if MONSTER.isHit(math.RTP10.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	case "300":
		if MONSTER.isHit(math.RTP11.TriggerWeight) {
			triggerIconId = math.RTP11.TriggerIconID
			bonusTypeId = math.RTP11.Type
		}

		if MONSTER.isHit(math.RTP11.HitWeight) {
			multiplier = h.pick(math)
			iconPay = math.IconPays[0]
		}

	default:
		return iconPay, triggerIconId, bonusTypeId, multiplier
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}

func (h *hotPot) pick(math *BsHotPotMath) int {
	multipliers := make([]rng.Option, 0, 7)
	for i := 0; i < len(math.IconPaysWeight); i++ {
		multipliers = append(multipliers, rng.Option{math.IconPaysWeight[i], math.Multiplier[i]})
	}

	return MONSTER.rng(multipliers).(int)
}
