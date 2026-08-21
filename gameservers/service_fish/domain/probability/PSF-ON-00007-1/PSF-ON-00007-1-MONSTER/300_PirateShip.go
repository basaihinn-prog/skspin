package PSF_ON_00007_1_MONSTER

import (
	"os"
	"strings"

	"serve/fish_comm/rng"
)

var PirateShip = &pirateShip{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "300"),
}

type pirateShip struct {
	bonusMode bool
}

type BsPirateShip struct {
	IconID         int                       `json:"IconID"`
	BonusID        string                    `json:"BonusID"`
	BonusTimes     []int                     `json:"BonusTimes"`
	IconPays       []int                     `json:"IconPays"`
	IconPaysWeight []int                     `json:"IconPaysWeight"`
	RTP1           struct{ HitBsPirateShip } `json:"RTP1"`
	RTP2           struct{ HitBsPirateShip } `json:"RTP2"`
	RTP3           struct{ HitBsPirateShip } `json:"RTP3"`
	RTP4           struct{ HitBsPirateShip } `json:"RTP4"`
	RTP5           struct{ HitBsPirateShip } `json:"RTP5"`
	RTP6           struct{ HitBsPirateShip } `json:"RTP6"`
	RTP7           struct{ HitBsPirateShip } `json:"RTP7"`
	RTP8           struct{ HitBsPirateShip } `json:"RTP8"`
	RTP9           struct{ HitBsPirateShip } `json:"RTP9"`
}

type HitBsPirateShip struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (p *pirateShip) Hit(rtpId string, math *BsPirateShip) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	multiplier = 0

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
		return 0, -1, -1, 0
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}

func (p *pirateShip) pick(math *BsPirateShip) int {
	iconPays := make([]rng.Option, 0, len(math.IconPaysWeight))

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPays = append(iconPays, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPays).(int)
}

func (p *pirateShip) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}

func (p *pirateShip) rtpHit(math *BsPirateShip, rtp *struct{ HitBsPirateShip }) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
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
		multiplier = p.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, multiplier
}
