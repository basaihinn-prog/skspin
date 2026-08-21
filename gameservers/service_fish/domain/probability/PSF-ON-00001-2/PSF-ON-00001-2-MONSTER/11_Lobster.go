package PSF_ON_00001_2_MONSTER

var Lobster = &lobster{}

type lobster struct {
}

func (s *lobster) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *lobster) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
