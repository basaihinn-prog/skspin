package PSF_ON_00001_1_MONSTER

var SeaHorse = &seaHorse{}

type seaHorse struct {
}

func (s *seaHorse) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
