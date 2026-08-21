package PSF_ON_00004_1_MONSTER

var PlatypusSeniorChef = &platypusSeniorChef{}

type platypusSeniorChef struct {
}

func (p *platypusSeniorChef) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (p *platypusSeniorChef) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (p *platypusSeniorChef) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
