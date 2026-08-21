package RKF_H5_00001_1_MONSTER

var Fish8 = &fish8{}

type fish8 struct{}

func (f *fish8) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
