package PSF_ON_00004_1_MONSTER

var GiantOarfish = &giantOarfish{}

type giantOarfish struct {
}

func (g *giantOarfish) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (g *giantOarfish) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
