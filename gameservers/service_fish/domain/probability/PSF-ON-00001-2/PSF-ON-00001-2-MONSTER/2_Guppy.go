package PSF_ON_00001_2_MONSTER

var Guppy = &guppy{}

type guppy struct {
}

func (s *guppy) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
