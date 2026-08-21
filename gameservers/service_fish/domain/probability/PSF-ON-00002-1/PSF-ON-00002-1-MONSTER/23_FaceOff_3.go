package PSF_ON_00002_1_MONSTER

var FaceOff3 = &faceOff3{}

type faceOff3 struct{}

func (f *faceOff3) Hit(rtpId string, math *BsFaceOffMath) (iconPay int, triggerIconId int, bonusTypeId int) {
	return FaceOff1.Hit(rtpId, math)
}

func (f *faceOff3) HitFs(math *FsFaceOffMath) (iconPay int) {
	return math.IconPays[0]
}
