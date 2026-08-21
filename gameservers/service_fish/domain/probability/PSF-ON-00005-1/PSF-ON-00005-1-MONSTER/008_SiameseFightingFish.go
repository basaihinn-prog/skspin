package PSF_ON_00005_1_MONSTER

var SiameseFightingFish = &siamesefightingfish{}

type siamesefightingfish struct{}

func (s *siamesefightingfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
