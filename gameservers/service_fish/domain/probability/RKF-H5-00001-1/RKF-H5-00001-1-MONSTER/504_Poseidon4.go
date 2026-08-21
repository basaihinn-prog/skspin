package RKF_H5_00001_1_MONSTER

import (
	"os"
	"strings"
)

var Poseidon4 = &poseidon4{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "504"),
}

type poseidon4 struct {
	bonusMode bool
}

func (p *poseidon4) Hit(rtpId string, math *BsPoseidonMath) (iconPay, triggerIconId, bonusTypeId int) {
	return Poseidon1.Hit(rtpId, math)
}
