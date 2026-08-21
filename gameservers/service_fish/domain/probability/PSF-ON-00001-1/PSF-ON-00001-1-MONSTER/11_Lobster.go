package PSF_ON_00001_1_MONSTER

var Lobster = &lobster{}

type lobster struct {
}

func (s *lobster) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *lobster) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
