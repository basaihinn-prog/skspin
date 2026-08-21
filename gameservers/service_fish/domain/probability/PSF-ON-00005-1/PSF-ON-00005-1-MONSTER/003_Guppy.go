package PSF_ON_00005_1_MONSTER

var Guppy = &guppy{}

type guppy struct{}

func (r *guppy) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
