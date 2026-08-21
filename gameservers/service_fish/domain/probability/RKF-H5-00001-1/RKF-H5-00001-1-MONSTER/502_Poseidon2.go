package RKF_H5_00001_1_MONSTER

import (
	"os"
	"strings"
)

var Poseidon2 = &poseidon2{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "502"),
}

type poseidon2 struct {
	bonusMode bool
}

func (p *poseidon2) Hit(rtpId string, math *BsPoseidonMath) (iconPay, triggerIconId, bonusTypeId int) {
	return Poseidon1.Hit(rtpId, math)
}
