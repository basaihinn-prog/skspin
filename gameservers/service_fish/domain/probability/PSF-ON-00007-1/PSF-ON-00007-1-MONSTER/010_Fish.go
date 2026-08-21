package PSF_ON_00007_1_MONSTER

var Fish10 = &fish10{}

type fish10 struct{}

func (f *fish10) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
