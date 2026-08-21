package PSF_ON_00002_1_MONSTER

var AsianArowana = &asianArowana{}

type asianArowana struct {
}

func (s *asianArowana) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
