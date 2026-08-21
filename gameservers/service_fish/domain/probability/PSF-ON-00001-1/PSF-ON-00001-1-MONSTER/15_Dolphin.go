package PSF_ON_00001_1_MONSTER

var Dolphin = &dolphin{}

type dolphin struct {
}

func (s *dolphin) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *dolphin) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
