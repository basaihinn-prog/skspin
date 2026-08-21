package PSF_ON_00002_1_MONSTER

var HermitCrab = &hermitCrab{}

type hermitCrab struct {
}

func (s *hermitCrab) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *hermitCrab) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
