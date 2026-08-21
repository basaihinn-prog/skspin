package PSF_ON_00005_1_MONSTER

var HermitCrab = &hermitcrab{}

type hermitcrab struct{}

func (h *hermitcrab) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (h *hermitcrab) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
