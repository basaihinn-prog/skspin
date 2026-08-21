package PSF_ON_00002_1_MONSTER

var FaceOff5 = &faceOff5{}

type faceOff5 struct{}

func (f *faceOff5) Hit(rtpId string, math *BsFaceOffMath) (iconPay int, triggerIconId int, bonusTypeId int) {
	return FaceOff1.Hit(rtpId, math)
}

func (f *faceOff5) HitFs(math *FsFaceOffMath) (iconPay int) {
	return math.IconPays[0]
}
