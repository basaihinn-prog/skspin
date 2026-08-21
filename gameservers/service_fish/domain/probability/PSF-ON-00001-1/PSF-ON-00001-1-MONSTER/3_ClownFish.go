package PSF_ON_00001_1_MONSTER

var ClownFish = &clownFish{}

type clownFish struct {
}

func (s *clownFish) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
