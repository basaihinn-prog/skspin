package PSF_ON_00005_1_MONSTER

import (
	"os"
	"strings"
)

var GiantWhale5 = &giantWhale5{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "505"),
}

type giantWhale5 struct {
	bonusMode bool
}

func (g *giantWhale5) Hit(rtpId string, math *BsGiantWhaleMath) (iconPay, triggerIconId, bonusTypeId int) {
	return GiantWhale1.Hit(rtpId, math)
}
