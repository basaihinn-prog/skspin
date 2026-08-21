package PSF_ON_00006_1_MONSTER

var Fish2 = &fish2{}

type fish2 struct{}

func (f *fish2) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
