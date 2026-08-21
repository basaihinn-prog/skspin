package PSF_ON_00004_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var WhiteTigerChef = &whiteTigerChef{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "400"),
}

type whiteTigerChef struct {
	bonusMode bool
}

type BsWhiteTigerChefMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	RTP1       struct{ HitBsWhiteTigerChefMath }
	RTP2       struct{ HitBsWhiteTigerChefMath }
	RTP3       struct{ HitBsWhiteTigerChefMath }
	RTP4       struct{ HitBsWhiteTigerChefMath }
	RTP5       struct{ HitBsWhiteTigerChefMath }
}

type HitBsWhiteTigerChefMath struct {
	ID             string
	AvgPay         []float32
	HitWeight      []int
	IconPays       []int
	IconPaysWeight []int
}

func (w *whiteTigerChef) Hit(rtpId string, math *BsWhiteTigerChefMath) (iconPay, avgPay int) {
	iconPay = 0

	switch rtpId {
	case math.RTP1.ID:
		iconPay, avgPay = w.rtpHit(&math.RTP1)

	case math.RTP2.ID:
		iconPay, avgPay = w.rtpHit(&math.RTP2)

	case math.RTP3.ID:
		iconPay, avgPay = w.rtpHit(&math.RTP3)

	case math.RTP4.ID:
		iconPay, avgPay = w.rtpHit(&math.RTP4)

	case math.RTP5.ID:
		iconPay, avgPay = w.rtpHit(&math.RTP5)

	default:
		return 0, 0
	}

	return iconPay, avgPay
}

func (w *whiteTigerChef) pick(math *struct{ HitBsWhiteTigerChefMath }) int {
	iconPays := make([]rng.Option, 0, len(math.IconPaysWeight))

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPays = append(iconPays, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPays).(int)
}

func (w *whiteTigerChef) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (w *whiteTigerChef) Rtp(order int, math *BsWhiteTigerChefMath) (rtpId string) {
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

func (w *whiteTigerChef) AvgPay(math *BsWhiteTigerChefMath) int {
	return int(math.RTP1.AvgPay[0])
}

func (w *whiteTigerChef) rtpHit(math *struct{ HitBsWhiteTigerChefMath }) (iconPay, avgPay int) {
	iconPay = 0
	avgPay = int(math.AvgPay[0])

	if MONSTER.isHit(math.HitWeight) {
		iconPay = w.pick(math)
	}

	return iconPay, avgPay
}
