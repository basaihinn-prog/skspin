package PSF_ON_00004_1_MONSTER

var SnaggletoothShark = &snaggletoothShark{}

type snaggletoothShark struct {
}

func (s *snaggletoothShark) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *snaggletoothShark) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (s *snaggletoothShark) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
