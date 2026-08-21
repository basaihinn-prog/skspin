package PSF_ON_00004_1_MONSTER

var HairtCrab = &hairtCrab{}

type hairtCrab struct {
}

func (h *hairtCrab) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (h *hairtCrab) rtpHit(iconPays int, bsFish *struct{ HitBsFishMath }) (iconPay int) {
	if MONSTER.isHit(bsFish.HitWeight) {
		iconPay = iconPays
	}

	return iconPay
}

func (h *hairtCrab) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (h *hairtCrab) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
