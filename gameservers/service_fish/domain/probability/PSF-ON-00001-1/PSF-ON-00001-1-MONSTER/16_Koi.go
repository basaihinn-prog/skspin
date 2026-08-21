package PSF_ON_00001_1_MONSTER

var Koi = &koi{}

type koi struct {
}

func (s *koi) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *koi) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
