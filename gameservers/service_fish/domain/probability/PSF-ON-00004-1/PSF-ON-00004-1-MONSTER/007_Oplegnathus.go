package PSF_ON_00004_1_MONSTER

var Oplegnathus = &oplegnathus{}

type oplegnathus struct {
}

func (s *oplegnathus) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *oplegnathus) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
