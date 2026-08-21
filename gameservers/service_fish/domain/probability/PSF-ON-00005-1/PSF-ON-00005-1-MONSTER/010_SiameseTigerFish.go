package PSF_ON_00005_1_MONSTER

var SiameseTigerFish = &siamesetigerfish{}

type siamesetigerfish struct{}

func (s *siamesetigerfish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
