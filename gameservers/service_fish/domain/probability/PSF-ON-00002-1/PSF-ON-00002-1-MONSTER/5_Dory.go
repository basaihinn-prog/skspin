package PSF_ON_00002_1_MONSTER

var Dory = &dory{}

type dory struct {
}

func (s *dory) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
