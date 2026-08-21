package PSF_ON_00005_1_MONSTER

var RiverShrimo = &rivershrimo{}

type rivershrimo struct{}

func (r *rivershrimo) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}
