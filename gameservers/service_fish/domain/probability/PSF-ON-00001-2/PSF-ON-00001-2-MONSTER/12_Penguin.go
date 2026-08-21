package PSF_ON_00001_2_MONSTER

var Penguin = &penguin{}

type penguin struct {
}

func (s *penguin) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
