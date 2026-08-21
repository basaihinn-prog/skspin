package PSF_ON_00001_1_MONSTER

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
	IsWildSubstitute bool
	BonusID          string
	BonusTimes       []int
	BonusTimesWeight []int
	BonusIconID      []int
	HitWeight        []int
	IconPays         []int
}

func (m *machineGun) pick(math *BsMachineGunMath) int {
	bonusTimes := make([]rng.Option, 0, 6)

	for i := 0; i < len(math.BonusTimesWeight); i++ {
		bonusTimes = append(bonusTimes, rng.Option{math.BonusTimesWeight[i], math.BonusTimes[i]})
	}
	return MONSTER.rng(bonusTimes).(int)
}
func (m *machineGun) Hit(math *BsMachineGunMath) (iconPay, bonusTimes int) {
	if MONSTER.isHit(math.HitWeight) {
		return math.IconPays[0], m.pick(math)
	}
	return 0, 0
}
