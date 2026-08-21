package PSF_ON_00005_1_MONSTER

var GoldShark = &goldshark{}

type goldshark struct{}

func (g *goldshark) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
