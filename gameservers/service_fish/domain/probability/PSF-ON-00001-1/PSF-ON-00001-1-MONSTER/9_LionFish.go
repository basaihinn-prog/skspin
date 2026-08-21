package PSF_ON_00001_1_MONSTER

var LionFish = &lionFish{}

type lionFish struct {
}

func (s *lionFish) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *lionFish) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
