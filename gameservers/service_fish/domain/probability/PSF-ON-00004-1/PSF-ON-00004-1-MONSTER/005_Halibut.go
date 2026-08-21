package PSF_ON_00004_1_MONSTER

var Halibut = &halibut{}

type halibut struct {
}

func (h *halibut) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (h *halibut) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
