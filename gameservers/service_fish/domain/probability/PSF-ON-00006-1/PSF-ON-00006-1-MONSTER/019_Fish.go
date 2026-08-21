package PSF_ON_00006_1_MONSTER

var Fish19 = &fish19{}

type fish19 struct{}

func (f *fish19) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
