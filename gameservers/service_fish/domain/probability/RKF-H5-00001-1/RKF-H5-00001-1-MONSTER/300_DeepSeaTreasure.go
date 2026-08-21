package RKF_H5_00001_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var DeepSeaTreasure = &deepSeaTreasure{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "300"),
}

type deepSeaTreasure struct {
	bonusMode bool
}

type BsDeepSeaTreasureMath struct {
	IconID         int                                `json:"IconID"`
	BonusID        string                             `json:"BonusID"`
	BonusTimes     []int                              `json:"BonusTimes"`
	IconPays       []int                              `json:"IconPays"`
	IconPaysWeight []int                              `json:"IconPaysWeight"`
	RTP1           struct{ HitBsDeepSeaTreasureMath } `json:"RTP1"`
	RTP2           struct{ HitBsDeepSeaTreasureMath } `json:"RTP2"`
	RTP3           struct{ HitBsDeepSeaTreasureMath } `json:"RTP3"`
	RTP4           struct{ HitBsDeepSeaTreasureMath } `json:"RTP4"`
	RTP5           struct{ HitBsDeepSeaTreasureMath } `json:"RTP5"`
	RTP6           struct{ HitBsDeepSeaTreasureMath } `json:"RTP6"`
	RTP7           struct{ HitBsDeepSeaTreasureMath } `json:"RTP7"`
	RTP8           struct{ HitBsDeepSeaTreasureMath } `json:"RTP8"`
	RTP9           struct{ HitBsDeepSeaTreasureMath } `json:"RTP9"`
}

type HitBsDeepSeaTreasureMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (d *deepSeaTreasure) Hit(rtpId string, math *BsDeepSeaTreasureMath) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	multiplier = 0

	switch rtpId {
	case math.RTP1.ID:
		return d.rtpHit(math, &math.RTP1)
	case math.RTP2.ID:
		return d.rtpHit(math, &math.RTP2)
	case math.RTP3.ID:
		return d.rtpHit(math, &math.RTP3)
	case math.RTP4.ID:
		return d.rtpHit(math, &math.RTP4)
	case math.RTP5.ID:
		return d.rtpHit(math, &math.RTP5)
	case math.RTP6.ID:
		return d.rtpHit(math, &math.RTP6)
	case math.RTP7.ID:
		return d.rtpHit(math, &math.RTP7)
	case math.RTP8.ID:
		return d.rtpHit(math, &math.RTP8)
	case math.RTP9.ID:
		return d.rtpHit(math, &math.RTP9)
	default:
		return 0, -1, -1, 0
	}
}

func (d *deepSeaTreasure) pick(math *BsDeepSeaTreasureMath) int {
	iconPays := make([]rng.Option, 0, len(math.IconPaysWeight))

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPays = append(iconPays, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPays).(int)
}

func (d *deepSeaTreasure) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}

func (d *deepSeaTreasure) rtpHit(math *BsDeepSeaTreasureMath, rtp *struct{ HitBsDeepSeaTreasureMath }) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
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
		multiplier = d.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}
