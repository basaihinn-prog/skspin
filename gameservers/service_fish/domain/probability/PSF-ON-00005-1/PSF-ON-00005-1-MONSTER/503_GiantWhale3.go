package PSF_ON_00005_1_MONSTER

import (
	"os"
	"strings"
)

var GiantWhale3 = &giantWhale3{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "503"),
}

type giantWhale3 struct {
	bonusMode bool
}

func (g *giantWhale3) Hit(rtpId string, math *BsGiantWhaleMath) (iconPay, triggerIconId, bonusTypeId int) {
	return GiantWhale1.Hit(rtpId, math)
}
