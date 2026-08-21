package RKF_H5_00001_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var Yggdrasil = &yggdrasil{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "100"),
}

type yggdrasil struct {
	bonusMode bool
}

type BsYggdrasilMath struct {
	IconID             int
	BonusID            string
	BonusTimes         []int
	LowIconPays        []int
	LowIconPaysWeight  []int
	HighIconPays       []int
	HighIconPaysWeight []int
	RTP1               struct{ HitBsYggdrasilMath } `json:"RTP1"`
	RTP2               struct{ HitBsYggdrasilMath } `json:"RTP2"`
	RTP3               struct{ HitBsYggdrasilMath } `json:"RTP3"`
	RTP4               struct{ HitBsYggdrasilMath } `json:"RTP4"`
	RTP5               struct{ HitBsYggdrasilMath } `json:"RTP5"`
	RTP6               struct{ HitBsYggdrasilMath } `json:"RTP6"`
	RTP7               struct{ HitBsYggdrasilMath } `json:"RTP7"`
	RTP8               struct{ HitBsYggdrasilMath } `json:"RTP8"`
	RTP9               struct{ HitBsYggdrasilMath } `json:"RTP9"`
}

type HitBsYggdrasilMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (y *yggdrasil) Hit(rtpId string, math *BsYggdrasilMath) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		return y.rtpHit(math, &math.RTP1)
	case math.RTP2.ID:
		return y.rtpHit(math, &math.RTP2)
	case math.RTP3.ID:
		return y.rtpHit(math, &math.RTP3)
	case math.RTP4.ID:
		return y.rtpHit(math, &math.RTP4)
	case math.RTP5.ID:
		return y.rtpHit(math, &math.RTP5)
	case math.RTP6.ID:
		return y.rtpHit(math, &math.RTP6)
	case math.RTP7.ID:
		return y.rtpHit(math, &math.RTP7)
	case math.RTP8.ID:
		return y.rtpHit(math, &math.RTP8)
	case math.RTP9.ID:
		return y.rtpHit(math, &math.RTP9)
	default:
		return -1, -1, nil
	}
}

func (y *yggdrasil) pick(math *BsYggdrasilMath) []int {
	lowPayOptions := y.setOptions(math.LowIconPaysWeight, math.LowIconPays)
	highPayOptions := y.setOptions(math.HighIconPaysWeight, math.HighIconPays)

	iconPays := []int{
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(highPayOptions).(int),
	}

	if len(iconPays) != 5 {
		panic("Error:PSFM_00013_MONSTER:pick:iconPays")
	}

	return iconPays
}

func (y *yggdrasil) setOptions(weight, iconPays []int) []rng.Option {
	payOptions := make([]rng.Option, 0, len(weight))
	for i := 0; i < len(weight); i++ {
		payOptions = append(payOptions, rng.Option{weight[i], iconPays[i]})
	}

	return payOptions
}
func (y *yggdrasil) rtpHit(bsMath *BsYggdrasilMath, rtp *struct{ HitBsYggdrasilMath }) (triggerIconId, bonusTypeId int, iconPays []int) {
	triggerIconId = -1
	bonusTypeId = -1
	iconPays = nil

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPays = y.pick(bsMath)
	}

	return triggerIconId, bonusTypeId, iconPays
}
