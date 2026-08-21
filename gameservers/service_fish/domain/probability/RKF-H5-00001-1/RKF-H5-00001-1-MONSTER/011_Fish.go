package RKF_H5_00001_1_MONSTER

var Fish11 = &fish11{}

type fish11 struct{}

func (f *fish11) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish11) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
