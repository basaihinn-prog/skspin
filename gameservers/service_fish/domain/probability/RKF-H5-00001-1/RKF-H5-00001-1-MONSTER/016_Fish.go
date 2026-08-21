package RKF_H5_00001_1_MONSTER

var Fish16 = &fish16{}

type fish16 struct{}

func (f *fish16) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish16) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
