package RKF_H5_00001_1_MONSTER

import (
	"os"
	"strings"
)

var Poseidon5 = &poseidon5{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "505"),
}

type poseidon5 struct {
	bonusMode bool
}

func (p *poseidon5) Hit(rtpId string, math *BsPoseidonMath) (iconPay, triggerIconId, bonusTypeId int) {
	return Poseidon1.Hit(rtpId, math)
}
