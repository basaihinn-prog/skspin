package PSF_ON_00001_1_MONSTER

var Dory = &dory{}

type dory struct {
}

func (s *dory) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *dory) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
