package PSF_ON_00001_1_MONSTER

var FaceOff4 = &faceOff4{}

type faceOff4 struct{}

func (f *faceOff4) Hit(math *BsFaceOffMath) (iconPay int) {
	return FaceOff1.Hit(math)
}
