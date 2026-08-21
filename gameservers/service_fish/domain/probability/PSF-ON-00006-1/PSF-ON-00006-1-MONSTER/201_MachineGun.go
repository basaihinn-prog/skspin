package PSF_ON_00006_1_MONSTER

import (
	"os"
	"strings"

	"serve/fish_comm/rng"
)

var MachineGun = &machineGun{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "201"),
}

type machineGun struct {
	bonusMode bool
}

type BsMachineGunMath struct {
	IconID           int                           `json:"IconID"`
	BonusID          string                        `json:"BonusID"`
	UseRTP           string                        `json:"UseRTP"`
	BonusTimes       []int                         `json:"BonusTimes"`
	BonusTimesWeight []int                         `json:"BonusTimesWeight"`
	BonusIconID      []int                         `json:"BonusIconID"`
	IconPays         []int                         `json:"IconPays"`
	RTP1             struct{ HitBsMachineGunMath } `json:"RTP1"`
	RTP2             struct{ HitBsMachineGunMath } `json:"RTP2"`
	RTP3             struct{ HitBsMachineGunMath } `json:"RTP3"`
	RTP4             struct{ HitBsMachineGunMath } `json:"RTP4"`
	RTP5             struct{ HitBsMachineGunMath } `json:"RTP5"`
	RTP6             struct{ HitBsMachineGunMath } `json:"RTP6"`
	RTP7             struct{ HitBsMachineGunMath } `json:"RTP7"`
	RTP8             struct{ HitBsMachineGunMath } `json:"RTP8"`
	RTP9             struct{ HitBsMachineGunMath } `json:"RTP9"`
}

type HitBsMachineGunMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (m *machineGun) Hit(rtpId string, math *BsMachineGunMath) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	bonusTimes = 0

	switch rtpId {
	case math.RTP1.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP1)

	case math.RTP2.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP2)

	case math.RTP3.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP3)

	case math.RTP4.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP4)

	case math.RTP5.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP5)

	case math.RTP6.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP6)

	case math.RTP7.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP7)

	case math.RTP8.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP8)

	case math.RTP9.ID:
		return m.rtpHit(math.IconPays, math, &math.RTP9)

	default:
		return 0, -1, -1, 0
	}

	return iconPay, triggerIconId, bonusTypeId, bonusTimes
}

func (m *machineGun) UseRtp(math *BsMachineGunMath) (rtpId string) {
	return math.UseRTP
}

func (m *machineGun) pick(math *BsMachineGunMath) int {
	bonusTimes := make([]rng.Option, 0, len(math.BonusTimes))

	for i := 0; i < len(math.BonusTimesWeight); i++ {
		bonusTimes = append(bonusTimes, rng.Option{math.BonusTimesWeight[i], math.BonusTimes[i]})
	}

	return MONSTER.rng(bonusTimes).(int)
}

func (m *machineGun) rtpHit(iconPays []int,
	math *BsMachineGunMath, rtp *struct{ HitBsMachineGunMath }) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	bonusTimes = 0

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = iconPays[0]
		bonusTimes = m.pick(math)
	}

	return iconPay, triggerIconId, bonusTypeId, bonusTimes
}
