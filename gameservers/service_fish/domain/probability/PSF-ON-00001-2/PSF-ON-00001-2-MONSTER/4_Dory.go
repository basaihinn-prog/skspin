package PSF_ON_00001_2_MONSTER

var Dory = &dory{}

type dory struct {
}

func (s *dory) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *dory) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
