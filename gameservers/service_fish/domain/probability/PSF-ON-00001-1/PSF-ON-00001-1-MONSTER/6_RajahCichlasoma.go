package PSF_ON_00001_1_MONSTER

var RajahCichlasoma = &rajahCichlasoma{}

type rajahCichlasoma struct {
}

func (s *rajahCichlasoma) Hit(math *BsFishMath) (iconPay, triggerIconId int) {
	return StarFish.Hit(math)
}
