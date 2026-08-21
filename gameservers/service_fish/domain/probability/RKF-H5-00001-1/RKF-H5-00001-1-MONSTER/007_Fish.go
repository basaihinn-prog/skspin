package RKF_H5_00001_1_MONSTER

var Fish7 = &fish7{}

type fish7 struct{}

func (f *fish7) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
