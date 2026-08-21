package PSF_ON_00002_1_MONSTER

var Lanternfish = &lanternfish{}

type lanternfish struct {
}

func (s *lanternfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
