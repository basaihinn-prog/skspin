package RKF_H5_00001_1_MONSTER

import (
	"os"
	"strings"
)

var Poseidon3 = &poseidon3{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "503"),
}

type poseidon3 struct {
	bonusMode bool
}

func (p *poseidon3) Hit(rtpId string, math *BsPoseidonMath) (iconPay, triggerIconId, bonusTypeId int) {
	return Poseidon1.Hit(rtpId, math)
}
