package PSF_ON_00001_1_MONSTER

var FaceOff2 = &faceOff2{}

type faceOff2 struct{}

func (f *faceOff2) Hit(math *BsFaceOffMath) (iconPay int) {
	return FaceOff1.Hit(math)
}
