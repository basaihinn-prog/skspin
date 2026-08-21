package PSF_ON_00004_1_MONSTER

var MahiMahi = &mahimahi{}

type mahimahi struct {
}

func (m *mahimahi) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (m *mahimahi) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (m *mahimahi) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
