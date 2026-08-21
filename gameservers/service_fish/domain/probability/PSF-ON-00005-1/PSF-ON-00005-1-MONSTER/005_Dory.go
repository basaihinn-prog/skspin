package PSF_ON_00005_1_MONSTER

var Dory = &dory{}

type dory struct{}

func (d *dory) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
