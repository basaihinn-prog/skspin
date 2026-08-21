package PSF_ON_00001_2_MONSTER

var Turtle = &turtle{}

type turtle struct {
}

func (s *turtle) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
