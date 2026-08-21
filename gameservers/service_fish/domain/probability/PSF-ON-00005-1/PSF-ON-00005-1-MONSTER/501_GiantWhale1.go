package PSF_ON_00005_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var GiantWhale1 = &giantWhale1{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "501"),
}

type giantWhale1 struct {
	bonusMode bool
}

type BsGiantWhaleMath struct {
	IconID         int                           `json:"IconID"`
	BonusID        string                        `json:"BonusID"`
	BonusTimes     []int                         `json:"BonusTimes"`
	IconPays       []int                         `json:"IconPays"`
	IconPaysWeight []int                         `json:"IconPaysWeight"`
	RTP1           struct{ HitBsGiantWhaleMath } `json:"RTP1"`
	RTP2           struct{ HitBsGiantWhaleMath } `json:"RTP2"`
	RTP3           struct{ HitBsGiantWhaleMath } `json:"RTP3"`
	RTP4           struct{ HitBsGiantWhaleMath } `json:"RTP4"`
	RTP5           struct{ HitBsGiantWhaleMath } `json:"RTP5"`
	RTP6           struct{ HitBsGiantWhaleMath } `json:"RTP6"`
	RTP7           struct{ HitBsGiantWhaleMath } `json:"RTP7"`
	RTP8           struct{ HitBsGiantWhaleMath } `json:"RTP8"`
	RTP9           struct{ HitBsGiantWhaleMath } `json:"RTP9"`
}

type HitBsGiantWhaleMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (g *giantWhale1) Hit(rtpId string, math *BsGiantWhaleMath) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		return g.rtpHit(math, &math.RTP1)

	case math.RTP2.ID:
		return g.rtpHit(math, &math.RTP2)

	case math.RTP3.ID:
		return g.rtpHit(math, &math.RTP3)

	case math.RTP4.ID:
		return g.rtpHit(math, &math.RTP4)

	case math.RTP5.ID:
		return g.rtpHit(math, &math.RTP5)

	case math.RTP6.ID:
		return g.rtpHit(math, &math.RTP6)

	case math.RTP7.ID:
		return g.rtpHit(math, &math.RTP7)

	case math.RTP8.ID:
		return g.rtpHit(math, &math.RTP8)

	case math.RTP9.ID:
		return g.rtpHit(math, &math.RTP9)

	default:
		return 0, -1, -1
	}

	return iconPay, triggerIconId, bonusTypeId
}

func (g *giantWhale1) pick(math *BsGiantWhaleMath) (iconPay int) {
	iconPays := make([]rng.Option, 0, 1)

	iconPays = append(iconPays, rng.Option{math.IconPaysWeight[0], math.IconPays[0]})
	iconPays = append(iconPays, rng.Option{math.IconPaysWeight[1], math.IconPays[1]})

	return MONSTER.rng(iconPays).(int)
}

func (g *giantWhale1) rtpHit(math *BsGiantWhaleMath, rtp *struct{ HitBsGiantWhaleMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = g.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId
}
