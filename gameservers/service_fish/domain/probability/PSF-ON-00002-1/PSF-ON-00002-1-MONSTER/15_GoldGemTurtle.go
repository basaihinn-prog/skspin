package PSF_ON_00002_1_MONSTER

var GoldGemTurtle = &goldGemTurtle{}

type goldGemTurtle struct {
}

func (s *goldGemTurtle) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *goldGemTurtle) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
