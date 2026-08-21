package RKF_H5_00001_1_MONSTER

var Fish1 = &fish1{}

type fish1 struct{}

func (f *fish1) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
