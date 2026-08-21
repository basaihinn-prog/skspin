package PSF_ON_00005_1_MONSTER

var GoldGemTurtle = &goldgemturtle{}

type goldgemturtle struct{}

func (g *goldgemturtle) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (g *goldgemturtle) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
