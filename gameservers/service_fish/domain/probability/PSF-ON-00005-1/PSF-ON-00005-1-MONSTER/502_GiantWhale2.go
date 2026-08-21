package PSF_ON_00005_1_MONSTER

import (
	"os"
	"strings"
)

var GiantWhale2 = &giantWhale2{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "502"),
}

type giantWhale2 struct {
	bonusMode bool
}

func (g *giantWhale2) Hit(rtpId string, math *BsGiantWhaleMath) (iconPay, triggerIconId, bonusTypeId int) {
	return GiantWhale1.Hit(rtpId, math)
}
