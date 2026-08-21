package PSF_ON_00001_2_MONSTER

var Shark = &shark{}

type shark struct {
}

func (s *shark) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
