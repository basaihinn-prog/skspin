package PSF_ON_00001_1_MONSTER

var Guppy = &guppy{}

type guppy struct {
}

func (s *guppy) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
