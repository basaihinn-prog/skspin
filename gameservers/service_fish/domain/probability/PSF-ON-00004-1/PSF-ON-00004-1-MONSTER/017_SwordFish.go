package PSF_ON_00004_1_MONSTER

var SwordFish = &swordFish{}

type swordFish struct {
}

func (s *swordFish) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (s *swordFish) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (s *swordFish) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
