package PSF_ON_00004_1_MONSTER

var ButterflyFish = &butterflyfish{}

type butterflyfish struct {
}

func (b *butterflyfish) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (b *butterflyfish) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
