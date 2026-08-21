package PSF_ON_00004_1_MONSTER

var LobsterWaiter = &lobsterWaiter{}

type lobsterWaiter struct {
}

func (l *lobsterWaiter) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	return Clam.Hit(rtpId, math)
}

func (l *lobsterWaiter) rtpHit(iconPays int, bsFish *struct{ HitBsFishMath }) (iconPay int) {
	if MONSTER.isHit(bsFish.HitWeight) {
		iconPay = iconPays
	}

	return iconPay
}

func (l *lobsterWaiter) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (l *lobsterWaiter) Rtp(order int, math *BsFishMath) (rtpId string) {
	return Clam.Rtp(order, math)
}
