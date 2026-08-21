package PSF_ON_00002_1_MONSTER

var FaceOff2 = &faceOff2{}

type faceOff2 struct{}

func (f *faceOff2) Hit(rtpId string, math *BsFaceOffMath) (iconPay int, triggerIconId int, bonusTypeId int) {
	return FaceOff1.Hit(rtpId, math)
}
