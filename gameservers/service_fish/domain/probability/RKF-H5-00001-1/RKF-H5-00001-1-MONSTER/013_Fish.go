package RKF_H5_00001_1_MONSTER

var Fish13 = &fish13{}

type fish13 struct{}

func (f *fish13) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish13) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
