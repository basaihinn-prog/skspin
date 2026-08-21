package PSF_ON_00004_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var LobsterDash = &lobsterDash{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "300"),
}

type lobsterDash struct {
	bonusMode bool
}

type BsDashMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	RTP1       struct{ HitBsDashMath }
	RTP2       struct{ HitBsDashMath }
	RTP3       struct{ HitBsDashMath }
	RTP4       struct{ HitBsDashMath }
	RTP5       struct{ HitBsDashMath }
}

type HitBsDashMath struct {
	ID             string
	AvgPay         []float32
	HitWeight      []int
	IconPays       []int
	IconPaysWeight []int
}

func (l *lobsterDash) Hit(rtpId string, math *BsDashMath) (avgPay, multiplier int) {
	avgPay = 0
	multiplier = 1

	switch rtpId {
	case math.RTP1.ID:
		avgPay, multiplier = l.rtpHit(&math.RTP1)

	case math.RTP2.ID:
		avgPay, multiplier = l.rtpHit(&math.RTP2)

	case math.RTP3.ID:
		avgPay, multiplier = l.rtpHit(&math.RTP3)

	case math.RTP4.ID:
		avgPay, multiplier = l.rtpHit(&math.RTP4)

	case math.RTP5.ID:
		avgPay, multiplier = l.rtpHit(&math.RTP5)

	default:
		return 0, 1
	}

	return avgPay, multiplier
}

func (l *lobsterDash) pick(math *struct{ HitBsDashMath }) int {
	iconPays := make([]rng.Option, 0, 11)

	for i := 0; i < len(math.IconPaysWeight); i++ {
		iconPays = append(iconPays, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	return MONSTER.rng(iconPays).(int)
}

func (l *lobsterDash) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (l *lobsterDash) Rtp(order int, math *BsDashMath) (rtpId string) {
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

func (l *lobsterDash) AvgPay(math *BsDashMath) int {
	return int(math.RTP1.AvgPay[0])
}
func (l *lobsterDash) rtpHit(math *struct{ HitBsDashMath }) (avgPay, multiplier int) {
	avgPay = int(math.AvgPay[0])
	multiplier = 0

	if MONSTER.isHit(math.HitWeight) {
		multiplier = l.pick(math)
	}

	return avgPay, multiplier
}
