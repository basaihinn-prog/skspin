package PSF_ON_00007_1_MONSTER

import (
	"os"
	"strings"

	"serve/fish_comm/rng"
)

var TreasureHermitCrab = &treasureHermitCrab{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "302"),
}

type treasureHermitCrab struct {
	bonusMode bool
}

type BsTreasureHermitCrabMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	RTP1       struct{ HitBsTreasureHermitCrabMath }
	RTP2       struct{ HitBsTreasureHermitCrabMath }
	RTP3       struct{ HitBsTreasureHermitCrabMath }
	RTP4       struct{ HitBsTreasureHermitCrabMath }
	RTP5       struct{ HitBsTreasureHermitCrabMath }
	RTP6       struct{ HitBsTreasureHermitCrabMath }
	RTP7       struct{ HitBsTreasureHermitCrabMath }
	RTP8       struct{ HitBsTreasureHermitCrabMath }
	RTP9       struct{ HitBsTreasureHermitCrabMath }
}

type HitBsTreasureHermitCrabMath struct {
	ID               string
	HitWeight        []int
	GainAmount       []int                                  `json:"Gain amount"`
	GainAmountWeight []int                                  `json:"Gain amount weight"`
	Gain1            struct{ GainBsTreasureHermitCrabMath } `json:"1 Gain "`
	Gain2            struct{ GainBsTreasureHermitCrabMath } `json:"2 Gain "`
	Gain3            struct{ GainBsTreasureHermitCrabMath } `json:"3 Gain "`
	Gain4            struct{ GainBsTreasureHermitCrabMath } `json:"4 Gain "`
	Gain5            struct{ GainBsTreasureHermitCrabMath } `json:"5 Gain "`
	TriggerIconID    int                                    `json:"TriggerIconID"`
	TriggerWeight    []int                                  `json:"TriggerWeight"`
	Type             int                                    `json:"Type"`
}

type GainBsTreasureHermitCrabMath struct {
	IconPays       []int
	IconPaysWeight []int
}

func (t *treasureHermitCrab) Hit(rtpId string, math *BsTreasureHermitCrabMath) (totalPay int, iconPays, iconPaysPick []int, triggerIconId, bonusTypeId int) {
	totalPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case math.RTP1.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP1)

	case math.RTP2.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP2)

	case math.RTP3.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP3)

	case math.RTP4.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP4)

	case math.RTP5.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP5)

	case math.RTP6.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP6)

	case math.RTP7.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP7)

	case math.RTP8.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP8)

	case math.RTP9.ID:
		totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId = t.rtpHit(&math.RTP9)

	default:
		return 0, nil, nil, -1, -1
	}

	return totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId
}

func (t *treasureHermitCrab) rngIconPay(gainAmount int, math *struct{ HitBsTreasureHermitCrabMath }) []int {
	iconPay := make([]int, 0, 4)

	for i := 0; i < 5; i++ {
		switch gainAmount {
		case 1:
			iconPay = append(iconPay, t.rngIconPays(&math.Gain1))

		case 2:
			iconPay = append(iconPay, t.rngIconPays(&math.Gain2))

		case 3:
			iconPay = append(iconPay, t.rngIconPays(&math.Gain3))

		case 4:
			iconPay = append(iconPay, t.rngIconPays(&math.Gain4))

		case 5:
			iconPay = append(iconPay, t.rngIconPays(&math.Gain5))

		default:
			return nil
		}
	}

	return iconPay
}

func (t *treasureHermitCrab) rngIconPays(math *struct{ GainBsTreasureHermitCrabMath }) int {
	iconPayOptions := make([]rng.Option, 0, len(math.IconPaysWeight))

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPayOptions = append(iconPayOptions, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPayOptions).(int)
}

func (t *treasureHermitCrab) rngGain(math *struct{ HitBsTreasureHermitCrabMath }) int {
	amountOptions := make([]rng.Option, 0, len(math.GainAmountWeight))

	for i := 0; i < len(math.GainAmountWeight); i++ {
		amountOptions = append(amountOptions, rng.Option{math.GainAmountWeight[i], math.GainAmount[i]})
	}

	return MONSTER.rng(amountOptions).(int)
}

func (t *treasureHermitCrab) Rtp(order int, math *BsTreasureHermitCrabMath) (rtpId string) {
	switch order {
	case 1:
		return math.RTP1.ID

	case 2:
		return math.RTP2.ID

	case 3:
		return math.RTP3.ID

	case 4:
		return math.RTP4.ID

	case 5:
		return math.RTP5.ID

	default:
		return ""
	}
}

func (t *treasureHermitCrab) rngPick(gainAmount int) []int {
	var pickOrder []int
	iconPaysPick := []int{0, 0, 0, 0, 0}

	options0 := t.remove(nil, -1)
	index0 := t.rngIndex(options0)
	pickOrder = append(pickOrder, index0)

	options1 := t.remove(options0, index0)
	index1 := t.rngIndex(options1)
	pickOrder = append(pickOrder, index1)

	options2 := t.remove(options1, index1)
	index2 := t.rngIndex(options2)
	pickOrder = append(pickOrder, index2)

	options3 := t.remove(options2, index2)
	index3 := t.rngIndex(options3)
	pickOrder = append(pickOrder, index3)

	options4 := t.remove(options3, index3)
	index4 := t.rngIndex(options4)
	pickOrder = append(pickOrder, index4)

	for i := 0; i < gainAmount; i++ {
		iconPaysPick[pickOrder[i]] = 1
	}

	return iconPaysPick
}

func (t *treasureHermitCrab) rngIndex(options []rng.Option) (index int) {
	return rng.New(options).Item.(int)
}

func (t *treasureHermitCrab) remove(indexOptions []rng.Option, indexOptionItem int) []rng.Option {
	if indexOptionItem == -1 {
		var indexOptions []rng.Option
		indexOptions = append(indexOptions, rng.Option{1, 0})
		indexOptions = append(indexOptions, rng.Option{1, 1})
		indexOptions = append(indexOptions, rng.Option{1, 2})
		indexOptions = append(indexOptions, rng.Option{1, 3})
		indexOptions = append(indexOptions, rng.Option{1, 4})
		return indexOptions
	}

	var newIndexOptions []rng.Option

	for _, v := range indexOptions {
		if v.Item.(int) != indexOptionItem {
			newIndexOptions = append(newIndexOptions, v)
		}
	}

	return newIndexOptions
}

func (t *treasureHermitCrab) rtpHit(rtp *struct{ HitBsTreasureHermitCrabMath }) (totalPay int, iconPays, iconPaysPick []int, triggerIconId, bonusTypeId int) {
	totalPay = 0
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.HitWeight) {
		gainAmount := t.rngGain(rtp)
		iconPays = t.rngIconPay(gainAmount, rtp)
		iconPaysPick = t.rngPick(gainAmount)
	}

	if len(iconPays) > 0 {
		for i := 0; i < len(iconPays); i++ {
			if iconPaysPick[i] == 1 {
				totalPay += iconPays[i]
			}
		}
	}

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	return totalPay, iconPays, iconPaysPick, triggerIconId, bonusTypeId
}
