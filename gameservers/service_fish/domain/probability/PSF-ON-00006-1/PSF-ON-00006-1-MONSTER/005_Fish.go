package PSF_ON_00006_1_MONSTER

var Fish5 = &fish5{}

type fish5 struct{}

func (f *fish5) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}
