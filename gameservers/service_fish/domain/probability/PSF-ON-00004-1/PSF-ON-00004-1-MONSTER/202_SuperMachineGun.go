package PSF_ON_00004_1_MONSTER

import (
	"os"
	"strings"
)

var SuperMachineGun = &superMachineGun{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "202"),
}

type superMachineGun struct {
	bonusMode bool
}

func (m *superMachineGun) Hit(rtpId string, math *BsMachineGunMath) (iconPay, bonusTimes, avgPay int) {
	return MachineGun.Hit(rtpId, math)
}

func (m *superMachineGun) Rtp(order int, math *BsMachineGunMath) (rtpId string) {
	return MachineGun.Rtp(order, math)
}

func (m *superMachineGun) UseRtp(math *BsMachineGunMath) (rtpId string) {
	return MachineGun.UseRtp(math)
}

func (m *superMachineGun) AvgPay(math *BsMachineGunMath) int {
	return MachineGun.AvgPay(math)
}
