package PSF_ON_00001_2_MONSTER

var PufferFish = &pufferFish{}

type pufferFish struct {
}

func (s *pufferFish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
