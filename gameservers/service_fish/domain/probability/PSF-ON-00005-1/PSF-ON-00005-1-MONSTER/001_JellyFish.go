package PSF_ON_00005_1_MONSTER

var JellyFish = &jellyfish{}

type jellyfish struct{}

func (j *jellyfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
