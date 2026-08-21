package PSF_ON_00007_1_MONSTER

var Fish8 = &fish8{}

type fish8 struct{}

func (s *fish8) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
