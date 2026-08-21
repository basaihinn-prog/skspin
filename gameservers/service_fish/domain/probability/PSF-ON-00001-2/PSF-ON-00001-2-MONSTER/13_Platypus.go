package PSF_ON_00001_2_MONSTER

var Platypus = &platypus{}

type platypus struct {
}

func (s *platypus) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *platypus) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
