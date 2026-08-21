package PSF_ON_00005_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var FruitDish = &fruitDish{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "300"),
}

type fruitDish struct {
	bonusMode bool
}

type BsFruitDish struct {
	IconID         int                      `json:"IconID"`
	BonusID        string                   `json:"BonusID"`
	BonusTimes     []int                    `json:"BonusTimes"`
	IconPays       []int                    `json:"IconPays"`
	IconPaysWeight []int                    `json:"IconPaysWeight"`
	RTP1           struct{ HitBsFruitDish } `json:"RTP1"`
	RTP2           struct{ HitBsFruitDish } `json:"RTP2"`
	RTP3           struct{ HitBsFruitDish } `json:"RTP3"`
	RTP4           struct{ HitBsFruitDish } `json:"RTP4"`
	RTP5           struct{ HitBsFruitDish } `json:"RTP5"`
	RTP6           struct{ HitBsFruitDish } `json:"RTP6"`
	RTP7           struct{ HitBsFruitDish } `json:"RTP7"`
	RTP8           struct{ HitBsFruitDish } `json:"RTP8"`
	RTP9           struct{ HitBsFruitDish } `json:"RTP9"`
}

type HitBsFruitDish struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (f *fruitDish) Hit(rtpId string, math *BsFruitDish) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	multiplier = 0

	switch rtpId {
	case math.RTP1.ID:
		return f.rtpHit(math, &math.RTP1)

	case math.RTP2.ID:
		return f.rtpHit(math, &math.RTP2)

	case math.RTP3.ID:
		return f.rtpHit(math, &math.RTP3)

	case math.RTP4.ID:
		return f.rtpHit(math, &math.RTP4)

	case math.RTP5.ID:
		return f.rtpHit(math, &math.RTP5)

	case math.RTP6.ID:
		return f.rtpHit(math, &math.RTP6)

	case math.RTP7.ID:
		return f.rtpHit(math, &math.RTP7)

	case math.RTP8.ID:
		return f.rtpHit(math, &math.RTP8)

	case math.RTP9.ID:
		return f.rtpHit(math, &math.RTP9)

	default:
		return 0, -1, -1, 0
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}

func (f *fruitDish) pick(math *BsFruitDish) int {
	iconPays := make([]rng.Option, 0, len(math.IconPaysWeight))

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPays = append(iconPays, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPays).(int)
}

func (f *fruitDish) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}

func (f *fruitDish) rtpHit(math *BsFruitDish, rtp *struct{ HitBsFruitDish }) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	multiplier = 0

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = 1
		multiplier = f.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}
