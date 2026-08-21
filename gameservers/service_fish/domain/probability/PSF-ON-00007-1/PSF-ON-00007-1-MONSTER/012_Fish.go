package PSF_ON_00007_1_MONSTER

var Fish12 = &fish12{}

type fish12 struct{}

func (f *fish12) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
