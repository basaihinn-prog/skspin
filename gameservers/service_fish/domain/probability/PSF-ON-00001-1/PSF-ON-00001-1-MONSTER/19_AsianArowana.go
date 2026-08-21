package PSF_ON_00001_1_MONSTER

var AsianArowana = &asianArowana{}

type asianArowana struct {
}

func (s *asianArowana) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
