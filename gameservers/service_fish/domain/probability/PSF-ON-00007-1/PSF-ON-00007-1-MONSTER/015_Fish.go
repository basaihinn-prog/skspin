package PSF_ON_00007_1_MONSTER

var Fish15 = &fish15{}

type fish15 struct{}

func (f *fish15) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return Fish0.Hit(rtpId, math)
}

func (f *fish15) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}
