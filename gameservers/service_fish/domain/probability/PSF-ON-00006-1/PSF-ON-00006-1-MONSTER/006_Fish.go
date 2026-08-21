package PSF_ON_00006_1_MONSTER

var Fish6 = &fish6{}

type fish6 struct{}

func (f *fish6) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
