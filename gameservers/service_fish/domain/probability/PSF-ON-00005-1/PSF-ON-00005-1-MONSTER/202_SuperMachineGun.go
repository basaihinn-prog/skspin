package PSF_ON_00005_1_MONSTER

var SuperMachineGun = &superMachineGun{}

type superMachineGun struct {
}

func (s *superMachineGun) UseRtp(math *BsMachineGunMath) (rtpId string) {
	return MachineGun.UseRtp(math)
}

func (s *superMachineGun) Hit(rtpId string, math *BsMachineGunMath) (iconPay, triggerIconId, bonusTypeId, bonusTimes int) {
	return MachineGun.Hit(rtpId, math)
}
