package PSF_ON_00001_1_MONSTER

var Manatee = &manatee{}

type manatee struct {
}

func (s *manatee) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *manatee) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
