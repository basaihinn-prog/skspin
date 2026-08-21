package PSF_ON_00002_1_MONSTER

var LionFish = &lionFish{}

type lionFish struct {
}

func (s *lionFish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *lionFish) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
