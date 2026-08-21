package PSF_ON_00004_1_MONSTER

var FlyingFish = &flyingfish{}

type flyingfish struct {
}

func (f *flyingfish) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (f *flyingfish) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (f *flyingfish) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
