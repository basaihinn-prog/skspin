package RKF_H5_00001_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var Poseidon1 = &poseidon1{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "501"),
}

type poseidon1 struct {
	bonusMode bool
}

type BsPoseidonMath struct {
	IconID         int                         `json:"IconID"`
	BonusID        string                      `json:"BonusID"`
	BonusTimes     []int                       `json:"BonusTimes"`
	IconPays       []int                       `json:"IconPays"`
	IconPaysWeight []int                       `json:"IconPaysWeight"`
	RTP1           struct{ HitBsPoseidonMath } `json:"RTP1"`
	RTP2           struct{ HitBsPoseidonMath } `json:"RTP2"`
	RTP3           struct{ HitBsPoseidonMath } `json:"RTP3"`
	RTP4           struct{ HitBsPoseidonMath } `json:"RTP4"`
	RTP5           struct{ HitBsPoseidonMath } `json:"RTP5"`
	RTP6           struct{ HitBsPoseidonMath } `json:"RTP6"`
	RTP7           struct{ HitBsPoseidonMath } `json:"RTP7"`
	RTP8           struct{ HitBsPoseidonMath } `json:"RTP8"`
	RTP9           struct{ HitBsPoseidonMath } `json:"RTP9"`
}

type HitBsPoseidonMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (p *poseidon1) Hit(rtpId string, math *BsPoseidonMath) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		return p.rtpHit(math, &math.RTP1)
	case math.RTP2.ID:
		return p.rtpHit(math, &math.RTP2)
	case math.RTP3.ID:
		return p.rtpHit(math, &math.RTP3)
	case math.RTP4.ID:
		return p.rtpHit(math, &math.RTP4)
	case math.RTP5.ID:
		return p.rtpHit(math, &math.RTP5)
	case math.RTP6.ID:
		return p.rtpHit(math, &math.RTP6)
	case math.RTP7.ID:
		return p.rtpHit(math, &math.RTP7)
	case math.RTP8.ID:
		return p.rtpHit(math, &math.RTP8)
	case math.RTP9.ID:
		return p.rtpHit(math, &math.RTP9)
	default:
		return 0, -1, -1
	}
}

func (p *poseidon1) pick(math *BsPoseidonMath) (iconPay int) {
	iconPays := make([]rng.Option, 0, 1)

	iconPays = append(iconPays, rng.Option{math.IconPaysWeight[0], math.IconPays[0]})
	iconPays = append(iconPays, rng.Option{math.IconPaysWeight[1], math.IconPays[1]})

	return MONSTER.rng(iconPays).(int)
}

func (p *poseidon1) rtpHit(math *BsPoseidonMath, rtp *struct{ HitBsPoseidonMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = p.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId
}
