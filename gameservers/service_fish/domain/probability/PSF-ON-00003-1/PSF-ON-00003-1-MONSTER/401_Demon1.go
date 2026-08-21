package PSF_ON_00003_1_MONSTER

import "serve/fish_comm/rng"

var Demon1 = &demon1{}

type demon1 struct {
}

type BsDemonMath struct {
	IconID     int                      `json:"IconID"`
	BonusID    string                   `json:"BonusID"`
	BonusTimes []int                    `json:"BonusTimes"`
	RTP1       struct{ HitBsDemonMath } `json:"RTP1"`
	RTP2       struct{ HitBsDemonMath } `json:"RTP2"`
	RTP3       struct{ HitBsDemonMath } `json:"RTP3"`
	RTP4       struct{ HitBsDemonMath } `json:"RTP4"`
	RTP5       struct{ HitBsDemonMath } `json:"RTP5"`
	RTP6       struct{ HitBsDemonMath } `json:"RTP6"`
	RTP7       struct{ HitBsDemonMath } `json:"RTP7"`
	RTP8       struct{ HitBsDemonMath } `json:"RTP8"`
	RTP9       struct{ HitBsDemonMath } `json:"RTP9"`
	RTP10      struct{ HitBsDemonMath } `json:"RTP10"`
}

type HitBsDemonMath struct {
	ID             string `json:"ID"`
	HitWeight      []int  `json:"HitWeight"`
	IconPays       []int  `json:"IconPays"`
	IconPaysWeight []int  `json:"IconPaysWeight"`
}

func (d *demon1) Hit(rtpId string, math *BsDemonMath) (iconPay int) {
	iconPay = 0

	switch rtpId {
	case math.RTP1.ID:
		iconPay = d.rtpHit(&math.RTP1)
	case math.RTP2.ID:
		iconPay = d.rtpHit(&math.RTP2)
	case math.RTP3.ID:
		iconPay = d.rtpHit(&math.RTP3)
	case math.RTP4.ID:
		iconPay = d.rtpHit(&math.RTP4)
	case math.RTP5.ID:
		iconPay = d.rtpHit(&math.RTP5)
	case math.RTP6.ID:
		iconPay = d.rtpHit(&math.RTP6)
	case math.RTP7.ID:
		iconPay = d.rtpHit(&math.RTP7)
	case math.RTP8.ID:
		iconPay = d.rtpHit(&math.RTP8)
	case math.RTP9.ID:
		iconPay = d.rtpHit(&math.RTP9)
	case math.RTP10.ID:
		iconPay = d.rtpHit(&math.RTP10)
	default:
		return 0
	}

	return iconPay
}

func (d *demon1) pick(math *struct{ HitBsDemonMath }) int {
	iconPays := make([]rng.Option, 0, len(math.IconPays))

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPays = append(iconPays, rng.Option{
			math.IconPaysWeight[i],
			math.IconPays[i],
		})
	}

	return MONSTER.rng(iconPays).(int)
}

func (d *demon1) Rtp(order int, math *BsDemonMath) (rtpId string) {
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
func (d *demon1) rtpHit(math *struct{ HitBsDemonMath }) (iconPay int) {
	if MONSTER.isHit(math.HitWeight) {
		iconPay = d.pick(math)
	}

	return iconPay
}
