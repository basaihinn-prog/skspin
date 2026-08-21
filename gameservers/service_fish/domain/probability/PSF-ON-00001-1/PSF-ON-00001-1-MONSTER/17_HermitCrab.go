package PSF_ON_00001_1_MONSTER

var HermitCrab = &hermitCrab{}

type hermitCrab struct {
}

func (s *hermitCrab) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}

func (s *hermitCrab) HitFs(math *FsFishMath) (iconPay int) {
	return StarFish.HitFs(math)
}
