package PSF_ON_00005_1_MONSTER

var APackOfBeer = &aPackOfBeer{}

type aPackOfBeer struct{}

func (a *aPackOfBeer) Hit(rtpId string, math *BsFruitDish) (iconPay, triggerIconId, bonusTypeId, multiplier int) {
	return FruitDish.Hit(rtpId, math)
}

func (a *aPackOfBeer) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
