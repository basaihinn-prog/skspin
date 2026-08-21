package PSF_ON_00001_1_MONSTER

var FaceOff3 = &faceOff3{}

type faceOff3 struct{}

func (f *faceOff3) Hit(math *BsFaceOffMath) (iconPay int) {
	return FaceOff1.Hit(math)
}

func (f *faceOff3) HitFs(math *FsFaceOffMath) (iconPay int) {
	return math.IconPays[0]
}
