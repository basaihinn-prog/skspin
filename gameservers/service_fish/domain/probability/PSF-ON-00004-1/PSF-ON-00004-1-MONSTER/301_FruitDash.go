package PSF_ON_00004_1_MONSTER

import (
	"os"
	"strings"
)

var FruitDash = &fruitDash{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "301"),
}

type fruitDash struct {
	bonusMode bool
}

func (f *fruitDash) Hit(rtpId string, math *BsDashMath) (avgPay, multiplier int) {
	return LobsterDash.Hit(rtpId, math)
}

func (f *fruitDash) HitFs(math *FsFishMath) (iconPays int) {
	return Clam.HitFs(math)
}

func (f *fruitDash) Rtp(order int, math *BsDashMath) (rtpId string) {
	return LobsterDash.Rtp(order, math)
}

func (f *fruitDash) AvgPay(math *BsDashMath) int {
	return LobsterDash.AvgPay(math)
}
