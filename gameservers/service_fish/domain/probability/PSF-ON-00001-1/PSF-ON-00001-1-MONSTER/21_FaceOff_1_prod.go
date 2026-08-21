//go:build prod
// +build prod

package PSF_ON_00001_1_MONSTER

func (f *faceOff1) Hit(math *BsFaceOffMath) (iconPay int) {
	if MONSTER.isHit(math.HitWeight) {
		return math.IconPays[0]
	}
	return 0
}

func (f *faceOff1) HitFs(math *FsFaceOffMath) (iconPay int) {
	return math.IconPays[0]
}
