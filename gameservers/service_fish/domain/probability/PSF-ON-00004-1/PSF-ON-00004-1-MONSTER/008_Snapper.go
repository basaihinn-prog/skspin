package PSF_ON_00004_1_MONSTER

var Snapper = &snapper{}

type snapper struct {
}

func (s *snapper) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *snapper) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
