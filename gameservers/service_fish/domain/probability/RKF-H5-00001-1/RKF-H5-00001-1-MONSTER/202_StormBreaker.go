package RKF_H5_00001_1_MONSTER

var StormBreaker = &stormBreaker{}

type stormBreaker struct {
}

func (s *stormBreaker) Hit(rtpId string, math *BsThunderHammerMath) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	return ThunderHammer.Hit(rtpId, math)
}

func (s *stormBreaker) UseRtp(math *BsThunderHammerMath) (rtpId string) {
	return ThunderHammer.UseRtp(math)
}
