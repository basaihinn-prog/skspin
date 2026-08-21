package PSF_ON_00007_1_MONSTER

var Fish9 = &fish9{}

type fish9 struct{}

func (f *fish9) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish9) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}
