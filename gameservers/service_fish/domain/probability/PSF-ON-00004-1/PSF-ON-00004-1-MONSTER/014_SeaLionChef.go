package PSF_ON_00004_1_MONSTER

var SeaLionChef = &seaLionChef{}

type seaLionChef struct {
}

func (s *seaLionChef) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *seaLionChef) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (s *seaLionChef) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
