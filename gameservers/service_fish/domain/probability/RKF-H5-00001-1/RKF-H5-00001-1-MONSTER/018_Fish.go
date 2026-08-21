package RKF_H5_00001_1_MONSTER

var Fish18 = &fish18{}

type fish18 struct{}

func (f *fish18) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
