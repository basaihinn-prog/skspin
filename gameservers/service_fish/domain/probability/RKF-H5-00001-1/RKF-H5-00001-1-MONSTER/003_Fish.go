package RKF_H5_00001_1_MONSTER

var Fish3 = &fish3{}

type fish3 struct{}

func (f *fish3) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
