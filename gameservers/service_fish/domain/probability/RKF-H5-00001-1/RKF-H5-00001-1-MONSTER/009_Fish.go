package RKF_H5_00001_1_MONSTER

var Fish9 = &fish9{}

type fish9 struct{}

func (f *fish9) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish9) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
