package PSF_ON_00005_1_MONSTER

var MangoFish = &mangofish{}

type mangofish struct{}

func (m *mangofish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
