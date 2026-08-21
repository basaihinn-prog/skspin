package RKF_H5_00001_1_MONSTER

var Fish12 = &fish12{}

type fish12 struct{}

func (f *fish12) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
