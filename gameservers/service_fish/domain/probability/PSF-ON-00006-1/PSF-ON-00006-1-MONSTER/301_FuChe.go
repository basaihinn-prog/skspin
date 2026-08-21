package PSF_ON_00006_1_MONSTER

var FuChe = &fuChe{}

type fuChe struct{}

func (t *fuChe) Hit(rtpId string, math *BsMagicShell) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	return MagicShell.Hit(rtpId, math)
}

func (t *fuChe) HitFs(math *FsFishMath, payIndex int) (iconPays int) {
	return Fish0.HitFs(math, payIndex)
}
