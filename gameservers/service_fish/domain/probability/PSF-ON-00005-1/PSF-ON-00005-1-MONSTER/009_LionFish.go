package PSF_ON_00005_1_MONSTER

var LionFish = &lionfish{}

type lionfish struct{}

func (l *lionfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (l *lionfish) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
