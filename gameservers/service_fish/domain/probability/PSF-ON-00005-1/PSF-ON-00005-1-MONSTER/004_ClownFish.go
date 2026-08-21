package PSF_ON_00005_1_MONSTER

var ClownFish = &clownfish{}

type clownfish struct{}

func (c *clownfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (c *clownfish) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
