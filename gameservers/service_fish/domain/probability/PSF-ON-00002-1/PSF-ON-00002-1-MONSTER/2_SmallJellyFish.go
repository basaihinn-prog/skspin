package PSF_ON_00002_1_MONSTER

var SmallJellyfish = &smallJellyfish{}

type smallJellyfish struct {
}

func (s *smallJellyfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
