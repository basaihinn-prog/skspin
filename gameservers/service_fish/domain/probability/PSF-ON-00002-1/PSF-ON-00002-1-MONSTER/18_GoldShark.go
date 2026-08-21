package PSF_ON_00002_1_MONSTER

var GoldShark = &goldShark{}

type goldShark struct {
}

func (s *goldShark) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
