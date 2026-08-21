package PSF_ON_00001_1_MONSTER

var PufferFish = &pufferFish{}

type pufferFish struct {
}

func (s *pufferFish) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
