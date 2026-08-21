package PSF_ON_00004_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var RedEnvelope = &redEnvelope{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "100"),
}

type redEnvelope struct {
	bonusMode bool
}

type BsRedEnvelopeMath struct {
	IconID    int
	BonusID   string
	BonusTime []int
	RTP1      struct{ HitBsRedEnvelopeMath }
	RTP2      struct{ HitBsRedEnvelopeMath }
	RTP3      struct{ HitBsRedEnvelopeMath }
	RTP4      struct{ HitBsRedEnvelopeMath }
	RTP5      struct{ HitBsRedEnvelopeMath }
}

type HitBsRedEnvelopeMath struct {
	ID                 string
	AvgPay             []float64
	HitWeight          []int
	LowIconPays        []int
	LowIconPaysWeight  []int
	HighIconPays       []int
	HighIconPaysWeight []int
}

func (r *redEnvelope) pick(math *struct{ HitBsRedEnvelopeMath }) []int {
	lowPayOptions := r.getOptions(math.LowIconPaysWeight, math.LowIconPays)
	highPayOptions := r.getOptions(math.HighIconPaysWeight, math.HighIconPays)

	iconPays := []int{
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(highPayOptions).(int),
	}

	return iconPays
}

func (r *redEnvelope) getOptions(weights, pays []int) []rng.Option {
	payOptions := make([]rng.Option, 0, len(weights))

	for i := 0; i < len(weights); i++ {
		payOptions = append(payOptions, rng.Option{weights[i], pays[i]})
	}

	return payOptions
}

func (r *redEnvelope) Hit(rtpId string, math *BsRedEnvelopeMath) (iconPays []int, avgPay int) {
	switch rtpId {
	case math.RTP1.ID:
		iconPays, avgPay = r.rtpHit(math.RTP1.HitWeight, &math.RTP1)

	case math.RTP2.ID:
		iconPays, avgPay = r.rtpHit(math.RTP2.HitWeight, &math.RTP2)

	case math.RTP3.ID:
		iconPays, avgPay = r.rtpHit(math.RTP3.HitWeight, &math.RTP3)

	case math.RTP4.ID:
		iconPays, avgPay = r.rtpHit(math.RTP4.HitWeight, &math.RTP4)

	case math.RTP5.ID:
		iconPays, avgPay = r.rtpHit(math.RTP5.HitWeight, &math.RTP5)

	default:
		return nil, 0
	}

	return iconPays, avgPay
}

func (r *redEnvelope) Rtp(order int, math *BsRedEnvelopeMath) (rtpId string) {
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

func (r *redEnvelope) AvgPay(math *BsRedEnvelopeMath) int {
	return int(math.RTP1.AvgPay[0])
}

func (r *redEnvelope) rtpHit(hitWeight []int, rtp *struct{ HitBsRedEnvelopeMath }) (iconPays []int, avgPay int) {
	avgPay = int(rtp.AvgPay[0])

	if MONSTER.isHit(hitWeight) {
		return r.pick(rtp), avgPay
	}

	return nil, avgPay
}
