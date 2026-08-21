package PSF_ON_00001_2_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var MachineGun = &machineGun{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "27"),
}

type machineGun struct {
	bonusMode bool
}

type BsMachineGunMath struct {
	IconID           int
	BonusID          string
	BonusTimes       []int
	BonusTimesWeight []int
	BonusIconID      []int
	IconPays         []int
	RTP1             struct{ HitBsMachineGunMath } `json:"RTP1"`
	RTP2             struct{ HitBsMachineGunMath } `json:"RTP2"`
	RTP3             struct{ HitBsMachineGunMath } `json:"RTP3"`
	RTP4             struct{ HitBsMachineGunMath } `json:"RTP4"`
	RTP5             struct{ HitBsMachineGunMath } `json:"RTP5"`
	RTP6             struct{ HitBsMachineGunMath } `json:"RTP6"`
	RTP7             struct{ HitBsMachineGunMath } `json:"RTP7"`
	RTP8             struct{ HitBsMachineGunMath } `json:"RTP8"`
	RTP9             struct{ HitBsMachineGunMath } `json:"RTP9"`
	RTP10            struct{ HitBsMachineGunMath } `json:"RTP10"`
	RTP11            struct{ HitBsMachineGunMath } `json:"RTP11"`
}

type HitBsMachineGunMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
	Type          int    `json:"Type"`
}

func (m *machineGun) pick(math *BsMachineGunMath) int {
	bonusTimes := make([]rng.Option, 0, 6)

	for i := 0; i < len(math.BonusTimesWeight); i++ {
		bonusTimes = append(bonusTimes, rng.Option{math.BonusTimesWeight[i], math.BonusTimes[i]})
	}

	return MONSTER.rng(bonusTimes).(int)
}

func (m *machineGun) Hit(rtpId string, math *BsMachineGunMath) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	bonusTimes = 0

	switch rtpId {
	case "0":
		if MONSTER.isHit(math.RTP1.TriggerWeight) {
			triggerIconId = math.RTP1.TriggerIconID
			bonusTypeId = math.RTP1.Type
		}

		if MONSTER.isHit(math.RTP1.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "20":
		if MONSTER.isHit(math.RTP2.TriggerWeight) {
			triggerIconId = math.RTP2.TriggerIconID
			bonusTypeId = math.RTP2.Type
		}

		if MONSTER.isHit(math.RTP2.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "40":
		if MONSTER.isHit(math.RTP3.TriggerWeight) {
			triggerIconId = math.RTP3.TriggerIconID
			bonusTypeId = math.RTP3.Type
		}

		if MONSTER.isHit(math.RTP3.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "60":
		if MONSTER.isHit(math.RTP4.TriggerWeight) {
			triggerIconId = math.RTP4.TriggerIconID
			bonusTypeId = math.RTP4.Type
		}

		if MONSTER.isHit(math.RTP4.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "80":
		if MONSTER.isHit(math.RTP5.TriggerWeight) {
			triggerIconId = math.RTP5.TriggerIconID
			bonusTypeId = math.RTP5.Type
		}

		if MONSTER.isHit(math.RTP5.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "901":
		if MONSTER.isHit(math.RTP6.TriggerWeight) {
			triggerIconId = math.RTP6.TriggerIconID
			bonusTypeId = math.RTP6.Type
		}

		if MONSTER.isHit(math.RTP6.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "902":
		if MONSTER.isHit(math.RTP7.TriggerWeight) {
			triggerIconId = math.RTP7.TriggerIconID
			bonusTypeId = math.RTP7.Type
		}

		if MONSTER.isHit(math.RTP7.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "100":
		if MONSTER.isHit(math.RTP8.TriggerWeight) {
			triggerIconId = math.RTP8.TriggerIconID
			bonusTypeId = math.RTP8.Type
		}

		if MONSTER.isHit(math.RTP8.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "150":
		if MONSTER.isHit(math.RTP9.TriggerWeight) {
			triggerIconId = math.RTP9.TriggerIconID
			bonusTypeId = math.RTP9.Type
		}

		if MONSTER.isHit(math.RTP9.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "200":
		if MONSTER.isHit(math.RTP10.TriggerWeight) {
			triggerIconId = math.RTP10.TriggerIconID
			bonusTypeId = math.RTP10.Type
		}

		if MONSTER.isHit(math.RTP10.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	case "300":
		if MONSTER.isHit(math.RTP11.TriggerWeight) {
			triggerIconId = math.RTP11.TriggerIconID
			bonusTypeId = math.RTP11.Type
		}

		if MONSTER.isHit(math.RTP11.HitWeight) {
			iconPay = math.IconPays[0]
			bonusTimes = m.pick(math)
		}

	default:
		return 0, -1, -1, 0
	}

	return iconPay, triggerIconId, bonusTypeId, bonusTimes
}
