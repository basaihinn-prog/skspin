package PSF_ON_00004_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var MachineGun = &machineGun{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "201"),
}

type machineGun struct {
	bonusMode bool
}

type BsMachineGunMath struct {
	IconID     int
	BonusID    string
	BonusTimes []int
	IconPays   []int
	UseRTP     string
	RTP1       struct{ HitBsMachineGunMath }
	RTP2       struct{ HitBsMachineGunMath }
	RTP3       struct{ HitBsMachineGunMath }
	RTP4       struct{ HitBsMachineGunMath }
	RTP5       struct{ HitBsMachineGunMath }
}

type HitBsMachineGunMath struct {
	ID               string
	AvgPay           []float32
	HitWeight        []int
	BonusTimes       []int
	BonusTimesWeight []int
	BonusIconID      []int
}

func (m *machineGun) Hit(rtpId string, math *BsMachineGunMath) (iconPay, bonusTimes, avgPay int) {
	iconPay = 0
	bonusTimes = 0

	switch rtpId {
	case math.RTP1.ID:
		iconPay, bonusTimes, avgPay = m.rtpHit(math.IconPays, &math.RTP1)

	case math.RTP2.ID:
		iconPay, bonusTimes, avgPay = m.rtpHit(math.IconPays, &math.RTP2)

	case math.RTP3.ID:
		iconPay, bonusTimes, avgPay = m.rtpHit(math.IconPays, &math.RTP3)

	case math.RTP4.ID:
		iconPay, bonusTimes, avgPay = m.rtpHit(math.IconPays, &math.RTP4)

	case math.RTP5.ID:
		iconPay, bonusTimes, avgPay = m.rtpHit(math.IconPays, &math.RTP5)

	default:
		return 0, 0, 0
	}

	return iconPay, bonusTimes, avgPay
}

func (m *machineGun) pick(rtp *struct{ HitBsMachineGunMath }) int {
	bonusTimes := make([]rng.Option, 0, len(rtp.BonusTimes))

	for i := 0; i < len(rtp.BonusTimesWeight); i++ {
		bonusTimes = append(bonusTimes, rng.Option{rtp.BonusTimesWeight[i], rtp.BonusTimes[i]})
	}

	return MONSTER.rng(bonusTimes).(int)
}

func (m *machineGun) Rtp(order int, math *BsMachineGunMath) (rtpId string) {
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

func (m *machineGun) UseRtp(math *BsMachineGunMath) (rtpId string) {
	return math.UseRTP
}

func (m *machineGun) AvgPay(math *BsMachineGunMath) int {
	return int(math.RTP1.AvgPay[0])
}

func (m *machineGun) rtpHit(iconPays []int, rtp *struct{ HitBsMachineGunMath }) (iconPay, bonusTimes, avgPay int) {
	avgPay = int(rtp.AvgPay[0])

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = iconPays[0]
		bonusTimes = m.pick(rtp)
	}

	return iconPay, bonusTimes, avgPay
}
