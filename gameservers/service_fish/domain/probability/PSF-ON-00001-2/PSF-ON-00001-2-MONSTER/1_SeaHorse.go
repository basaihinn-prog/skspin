package PSF_ON_00001_2_MONSTER

var SeaHorse = &seaHorse{}

type seaHorse struct {
}

func (s *seaHorse) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
