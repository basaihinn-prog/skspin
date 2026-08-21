package PSF_ON_00002_1_MONSTER

var RajahCichlasoma = &rajahCichlasoma{}

type rajahCichlasoma struct {
}

func (s *rajahCichlasoma) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
