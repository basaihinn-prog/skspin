package PSF_ON_00001_2_MONSTER

var Koi = &koi{}

type koi struct {
}

func (s *koi) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *koi) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
