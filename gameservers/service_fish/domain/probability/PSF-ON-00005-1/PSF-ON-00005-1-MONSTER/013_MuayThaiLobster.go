package PSF_ON_00005_1_MONSTER

var MuayThaiLobster = &muaythailobster{}

type muaythailobster struct{}

func (m *muaythailobster) Hit(rtpId string, math *BsFishMath) (iconPays int, triggerIconId int, bonusTypeId int) {
	return StarFish.Hit(rtpId, math)
}

func (m *muaythailobster) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
