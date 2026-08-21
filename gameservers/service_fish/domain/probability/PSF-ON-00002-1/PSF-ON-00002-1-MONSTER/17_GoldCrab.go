package PSF_ON_00002_1_MONSTER

var GoldCrab = &goldCrab{}

type goldCrab struct {
}

func (s *goldCrab) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *goldCrab) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
