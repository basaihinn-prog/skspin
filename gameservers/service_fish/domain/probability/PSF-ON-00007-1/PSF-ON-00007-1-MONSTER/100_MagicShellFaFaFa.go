package PSF_ON_00007_1_MONSTER

import (
	"os"
	"strings"

	"serve/fish_comm/rng"
)

var MagicShellFaFaFa = &magicShellFaFaFa{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "100"),
}

type magicShellFaFaFa struct {
	bonusMode bool
}

type BsMagicShellFaFaFaMath struct {
	IconID             int
	BonusID            string
	BonusTimes         []int
	LowIconPays        []int
	LowIconPaysWeight  []int
	HighIconPays       []int
	HighIconPaysWeight []int
	RTP1               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP1"`
	RTP2               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP2"`
	RTP3               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP3"`
	RTP4               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP4"`
	RTP5               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP5"`
	RTP6               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP6"`
	RTP7               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP7"`
	RTP8               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP8"`
	RTP9               struct{ HitBsMagicShellFaFaFaMath } `json:"RTP9"`
}

type HitBsMagicShellFaFaFaMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (m *magicShellFaFaFa) Hit(rtpId string, math *BsMagicShellFaFaFaMath) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		return m.rtpHit(math, &math.RTP1)

	case math.RTP2.ID:
		return m.rtpHit(math, &math.RTP2)

	case math.RTP3.ID:
		return m.rtpHit(math, &math.RTP3)

	case math.RTP4.ID:
		return m.rtpHit(math, &math.RTP4)

	case math.RTP5.ID:
		return m.rtpHit(math, &math.RTP5)

	case math.RTP6.ID:
		return m.rtpHit(math, &math.RTP6)

	case math.RTP7.ID:
		return m.rtpHit(math, &math.RTP7)

	case math.RTP8.ID:
		return m.rtpHit(math, &math.RTP8)

	case math.RTP9.ID:
		return m.rtpHit(math, &math.RTP9)

	default:
		return -1, -1, nil
	}

	return triggerIconId, bonusTypeId, iconPays
}

func (m *magicShellFaFaFa) pick(math *BsMagicShellFaFaFaMath) []int {
	lowPayOptions := m.setOptions(math.LowIconPaysWeight, math.LowIconPays)
	highPayOptions := m.setOptions(math.HighIconPaysWeight, math.HighIconPays)

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

func (m *magicShellFaFaFa) setOptions(weight, iconPays []int) []rng.Option {
	payOptions := make([]rng.Option, 0, len(weight))
	for i := 0; i < len(weight); i++ {
		payOptions = append(payOptions, rng.Option{weight[i], iconPays[i]})
	}

	return payOptions
}

func (m *magicShellFaFaFa) rtpHit(bsMath *BsMagicShellFaFaFaMath, rtp *struct{ HitBsMagicShellFaFaFaMath }) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1
	iconPays = nil

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPays = m.pick(bsMath)
	}

	return triggerIconId, bonusTypeId, iconPays
}
