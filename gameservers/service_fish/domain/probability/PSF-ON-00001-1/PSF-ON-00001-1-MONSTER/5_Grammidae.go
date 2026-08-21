package PSF_ON_00001_1_MONSTER

var Grammidae = &grammidae{}

type grammidae struct {
}

func (s *grammidae) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
