package RKF_H5_00001_1_MONSTER

var Fish19 = &fish19{}

type fish19 struct{}

func (f *fish19) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
