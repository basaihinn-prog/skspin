package PSF_ON_00001_1_MONSTER

var LanternFish = &lanternFish{}

type lanternFish struct {
}

func (s *lanternFish) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
