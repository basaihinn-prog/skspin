package RKF_H5_00001_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var ThunderHammer = &thunderHammer{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "201"),
}

type thunderHammer struct {
	bonusMode bool
}

type BsThunderHammerMath struct {
	IconID           int                              `json:"IconID"`
	BonusID          string                           `json:"BonusID"`
	UseRTP           string                           `json:"UseRTP"`
	BonusTimes       []int                            `json:"BonusTimes"`
	BonusTimesWeight []int                            `json:"BonusTimesWeight"`
	BonusIconID      []int                            `json:"BonusIconID"`
	IconPays         []int                            `json:"IconPays"`
	RTP1             struct{ HitBsThunderHammerMath } `json:"RTP1"`
	RTP2             struct{ HitBsThunderHammerMath } `json:"RTP2"`
	RTP3             struct{ HitBsThunderHammerMath } `json:"RTP3"`
	RTP4             struct{ HitBsThunderHammerMath } `json:"RTP4"`
	RTP5             struct{ HitBsThunderHammerMath } `json:"RTP5"`
	RTP6             struct{ HitBsThunderHammerMath } `json:"RTP6"`
	RTP7             struct{ HitBsThunderHammerMath } `json:"RTP7"`
	RTP8             struct{ HitBsThunderHammerMath } `json:"RTP8"`
	RTP9             struct{ HitBsThunderHammerMath } `json:"RTP9"`
}

type HitBsThunderHammerMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (t *thunderHammer) Hit(rtpId string, math *BsThunderHammerMath) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	bonusTimes = 0

	switch rtpId {
	case math.RTP1.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP1)
	case math.RTP2.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP2)
	case math.RTP3.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP3)
	case math.RTP4.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP4)
	case math.RTP5.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP5)
	case math.RTP6.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP6)
	case math.RTP7.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP7)
	case math.RTP8.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP8)
	case math.RTP9.ID:
		return t.rtpHit(math.IconPays, math, &math.RTP9)
	default:
		return 0, -1, -1, 0
	}
}

func (t *thunderHammer) UseRtp(math *BsThunderHammerMath) (rtpId string) {
	return math.UseRTP
}

func (t *thunderHammer) pick(math *BsThunderHammerMath) int {
	bonusTimes := make([]rng.Option, 0, len(math.BonusTimes))

	for i := 0; i < len(math.BonusTimesWeight); i++ {
		bonusTimes = append(bonusTimes, rng.Option{math.BonusTimesWeight[i], math.BonusTimes[i]})
	}

	return MONSTER.rng(bonusTimes).(int)
}

func (t *thunderHammer) rtpHit(iconPays []int,
	math *BsThunderHammerMath, rtp *struct{ HitBsThunderHammerMath }) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	bonusTimes = 0

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = iconPays[0]
		bonusTimes = t.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, bonusTimes
}
