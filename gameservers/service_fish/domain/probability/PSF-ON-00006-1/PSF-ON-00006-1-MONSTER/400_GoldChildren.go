package PSF_ON_00006_1_MONSTER

import (
	"os"
	"strings"

	"serve/fish_comm/rng"
)

var GoldChildren = &goldChildren{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "400"),
}

type goldChildren struct {
	bonusMode bool
}

type BsGoldChildrenMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	RTP1       struct{ HitBsGoldChildrenMath }
	RTP2       struct{ HitBsGoldChildrenMath }
	RTP3       struct{ HitBsGoldChildrenMath }
	RTP4       struct{ HitBsGoldChildrenMath }
	RTP5       struct{ HitBsGoldChildrenMath }
	RTP6       struct{ HitBsGoldChildrenMath }
	RTP7       struct{ HitBsGoldChildrenMath }
	RTP8       struct{ HitBsGoldChildrenMath }
	RTP9       struct{ HitBsGoldChildrenMath }
}

type HitBsGoldChildrenMath struct {
	ID             string `json:"ID"`
	HitWeight      []int  `json:"HitWeight"`
	IconPays       []int  `json:"IconPays"`
	IconPaysWeight []int  `json:"IconPaysWeight"`
	TriggerIconID  int    `json:"TriggerIconID"`
	TriggerWeight  []int  `json:"TriggerWeight"`
	Type           int    `json:"Type"`
}

func (o *goldChildren) Hit(rtpId string, math *BsGoldChildrenMath) (iconPay, triggerIconId, bonusTypeId int) {
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

func (o *goldChildren) pick(math *struct{ HitBsGoldChildrenMath }) int {
	options := make([]rng.Option, 0, len(math.IconPays))

	for i := 0; i < len(math.IconPays); i++ {
		options = append(options, rng.Option{Weight: math.IconPaysWeight[i], Item: math.IconPays[i]})
	}

	return MONSTER.rng(options).(int)
}

func (o *goldChildren) rtpHit(rtp *struct{ HitBsGoldChildrenMath }) (iconPay, triggerIconId, bonusTypeId int) {
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
