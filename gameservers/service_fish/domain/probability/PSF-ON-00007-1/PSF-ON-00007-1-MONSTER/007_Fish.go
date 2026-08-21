package PSF_ON_00007_1_MONSTER

var Fish7 = &fish7{}

type fish7 struct{}

func (f *fish7) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
