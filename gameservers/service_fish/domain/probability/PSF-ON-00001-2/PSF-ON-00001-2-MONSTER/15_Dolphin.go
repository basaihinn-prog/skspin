package PSF_ON_00001_2_MONSTER

var Dolphin = &dolphin{}

type dolphin struct {
}

func (s *dolphin) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *dolphin) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
