package PSF_ON_00005_1_MONSTER

var GoldCoinToad = &goldcointoad{}

type goldcointoad struct{}

func (g *goldcointoad) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (g *goldcointoad) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
