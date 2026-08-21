package PSF_ON_00001_1_MONSTER

var Platypus = &platypus{}

type platypus struct {
}

func (s *platypus) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *platypus) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
