package PSF_ON_00005_1_MONSTER

import (
	"os"
	"strings"
)

var GiantWhale4 = &giantWhale4{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "504"),
}

type giantWhale4 struct {
	bonusMode bool
}

func (g *giantWhale4) Hit(rtpId string, math *BsGiantWhaleMath) (iconPay, triggerIconId, bonusTypeId int) {
	return GiantWhale1.Hit(rtpId, math)
}
