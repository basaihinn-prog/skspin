package PSF_ON_00003_1_MONSTER

var Demon3 = &demon3{}

type demon3 struct {
}

func (d *demon3) Hit(rtpId string, math *BsDemonMath) (iconPay int) {
	return Demon1.Hit(rtpId, math)
}

func (d *demon3) HitFs(math *FsZombieMath) (iconPays int) {
	return math.IconPays[0]
}

func (d *demon3) Rtp(order int, math *BsDemonMath) (rtpId string) {
	return Demon1.Rtp(order, math)
}
