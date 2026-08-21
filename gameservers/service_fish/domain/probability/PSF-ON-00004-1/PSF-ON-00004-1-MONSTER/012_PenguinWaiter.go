package PSF_ON_00004_1_MONSTER

var PenguinWaiter = &penguinWaiter{}

type penguinWaiter struct {
}

func (p *penguinWaiter) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (p *penguinWaiter) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
