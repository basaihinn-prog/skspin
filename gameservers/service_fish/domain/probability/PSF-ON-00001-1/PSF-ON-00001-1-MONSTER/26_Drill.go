package PSF_ON_00001_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var Drill = &drill{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "26"),
}

type drill struct {
	bonusMode bool
}

type BsDrillMath struct {
	IconID           int
	IsWildSubstitute bool
	BonusID          string
	BonusTimes       []int
	BonusTimesWeight []int
	BonusIconID      []int
	HitWeight        []int
	IconPays         []int
}

func (d *drill) pick(math *BsDrillMath) int {
	bonusTimes := make([]rng.Option, 0, 6)

	for i := 0; i < len(math.BonusTimesWeight); i++ {
		bonusTimes = append(bonusTimes, rng.Option{math.BonusTimesWeight[i], math.BonusTimes[i]})
	}
	return MONSTER.rng(bonusTimes).(int)
}
func (d *drill) Hit(math *BsDrillMath) (iconPay, bonusTimes int) {
	if MONSTER.isHit(math.HitWeight) {
		return math.IconPays[0], d.pick(math)
	}
	return 0, 0
}
