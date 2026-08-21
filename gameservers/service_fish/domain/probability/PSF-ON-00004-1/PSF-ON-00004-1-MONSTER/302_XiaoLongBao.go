package PSF_ON_00004_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var XiaoLongBao = &xiaoLongBao{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "302"),
}

type xiaoLongBao struct {
	bonusMode bool
}

type BsXiaoLongBaoMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	RTP1       struct{ HitBsXiaoLongBaoMath }
	RTP2       struct{ HitBsXiaoLongBaoMath }
	RTP3       struct{ HitBsXiaoLongBaoMath }
	RTP4       struct{ HitBsXiaoLongBaoMath }
	RTP5       struct{ HitBsXiaoLongBaoMath }
}

type HitBsXiaoLongBaoMath struct {
	ID               string
	HitWeight        []int
	AvgPay           []int
	GainAmount       []int                           `json:"Gain amount"`
	GainAmountWeight []int                           `json:"Gain amount weight"`
	Gain1            struct{ GainBsXiaoLongBaoMath } `json:"1 Gain "`
	Gain2            struct{ GainBsXiaoLongBaoMath } `json:"2 Gain "`
	Gain3            struct{ GainBsXiaoLongBaoMath } `json:"3 Gain "`
	Gain4            struct{ GainBsXiaoLongBaoMath } `json:"4 Gain "`
	Gain5            struct{ GainBsXiaoLongBaoMath } `json:"5 Gain "`
}

type GainBsXiaoLongBaoMath struct {
	IconPays       []int
	IconPaysWeight []int
}

func (x *xiaoLongBao) Hit(rtpId string, math *BsXiaoLongBaoMath) (totalPay int, iconPays, iconPaysPick []int, avgPay int) {
	totalPay = 0

	switch rtpId {
	case math.RTP1.ID:
		totalPay, iconPays, iconPaysPick, avgPay = x.rtpHit(&math.RTP1)

	case math.RTP2.ID:
		totalPay, iconPays, iconPaysPick, avgPay = x.rtpHit(&math.RTP2)

	case math.RTP3.ID:
		totalPay, iconPays, iconPaysPick, avgPay = x.rtpHit(&math.RTP3)

	case math.RTP4.ID:
		totalPay, iconPays, iconPaysPick, avgPay = x.rtpHit(&math.RTP4)

	case math.RTP5.ID:
		totalPay, iconPays, iconPaysPick, avgPay = x.rtpHit(&math.RTP5)

	default:
		return 0, nil, nil, 0
	}

	return totalPay, iconPays, iconPaysPick, avgPay
}

func (x *xiaoLongBao) rngIconPay(gainAmount int, math *struct{ HitBsXiaoLongBaoMath }) []int {
	iconPay := make([]int, 0, 4)

	for i := 0; i < 5; i++ {
		switch gainAmount {
		case 1:
			iconPay = append(iconPay, x.rngIconPays(&math.Gain1))

		case 2:
			iconPay = append(iconPay, x.rngIconPays(&math.Gain2))

		case 3:
			iconPay = append(iconPay, x.rngIconPays(&math.Gain3))

		case 4:
			iconPay = append(iconPay, x.rngIconPays(&math.Gain4))

		case 5:
			iconPay = append(iconPay, x.rngIconPays(&math.Gain5))

		default:
			return nil
		}
	}

	return iconPay
}

func (x *xiaoLongBao) rngIconPays(math *struct{ GainBsXiaoLongBaoMath }) int {
	iconPayOptions := make([]rng.Option, 0, len(math.IconPaysWeight))

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPayOptions = append(iconPayOptions, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPayOptions).(int)
}

func (x *xiaoLongBao) rngGain(math *struct{ HitBsXiaoLongBaoMath }) int {
	amountOptions := make([]rng.Option, 0, len(math.GainAmountWeight))

	for i := 0; i < len(math.GainAmountWeight); i++ {
		amountOptions = append(amountOptions, rng.Option{math.GainAmountWeight[i], math.GainAmount[i]})
	}

	return MONSTER.rng(amountOptions).(int)
}

func (x *xiaoLongBao) Rtp(order int, math *BsXiaoLongBaoMath) (rtpId string) {
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

func (x *xiaoLongBao) rngPick(gainAmount int) []int {
	var pickOrder []int
	iconPaysPick := []int{0, 0, 0, 0, 0}

	options0 := x.remove(nil, -1)
	index0 := x.rngIndex(options0)
	pickOrder = append(pickOrder, index0)

	options1 := x.remove(options0, index0)
	index1 := x.rngIndex(options1)
	pickOrder = append(pickOrder, index1)

	options2 := x.remove(options1, index1)
	index2 := x.rngIndex(options2)
	pickOrder = append(pickOrder, index2)

	options3 := x.remove(options2, index2)
	index3 := x.rngIndex(options3)
	pickOrder = append(pickOrder, index3)

	options4 := x.remove(options3, index3)
	index4 := x.rngIndex(options4)
	pickOrder = append(pickOrder, index4)

	for i := 0; i < gainAmount; i++ {
		iconPaysPick[pickOrder[i]] = 1
	}

	return iconPaysPick
}

func (x *xiaoLongBao) rngIndex(options []rng.Option) (index int) {
	return rng.New(options).Item.(int)
}

func (x *xiaoLongBao) remove(indexOptions []rng.Option, indexOptionItem int) []rng.Option {
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

func (x *xiaoLongBao) AvgPay(math *BsXiaoLongBaoMath) int {
	return int(math.RTP1.AvgPay[0])
}

func (x *xiaoLongBao) rtpHit(rtp *struct{ HitBsXiaoLongBaoMath }) (totalPay int, iconPays, iconPaysPick []int, avgPay int) {
	totalPay = 0
	avgPay = rtp.AvgPay[0]

	if MONSTER.isHit(rtp.HitWeight) {
		gainAmount := x.rngGain(rtp)
		iconPays = x.rngIconPay(gainAmount, rtp)
		iconPaysPick = x.rngPick(gainAmount)
	}

	if len(iconPays) > 0 {
		for i := 0; i < len(iconPays); i++ {
			if iconPaysPick[i] == 1 {
				totalPay += iconPays[i]
			}
		}
	}

	return totalPay, iconPays, iconPaysPick, avgPay
}
