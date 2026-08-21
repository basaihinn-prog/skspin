package PSF_ON_00007_1_MONSTER

var Fish17 = &fish17{}

type fish17 struct{}

func (f *fish17) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish17) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}
