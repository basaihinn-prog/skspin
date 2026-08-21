package PSF_ON_00005_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var RedEnvelopeTree = &redEnvelopeTree{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "100"),
}

type redEnvelopeTree struct {
	bonusMode bool
}

type BsRedEnvelopeTreeMath struct {
	IconID             int
	BonusID            string
	BonusTimes         []int
	LowIconPays        []int
	LowIconPaysWeight  []int
	HighIconPays       []int
	HighIconPaysWeight []int
	RTP1               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP1"`
	RTP2               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP2"`
	RTP3               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP3"`
	RTP4               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP4"`
	RTP5               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP5"`
	RTP6               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP6"`
	RTP7               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP7"`
	RTP8               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP8"`
	RTP9               struct{ HitBsRedEnvelopeTreeMath } `json:"RTP9"`
}

type HitBsRedEnvelopeTreeMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (r *redEnvelopeTree) Hit(rtpId string, math *BsRedEnvelopeTreeMath) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		return r.rtpHit(math, &math.RTP1)

	case math.RTP2.ID:
		return r.rtpHit(math, &math.RTP2)

	case math.RTP3.ID:
		return r.rtpHit(math, &math.RTP3)

	case math.RTP4.ID:
		return r.rtpHit(math, &math.RTP4)

	case math.RTP5.ID:
		return r.rtpHit(math, &math.RTP5)

	case math.RTP6.ID:
		return r.rtpHit(math, &math.RTP6)

	case math.RTP7.ID:
		return r.rtpHit(math, &math.RTP7)

	case math.RTP8.ID:
		return r.rtpHit(math, &math.RTP8)

	case math.RTP9.ID:
		return r.rtpHit(math, &math.RTP9)

	default:
		return -1, -1, nil
	}

	return triggerIconId, bonusTypeId, iconPays
}

func (r *redEnvelopeTree) pick(math *BsRedEnvelopeTreeMath) []int {
	lowPayOptions := r.setOptions(math.LowIconPaysWeight, math.LowIconPays)
	highPayOptions := r.setOptions(math.HighIconPaysWeight, math.HighIconPays)

	iconPays := []int{
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(highPayOptions).(int),
	}

	if len(iconPays) != 5 {
		panic("Error:PSFM_00006_MONSTER:pick:iconPays")
	}

	return iconPays
}

func (r *redEnvelopeTree) setOptions(weight, iconPays []int) []rng.Option {
	payOptions := make([]rng.Option, 0, len(weight))
	for i := 0; i < len(weight); i++ {
		payOptions = append(payOptions, rng.Option{weight[i], iconPays[i]})
	}

	return payOptions
}

func (r *redEnvelopeTree) rtpHit(bsMath *BsRedEnvelopeTreeMath, rtp *struct{ HitBsRedEnvelopeTreeMath }) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1
	iconPays = nil

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPays = r.pick(bsMath)
	}

	return triggerIconId, bonusTypeId, iconPays
}
