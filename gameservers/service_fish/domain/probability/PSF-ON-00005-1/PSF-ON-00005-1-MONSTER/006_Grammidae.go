package PSF_ON_00005_1_MONSTER

var Grammidae = &grammidae{}

type grammidae struct{}

func (g *grammidae) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
