package PSF_ON_00004_1_MONSTER

var Shrimp = &shrimp{}

type shrimp struct {
}

func (s *shrimp) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *shrimp) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
