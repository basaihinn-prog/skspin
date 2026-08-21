package PSF_ON_00001_2_MONSTER

var FaceOff4 = &faceOff4{}

type faceOff4 struct{}

func (f *faceOff4) Hit(rtpId string, math *BsFaceOffMath) (iconPay int, triggerIconId int, bonusTypeId int) {
	return FaceOff1.Hit(rtpId, math)
}
