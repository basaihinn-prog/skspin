package PSF_ON_00002_1_MONSTER

var Clownfish = &clownFish{}

type clownFish struct {
}

func (s *clownFish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *clownFish) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
