package PSF_ON_00001_1_MONSTER

var FaceOff5 = &faceOff5{}

type faceOff5 struct{}

func (f *faceOff5) Hit(math *BsFaceOffMath) (iconPay int) {
	return FaceOff1.Hit(math)
}

func (f *faceOff5) HitFs(math *FsFaceOffMath) (iconPay int) {
	return math.IconPays[0]
}
