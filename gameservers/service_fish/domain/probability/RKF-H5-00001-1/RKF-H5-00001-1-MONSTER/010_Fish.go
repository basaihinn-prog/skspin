package RKF_H5_00001_1_MONSTER

var Fish10 = &fish10{}

type fish10 struct{}

func (f *fish10) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
