package PSF_ON_00001_2_MONSTER

var LanternFish = &lanternFish{}

type lanternFish struct {
}

func (s *lanternFish) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
