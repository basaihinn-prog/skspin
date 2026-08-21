package PSF_ON_00003_1_MONSTER

var Demon2 = &demon2{}

type demon2 struct {
}

func (d *demon2) Hit(rtpId string, math *BsDemonMath) (iconPay int) {
	return Demon1.Hit(rtpId, math)
}

func (d *demon2) Rtp(order int, math *BsDemonMath) (rtpId string) {
	return Demon1.Rtp(order, math)
}
