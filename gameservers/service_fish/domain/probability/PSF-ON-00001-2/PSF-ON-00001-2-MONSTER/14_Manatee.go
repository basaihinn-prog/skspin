package PSF_ON_00001_2_MONSTER

var Manatee = &manatee{}

type manatee struct {
}

func (s *manatee) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *manatee) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
