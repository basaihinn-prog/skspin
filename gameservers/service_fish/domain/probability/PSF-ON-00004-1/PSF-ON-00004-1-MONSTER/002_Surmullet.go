package PSF_ON_00004_1_MONSTER

var Surmullet = &surmullet{}

type surmullet struct {
}

func (s *surmullet) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *surmullet) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
