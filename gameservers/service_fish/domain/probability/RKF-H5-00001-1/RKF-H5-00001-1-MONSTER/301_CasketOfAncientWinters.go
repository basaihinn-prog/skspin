package RKF_H5_00001_1_MONSTER

var CasketOfAncientWinters = &casketOfAncientWinters{}

type casketOfAncientWinters struct{}

func (c *casketOfAncientWinters) Hit(rtpId string, math *BsDeepSeaTreasureMath) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	return DeepSeaTreasure.Hit(rtpId, math)
}

func (c *casketOfAncientWinters) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
