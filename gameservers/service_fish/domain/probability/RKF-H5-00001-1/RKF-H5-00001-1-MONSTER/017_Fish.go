package RKF_H5_00001_1_MONSTER

var Fish17 = &fish17{}

type fish17 struct{}

func (f *fish17) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish17) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
