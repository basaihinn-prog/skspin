package PSF_ON_00002_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var RedEnvelope = &redEnvelope{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "28"),
}

type redEnvelope struct {
	bonusMode bool
}

type BsRedEnvelopeMath struct {
	IconID             int
	BonusID            string
	BonusTimes         []int
	LowIconPays        []int
	LowIconPaysWeight  []int
	HighIconPays       []int
	HighIconPaysWeight []int
	RTP1               struct{ HitBsRedEnvelopeMath } `json:"RTP1"`
	RTP2               struct{ HitBsRedEnvelopeMath } `json:"RTP2"`
	RTP3               struct{ HitBsRedEnvelopeMath } `json:"RTP3"`
	RTP4               struct{ HitBsRedEnvelopeMath } `json:"RTP4"`
	RTP5               struct{ HitBsRedEnvelopeMath } `json:"RTP5"`
	RTP6               struct{ HitBsRedEnvelopeMath } `json:"RTP6"`
	RTP7               struct{ HitBsRedEnvelopeMath } `json:"RTP7"`
	RTP8               struct{ HitBsRedEnvelopeMath } `json:"RTP8"`
	RTP9               struct{ HitBsRedEnvelopeMath } `json:"RTP9"`
	RTP10              struct{ HitBsRedEnvelopeMath } `json:"RTP10"`
	RTP11              struct{ HitBsRedEnvelopeMath } `json:"RTP11"`
}

type HitBsRedEnvelopeMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (r *redEnvelope) pick(math *BsRedEnvelopeMath) []int {
	lowPayOptions := make([]rng.Option, 0, 6)
	for i := 0; i < len(math.LowIconPaysWeight); i++ {
		lowPayOptions = append(lowPayOptions, rng.Option{math.LowIconPaysWeight[i], math.LowIconPays[i]})
	}

	highPayOptions := make([]rng.Option, 0, 5)
	for i := 0; i < len(math.HighIconPaysWeight); i++ {
		highPayOptions = append(highPayOptions, rng.Option{math.HighIconPaysWeight[i], math.HighIconPays[i]})
	}

	iconPays := []int{
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(highPayOptions).(int),
	}

	if len(iconPays) != 5 {
		panic("Error:PSFM_00001_MONSTER:pick:iconPays")
	}

	return iconPays
}

func (r *redEnvelope) Hit(rtpId string, math *BsRedEnvelopeMath) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case "0":
		if MONSTER.isHit(math.RTP1.TriggerWeight) {
			triggerIconId = math.RTP1.TriggerIconID
			bonusTypeId = math.RTP1.Type
		}

		if MONSTER.isHit(math.RTP1.HitWeight) {
			iconPays = r.pick(math)
		}

	case "20":
		if MONSTER.isHit(math.RTP2.TriggerWeight) {
			triggerIconId = math.RTP2.TriggerIconID
			bonusTypeId = math.RTP2.Type
		}

		if MONSTER.isHit(math.RTP2.HitWeight) {
			iconPays = r.pick(math)
		}

	case "40":
		if MONSTER.isHit(math.RTP3.TriggerWeight) {
			triggerIconId = math.RTP3.TriggerIconID
			bonusTypeId = math.RTP3.Type
		}

		if MONSTER.isHit(math.RTP3.HitWeight) {
			iconPays = r.pick(math)
		}

	case "60":
		if MONSTER.isHit(math.RTP4.TriggerWeight) {
			triggerIconId = math.RTP4.TriggerIconID
			bonusTypeId = math.RTP4.Type
		}

		if MONSTER.isHit(math.RTP4.HitWeight) {
			iconPays = r.pick(math)
		}

	case "80":
		if MONSTER.isHit(math.RTP5.TriggerWeight) {
			triggerIconId = math.RTP5.TriggerIconID
			bonusTypeId = math.RTP5.Type
		}

		if MONSTER.isHit(math.RTP5.HitWeight) {
			iconPays = r.pick(math)
		}

	case "901":
		if MONSTER.isHit(math.RTP6.TriggerWeight) {
			triggerIconId = math.RTP6.TriggerIconID
			bonusTypeId = math.RTP6.Type
		}

		if MONSTER.isHit(math.RTP6.HitWeight) {
			iconPays = r.pick(math)
		}

	case "902":
		if MONSTER.isHit(math.RTP7.TriggerWeight) {
			triggerIconId = math.RTP7.TriggerIconID
			bonusTypeId = math.RTP7.Type
		}

		if MONSTER.isHit(math.RTP7.HitWeight) {
			iconPays = r.pick(math)
		}

	case "100":
		if MONSTER.isHit(math.RTP8.TriggerWeight) {
			triggerIconId = math.RTP8.TriggerIconID
			bonusTypeId = math.RTP8.Type
		}

		if MONSTER.isHit(math.RTP8.HitWeight) {
			iconPays = r.pick(math)
		}

	case "150":
		if MONSTER.isHit(math.RTP9.TriggerWeight) {
			triggerIconId = math.RTP9.TriggerIconID
			bonusTypeId = math.RTP9.Type
		}

		if MONSTER.isHit(math.RTP9.HitWeight) {
			iconPays = r.pick(math)
		}

	case "200":
		if MONSTER.isHit(math.RTP10.TriggerWeight) {
			triggerIconId = math.RTP10.TriggerIconID
			bonusTypeId = math.RTP10.Type
		}

		if MONSTER.isHit(math.RTP10.HitWeight) {
			iconPays = r.pick(math)
		}

	case "300":
		if MONSTER.isHit(math.RTP11.TriggerWeight) {
			triggerIconId = math.RTP11.TriggerIconID
			bonusTypeId = math.RTP11.Type
		}

		if MONSTER.isHit(math.RTP11.HitWeight) {
			iconPays = r.pick(math)
		}

	default:
		return -1, -1, nil
	}

	return triggerIconId, bonusTypeId, iconPays
}
