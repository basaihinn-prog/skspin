package PSF_ON_00005_1_MONSTER

var GoldCrab = &goldcrab{}

type goldcrab struct{}

func (g *goldcrab) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (g *goldcrab) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
