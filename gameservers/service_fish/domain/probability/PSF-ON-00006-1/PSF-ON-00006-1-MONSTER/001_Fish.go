package PSF_ON_00006_1_MONSTER

var Fish1 = &fish1{}

type fish1 struct{}

func (f *fish1) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
