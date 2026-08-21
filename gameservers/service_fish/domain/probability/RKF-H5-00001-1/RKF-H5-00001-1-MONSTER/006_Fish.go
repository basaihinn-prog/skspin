package RKF_H5_00001_1_MONSTER

var Fish6 = &fish6{}

type fish6 struct{}

func (f *fish6) Hit(rtpId string, math *BsFishMath) (iconPays, triggerIconId, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
