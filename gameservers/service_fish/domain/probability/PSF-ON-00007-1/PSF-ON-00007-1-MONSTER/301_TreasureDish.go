package PSF_ON_00007_1_MONSTER

var TreasureDish = &treasureDish{}

type treasureDish struct{}

func (t *treasureDish) Hit(rtpId string, math *BsPirateShip) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	return PirateShip.Hit(rtpId, math)
}

func (t *treasureDish) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}
