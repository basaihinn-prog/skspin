package PSF_ON_00006_1_MONSTER

var Fish11 = &fish11{}

type fish11 struct{}

func (f *fish11) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish11) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}
