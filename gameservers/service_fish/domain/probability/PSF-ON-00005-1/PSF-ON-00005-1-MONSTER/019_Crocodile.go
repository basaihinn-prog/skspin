package PSF_ON_00005_1_MONSTER

var Crocodile = &crocodile{}

type crocodile struct{}

func (c *crocodile) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
