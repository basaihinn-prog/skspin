package PSF_ON_00002_1_MONSTER

var GiantJellyfish = &giantJellyfish{}

type giantJellyfish struct {
}

func (s *giantJellyfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *giantJellyfish) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
