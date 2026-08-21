package PSF_ON_00001_1_MONSTER

var Penguin = &penguin{}

type penguin struct {
}

func (s *penguin) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
