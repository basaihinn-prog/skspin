package PSF_ON_00006_1_MONSTER

var Fish3 = &fish3{}

type fish3 struct{}

func (f *fish3) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
