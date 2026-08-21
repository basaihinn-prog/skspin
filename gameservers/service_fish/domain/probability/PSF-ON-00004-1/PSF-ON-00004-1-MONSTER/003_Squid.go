package PSF_ON_00004_1_MONSTER

var Squid = &squid{}

type squid struct {
}

func (s *squid) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *squid) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
