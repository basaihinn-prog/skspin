package PSF_ON_00007_1_MONSTER

var Fish16 = &fish16{}

type fish16 struct{}

func (f *fish16) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish16) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}
