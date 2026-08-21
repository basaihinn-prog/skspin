package PSF_ON_00001_1_MONSTER

var Shark = &shark{}

type shark struct {
}

func (s *shark) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
