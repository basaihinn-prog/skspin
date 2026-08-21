package PSF_ON_00004_1_MONSTER

var Stingray = &stingray{}

type stingray struct {
}

func (s *stingray) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *stingray) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
