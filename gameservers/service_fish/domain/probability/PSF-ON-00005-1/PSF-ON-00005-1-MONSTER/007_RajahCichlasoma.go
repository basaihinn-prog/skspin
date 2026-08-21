package PSF_ON_00005_1_MONSTER

var RajahCichlasoma = &rajahcichlasoma{}

type rajahcichlasoma struct{}

func (r *rajahcichlasoma) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
