package PSF_ON_00004_1_MONSTER

var WhaleShark = &whaleShark{}

type whaleShark struct {
}

func (w *whaleShark) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (w *whaleShark) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
