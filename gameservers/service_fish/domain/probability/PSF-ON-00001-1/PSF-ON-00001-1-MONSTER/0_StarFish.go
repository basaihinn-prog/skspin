package PSF_ON_00001_1_MONSTER

import (
	"os"
	"strings"
)

var StarFish = &starFish{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "31"),
}

type starFish struct {
	bonusMode bool
}

type BsFishMath struct {
	IconID           int    `json:"IconID"`
	IsWildSubstitute bool   `json:"IsWildSubstitute"`
	BonusID          string `json:"BonusID"`
	BonusTimes       []int  `json:"BonusTimes"`
	HitWeight        []int  `json:"HitWeight"`
	IconPays         []int  `json:"IconPays"`
	TriggerIconID    int    `json:"TriggerIconID"`
	TriggerWeight    []int  `json:"TriggerWeight"`
}

type FsFishMath struct {
	IconID   int   `json:"IconID"`
	IconPays []int `json:"IconPays"`
}

func (s *starFish) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	if MONSTER.isHit(math.HitWeight) {
		if MONSTER.isHit(math.TriggerWeight) {
			return math.IconPays[0], math.TriggerIconID
		}
		return math.IconPays[0], -1
	}
	return 0, -1
}

func (s *starFish) HitFs(math *FsFishMath) (iconPays int) {
	return math.IconPays[0]
}
