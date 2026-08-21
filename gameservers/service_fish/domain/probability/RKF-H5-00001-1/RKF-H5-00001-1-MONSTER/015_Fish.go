package RKF_H5_00001_1_MONSTER

var Fish15 = &fish15{}

type fish15 struct{}

func (f *fish15) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish15) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
