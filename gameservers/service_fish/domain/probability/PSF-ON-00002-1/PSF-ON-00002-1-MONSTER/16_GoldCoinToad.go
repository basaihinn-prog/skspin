package PSF_ON_00002_1_MONSTER

var GoldCoinToad = &goldCoinToad{}

type goldCoinToad struct {
}

func (s *goldCoinToad) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (s *goldCoinToad) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
