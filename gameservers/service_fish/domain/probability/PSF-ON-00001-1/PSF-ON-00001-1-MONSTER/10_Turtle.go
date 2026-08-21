package PSF_ON_00001_1_MONSTER

var Turtle = &turtle{}

type turtle struct {
}

func (s *turtle) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
