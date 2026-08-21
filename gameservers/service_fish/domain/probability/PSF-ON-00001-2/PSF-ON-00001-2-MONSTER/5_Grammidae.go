package PSF_ON_00001_2_MONSTER

var Grammidae = &grammidae{}

type grammidae struct {
}

func (s *grammidae) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
