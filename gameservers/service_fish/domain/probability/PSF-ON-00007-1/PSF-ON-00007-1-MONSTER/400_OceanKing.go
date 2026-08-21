package PSF_ON_00007_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var OceanKing = &oceanKing{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "400"),
}

type oceanKing struct {
	bonusMode bool
}

type BsOceanKingMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	RTP1       struct{ HitBsOceanKingMath }
	RTP2       struct{ HitBsOceanKingMath }
	RTP3       struct{ HitBsOceanKingMath }
	RTP4       struct{ HitBsOceanKingMath }
	RTP5       struct{ HitBsOceanKingMath }
	RTP6       struct{ HitBsOceanKingMath }
	RTP7       struct{ HitBsOceanKingMath }
	RTP8       struct{ HitBsOceanKingMath }
	RTP9       struct{ HitBsOceanKingMath }
}

type HitBsOceanKingMath struct {
	ID             string `json:"ID"`
	HitWeight      []int  `json:"HitWeight"`
	IconPays       []int  `json:"IconPays"`
	IconPaysWeight []int  `json:"IconPaysWeight"`
	TriggerIconID  int    `json:"TriggerIconID"`
	TriggerWeight  []int  `json:"TriggerWeight"`
	Type           int    `json:"Type"`
}

func (o *oceanKing) Hit(rtpId string, math *BsOceanKingMath) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		return o.rtpHit(&math.RTP1)
	case math.RTP2.ID:
		return o.rtpHit(&math.RTP2)
	case math.RTP3.ID:
		return o.rtpHit(&math.RTP3)
	case math.RTP4.ID:
		return o.rtpHit(&math.RTP4)
	case math.RTP5.ID:
		return o.rtpHit(&math.RTP5)
	case math.RTP6.ID:
		return o.rtpHit(&math.RTP6)
	case math.RTP7.ID:
		return o.rtpHit(&math.RTP7)
	case math.RTP8.ID:
		return o.rtpHit(&math.RTP8)
	case math.RTP9.ID:
		return o.rtpHit(&math.RTP9)
	default:
		return iconPay, triggerIconId, bonusTypeId
	}
}

func (o *oceanKing) pick(math *struct{ HitBsOceanKingMath }) int {
	options := make([]rng.Option, 0, len(math.IconPays))

	for i := 0; i < len(math.IconPays); i++ {
		options = append(options, rng.Option{Weight: math.IconPaysWeight[i], Item: math.IconPays[i]})
	}

	return MONSTER.rng(options).(int)
}

func (o *oceanKing) rtpHit(rtp *struct{ HitBsOceanKingMath }) (iconPay, triggerIconId, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = o.pick(rtp)
	}

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	return iconPay, triggerIconId, bonusTypeId
}
