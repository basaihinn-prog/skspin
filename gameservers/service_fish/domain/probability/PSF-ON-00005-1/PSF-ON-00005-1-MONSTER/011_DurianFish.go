package PSF_ON_00005_1_MONSTER

var DurianFish = &durianfish{}

type durianfish struct{}

func (d *durianfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (d *durianfish) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
