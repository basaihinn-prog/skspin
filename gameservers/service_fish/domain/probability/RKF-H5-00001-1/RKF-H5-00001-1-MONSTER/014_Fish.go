package RKF_H5_00001_1_MONSTER

var Fish14 = &fish14{}

type fish14 struct{}

func (f *fish14) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish14) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
