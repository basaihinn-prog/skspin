package RKF_H5_00001_1_MONSTER

var Fish4 = &fish4{}

type fish4 struct{}

func (f *fish4) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish4) HitFs(math *FsFishMath) (iconPays int) {
	return Fish0.HitFs(math)
}
