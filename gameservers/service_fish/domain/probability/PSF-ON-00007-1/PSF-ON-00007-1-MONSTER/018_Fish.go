package PSF_ON_00007_1_MONSTER

var Fish18 = &fish18{}

type fish18 struct{}

func (f *fish18) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
